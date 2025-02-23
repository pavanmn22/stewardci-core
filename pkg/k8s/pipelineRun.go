package k8s

import (
	"context"
	"fmt"
	"net/url"
	"time"

	api "github.com/SAP/stewardci-core/pkg/apis/steward/v1alpha1"
	stewardv1alpha1 "github.com/SAP/stewardci-core/pkg/client/clientset/versioned/typed/steward/v1alpha1"
	"github.com/SAP/stewardci-core/pkg/metrics"
	utils "github.com/SAP/stewardci-core/pkg/utils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	klog "k8s.io/klog/v2"
)

// PipelineRun is a wrapper for the K8s PipelineRun resource
type PipelineRun interface {
	fmt.Stringer
	GetAPIObject() *api.PipelineRun
	GetStatus() *api.PipelineStatus
	GetSpec() *api.PipelineSpec
	GetName() string
	GetKey() string
	GetRunNamespace() string
	GetAuxNamespace() string
	GetNamespace() string
	GetPipelineRepoServerURL() (string, error)
	HasDeletionTimestamp() bool
	AddFinalizer(ctx context.Context) error
	CommitStatus(ctx context.Context) ([]*api.StateItem, error)
	DeleteFinalizerIfExists(ctx context.Context) error
	InitState() error
	UpdateState(api.State, metav1.Time) error
	UpdateResult(api.Result, metav1.Time)
	UpdateContainer(*corev1.ContainerState)
	StoreErrorAsMessage(error, string) error
	UpdateRunNamespace(string)
	UpdateAuxNamespace(string)
	UpdateMessage(string)
}

type pipelineRun struct {
	client          stewardv1alpha1.PipelineRunInterface
	apiObj          *api.PipelineRun
	copied          bool
	changes         []changeFunc
	commitRecorders []commitRecorderFunc
}

type changeFunc func(*api.PipelineStatus) (commitRecorderFunc, error)
type commitRecorderFunc func() *api.StateItem

// NewPipelineRun creates a managed pipeline run object.
// If a factory is provided a new version of the pipelinerun is fetched.
// All changes are done on the fetched object.
// If no pipeline run can be found matching the apiObj, nil,nil is returned.
// An error is only returned if a Get for the pipelinerun returns an error other than a NotFound error.
// If you call with factory nil you can only use the Get* functions
// If you use functions changing the pipeline run without factroy set you will get an error.
// The provided PipelineRun object is never modified and copied as late as possible.
func NewPipelineRun(ctx context.Context, apiObj *api.PipelineRun, factory ClientFactory) (PipelineRun, error) {
	if factory == nil {
		return &pipelineRun{
			apiObj: apiObj,
			copied: false,
		}, nil
	}
	client := factory.StewardV1alpha1().PipelineRuns(apiObj.GetNamespace())
	obj, err := client.Get(ctx, apiObj.GetName(), metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &pipelineRun{
		apiObj:          obj,
		copied:          true,
		client:          client,
		changes:         []changeFunc{},
		commitRecorders: []commitRecorderFunc{},
	}, nil
}

// GetAPIObject implements part of interface `PipelineRun`.
func (r *pipelineRun) GetAPIObject() *api.PipelineRun {
	return r.apiObj
}

// GetRunNamespace returns the namespace in which the build takes place
func (r *pipelineRun) GetRunNamespace() string {
	return r.apiObj.Status.Namespace
}

// GetAuxNamespace returns the namespace hosting auxiliary services
// for the pipeline run.
func (r *pipelineRun) GetAuxNamespace() string {
	return r.apiObj.Status.AuxiliaryNamespace
}

// GetKey returns the key of the pipelineRun
func (r *pipelineRun) GetKey() string {
	key, _ := cache.MetaNamespaceKeyFunc(r.apiObj)
	return key
}

// GetNamespace returns the namespace of the underlying pipelineRun object
func (r *pipelineRun) GetNamespace() string {
	return r.apiObj.GetNamespace()
}

// GetPipelineRepoServerURL returns the server hosting the Jenkinsfile repository
func (r *pipelineRun) GetPipelineRepoServerURL() (string, error) {
	urlString := r.GetSpec().JenkinsFile.URL
	repoURL, err := url.Parse(urlString)
	if err != nil {
		return "", errors.Wrapf(err, "value %q of field spec.jenkinsFile.url is invalid [%s]", urlString, r.String())
	}
	if !(repoURL.Scheme == "http") && !(repoURL.Scheme == "https") {
		return "", fmt.Errorf("value %q of field spec.jenkinsFile.url is invalid [%s]: scheme not supported: %q", urlString, r.String(), repoURL.Scheme)
	}
	return fmt.Sprintf("%s://%s", repoURL.Scheme, repoURL.Host), nil
}

func (r *pipelineRun) GetName() string {
	return r.apiObj.GetName()
}

// GetStatus return the Status
// the returned PipelineStatus MUST NOT be modified
// use the prodided Update* functions instead
func (r *pipelineRun) GetStatus() *api.PipelineStatus {
	return &r.apiObj.Status
}

// GetSpec return the spec part of the PipelineRun resource
// the returned PipelineSpec MUST NOT be modified
func (r *pipelineRun) GetSpec() *api.PipelineSpec {
	return &r.apiObj.Spec
}

// InitState initializes the state as 'new' if state was undefined (empty) before.
// The state's start time will be set to the object's creation time.
// Fails if a state is set already.
func (r *pipelineRun) InitState() error {
	r.ensureCopy()
	klog.V(3).Infof("Set State to New [%s]", r.String())
	return r.changeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {

		if s.State != api.StateUndefined {
			return nil, fmt.Errorf("Cannot initialize multiple times")
		}

		newStateDetails := api.StateItem{
			State:     api.StateNew,
			StartedAt: r.apiObj.ObjectMeta.CreationTimestamp,
		}
		s.StateDetails = newStateDetails
		s.State = api.StateNew
		return nil, nil
	})
}

// UpdateState set end time of current (defined) state (A) and store it to the history.
// if no current state is defined a new state (A) with cretiontime of the pipelinerun as start time is created.
// It also creates a new current state (B) with start time.
func (r *pipelineRun) UpdateState(state api.State, ts metav1.Time) error {
	if r.apiObj.Status.State == api.StateUndefined {
		if err := r.InitState(); err != nil {
			return err
		}
	}
	r.ensureCopy()
	klog.V(3).Infof("Update State to %s [%s]", state, r.String())
	oldStateDetails := r.apiObj.Status.StateDetails

	return r.changeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {
		currentStateDetails := s.StateDetails
		if currentStateDetails.State != oldStateDetails.State {
			return nil, fmt.Errorf("State cannot be updated as it was changed concurrently from %q to %q", oldStateDetails.State, currentStateDetails.State)
		}
		if state == api.StatePreparing {
			s.StartedAt = &ts
		}
		currentStateDetails.FinishedAt = ts
		his := s.StateHistory
		his = append(his, currentStateDetails)

		commitRecorderFunc := func() *api.StateItem {
			return &currentStateDetails
		}
		newStateDetails := api.StateItem{State: state, StartedAt: ts}
		if state == api.StateFinished {
			newStateDetails.FinishedAt = ts
		}

		s.StateDetails = newStateDetails
		s.StateHistory = his
		s.State = state
		return commitRecorderFunc, nil
	})
}

// String returns the full qualified name of the pipeline run
func (r *pipelineRun) String() string {
	return fmt.Sprintf("PipelineRun{name: %s, namespace: %s, state: %s}", r.GetName(), r.GetNamespace(), string(r.GetStatus().State))
}

// UpdateResult of the pipeline run
func (r *pipelineRun) UpdateResult(result api.Result, ts metav1.Time) {
	r.ensureCopy()
	r.mustChangeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {
		s.Result = result
		s.FinishedAt = &ts
		return nil, nil
	})
}

// UpdateContainer ...
func (r *pipelineRun) UpdateContainer(c *corev1.ContainerState) {
	if c == nil {
		return
	}
	r.ensureCopy()
	r.mustChangeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {
		s.Container = *c
		return nil, nil
	})
}

// StoreErrorAsMessage stores the error as message in the status
func (r *pipelineRun) StoreErrorAsMessage(err error, message string) error {
	if err != nil {
		text := fmt.Sprintf("ERROR: %s [%s]: %s", utils.Trim(message), r.String(), err.Error())
		klog.V(3).Infof(text)
		r.UpdateMessage(text)
	}
	return nil
}

// UpdateMessage stores string as message in the status
func (r *pipelineRun) UpdateMessage(message string) {
	r.ensureCopy()

	r.mustChangeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {
		old := s.Message
		if old != "" {
			his := s.History
			his = append(his, old)
			s.History = his
		}
		s.Message = utils.Trim(message)
		s.MessageShort = utils.ShortenMessage(message, 100)
		return nil, nil
	})
}

// UpdateRunNamespace overrides the namespace in which the builds happens
func (r *pipelineRun) UpdateRunNamespace(ns string) {
	r.ensureCopy()
	r.mustChangeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {
		s.Namespace = ns
		return nil, nil
	})
}

// UpdateAuxNamespace overrides the namespace hosting auxiliary services
// for the pipeline run.
func (r *pipelineRun) UpdateAuxNamespace(ns string) {
	r.ensureCopy()
	r.mustChangeStatusAndStoreForRetry(func(s *api.PipelineStatus) (commitRecorderFunc, error) {
		s.AuxiliaryNamespace = ns
		return nil, nil
	})
}

//HasDeletionTimestamp returns true if deletion timestamp is set
func (r *pipelineRun) HasDeletionTimestamp() bool {
	return !r.apiObj.ObjectMeta.DeletionTimestamp.IsZero()
}

// AddFinalizer adds a finalizer to pipeline run
func (r *pipelineRun) AddFinalizer(ctx context.Context) error {
	changed, finalizerList := utils.AddStringIfMissing(r.apiObj.ObjectMeta.Finalizers, FinalizerName)
	if changed {
		err := r.updateFinalizers(ctx, finalizerList)
		return err
	}
	return nil
}

// DeleteFinalizerIfExists deletes a finalizer from pipeline run
func (r *pipelineRun) DeleteFinalizerIfExists(ctx context.Context) error {
	changed, finalizerList := utils.RemoveString(r.apiObj.ObjectMeta.Finalizers, FinalizerName)
	if changed {
		return r.updateFinalizers(ctx, finalizerList)
	}
	return nil
}

func (r *pipelineRun) updateFinalizers(ctx context.Context, finalizerList []string) error {
	if len(r.changes) > 0 {
		return fmt.Errorf("cannot add finalizers when we have uncommited status updates")
	}
	if r.client == nil {
		panic(fmt.Errorf("No factory provided to store updates [%s]", r.String()))
	}
	r.ensureCopy()
	start := time.Now()
	r.apiObj.ObjectMeta.Finalizers = finalizerList
	result, err := r.client.Update(ctx, r.apiObj, metav1.UpdateOptions{})
	end := time.Now()
	elapsed := end.Sub(start)
	klog.V(4).Infof("finish update finalizer after %s in %s", elapsed, r.apiObj.Name)
	if err != nil {
		return errors.Wrap(err,
			fmt.Sprintf("Failed to update finalizers [%s]", r.String()))
	}
	r.apiObj = result
	return nil
}

// mustChangeStatusAndStoreForRetry calls changeStatusAndStoreForRetry and
// panics in case of an error.
func (r *pipelineRun) mustChangeStatusAndStoreForRetry(change changeFunc) {
	err := r.changeStatusAndStoreForRetry((change))
	if err != nil {
		panic(err)
	}
}

// changeStatusAndStoreForRetry receives a function applying changes to pipelinerun.Status
// This function get executed on the current memory representation of the pipeline run
// and remembered so that it can be re-applied later in case of a re-try. The change function
// must only apply changes to pipelinerun.Status.
//
func (r *pipelineRun) changeStatusAndStoreForRetry(change changeFunc) error {
	commitRecorder, err := change(r.GetStatus())
	if err == nil {
		r.changes = append(r.changes, change)
		r.commitRecorders = append(r.commitRecorders, commitRecorder)
	}

	return err
}

// CommitStatus executes `change` and writes the
// status of the underlying PipelineRun object to storage afterwards.
// `change` is expected to mutate only the status of the underlying
// object, not more.
// In case of a conflict (object in storage is different version than
// ours), the update is retried with backoff:
//     - wait
//     - fetch object from storage
//     - run `change`
//     - write object status to storage
// After too many conflicts retrying is aborted, in which case an
// error is returned.
// Non-conflict errors are returned without retrying.
//
// Pitfall: If the underlying PipelineRun object was changed in memory
// compared to the version in storage _before calling this function_,
// that change _gets_ persisted in case there's _no_ update conflict, but
// gets _lost_ in case there _is_ an update conflict! This is hard to find
// by tests, as those typically do not encounter update conflicts.
func (r *pipelineRun) CommitStatus(ctx context.Context) ([]*api.StateItem, error) {
	if r.client == nil {
		panic(fmt.Errorf("No factory provided to store updates [%s]", r.String()))
	}

	klog.V(5).Infof("enter commitStatus for pipeline run %q ...", r.String())
	if len(r.changes) == 0 {
		klog.V(5).Infof("commitStatus no changes found for pipeline run %q.", r.String())
		return nil, nil
	}

	retryCount := uint64(0)
	defer func(start time.Time) {
		if retryCount > 0 {
			codeLocationSkipFrames := uint16(1)
			codeLocation := metrics.CodeLocation(codeLocationSkipFrames)
			latency := time.Since(start)
			metrics.Retries.Observe(codeLocation, retryCount, latency)
			klog.V(5).InfoS("retry was required",
				"location", codeLocation,
				"count", retryCount,
				"latency", latency,
				"namespace", r.GetNamespace(),
				"pipelineRun", r.GetName(),
			)
		}
	}(time.Now())

	var changeError error
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		var err error

		if retryCount > 0 {
			klog.V(5).Infof("commitStatus reload pipeline run for retry %q ...", r.String())
			new, err := r.client.Get(ctx, r.apiObj.GetName(), metav1.GetOptions{})
			if err != nil {
				return errors.Wrap(err,
					"failed to fetch pipeline after update conflict")
			}
			r.apiObj = new
			r.copied = true
			r.commitRecorders = []commitRecorderFunc{}
			var commitRecorder func() *api.StateItem
			klog.V(5).Infof("commitStatus applies %d change(s)", len(r.changes))
			for i, change := range r.changes {
				commitRecorder, changeError = change(r.GetStatus())
				if changeError != nil {
					klog.V(5).Infof("applying change %d failed with error: %s", i, changeError.Error())
					return nil
				}
				r.commitRecorders = append(r.commitRecorders, commitRecorder)
			}
		}

		result, err := r.client.UpdateStatus(ctx, r.apiObj, metav1.UpdateOptions{})
		if err == nil {
			r.apiObj = result
			return nil
		}
		retryCount++
		return err
	})
	r.changes = []changeFunc{}
	if changeError != nil {
		return nil, changeError
	}
	result := []*api.StateItem{}
	for _, recorder := range r.commitRecorders {
		if recorder != nil {
			result = append(result, recorder())
		}
	}
	return result, errors.Wrapf(err, "failed to update status [%s]", r.String())
}

func (r *pipelineRun) ensureCopy() {
	if !r.copied {
		r.apiObj = r.apiObj.DeepCopy()
		r.copied = true
	}
}

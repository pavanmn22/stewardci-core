module github.com/SAP/stewardci-core

go 1.16

require (
	github.com/benbjohnson/clock v1.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/lithammer/dedent v1.1.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/openzipkin/zipkin-go v0.3.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/statsd_exporter v0.22.4 // indirect
	github.com/tektoncd/pipeline v0.34.0
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v1.5.2
	k8s.io/klog/v2 v2.40.1
	knative.dev/pkg v0.0.0-20220131144930-f4b57aef0006
)

replace (
	k8s.io/api => k8s.io/api v0.21.7 // kubernetes-1.21.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.7 // kubernetes-1.21.7
	k8s.io/apiserver => k8s.io/apiserver v0.21.7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.7
	k8s.io/client-go => k8s.io/client-go v0.21.7 // kubernetes-1.21.7
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.7
	k8s.io/code-generator => k8s.io/code-generator v0.21.7
	k8s.io/component-base => k8s.io/component-base v0.21.7
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.7
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.7
	k8s.io/cri-api => k8s.io/cri-api v0.21.7
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.7
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.7
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.7
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.7
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.7
	k8s.io/kubectl => k8s.io/kubectl v0.21.7
	k8s.io/kubelet => k8s.io/kubelet v0.21.7
	k8s.io/kubernetes => k8s.io/kubernetes v1.21.7
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.7
	k8s.io/metrics => k8s.io/metrics v0.21.7
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.7
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.7
	knative.dev/pkg => knative.dev/pkg v0.0.0-20211101212339-96c0204a70dc // release-1.0
)

module github.com/SAP/stewardci-core/charts/steward/test

go 1.16

require (
	github.com/gruntwork-io/terratest v0.38.5
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.21.7 // kubernetes-1.21.7
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
)

module github.com/warroyo/nsxt-tag-controller

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.3
	github.com/vmware/vsphere-automation-sdk-go/runtime v0.3.1
	github.com/vmware/vsphere-automation-sdk-go/services/nsxt v0.5.0
	gitlab.eng.vmware.com/core-build/guest-cluster-controller v0.0.0-20210413224313-fa16808e7391
	gitlab.eng.vmware.com/core-build/nsx-ujo v0.0.0-20210412204140-b9868e395e5d
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	sigs.k8s.io/controller-runtime v0.5.14
)

replace (
	dmitri.shuralyov.com/gpu/mtl => dmitri.shuralyov.com/gpu/mtl v0.0.0-20201218220906-28db891af037
	k8s.io/api => k8s.io/api v0.17.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.8
	k8s.io/client-go => k8s.io/client-go v0.17.8
)

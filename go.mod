module kubecp

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alecthomas/units v0.0.0-20201120081800-1786d5ef83d4 // indirect
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.20.1 // indirect
	github.com/go-openapi/swag v0.19.13 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.7.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/kubernetes v1.20.2
	k8s.io/metrics v0.20.2
)

replace k8s.io/api => k8s.io/api v0.20.2

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.2

replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.0-alpha.0

replace k8s.io/apiserver => k8s.io/apiserver v0.20.2

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.2

replace k8s.io/client-go => k8s.io/client-go v0.20.2

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.2

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.2

replace k8s.io/code-generator => k8s.io/code-generator v0.20.3-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.20.2

replace k8s.io/component-helpers => k8s.io/component-helpers v0.20.2

replace k8s.io/controller-manager => k8s.io/controller-manager v0.20.2

replace k8s.io/cri-api => k8s.io/cri-api v0.20.3-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.2

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.2

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.2

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.2

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.2

replace k8s.io/kubectl => k8s.io/kubectl v0.20.2

replace k8s.io/kubelet => k8s.io/kubelet v0.20.2

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.2

replace k8s.io/metrics => k8s.io/metrics v0.20.2

replace k8s.io/mount-utils => k8s.io/mount-utils v0.20.3-rc.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.2
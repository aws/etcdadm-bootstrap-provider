module github.com/aws/etcdadm-bootstrap-provider

go 1.16

require (
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/go-logr/logr v1.2.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v3 v3.0.0 // indirect
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	sigs.k8s.io/cluster-api v1.0.1
	sigs.k8s.io/controller-runtime v0.11.1

)

replace (
	github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible
	sigs.k8s.io/cluster-api => github.com/mrajashree/cluster-api v1.1.3-custom
)

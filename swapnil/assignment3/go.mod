module swapnil/k8s-dev-training/swapnil/assignment3

go 1.13

//replace (
//	github.com/swapnil-velotio/k8s-dev-training/swapnil/assignment3 => ../helpers
//)
require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/robfig/cron v1.2.0
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)

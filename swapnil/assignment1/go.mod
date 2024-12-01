module github.com/swapnil-velotio/k8s-dev-training/swapnil/assignment1

go 1.15

replace (
	github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers => ../helpers
)
require (
	github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers v0.0.0-20210407130751-b333e51a55df
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	sigs.k8s.io/controller-runtime v0.8.3
)
//replace (
//	 github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers =>  /home/swapnil/go/src/swapnil/k8s-dev-training/swapnil/helpers
//)

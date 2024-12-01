## How to use
It is assumed that a minikube like k8s-cluster is already running and docker is installed on your machine.

1. Make sure python is installed on your machine. Run the following command on CLI.
``` python3 py_script.py ```
This will build docker image, create cluster roles and bindings and run a pod on k8s cluster. this pod will in-turn create k8s resources in the cluster.
2. 'pod-asgn' will be cerated in the k8s-cluster. 
Run ```kubectl logs pod-asgn``` to view the output of the pod.
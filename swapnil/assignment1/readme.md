# CRUD operation for configmap, pod and deployment using go client and controller runtime

solution is in dedicated folder for each resource with go client and controller runtime go file

Findings

go client seemed more easy for CRUD operation on kubernetes resources 
It provides dynamic client where we just need to provide appropriate pointer to the resource and it figures out how to do that operation on that type of resource
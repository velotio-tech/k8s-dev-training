## Markers

Markers are single line comments  that start with a plus, represent some addtional information of our implementation (like fields) and used by `controller-runtime` and `make` to make changes in CRDs for additional recommedations or restrictions.
Sample Syntax: `+path:to:marker:arg1=val,arg2=val2`
The arguments can be of strings, ints, bools, slices, or maps data type.

Marker Types:
1. Empty: This is a similar to boolean flags on CLI (--). Eg. `+kubebuilder:validation:Optional`
2. Anonymous: This take only one value as argument. Eg. `+kubebuilder:validation:MaxItems=2`
3. Multi-Option: This takes multiple arguments. Ordering is not compulsory. Eg. `+kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string`



## Subresources

Subresources are separate kind of endpoints that have a suffix appended to their path of the normal resource. For example, the pod standard HTTP path is `/api/v1/namespace/namespace/pods/name`. & status subresource can be accessed using API endpoint `/api/v1/namespace/<namespace>/pods/<name>/status`.

There are mainly two types of subresources : `/scale` and `/status`. Both are not enabled by default. We can enable the `status` subresource by `// +kubebuilder:subresource:status`. When I write this genarator (marker) within the types.go file near `status` struct and perform `make install`, the CRD gets updated with these fields:
`subresources:
    status: {}`


**When I change the fields and values in `spec` of types.go OR .yaml of CRD, if I change the contents of `status`, it will throw an error and won't change. i.e The PUT and POST verbs on objects MUST ignore the "status" values, to avoid accidentally overwriting the status in read-modify-write scenarios. A /status subresource MUST be provided to enable system components to update statuses of resources they manage.**



References:
* https://book.kubebuilder.io
* https://book-v1.book.kubebuilder.io/basics/status_subresource.html
* https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#status-subresource
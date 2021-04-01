This code has CRD named My Job. It runs job on the basics of resource Type.

Controller  has mainly defined two resources structs that contain what to be checked and what needs to be maintained. Those structs are: `spec{}` and `status{}`. 
* `spec{}` holds fields that define what the desired condition is for the resource.
* `status{}` holds fields that define what the current condition is of the resource.

A controller tries to find out what the current status of resource is and tries to match it with what mentioned in spec.

if the error is returned from `Reconcile()`, then it  will be requeued (using logging). There are also two components of a controller: Informer/SharedInformer and Workqueue.

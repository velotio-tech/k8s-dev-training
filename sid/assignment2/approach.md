## CodeSanity CR
It is responsible for maintaining a test coverage criteria across a cluster. It allows defining a test specification and specify which pods and images should be coveraged and then it's the controller's job to do the assessment everytime a new pod is created or updated.

### Spec
 - A `TestFilesRegexStr` field to specify which files should be covered.
 - An optional `RequiredCoverage` field to specify what should be required coverage criteria.
 - An optional `PodNames` field to specify which pods should be covered based on their names. By default, it will assess all the pods.
 - An optional `Images` field to specify which pods should be covered based on their images. By default, it will assess all the pods.

### Status
 - A `CoverageMap` field that will map different image's name-tag slug to it's latest calculated coverage.
 - A `HealthyPods` field that will map different healthy pod names to their calculated coverage. If no coverage criteria is present in CR's spec then all the pods whose images were able to run test command successfully will be considered healthy.
 - A `UnhealthyPods` field that will map different unhealthy pod names to their calculated coverage.
 - A `LastRunAt` field will tell when the latest assessment was made.

The assumption here is that every image involved can be invoked with a `--test` flag along with files regex and coverage to run tests on the said image. This can be mocked for the sake of this exercise.

### Controller Logic

 - When a new CR is created, cycle through all the pods and spawn jobs to run tests and populate different fields in the status. Attach functions on Pods for create, update and delete using shared informer.
 - When a new Pod is created, controller will spawn an new k8s job to run tests and update the CR's status.
 - When a pod is deleted, controller will just updated it's status to remove the pod's data as well as image it there are no pods corresponding to it.
 - When a pod is updated, controller will check of the image was updated and spawn a new job to run tests if the coverage information for this new image was not present in `CoverageMap` field in CR's status.
 - For every run, it will update the `LastRunAt` field.
 - When an existing CR is updated, let's say coverage criteria was changes. Then, rerun the jobs for different images and updated the status accordingly.

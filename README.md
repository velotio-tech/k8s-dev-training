<<<<<<< HEAD
# k8s-assignment2
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/k8s-assignment2:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/k8s-assignment2:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

=======
# k8s-dev-training
Training programme covering development exercise for kubernetes native development


## Assignment 1

Aim: A small task to learn how to create the clients and their usages with `client-go` and `controller-runtime` and their differences

Requirements:

- Create the clients for any 3 k8s resources using incluster config

- perform CRUD operations

- Similarly, explore controller-runtime library and perform CRUD using it

- At the end write up what non-general things you learn while exploring the library. It can be any new utility functions, any new efficient ways to do something or may be some helpful blog around it


## Assignment 2

Aim: Learn to create a CRD (k8s API) using kubebuilder 2.0 and implement the validations, subResources

- Define 1 API of your own.

- Define a logical set of spec and status fields in the CRD, atleast 3-4 fields in both of them.

- You should use the markers to have the proper restrictions in each field (Explore the markers and their types and make use of a few of them atleast)
  use the `controller-gen` tools to generate the crd yaml for the same.
  
- Do some tryouts by changing the fields and adding the fields of the custom resource of that CRD, with both positive and negative scenarios wrt the restrictions and field types

- Make the `status` as Subresource and try to update status by putting it in yaml, it won't, try to read up on why is it not possible

- Post what is the subresources and why you cannot use the general CRUD on them as you can do on `spec`


## Assignment 3

Aim: Learn to create a controller using the API and manage its reconcilation based on the predicates, ownerships and field indexers

- Explore and learn about the custom controllers

- Understand the architecture of `shared-informers` in extension k8s

- In the API from `Assignemt 2`, have a field in the spec which will take the input as GVK and their names.

- Make controller for the above CRD which should read the crd and make the resource as specified. Do the crud on the generated resource and make sure that the events are handled using predicates in such a way that the current state is always matched with the desired state as specified in the custom resource

- Also do the crud on the spec of your custom resource by adding/removing/changing the resource gvk field as well as other fields and make sure the reconcile works. Also handle the  events  for the custom resource also.

- Make sure you also push the owner reference on the created resource
  
## Assignment 4

Aim: Learn to perform the reconcilation in a bit advanced set of restrictions. Such as in one namespace, CRUD on children based on owner references, also on the children with only field indexers, 3 level ownerships.

- Using the outcome of `Assignment 3`, create a new API or a new version of the older API and reconcile it in such a way that controller should create resources with 3 levels.
Ex: Your API should create a deployment, it will create a pod (after replica set) and add a script in pod command to create a `Job` resource

- Perform CRUD on all the levels of resources. Any one level of resource should not have the `OwnerReference` on it. It can be the first one or the last one.

- Use Owner reference to get the child resources and use field indexer to get the resource which donot have owner ref and then perform some meaningful operations

- Your controller should be reconciling the resources only in its own namespace not from any other. Use predicate here.  


## Assignment 5

Aim: To learn about the Qemu-img utility and datamover

- Explore the `qemu-img` utility and `qcow2` image format. Understand What, Why and How of this utility.

- Try out a few `qemu` commands, check for their sizes. You should be able to tell the difference between `Virtual Size` and `Disk Size`

- Try out the `libguestfs` library commands being used in datamover. To create qcow2 images from block device, create qcow2 from file system

- How to restore the data from qcow2 images 

- Significance of `Bacing Chain` and how to explicitly put a chain to an image as well as how to create a new `Diff` image from 2 sources of data and referring another image using the chain.
Use `-D` flag in the qemu-img utility command.


### General Requirement

- Comment all the things you are doing

- Code should be structured even though it is an assignment

- write up a simple summary of what you tried to do and how did you do in the commit comments


### Submission

- Fork this repository and work on it

- Raise a PR against this repository

- Add/ Tag the reviewer for that PR
>>>>>>> 9d82fa2aee09f59afac870b27289415abcb46f35

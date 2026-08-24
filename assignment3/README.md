## Assignment 3

In Assignment 2, we built the *brain blueprint* (the CRD). Now, in Assignment 3, you build the *engine* (the Controller) that watches your Custom Resource and makes real things happen in your Kubernetes cluster.

## The Big Picture: What Are You Building in Assignment 3?

In this assignment, you will build an **Operator/Controller** for your `AppService` CRD.

Specifically, you will:

1. Extend `AppService` so users can specify desired Kubernetes resources dynamically via Group-Version-Kind (GVK) and resource names.
2. Build a **Reconciler Loop** that reads the CRD, checks what resources exist, and performs CRUD operations to sync reality with the desired state.
3. Attach **Owner References** so Kubernetes knows your controller owns the generated resources.
4. Filter cluster events using **Predicates** so your controller only reacts when relevant changes occur.
5. Understand the underlying engine powering all of this: **SharedInformers**.

# Deep Dive: Assignment 3

## 1. Core Concepts Explained

### A. The Reconcile Loop (`Reconcile`)

The heart of any Kubernetes Operator is the **Reconciliation Loop**. Its job is simple: **Compare Desired State (`spec`) with Actual State (the cluster), and make changes until they match.**

$$\text{Desired State } (\text{spec}) \quad \Longleftrightarrow \quad \text{Actual State } (\text{cluster})$$

The `Reconcile` function receives a request with just a `NamespacedName` (namespace and name of your object). It must be **idempotent**—meaning running it once or running it 100 times with the same state produces the exact same cluster outcome.

### B. SharedInformers & Caches

How does a controller know when something changes without hammering the Kubernetes API server with constant HTTP polling? **SharedInformers**.

1. **List-Watch Mechanism:** An Informer connects to the K8s API server using a `List` call on startup, followed by an HTTP `Watch` stream for updates (Added, Updated, Deleted).
2. **Local Cache (IndexerStore):** The Informer stores all observed objects in a fast, in-memory local cache.
3. **Shared Memory:** Multiple controllers inside the same process share the same Informer cache (hence *Shared*Informer) to save RAM and network bandwidth.

When you execute `r.Get()` or `r.List()` inside a `controller-runtime` reconciler, you are reading from this fast local cache, **not** making a network call to `etcd`!

### C. Predicates (Event Filtering)

By default, whenever a watched resource is created, updated, or deleted, an event triggers the `Reconcile` loop. **Predicates** act as filters/guards before an event hits your controller logic.

For example, if a Pod's `resourceVersion` or `metadata.managedFields` changes, K8s triggers an Update event. A `GenerationChangedPredicate` filters out noise so your controller *only* runs when the `.spec` actually changes, ignoring harmless background status updates.

### D. Owner References (`SetControllerReference`)

When your controller creates a child resource (like a `Deployment` or `Service`), you must attach an **Owner Reference** pointing back to your `AppService` Custom Resource.

This serves two crucial roles:

1. **Garbage Collection (Cascading Delete):** When a user deletes the parent `AppService`, Kubernetes automatically deletes all child resources.
2. **Event Routing:** When a child resource changes (e.g., someone manually deletes a managed `Service`), the controller framework uses the Owner Reference to identify the parent `AppService` and enqueue a Reconcile request for it.

## 2. Architecture & Event Flow Diagram

```
[ K8s API Server ]
       │
       │ HTTP Watch Stream (Events: Create/Update/Delete)
       ▼
┌─────────────────────────────────────────────────────────────┐
│                       SharedInformer                        │
│  ┌─────────────────────────┐   ┌─────────────────────────┐  │
│  │ Local In-Memory Cache   │   │  Reflector & Delta FIFO │  │
│  └─────────────────────────┘   └─────────────────────────┘  │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      Event Predicates                       │
│       (Filter out unneeded updates / spec changes only)     │
└──────────────────────────────┬──────────────────────────────┘
                               │ (Passes Filter)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                        WorkQueue                            │
│           (Queues NamespacedName: "default/my-app")         │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                     Reconciler Loop                         │
│  1. Fetch AppService from Local Cache                       │
│  2. Read desired GVK/Resources from spec                    │
│  3. Check actual state in cluster                           │
│  4. Perform CRUD to reach desired state                     │
│  5. Set OwnerReference on child resources                   │
│  6. Update AppService.status                                │
└─────────────────────────────────────────────────────────────┘
```

## 3. Step-by-Step Implementation Guide

## Standalone Directory Layout

```bash
k8s-dev-training/
├── assignment1/
├── assignment2/
└── assignment3/
    ├── Makefile
    ├── PROJECT
    ├── go.mod
    ├── go.sum
    ├── main.go / cmd/main.go
    ├── api/
    │   └── v1alpha1/
    │       ├── dynamicapp_types.go
    │       ├── groupversion_info.go
    │       └── zz_generated.deepcopy.go
    ├── config/
    │   ├── crd/
    │   ├── rbac/
    │   └── samples/
    │       └── apps_v1alpha1_dynamicapp.yaml
    └── internal/
        └── controller/
            └── dynamicapp_controller.go
```

## Step 3.1: Initialize the Standalone Assignment 3 Project

Navigate to your workspace root and initialize Kubebuilder inside a fresh `assignment3` folder:

```bash
cd k8s-dev-training
mkdir -p assignment3
cd assignment3

# Initialize new Kubebuilder project
kubebuilder init --domain example.com --repo github.com/your-username/k8s-dev-training/assignment3

# Create a brand new API for Assignment 3
kubebuilder create api --group apps --version v1alpha1 --kind DynamicApp --resource=true --controller=true
```

## Step 3.2: Define the Schema (`api/v1alpha1/dynamicapp_types.go`)

Open `api/v1alpha1/dynamicapp_types.go` and replace its contents with the schema This defines target resources (GVKs and names) for dynamic management.


## Step 3.3: Implement the Reconciler & Predicates (`internal/controller/dynamicapp_controller.go`)

Open `internal/controller/dynamicapp_controller.go` and replace its contents with the reconciliation loop logic:


## Step 3.4: Generate Code & Manifests

Execute the generator commands inside the `assignment3` folder:

```
# Generate deepcopy methods
make generate

# Generate CRD YAML manifests
make manifests

# Install Assignment 3 CRD into your active kind cluster
make install
```

Verify CRD installation:

```
kubectl get crd dynamicapps.apps.example.com
```

## Step 3.5: Test & Validate Execution

### 1. Run Controller Locally

In Terminal 1:

```
make run
```

### 2. Create Sample Custom Resource

In Terminal 2, create `config/samples/apps_v1alpha1_dynamicapp.yaml`:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: DynamicApp
metadata:
  name: demo-dynamic-app
  namespace: default
spec:
  replicas: 2
  targetResources:
    - group: ""
      version: "v1"
      kind: "ConfigMap"
      name: "dynamic-cm-sample"
```

Apply the Custom Resource:

```
kubectl apply -f config/samples/apps_v1alpha1_dynamicapp.yaml
```

### 3. Verify Child Creation & Owner Reference

Check if the controller created the child ConfigMap:

```
kubectl get configmap dynamic-cm-sample -o yaml
```

**Expected Output Snippet:**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: dynamic-cm-sample
  namespace: default
  ownerReferences:
  - apiVersion: apps.example.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: DynamicApp
    name: demo-dynamic-app
```

### 4. Test Cascading Deletion

Delete the parent resource and verify automatic cleanup:

```
kubectl delete dynamicapp demo-dynamic-app
kubectl get configmap dynamic-cm-sample
```

**Expected Output:** `Error from server (NotFound): configmaps "dynamic-cm-sample" not found`

### **Aim**

Build a control loop (Reconciler) that continuously watches your Custom Resource (`DynamicApp`), automatically provisions target child resources (e.g., `ConfigMap`), attaches **Owner References**, and uses **Predicates** to filter unnecessary events.

### **Files Created**

```bash
assignment3/
├── api/v1alpha1/
│   └── dynamicapp_types.go          <-- CRD schema defining target GVK array
├── config/crd/bases/
│   └── apps.example.com_dynamicapps.yaml <-- Generated CRD blueprint
├── config/samples/
│   └── apps_v1alpha1_dynamicapp.yaml <-- Sample CR defining target child resources
└── internal/controller/
    └── dynamicapp_controller.go     <-- Reconcile loop logic & Predicate filters
```

### **Execution Commands**

```
# 1. Initialize project with controller logic enabled
kubebuilder init --domain example.com --repo github.com/your-username/assignment3
kubebuilder create api --group apps --version v1alpha1 --kind DynamicApp --resource=true --controller=true

# 2. Build code & install CRD
make generate && make manifests
make install

# 3. Run controller locally against your kind cluster
make run

# 4. In a separate terminal, apply the CR
kubectl apply -f config/samples/apps_v1alpha1_dynamicapp.yaml
```

### **End-to-End Reconciliation Flow**

```bash
 1. User applies CR (apps_v1alpha1_dynamicapp.yaml)
                        │
                        ▼
 2. API Server receives CR ──► 3. SharedInformer receives watch event
                                              │
                                              ▼
                                 4. Custom Predicate filters event
                                 (Checks if metadata.generation changed)
                                              │
                                              ▼
                                 5. Reconcile Loop Executes
                                 (internal/controller/dynamicapp_controller.go)
                                              │
                                              ├── Reads CR spec.targetResources
                                              ├── Checks if child ConfigMap exists
                                              ├── Creates missing ConfigMap
                                              ├── Injects OwnerReference into ConfigMap
                                              └── Updates CR.status.phase = "Running"
                                              │
                                              ▼
 6. If CR is deleted ────────► 7. K8s Garbage Collector deletes child ConfigMap
                                 (Enabled by Owner Reference)
```

#### **Term Glossary:**

- **Controller:** A loop that reads desired state (`spec`) from a CR, checks actual cluster state, and performs CRUD operations until desired state matches actual state.
- **SharedInformer:** Background worker mechanism that maintains a fast local in-memory cache of cluster objects synced via HTTP watch streams.
- **Owner Reference:** Metadata attached to a child object pointing back to its parent CR. Enables **cascading deletion** (deleting parent automatically cleans up child) and event routing.
- **Predicate:** Event filter logic used to prevent unnecessary reconcile loops (e.g., ignoring status-only update events).

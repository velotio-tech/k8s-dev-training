# Assignment 4

We will step up from basic single-child reconciliation to building production-grade operator logic handling **multi-tiered hierarchies**, **fast lookups using indexers**, and **strict multi-tenant namespace scoping**.

## The Big Picture: What Are You Building in Assignment 4?

In this assignment, you will extend your controller to manage a **3-Level Cascading Hierarchy** while introducing advanced `controller-runtime` patterns:

1. **3-Level Resource Chain:**
    - **Level 1 (Parent):** Your Custom Resource (`MultiTierApp`).
    - **Level 2 (Child):** A standard Kubernetes `Deployment` owned by the `MultiTierApp` (Level 1 $\rightarrow$ Level 2 Owner Reference).
    - **Level 3 (Grandchild):** A `Job` resource created by a script injected into a Pod/Deployment. **Crucially, one resource level in the hierarchy will NOT have an Owner Reference!**
2. **Field Indexers:** Use an in-memory indexer (`SetFieldsFromObject`) to instantly look up unowned resources in the local informer cache without expensive API server scans.
3. **Namespace-Scoped Predicates:** Enforce strict namespace isolation so your controller only watches and reconciles resources within its designated active namespace.

# Deep Dive: Concepts & Architecture

## 1. Core Concepts Explained

### A. 3-Level Ownership & Unowned Resources

In production operators, not all resources can hold an Owner Reference:

- **Cross-Namespace Resources:** A child resource in Namespace `B` **cannot** have an `ownerReference` pointing to a parent in Namespace `A` (Kubernetes garbage collection forbids cross-namespace ownership).
- **Unowned Level Pattern:** In this assignment, Level 3 (`Job`) will **not** have an Owner Reference pointing back to Level 1 or 2.

```
 [ Level 1: Parent CR ] (MultiTierApp)
         │
         │  (Holds OwnerReference)
         ▼
 [ Level 2: Child ] (Deployment)
         │
         │  (Creates Pods running script)
         ▼
 [ Level 3: Grandchild ] (Job)  <-- NO Owner Reference!
```

### B. Field Indexers

When a child resource does **not** have an Owner Reference, how does the controller find it in the cache without querying the API server repeatedly?

**Field Indexers** create an in-memory lookup index in the `controller-runtime` cache—similar to a database index on a column.

- You index resources by a custom field (e.g., `metadata.annotations.parent-cr`).
- When reconciling, you query the local cache using `client.MatchingFields{"metadata.annotations.parent-cr": "cr-name"}`.

### C. Namespace-Scoped Predicates

In multi-tenant clusters, operators are often deployed to watch a single namespace. A **Namespace Predicate** intercepts incoming watch events at the Informer boundary and drops any event originating outside the controller's designated namespace.

## 2. Structural Architecture & Flow Diagram

```
 [ K8s API Event Watch Stream ]
                │
                ▼
┌──────────────────────────────────────────┐
│      Namespace-Scoped Predicate          │
│   (Drop events if namespace != "dev")   │
└───────────────────┬──────────────────────┘
                    │ (Allowed)
                    ▼
┌──────────────────────────────────────────┐
│             Reconciler Loop              │
└───────────────────┬──────────────────────┘
                    │
                    ├── 1. Get MultiTierApp (Level 1)
                    │
                    ├── 2. Reconcile Deployment (Level 2)
                    │      └── Attached OwnerReference -> MultiTierApp
                    │
                    └── 3. Reconcile Job (Level 3)
                           ├── NO OwnerReference attached!
                           ├── Attached Annotation: "app.example.com/parent": "my-app"
                           └── Query Cache via Field Indexer: "spec.parentRef"
```

## 3. Step-by-Step Implementation Guide

Let's set up **Assignment 4** in its own standalone directory: `assignment4`.

### Step 3.1: Directory & Kubebuilder Setup

```
cd k8s-dev-training
mkdir -p assignment4
cd assignment4

# Initialize standalone project
kubebuilder init --domain example.com --repo github.com/Ajay-Raut-Rsystems/k8s-dev-training/assignment4

# Create API & Controller
kubebuilder create api --group apps --version v1alpha1 --kind MultiTierApp --resource=true --controller=true
```

### Step 3.2: Define Schema (`api/v1alpha1/multitierapp_types.go`)

Open `api/v1alpha1/multitierapp_types.go` and define the schema:

Update the /api/v1alpha1/groupversion_info.go :




### Step 3.3: Implement Controller Logic (`internal/controller/multitierapp_controller.go`)

Open `internal/controller/multitierapp_controller.go` and implement 3-level management, field indexing, and predicates:


### Step 3.4: Configure `cmd/main.go` for Active Namespace

In `cmd/main.go` (or `main.go`), pass the target namespace (e.g., `default`) when setting up the reconciler:

```go
if err = (&controller.MultiTierAppReconciler{
    Client:          mgr.GetClient(),
    Scheme:          mgr.GetScheme(),
    WatchNamespace: "default", // Restricts reconciler to "default" namespace
}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "unable to create controller", "controller", "MultiTierApp")
    os.Exit(1)
}
```

### Step 3.5: Execution & Validation

1. Generate deepcopy code & CRD manifests:Bash
    
    ```
    make generate
    make manifests
    make install
    ```
    
2. Run the controller locally:Bash
    
    ```
    make run
    ```
    
3. Create sample CR `config/samples/apps_v1alpha1_multitierapp.yaml`:YAML
    
    ```yaml
    apiVersion: apps.example.com/v1alpha1
    kind: MultiTierApp
    metadata:
      name: sample-multitier
      namespace: default
    spec:
      replicas: 2
      targetNamespace: "default"
      jobCommand: "echo 'Level 3 Job executed successfully'"
    ```
    
4. Apply the CR:Bash
    
    ```
    kubectl apply -f config/samples/apps_v1alpha1_multitierapp.yaml
    ```
    
5. Verify 3-Level Resources:
    - **Level 2 Deployment (Owned):**Bash
        
        ```
        kubectl get deployment sample-multitier-deploy -o yaml | grep -A 5 ownerReferences
        ```
        
        *(Shows OwnerReference pointing to `sample-multitier`)*
        
    - **Level 3 Job (Unowned, Indexed):**Bash
        
        ```
        kubectl get job sample-multitier-job -o yaml | grep ownerReferences
        ```
        
        *(Returns nothing—no OwnerReference exists!)*
        
    - Verify annotation on Job:Bash
        
        ```
        kubectl get job sample-multitier-job -o yaml | grep -A 2 annotations
        ```
        
        *(Shows `apps.example.com/parent-cr: sample-multitier`)*
        
6. Test Namespace Scoping:
    
    Create a `MultiTierApp` instance in a different namespace (e.g., `kube-system`). Observe that the controller **ignores** it completely due to the namespace predicate filter!
    

## 4. Key Engineering Takeaways

- **Field Indexers enable high-performance lookups:** When objects lack Owner References, Field Indexers prevent costly $O(N)$ API server scans by indexing fields in local cache memory.
- **Cascading deletion requires manual handling for unowned children:** Deleting a `MultiTierApp` automatically cleans up Level 2 Deployment via K8s Garbage Collection, but Level 3 Job remains until explicitly cleaned up by the reconciler via finalizers or custom delete logic.
- **Namespace Predicates ensure multi-tenant safety:** Filtering events early at the predicate level prevents unnecessary work queues and resource leakage across namespace boundaries.
- 
    
    A multi-tenant cluster is a single shared computing environment (such as a Kubernetes cluster) used by multiple distinct groups—called tenants—while keeping their applications and data isolated.
    
    Think of it like an apartment building. Each tenant (a different internal team, department, or external customer) has their own private apartment unit, but everyone shares the same foundation, plumbing, electricity, and security system.
    
    **How It Works**
    
    - **Shared Infrastructure:** All tenants run their workloads on the exact same pool of physical or virtual servers, using the same operating system and control plane.
    - **Logical Isolation:** Instead of giving every team their own private cluster (single-tenancy), administrators use software boundaries to separate them.
    - **Policy Enforcement:** Tools like Role-Based Access Control (RBAC), resource limits, and network rules ensure one tenant cannot see, modify, or overwhelm another tenant's resources.
    
    **Types of Multi-Tenancy**
    
    - **Soft Multi-Tenancy:** Used for trusted internal teams within the same company. Separation is logical using Kubernetes Namespaces and access rules.
    - **Hard Multi-Tenancy:** Used for untrusted external clients or strict regulatory environments. This requires heavy isolation, often using virtual clusters (vCluster) or dedicated security barriers.
    
    ---
    
    In the Kubebuilder command,  **tells the scaffolding tool to generate the Custom Resource Definition (CRD) source code and API structures for your custom type**. [1,Kubernetes.%20These%20are%20used%20by%20the%20operator)]
    
    Here is exactly what that means and what it creates:
    
    **What  Does**
    
    - **Generates the Schema**: It creates the Go structs ( and ) where you define the configuration fields for your custom object.
    - **Registers the API**: It adds the code required to register your new object type () with the Kubernetes API server.
    - **Enables Manifest Generation**: It sets up the markers () so that the tool can automatically build the YAML files for your CRD.
    
    Code Generated by This Flag
    
    **When you set this flag to true, Kubebuilder creates a file typically located at . Inside this file, you will find:**
    
    - **MultiTierAppSpec**: The struct where you define the *desired* state (e.g., number of replicas, image names).
    - **MultiTierAppStatus**: The struct where you define the *observed* state (e.g., number of active pods, current conditions).
    - **MultiTierApp**: The top-level object that integrates into the Kubernetes API.
    
    **Why You Might Toggle It**
    
    - **—resource=true (Default)**: Use this when you want to store data in Kubernetes using a brand new custom YAML object.
    - **—resource=false** : Use this if you only want to create a controller to watch *existing* core Kubernetes resources (like Pods or Deployments) without inventing a new object type.
    
    ---

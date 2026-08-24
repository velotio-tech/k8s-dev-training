# Deep Dive: Assignment 1

## 1. Task I need to do?

Need to write Go programs that run inside or against your `kind` cluster to perform **Create, Read, Update, Delete (CRUD)** operations on **3 standard K8s resources** (e.g., `ConfigMap`, `Pod`, `Service`).

You must do this twice:

1. Using **`client-go`**: The official, low-level Kubernetes client SDK.
2. Using **`controller-runtime`**: The higher-level abstraction framework built on top of `client-go` (used by Kubebuilder and Operator SDK).
3. **Write-up**: Document non-obvious differences, performance traits, and utility functions discovered.

---

## 2. Key Concepts & Differences

### `client-go`

* **What it is:** The foundational Go library maintained by the Kubernetes project.
* **How it talks to K8s:** Uses specific, typed clientsets (e.g., `clientset.CoreV1().Pods(namespace)`).
* **Behavior:** By default, every `.Get()` or `.List()` call sends a **direct HTTP request** to the K8s API server.
* **When to use:** CLI tools (like `kubectl` plugins), simple scripts, or lightweight utilities.

### `controller-runtime`

* **What it is:** A higher-level framework created by the Kubernetes Special Interest Group (SIG) API Machinery.
* **How it talks to K8s:** Uses a unified, generic `client.Client` interface capable of handling any resource (standard or Custom Resource) using Go types.
* **Behavior:** Uses **Informers and Caches** under the hood. Reading data (`Get`/`List`) reads from a **local in-memory cache** synced via watches, drastically reducing API server load. Writes (`Create`/`Update`/`Delete`) go directly to the API server.
* **When to use:** Custom Controllers, Operators, and Kubebuilder projects.

---

## 3. Architecture & Data Flow Comparison

```
[ Your Program ]
      │
      ├── (A) client-go Direct API Approach
      │     │
      │     └── Direct REST HTTP Call ───────────────► [ K8s API Server ] ──► [ etcd ]
      │
      └── (B) controller-runtime Cached Approach
            │
            ├── READS (Get/List) ──► [ Local Cache (Informer) ] (Fast, 0 API Server load)
            │                                ▲
            │                        Watch Sync Stream
            │                                │
            └── WRITES (Create/Update) ──────┴────────► [ K8s API Server ] ──► [ etcd ]

```

---

## 4. Step-by-Step Implementation Guide

Follow these steps to set up and run the code on your Linux machine with `kind` and `docker`.

### Step 4.1: Setup Go Module

Create a clean directory for Assignment 1:

```bash
mkdir -p k8s-dev-training/assignment1
cd k8s-dev-training/assignment1
go mod init assignment1

```

Install the required dependencies:

```bash
go get k8s.io/api/core/v1@v0.30.0
go get k8s.io/apimachinery/pkg/apis/meta/v1@v0.30.0
go get k8s.io/client-go@v0.30.0
go get sigs.k8s.io/controller-runtime@v0.18.0

```

---

### Step 4.2: Part A — `client-go` Implementation

Created a file named `client_go_crud.go`:


---

### Step 4.3: Part B — `controller-runtime` Implementation

Create a file named `controller_runtime_crud.go`:

---

### Step 4.4: Execution & Verification

Make sure your `kind` cluster is active:

```bash
kubectl cluster-info

```

Run both programs directly against your local cluster:

```bash
go run client_go_demo.go
go run controller_runtime_demo.go

```

---

## 5. Write-up: Key Learnings & Non-General Findings

Summary of my learnings :

1. **Client Ergonomics:**
* `client-go` requires accessing resources via versioned interfaces (`clientset.CoreV1().Pods(ns)`).
* `controller-runtime` uses a single `client.Client` interface where you pass generic Go structs representing K8s objects (`k8sClient.Get(ctx, key, obj)`).


2. **In-Cluster vs Out-of-Cluster Configuration:**
* Both scripts utilize `rest.InClusterConfig()`, which looks for service account tokens mounted at `/var/run/secrets/kubernetes.io/serviceaccount` inside a Pod.
* `clientcmd.BuildConfigFromFlags("", kubeconfig)` serves as an essential fallback when developing locally on a laptop.


3. **Status Subresource Handling:**
* In `controller-runtime`, modifying object status requires calling `k8sClient.Status().Update()`. Regular `k8sClient.Update()` ignores changes made to `.Status` fields when the status subresource is enabled.


4. **Useful Utilities Discovered:**
* `types.NamespacedName`: A clean utility struct in `k8s.io/apimachinery/pkg/types` to identify objects by namespace and name without string concatenation.
* `client.Object`: A universal interface in `controller-runtime` implemented by all Kubernetes API objects, enabling write helpers that accept any K8s object type.
* `controller-runtime/pkg/client/fake`: A handy utility package for writing unit tests without needing a running API server or `envtest` setup.

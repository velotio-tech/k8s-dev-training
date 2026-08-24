# Assignment 2

Extending the Kubernetes API itself by creating Custom Resource Definitions (CRDs) using **Kubebuilder v2.0/v3+** patterns, schema markers, and status subresources.

## The Big Picture: What Are You Building in Assignment 2?

In Kubernetes, standard resources are things like `Pod`, `Service`, or `Deployment`. A **CRD (Custom Resource Definition)** lets you register your *own* custom object with the Kubernetes API server—for example, `CronJob`, `MyDatabase`, or `AppDeployment`.

In this assignment, you will:

1. Initialize a Kubebuilder project.
2. Define a custom API schema with 3–4 spec fields and 3–4 status fields.
3. Apply schema validation markers (`+kubebuilder:validation:...`).
4. Generate CRD YAML manifests using `controller-gen`.
5. Experiment with resource creation (valid vs. invalid specs).
6. Enable and explore the **Status Subresource** to understand why standard CRUD won't update `.status`.

# Deep Dive: Assignment 2

## 1. What is Being Asked?

- **Define 1 API:** Choose a logical domain (e.g., `AppService`, `DatabaseCluster`, or `BackupJob`).
- **Define Spec & Status Fields:**
    - **Spec fields:** Desired state supplied by the user (at least 3–4 fields).
    - **Status fields:** Observed state updated by the system/controller (at least 3–4 fields).
- **Apply Markers:** Use Go comment tags (markers) to enforce validation rules like minimum/maximum values, enum constraints, optional/required flags, and array lengths.
- **Generate & Experiment:** Run `make manifests` (using `controller-gen`) to emit YAML. Create valid and invalid custom resource instances to test K8s OpenAPI validation.
- **Status Subresource Investigation:** Enable `/status` subresource, test modifying `.status` directly via `kubectl apply` vs. status API endpoints, and explain why standard CRUD cannot mutate status fields.

## 2. Key Concepts & Architecture

### What is a CRD & OpenAPI v3 Schema?

When you define a Go struct for your CRD, `controller-gen` reads your code and comment markers to produce an OpenAPI v3 validation schema in YAML. When a user sends `kubectl apply -f my-cr.yaml`, the K8s API server validates the request payload against this schema **before** storing it in `etcd`.

### Spec vs. Status Separation

- **`spec`**: Represents the **Desired State** set by the human user or client.
- **`status`**: Represents the **Observed/Actual State** populated exclusively by a Controller or Operator.

### The Status Subresource (`/status`)

By default, Kubernetes treats an entire YAML file as a single object. When you add the `+kubebuilder:subresource:status` marker:

1. The API server creates a dedicated REST endpoint: `/apis/<group>/<version>/namespaces/<ns>/<plural>/<name>/status`.
2. **Main Endpoint (`/apis/.../<name>`):** Ignores all modifications to the `.status` block during standard `Update`/`Apply` calls.
3. **Status Endpoint (`/apis/.../<name>/status`):** Only modifies the `.status` block and **ignores changes to `.spec`**.

## 3. Data Flow & Subresource Mechanism

```
                       ┌─────────────────────────────────────────────────────────────┐
                       │                   Kubernetes API Server                     │
                       └─────────────────────────────────────────────────────────────┘
                                                      │
         ┌────────────────────────────────────────────┴────────────────────────────────────────────┐
         │                                                                                         │
         ▼                                                                                         ▼
  [ Standard REST Endpoint ]                                                      [ /status Subresource Endpoint ]
  POST/PUT /apis/apps.example.com/v1alpha1/.../my-app                               PUT /apis/apps.example.com/v1alpha1/.../my-app/status
         │                                                                                         │
         ├── Modifies: .spec, .metadata                                                            ├── Modifies: .status ONLY
         └── Strips/Ignores: .status changes                                                       └── Strips/Ignores: .spec changes
         │                                                                                         │
         └────────────────────────────────────────────┬────────────────────────────────────────────┘
                                                      │
                                                      ▼
                                                [ etcd Storage ]
```

## 4. Step-by-Step Implementation Guide

### Step 4.1: Directory Setup & Kubebuilder Initialization

On your Linux machine with Docker and `kind` running:

```bash
# Create directory for Assignment 2
mkdir -p k8s-dev-training/assignment2
cd k8s-dev-training/assignment2

# Initialize Kubebuilder project
kubebuilder init --domain example.com --repo github.com/your-username/k8s-dev-training/assignment2

# Create API (Group: apps, Version: v1alpha1, Kind: AppService)
kubebuilder create api --group apps --version v1alpha1 --kind AppService --resource=true --controller=false
```

### Step 4.2: Define API Spec & Status (`api/v1alpha1/appservice_types.go`)

Open `api/v1alpha1/appservice_types.go` and replace the struct definitions with schema featuring rich markers:



### Step 4.3: Generate Manifests & Install CRD

1. Generate deepcopy code and CRD YAML manifests using `controller-gen`:
    
    ```
    make manifests
    ```
    
2. Inspect the generated CRD file at `config/crd/bases/apps.example.com_appservices.yaml`. Notice how markers were translated into OpenAPI v3 validation rules.
3. Install the CRD into your `kind` cluster:
    
    ```
    make install
    ```
    
4. Verify installation:
    
    ```
    kubectl get crd appservices.apps.example.com
    ```
    

### Step 4.4: Tryout Positive & Negative Scenarios

Create a test folder `config/samples/testing/` to perform validation tests:

```
mkdir -p config/samples/testing
```

#### Test 1: Valid Custom Resource (`valid-cr.yaml`)

Create `config/samples/testing/valid-cr.yaml`:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: valid-app
  namespace: default
spec:
  replicas: 3
  environment: "staging"
  image: "nginx:latest"
  ports:
    - 80
    - 443
```

Apply it:

```
kubectl apply -f config/samples/testing/valid-cr.yaml
# Result: appservice.apps.example.com/valid-app created
```

#### Test 2: Invalid Replica Count (`invalid-replicas.yaml`)

Create `config/samples/testing/invalid-replicas.yaml`:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: invalid-replicas-app
  namespace: default
spec:
  replicas: 15 # Exceeds Maximum=10 marker
  environment: "dev"
  image: "nginx:latest"
```

Apply it:

```
kubectl apply -f config/samples/testing/invalid-replicas.yaml
```

**Expected Error Output:**

```
The AppService "invalid-replicas-app" is invalid: spec.replicas: Invalid value: 15: spec.replicas in body should be less than or equal to 10
```

We enforce the limits in your Go code using **Kubebuilder OpenAPI Validation Markers** placed directly above the `Replicas` field in `api/v1alpha1/appservice_types.go`.

### Where it is configured in your project

#### 1. In `api/v1alpha1/appservice_types.go`

In your Go struct, you added the `+kubebuilder:validation:Maximum` comment tag:

```go
type AppServiceSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// Replicas is the desired number of pod instances
	Replicas int32 `json:"replicas,omitempty"`

	// ... other fields
}
```

#### 2. How it translates into your CRD (`config/crd/bases/apps.example.com_appservices.yaml`)

When you ran `make manifests`, `controller-gen` parsed those comment tags and generated OpenAPI v3 schema validation rules inside your CRD manifest:

```yaml
properties:
  spec:
    properties:
      replicas:
        type: integer
        minimum: 1
        maximum: 10  # <-- controller-gen generated this rule
```

### How the validation works at runtime

When you run `kubectl apply -f config/samples/testing/invalid-replicas.yaml`:

1. Your request hits the **Kubernetes API Server**.
2. Before saving the object to etcd, the API server checks the incoming `spec.replicas: 15` against the **OpenAPI v3 schema** defined in the registered CRD.
3. Since `15 > 10`, the API Server immediately rejects the request with `Invalid value: 15: spec.replicas in body should be less than or equal to 10`.

This happens entirely at the API server layer **before your custom controller code even sees or receives the object**.

#### Test 3: Invalid Enum Environment (`invalid-env.yaml`)

Create `config/samples/testing/invalid-env.yaml`:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: invalid-env-app
  namespace: default
spec:
  replicas: 2
  environment: "production" # Invalid! Enum restricts to [dev, staging, prod]
  image: "nginx:latest"
```

Apply it:

```bash
kubectl apply -f config/samples/testing/invalid-env.yaml
```

**Expected Error Output:**

```
The AppService "invalid-env-app" is invalid: spec.environment: Unsupported value: "production": supported values: "dev", "staging", "prod"
```

### Step 4.5: Status Subresource Experimentation

Now, let's explore why direct CRUD on status fails via standard YAML applies.

1. Create `config/samples/testing/cr-with-status.yaml`:

```yaml
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: status-test-app
  namespace: default
spec:
  replicas: 2
  environment: "dev"
  image: "nginx:alpine"
status:
  phase: "Running"
  readyReplicas: 2
```

1. Apply the manifest:
    
    ```bash
    kubectl apply -f config/samples/testing/cr-with-status.yaml
    ```
    
2. Fetch the object back from the cluster:
    
    ```
    kubectl get appservice status-test-app -o yaml
    ```
    

**Observation:**

The `.status` block in the cluster output is **empty** or defaults to zero values! The `phase: "Running"` and `readyReplicas: 2` fields supplied in the YAML were completely ignored by the API server.

1. **Updating Status Correctly (via API / curl):**
    
    To update status, you must hit the `/status` subresource directly:
    

```bash
# Proxy local port to API server
kubectl proxy --port=8001 &

# Update status via raw HTTP subresource endpoint
curl -X PUT http://localhost:8001/apis/apps.example.com/v1alpha1/namespaces/default/appservices/status-test-app/status \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion": "apps.example.com/v1alpha1",
    "kind": "AppService",
    "metadata": {
      "name": "status-test-app",
      "namespace": "default",
      "resourceVersion": "<GET_RESOURCE_VERSION_FROM_KUBECTL>"
    },
    "status": {
      "phase": "Running",
      "readyReplicas": 2
    }
  }'
```

## 5. Write-up: Subresources & Why General CRUD Cannot Mutate Status

### What is a Subresource?

A subresource is an auxiliary endpoint hanging off a main Kubernetes resource REST path (e.g., `/status`, `/scale`, `/exec`). It targets a specific subset of an object's state rather than the whole document.

### Why standard CRUD (`kubectl apply` / REST `.Update()`) ignores status:

1. **Security & Role-Based Access Control (RBAC):** Users need permission to update `spec` (desired configuration), but shouldn't directly fake system readiness. Operators/Controllers are granted RBAC rights specifically for `/status` subresources.
2. **Concurrency & Optimistic Locking:** If users and automated controllers continuously modified the same document endpoint, they would constantly overwrite each other's changes (`ResourceVersion` conflict errors).
3. **API Server Isolation:** When the `+kubebuilder:subresource:status` marker is defined, the K8s API server handler explicitly ignores changes to `.status` on standard paths (`PUT/PATCH /apis/apps.example.com/.../appservices/name`), requiring clients to use `PUT /apis/apps.example.com/.../appservices/name/status`.

**An external user or admin *can* use the `/status` subresource endpoint directly**, but standard tools like `kubectl apply` do not hit that endpoint by default.

Here is why it feels confusing and how the permissions and architecture actually work:

### 1. `kubectl apply` vs. The `/status` Endpoint

When you run `kubectl apply -f my-cr.yaml`, `kubectl` sends an HTTP `PUT` or `PATCH` request strictly to the **main object endpoint**:

$$\text{Main Endpoint: } \texttt{/apis/apps.example.com/v1alpha1/namespaces/default/appservices/status-test-app}$$

When the Status Subresource is enabled on a CRD, Kubernetes automatically splits the resource into two separate API paths:

| **API Endpoint** | **What it controls** | **What happens if you send .spec** | **What happens if you send .status** |
| --- | --- | --- | --- |
| **Main Endpoint** (`/appservices/name`) | `.metadata`, `.spec` | **Saved** | **Silently Ignored / Erased** |
| **Status Endpoint** (`/appservices/name/status`) | `.status` | **Ignored** | **Saved** |

Because `kubectl apply` only calls the **Main Endpoint**, any `.status` block defined inside your YAML file is stripped out before being saved.

### 2. Can an external user modify the status manually?

**Yes.** Anyone or any tool (like `curl`, a custom script, or `kubectl`) that sends a request **specifically targeted at the `/status` subresource path** can update the status, provided their RBAC permissions allow it.

To manually update the status of `status-test-app` right now from your terminal using `kubectl`:

```
kubectl replace --raw /apis/apps.example.com/v1alpha1/namespaces/default/appservices/status-test-app/status -f - <<EOF
{
  "apiVersion": "apps.example.com/v1alpha1",
  "kind": "AppService",
  "metadata": {
    "name": "status-test-app",
    "namespace": "default"
  },
  "status": {
    "phase": "Running",
    "readyReplicas": 2
  }
}
EOF
```

If you run `kubectl get appservice status-test-app -n default -o yaml` after executing the command above, you will see `.status` successfully updated!

### 3. Why Kubernetes designed it this way

Separating these endpoints solves two major issues in distributed systems:

- **RBAC Granularity:** In production, you can grant human developers access to edit `.spec` (to request changes), while restricting permission on `/status` so that *only* the controller's `ServiceAccount` can write to `/status` (to report actual state).
- **Concurrency Conflicts (Optimistic Locking):** Controllers frequently update status (e.g., reporting pod readiness). If status updates were routed through the main endpoint, a controller updating `.status` would constantly collide with users updating `.spec`, causing constant `409 Conflict` errors.

!image.png

Extra findings :

- 

The YAML is generated from the Go markers by `make manifests`, configured in `Makefile:43-47`. Therefore, you should normally edit only the Go markers and regenerate:

```bash
ajayraut@EMPID20144:~/Documents/k8s-assignments/assignment2$ tree
.
├── AGENTS.md
├── api
│   └── v1alpha1
│       ├── appservice_types.go [This file is the main for defining the schema of the API, speciifyin the rules using go markers, read by the controllen-gen , and this file helps to create the .yaml i.e the CRD]
│       ├── groupversion_info.go
│       └── zz_generated.deepcopy.go
├── cmd
│   └── main.go
├── config
│   ├── crd
│   │   ├── bases
│   │   │   └── apps.example.com_appservices.yaml   [This file is generated using the make manifest command , and installed in the cluster using the make install command]
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   └── samples
│       ├── apps_v1alpha1_appservice.yaml
│       ├── kustomization.yaml
│       └── testing
│           ├── cr-with-status.yaml
│           ├── invalid-env.yaml
│           ├── invalid-replicas.yaml
│           └── valid-cr.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
├── Makefile
```

│       ├── **`appservice_types.go`** 

- [This file is the main for defining the schema of the API, speciifyin the rules using go markers, read by the controllen-gen , and this file helps to create the .yaml i.e the CRD]
- so if we want to update the rules / enforce new rules , edit the markers of this file, then run make manifestes again to the update the /**`config/crd/bases/apps.example.com_appservices.yaml`** CRD file which is actually being consumed by kubernetes and enforces the rules using this
- So to avoid mistakes of YAML file creation, we have seprattion, of the rule defined using the go file that generate the CRD yaml file.

```bash
assignment2/
├── Makefile                          <-- Automation scripts generated by Kubebuilder
├── api/v1alpha1/
│   └── appservice_types.go          <-- API struct definition with schema markers
├── config/crd/bases/
│   └── apps.example.com_appservices.yaml <-- CRD Manifest (Generated)
└── config/samples/testing/
├── valid-cr.yaml                <-- Valid Custom Resource instance
├── invalid-replicas.yaml        <-- Triggers OpenAPI schema validation error
└── cr-with-status.yaml          <-- Proves status subresource behavior
```

### **Resource Flow & Key Terms**

```
                       [ Go Structs + Markers ]
                       (api/v1alpha1/appservice_types.go)
                                  │
                                  │ make manifests (controller-gen)
                                  ▼
                   [ Custom Resource Definition (CRD) ]
             (config/crd/bases/apps.example.com_appservices.yaml)
                                  │
                                  │ kubectl apply -f (Installed into Cluster)
                                  ▼
                     [ API Server registers new API ]
                                  │
                                  │ kubectl apply -f valid-cr.yaml
                                  ▼
                     [ Custom Resource (CR) Instance ]
                                (valid-app)
```

#### **Term Glossary:**

- **CRD (Custom Resource Definition):** The *blueprint/schema* installed cluster-wide. It instructs the API server how to validate and store a new resource type (`AppService`).
- **CR (Custom Resource):** An *instance* created by a user based on the CRD schema (analogous to creating a Pod from the Pod specification).
- **Markers:** Go comment tags (e.g., `// +kubebuilder:validation:Minimum=1`) converted by `controller-gen` into OpenAPI v3 schema validation rules inside the CRD.
- **Status Subresource (`/status`):** A dedicated REST endpoint. When enabled via `+kubebuilder:subresource:status`, standard `kubectl apply` updates on `.status` are **ignored by design**. Status can only be mutated via dedicated status API endpoints by controllers.

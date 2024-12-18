# Action Plan for Custom Controller Implementation

## Step 1: Understand Custom Controllers and Shared Informers in Kubernetes

### Learn the Architecture of Shared Informers
- Study how shared informers work to cache and watch Kubernetes resources.
- Understand their role in reducing API server calls and maintaining a local cache of resource states.
- Key concepts to review:
  - **SharedIndexInformer**: Used for caching and index-based lookup.
  - **Event Handlers**: Add, Update, Delete handlers.
  - **Resync Period**: Configurable periodic resync of resources.

### Study Custom Controllers
- Learn how Kubernetes controllers operate in a control loop.
- Understand reconciliation logic to align the actual state with the desired state.
- Review the usage of the **controller-runtime** library for implementing controllers.

## Step 2: Design the CRD

### CRD Definition
- Define a Custom Resource Definition (CRD) that includes:
  - **Spec** field for inputting `GroupVersionKind` (GVK) and resource names.
- Example CRD Spec:
  ```yaml
  apiVersion: example.com/v1
  kind: CustomResource
  metadata:
    name: my-custom-resource
  spec:
    resources:
      - group: apps
        version: v1
        kind: Deployment
        name: my-deployment
      - group: core
        version: v1
        kind: ConfigMap
        name: my-configmap
  ```

## Step 3: Implement the Controller

### Controller Setup
- Use the `controller-runtime` library to simplify the implementation.
- Initialize the manager and the custom controller.

### Watch Resources
- Configure the controller to watch:
  - The Custom Resource for changes to its `spec`.
  - Relevant Kubernetes resources specified by the `GroupVersionKind` in the CRD.
- Ensure that reconciliation logic triggers upon relevant events.

### Reconcile Loop
- Implement the reconciliation loop to ensure the desired state matches the actual state:
  - Read the `resources` field in the CRD's `spec`.
  - For each specified GVK and resource name:
    - Check if the resource exists in the cluster. If not, create it.
    - If the resource exists but is not in the desired state, update it.
    - Delete any resources not specified in the current `spec` but managed by the CRD.
  - Maintain consistency between the custom resource and the actual resources in the cluster.

### Owner References
- Add the custom resource as the **OwnerReference** to all managed resources.
- This ensures automatic garbage collection of managed resources if the custom resource is deleted.

## Step 4: Handle CRUD Operations on the Custom Resource Spec

### Spec Updates
- Design the controller to handle changes to the `spec` of the custom resource:
  - Additions: Create new resources for newly added GVKs.
  - Modifications: Update existing resources to match the new desired state.
  - Deletions: Remove resources no longer specified in the `spec`.

### Event Handling with Predicates
- Use predicates to filter events for efficiency:
  - Ignore events that do not change the desired state.
  - Process events for additions, updates, and deletions in the `resources` list.

### Reconcile Logic for Spec Changes
- Re-evaluate the `spec` during each reconciliation cycle to:
  - Dynamically adjust resources based on the updated desired state.
  - Ensure the cluster state reflects the changes in the custom resource.

## Step 5: Test Reconciliation and Event Handling

### Unit Testing
- Write comprehensive unit tests for the controller logic:
  - Mock Kubernetes API calls to simulate the cluster state.
  - Verify that the controller performs the correct CRUD operations for resources.

### Integration Testing
- Deploy the controller in a test Kubernetes cluster.
- Apply test cases with different CRD instances and validate:
  - Resources are created, updated, or deleted as per the custom resource `spec`.
  - Events are correctly filtered and handled.
  - Owner references are appropriately set.

## Step 6: Documentation and Deployment

### Document the CRD and Controller
- Provide clear instructions for:
  - Installing the CRD.
  - Deploying the custom controller.
  - Writing custom resources based on the CRD schema.
  - Verifying and troubleshooting the controller's behavior.

### Deployment
- Build and package the controller as a Docker image.
- Create Kubernetes manifests (e.g., `Deployment`, `ServiceAccount`, `Role`, `RoleBinding`) for deploying the controller in a cluster.
- Deploy the controller and test it in a staging environment before production rollout.
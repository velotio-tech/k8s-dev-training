# Assignment 2 setup — run these on your own laptop


## 1. Install kubebuilder (one-time)

```bash
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/linux/amd64"
chmod +x kubebuilder
sudo mv kubebuilder /usr/local/bin/
kubebuilder version
```

## 2. Scaffold the project

```bash
mkdir webapp-operator && cd webapp-operator
go mod init github.com/<you>/webapp-operator
kubebuilder init --domain example.com --repo github.com/<you>/webapp-operator
```

This generates a standard project skeleton: `main.go`, `PROJECT`, `Makefile`,
`config/` (kustomize manifests), `hack/`.

## 3. Scaffold the API (this is the important step)

```bash
kubebuilder create api --group apps --version v1alpha1 --kind WebApp
```

When it asks:
```
Create Resource [y/n]  -> y
Create Controller [y/n] -> y
```

This generates:
- `api/v1alpha1/webapp_types.go`  ← you edit this (spec/status structs + markers)
- `internal/controller/webapp_controller.go` ← Assignment 3 territory, ignore for now
- `api/v1alpha1/zz_generated.deepcopy.go` ← auto-generated, never hand-edit

## 4. Replace the generated types.go

Delete the placeholder content of `api/v1alpha1/webapp_types.go` and replace
it with `webapp_types.go` from this folder (I'm giving you the full file
below with markers and comments explaining each one).

## 5. Generate CRD YAML + deepcopy code

```bash
make manifests   # reads your markers, writes config/crd/bases/*.yaml
make generate    # regenerates zz_generated.deepcopy.go
```

`make manifests` is doing the "use controller-gen tools" step the
assignment asks for — it's already wired into the Makefile kubebuilder
scaffolded for you.

## 6. Install the CRD into your kind cluster

```bash
make install
kubectl get crd webapps.apps.example.com
kubectl explain webapp.spec   # this reads straight from your markers!
```

`kubectl explain` reading your Go comments back to you is a nice "it's
working" signal — go try it once you're here.

## 7. Try the sample CRs (positive and negative scenarios)

See `samples/` in this folder — apply the valid one, then each invalid
one, and read the rejection messages carefully; they map 1:1 to the
markers you wrote.

```bash
kubectl apply -f config/samples/valid-webapp.yaml
kubectl apply -f config/samples/invalid-replicas.yaml
kubectl apply -f config/samples/invalid-environment.yaml
kubectl apply -f config/samples/missing-image.yaml
```

## 8. The status subresource experiment (the actual point of this assignment)

```bash
kubectl apply -f config/samples/valid-webapp.yaml
kubectl get webapp demo-webapp -o yaml
```

Now try editing status directly:

```bash
kubectl edit webapp demo-webapp
```

Add a `status:` block with some fields and save. **Watch what happens on
re-read** — your status edit gets silently dropped. That's the proof.

Then do it the *correct* way, via the status subresource explicitly:

```bash
kubectl patch webapp demo-webapp --subresource=status --type=merge \
  -p '{"status":{"phase":"Running","availableReplicas":1,"message":"manually set for testing"}}'

kubectl get webapp demo-webapp -o yaml
```

This time it sticks — because you wrote to `/apis/apps.example.com/v1alpha1/.../webapps/demo-webapp/status`
instead of the main resource endpoint.
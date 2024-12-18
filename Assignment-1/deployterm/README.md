# Deployterm
A simple TUI (Text User Interface) tool for managing Kubernetes deployments.

## Features
- **List all namespaces**: View all available namespaces in your Kubernetes cluster.
- **Create new deployments**: Easily create a new deployments under a specific namespace.
- **List deployments in each namespace**: Easily see all deployments under a specific namespace.
- **Describe deployment**: Get detailed information about a specific deployment.
- **Edit deployment**: Modify deployment configurations (e.g., replicas, containers, etc.).
- **Delete deployment**: Remove a deployment from a specified namespace.

## Usage
```sh
$ ./deployterm --help
```

```
Usage of ./deployterm:
  -kubeconfig string
        path to kubeconfig file
  -use-controller-runtime
        use controller-runtime library instead of client-go
```

## Demo
![Demo](screen_recording.gif)

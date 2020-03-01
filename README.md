# ReplicatedConfigMap

`ReplicatedConfigMap` is an Operator that comes bundled with a CRD such that you
can created a `ReplicatedConfigMap` resource, and the data that's placed into a
`ReplicatedConfigMap` will be duplicated to namespaces with a specific
annotation.

`ReplicatedConfigMap` is a simple project to learn the basics around developing
an Operator for the Kubernetes ecosystem.

## Requirements / Goals
There is a small number of requirements that the `ReplicatedConfigMap` needs to
address.

1. `ReplicatedConfigMap` is implemented via a CRD. The provided functionality
   could have been implemented with just a controller and labels or annotations
   on a `ConfigMap`, but we explicitly wanted to create a CRD.
2. The `ReplicatedConfigMap` controller is to only pay attention to namespaces
   with a specific annotation.
3. The `ReplicatedConfigMap` is to be built using `kubebuilder` as opposed to
   other Operator frameworks.
4. The `ReplicatedConfigMap` needs to immediately react to new namespaces that
   are created with the appropriate annotations.

## Open Questions

### Do I need to run the controller in each namespace that I want it to listen to?

No, the `ReplicatedConfigMap` controller was built using `kubebuilder`'s
cluster wide mode. This means that when you create a `ReplicatedConfigMap`, it
actually doesn't live in any namespace at all.

Allowing for cross namespace owner references without being a non-namespaced
resource is [explicitly not
allowed](https://github.com/kubernetes-sigs/controller-runtime/pull/675) by the
`controller-runtime` project.

### What's the behavior of removing managed ConfigMaps?

ConfigMaps that were created by the `ReplicatedConfigMap` controller that are
deleted manually will get instantly re-created by the controller.

### What's the behavior of removing the annotation set on a namespace?

The ConfigMaps that were created will stay there. There was no requirement to
reconcile deletion of ConfigMaps when a namespace is no longer being watched.

### What's the behavior of removing the source `ReplicatedConfigMap`

All children of the source `ReplicatedConfigMap` CRD will get deleted via
Kubernetes' ownership model.

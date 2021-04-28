---
title: Quarkus Operator Tutorial
linkTitle: Tutorial
weight: 30
description: An in-depth walkthough of building and running a Quarkus-based operator.
---

## Prerequisites

- Java through the [installation guide](https://java.com/en/download/help/download_options.html).
- User authorized with `cluster-admin` permissions.
- Maven installation [installation guide](https://maven.apache.org/install.html)

## Overview

We will create a sample project to let you know how it works and this sample will:

- Create a Memcached Deployment if it doesn't exist
- Ensure that the Deployment size is the same as specified by the Memcached Custom Resource spec
- Update the Memcached Custom Resource status using the status writer with the names of the Custom Resource's pods

## Create a new project

Use the CLI to create a new memcached-quarkus-operator project:

```sh
mkdir memcached-quarkus-operator
cd memcached-quarkus-operator
# we'll use a domain of example.com
# so all API groups will be <group>.example.com
operator-sdk init --plugins quarkus --domain example.com --project-name memcached-quarkus-operator
```

#### A note on dependency management

`operator-sdk init` generates `pom.xml` file. This file contains all the dependencies required to run the operator.

### MemcachedQuarkusOperator

The main program of the operator `MemcachedQuarkusOperator.java` initializes and runs the operator. The operator uses java-operator-sdk which is similar to the Go lang version of controller-runtime.

The code below will initialize and define the informers/watches for your operator.

```
  @Override
  public int run(String... args) throws Exception {
    operator.start();

    Quarkus.waitForExit();
    return 0;
  }
```

## Create a new API and Controller

Create a new Custom Resource Definition (CRD) API with group `cache` version `v1` and Kind `Memcached`. The plugin, still in its alpha state, will output debug messages which are normal.

`create api` command will  scaffold the `MemcachedController`, `MemcachedSpec`, `MemcachedStatus` and `Memcached`. 

```console
$ operator-sdk create api --plugins quarkus --group cache --version v1 --kind Memcached
InjectResource called
UpdateResource called
Scaffold called
NewCreateAPIScaffolder called
InjectFS called
api.Scaffold()
PostScaffold called
...
```

After the `create api` command the file structure will be shown as below.

$ tree
.
├── pom.xml
├── PROJECT
└── src
    └── main
        ├── java
        │   └── com
        │       └── example
        │           ├── MemcachedController.java
        │           ├── Memcached.java
        │           ├── MemcachedQuarkusOperator.java
        │           ├── MemcachedSpec.java
        │           └── MemcachedStatus.java
        └── resources
            └── application.properties

6 directories, 8 files

#### Understanding Kubernetes APIs

For an in-depth explanation of Kubernetes APIs and the group-version-kind model, check out these [kubebuilder docs](https://book.kubebuilder.io/cronjob-tutorial/gvks.html).

The `java-operator-plugins` project uses the APIs from [java-operator-sdk](https://github.com/java-operator-sdk/java-operator-sdk). The java-operator-sdk library is to Java what `controller-runtime` is to Go. It provides a framework to help Java projects interact with the Kubernetes API.

### Define the API

#### `MemcachedSpec`

Initially, the scaffolded Spec file will be empty. The operator developer needs to add attributes to this file according to his need. For the Memcached example, we added the size field as shown below example.

```
// MemcachedSpec defines the desired state of Memcached
public class MemcachedSpec {

    // Add Spec information here
    // Size is the size of the memcached deployment
    private Integer size;

    public Integer getSize() {
        return size;
    }

    public void setSize(Integer size) {
        this.size = size;
    }
}
```

#### `MemcachedStatus`

Similar to the Spec file `MemcachedStatus` file got scaffolded as part of the `create api` command. The user has to modify a Status file. For the Memcached example, we too list of nodes as shown below.
The nodes field is a list of string values and it contains the name of the Memcached pods.


```
import java.util.ArrayList;
import java.util.List;

// MemcachedStatus defines the observed state of Memcached
public class MemcachedStatus {

    // Add Status information here
	// Nodes are the names of the memcached pods
    private List<String> nodes;

    public List<String> getNodes() {
        if (nodes == null) {
            nodes = new ArrayList<>();
        }
        return nodes;
    }

    public void setNodes(List<String> nodes) {
        this.nodes = nodes;
    }
}
```

**Note** The Node field is just to illustrate an example of a Status field. In real cases, it would be recommended to use [Conditions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties).

#### `Memcached`

`Memcached` file scaffolded via `create api` command. It extends the properties of `MemcachedSpec` and `MemcachedStatus`.

```
// Memcached is the Schema for the memcacheds API

@Version("v1alpha1")
@Group("cache.example.com")
public class Memcached extends CustomResource<MemcachedSpec, MemcachedStatus>
    implements Namespaced {}
```

### Apply Custom Resource and CRD's using below command

#### CRD - 

Create a file with the name `crd.yaml`. CRD enables users to add their own/custom objects to the Kubernetes cluster.

This file will contain the below content. 

```
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.4.1
  name: memcacheds.cache.example.com
spec:
  group: cache.example.com
  names:
    kind: Memcached
    listKind: MemcachedList
    plural: memcacheds
    singular: memcached
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        ...
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
```

Create the CRD:

`kubectl apply -f crd.yaml`

### Create Memcached Custom Resource - memcached-sample.yaml

Create the sample Memcached Custom Resource manifest at k8s/samples/memcached-sample.yaml and define the spec as the following.

This file will contain the below content. 

```
apiVersion: cache.example.com/v1
kind: Memcached
metadata:
  name: memcached-sample
spec:
  # Add fields here
  size: 1
```

Create the Custom Resource:

`kubectl apply -f memcached-sample.yaml`

## Implement the Controller

Add the below-mentioned code snippet in `MemcachedController.java` file. Initially, this file will contain the empty methods `createOrUpdateResource` and `deleteResource`. Please add the below code in respective methods as part of controller logic. Also, add the `createMemcachedDeployment` method that will create the Deployment for your operator.

**Note**: Next two subsections explain the two methods `createOrUpdateResource` and `deleteResource`. These two methods get called whenever some update/create/delete event occurs in the cluster.

These methods are already scaffolded as part of the `create api` command. In this memcached example, we will need to watch the deployment so we can react to size changes. We will accomplish this in the steps below.

### createOrUpdateResource

First, let's get the Deployment. In the MemachedController.java, find the createOrUpdateResource method and add the following code below the `// TODO: fill in logic` comment. It should look like:

```
        Deployment deployment =
                client
                        .apps()
                        .deployments()
                        .inNamespace(resource.getMetadata().getNamespace())
                        .withName(resource.getMetadata().getName())
                        .get();
```

If the deployment is `null` that means we need to create the deployment for it. In the `MemachedController.java`, find the `createOrUpdateResource` method and add the following code.

Below code will verify that Deployment within the cluster got created or not. If deployment is null then it will create deployment. `createMemcachedDeployment(resource)` creates the Deployment and then it will be applied by using `client.apps().deployments().create(newDeployment);` code. `createMemcachedDeployment(resource)` method explained in the next part.

```
        if (deployment == null) {
            Deployment newDeployment = createMemcachedDeployment(resource);
            client.apps().deployments().create(newDeployment);
            return UpdateControl.noUpdate();
        }
```

Once we create the deployment, we need to decide whether we have to reconcile it or not.
If there is no need of reconciliation then return `UpdateControl.noUpdate()` else we need to return `UpdateControl.updateStatusSubResource(resource)`

After the Creation of the Deployment, get the current and required replicas by using the below code.

```
        int currentReplicas = deployment.getSpec().getReplicas();
        int requiredReplicas = resource.getSpec().getSize();
```

If currentReplicas does not match with requiredReplicas then we need to update the `Deployment`. This process will be done by using the below code.

```
        if (currentReplicas != requiredReplicas) {
            deployment.getSpec().setReplicas(requiredReplicas);
            client.apps().deployments().createOrReplace(deployment);
            return UpdateControl.noUpdate();
        }
```

Then, let's get the list of pods and pod names. In the `MemachedController.java`, find the `createOrUpdateResource` method and add the following code. It should look like:

```
        List<Pod> pods =
                client
                        .pods()
                        .inNamespace(resource.getMetadata().getNamespace())
                        .withLabels(labelsForMemcached(resource))
                        .list()
                        .getItems();

        List<String> podNames =
                pods.stream().map(p -> p.getMetadata().getName()).collect(Collectors.toList());
```

Now, check whether resources get created or not. Then, it verifies podnames with the Memcached resources. If there is a mismatch in either of these conditions then we need to do a reconciliation.

```
        if (resource.getStatus() == null
                || !CollectionUtils.isEqualCollection(podNames, resource.getStatus().getNodes())) {
            if (resource.getStatus() == null) resource.setStatus(new MemcachedStatus());
            resource.getStatus().setNodes(podNames);
            return UpdateControl.updateStatusSubResource(resource);
        }
```

The complete `createOrUpdateResource` function will look like below in `MemachedController.java`.

```
    @Override
    public UpdateControl<Memcached> createOrUpdateResource(
        Memcached resource, Context<Memcached> context) {
        // TODO: fill in logic
        System.out.println("Create or Update Control");
        Deployment deployment =
                client
                        .apps()
                        .deployments()
                        .inNamespace(resource.getMetadata().getNamespace())
                        .withName(resource.getMetadata().getName())
                        .get();

        System.out.println(deployment);


        if (deployment == null) {
            Deployment newDeployment = createMemcachedDeployment(resource);
            client.apps().deployments().create(newDeployment);
            return UpdateControl.noUpdate();
        }

        int currentReplicas = deployment.getSpec().getReplicas();
        int requiredReplicas = resource.getSpec().getSize();
        if (currentReplicas != requiredReplicas) {
            deployment.getSpec().setReplicas(requiredReplicas);
            client.apps().deployments().createOrReplace(deployment);
            return UpdateControl.noUpdate();
        }

        List<Pod> pods =
                client
                        .pods()
                        .inNamespace(resource.getMetadata().getNamespace())
                        .withLabels(labelsForMemcached(resource))
                        .list()
                        .getItems();

        List<String> podNames =
                pods.stream().map(p -> p.getMetadata().getName()).collect(Collectors.toList());

        if (resource.getStatus() == null
                || !CollectionUtils.isEqualCollection(podNames, resource.getStatus().getNodes())) {
            if (resource.getStatus() == null) resource.setStatus(new MemcachedStatus());
            resource.getStatus().setNodes(podNames);
            return UpdateControl.updateStatusSubResource(resource);
        }

        return UpdateControl.noUpdate();
    }
```

### createMemcachedDeployment

Create `createMemcachedDeployment` method and add the below code snippet. This method simply creates the Deployment.

```
    private Deployment createMemcachedDeployment(Memcached m) {
        return new DeploymentBuilder()
            .withMetadata(
                new ObjectMetaBuilder()
                    .withName(m.getMetadata().getName())
                    .withNamespace(m.getMetadata().getNamespace())
                    .withOwnerReferences(
                        new OwnerReferenceBuilder()
                            .withApiVersion("v1alpha1")
                            .withKind("Memcached")
                            .withName(m.getMetadata().getName())
                            .withUid(m.getMetadata().getUid())
                            .build())
                    .build())
            .withSpec(
                new DeploymentSpecBuilder()
                    .withReplicas(m.getSpec().getSize())
                    .withSelector(
                        new LabelSelectorBuilder().withMatchLabels(labelsForMemcached(m)).build())
                    .withTemplate(
                        new PodTemplateSpecBuilder()
                            .withMetadata(
                                new ObjectMetaBuilder().withLabels(labelsForMemcached(m)).build())
                            .withSpec(
                                new PodSpecBuilder()
                                    .withContainers(
                                        new ContainerBuilder()
                                            .withImage("memcached:1.4.36-alpine")
                                            .withName("memcached")
                                            .withCommand("memcached", "-m=64", "-o", "modern", "-v")
                                            .withPorts(
                                                new ContainerPortBuilder()
                                                    .withContainerPort(11211)
                                                    .withName("memcached")
                                                    .build())
                                            .build())
                                    .build())
                            .build())
                    .build())
            .build();
    }
```

### deleteResource

The `deleteResource` method is an implemented method. The deletion part will be taken care of by the Java Operator SDK  library, that's why it is empty.
The code snippet for `deleteResource` is as shown below.

```
    @Override
    public DeleteControl deleteResource(Memcached resource, Context<Memcached> context) {
        // nothing to do here...
        // framework takes care of deleting the resource object
        // k8s takes care of deleting deployment and pods because of ownerreference set
        System.out.println("Delete Control");
        return DeleteControl.DEFAULT_DELETE;
    }
```

## Run the Operator Locally

### Run locally outside the cluster

The following steps will show how to run your operator locally.

Compile your operator with the below command

`mvn clean install`

It will create a `.jar` file for your operator in `target/quarkus-app`.  Now, run the `jar` file using the below command.

`java -jar quarkus-run.jar`

This command will run your operator locally. You can check cluster pods and deployment with the below commands.

`kubectl get deployment`

`kubectl get pods`

Delete One of the Pod forcefully then Memcached operator will create new automatically.

`kubectl delete pod pod-name`

In the end, change the size from the memcached-sample.yaml file and apply it to the cluster. After these steps, an operator will make sure that the cluster has an updated number of pods in it.
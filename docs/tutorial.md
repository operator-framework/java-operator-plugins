# Java Operator Tutorial
### An in-depth walkthrough of building and running a Java-based operator.

## Prerequisites

- [Operator SDK](https://sdk.operatorframework.io/docs/installation/) v1.8.0 or newer
- [Java](https://java.com/en/download/help/download_options.html) 11
- [Maven 3.6.3](https://maven.apache.org/install.html) or newer
- User authorized with `cluster-admin` permissions.

## Overview

We will create a sample project to let you know how it works and this sample will:

- Create a Memcached Deployment if it doesn't exist
- Ensure that the Deployment size is the same as specified by the Memcached Custom Resource (CR) spec
- Update the Memcached CR status using the status writer with the names of the CR's pods

## Create a new project

Use the [Operator SDK](https://sdk.operatorframework.io/docs/installation/) CLI to create a
new memcached-quarkus-operator project:

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

The quarkus plugin will scaffold out several files during the `init` phase. One
of these files is the operator's main program, `MemcachedQuarkusOperator.java`.
This file initializes and runs the operator. The operator uses java-operator-sdk,
which is similar to
[controller-runtime](https://github.com/kubernetes-sigs/controller-runtime), to
make operator development easier.

The important part of the `MemcachedQuarkusOperator.java` is the `run` method
which will start the operator and initializes the informers and watches for your
operator.

Here is an example of the `run` method that will typically be scaffolded out by
this plugin:

```
  @Override
  public int run(String... args) throws Exception {
    operator.start();

    Quarkus.waitForExit();
    return 0;
  }
```

## Create a new API and Controller

An operator isn't much good without an API to work with. Create a new Custom
Resource Definition (CRD) API with group `cache`, version `v1`, and Kind
`Memcached`.

Use the `create api` command to scaffold the `MemcachedController`,
`MemcachedSpec`, `MemcachedStatus` and `Memcached`. These files represent the
API. The plugin may show some debug statements which is normal as it is still in
the alpha state.

```console
$ operator-sdk create api --plugins quarkus --group cache --version v1 --kind Memcached
```

After running the `create api` command the file structure will change to match the
one shown as below.

```
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
```

#### Understanding Kubernetes APIs

For an in-depth explanation of Kubernetes APIs and the group-version-kind model, check out these [kubebuilder docs](https://book.kubebuilder.io/cronjob-tutorial/gvks.html).

The `java-operator-plugins` project uses the APIs from [java-operator-sdk](https://github.com/java-operator-sdk/java-operator-sdk). The java-operator-sdk library is to Java what `controller-runtime` is to Go. It provides a framework to help Java projects interact with the Kubernetes API.

### Define the API

#### `MemcachedSpec`

Initially, the scaffolded Spec file will be empty. The operator developer needs
to add attributes to this file according to their needs. For the `Memcached`
example, we will add the size field as shown in the example below.

The `MemcachedSpec` class defines the desired state of `Memcached`.

```
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

As you can see, we added the `size` attribute along with corresponding getter
and setter.

#### `MemcachedStatus`

Similar to the Spec file above, the `MemcachedStatus` file was scaffolded as
part of the `create api` command. The user will need to modify a Status file
in order to add any desired attributes. For this `Memcached` example, we will
add a list of nodes as shown below.

The nodes field is a list of string values and it contains the name of
the Memcached pods.  The `MemcachedStatus` defines the observed state
of `Memcached`.


```
import java.util.ArrayList;
import java.util.List;

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

**Note** The Node field is just to illustrate an example of a Status field. In
real use cases, it is recommended that you use
[Conditions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties).

#### `Memcached`

Now that we have Spec and Status classes, let's look at the `Memcached` class.
This file was also scaffolded via `create api` command. Notice it extends both
`MemcachedSpec` and `MemcachedStatus`.

The `Memcached` is the Schema for the Memcacheds API.

```

@Version("v1alpha1")
@Group("cache.example.com")
public class Memcached extends CustomResource<MemcachedSpec, MemcachedStatus>
    implements Namespaced {}
```

You have now created the necessary classes for the API.

### Apply Custom Resource and CRD's using below command

There are a couple of ways to create the CRD. You can either create the file
manually. OR let the quarkus extensions defined in pom.xml use the annotations
on your Spec/Status classes to create the crd files for you.

#### Manually create `crd.yaml`

Create a file with the name `crd.yaml`. A CRD enables users to add their
own/custom objects to the Kubernetes cluster. Below you will find an example
`Memcached` CRD.

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
        type: object
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

### TODO: Creating the CRD seems out of place
Create the CRD:

`kubectl apply -f crd.yaml`
#########

### Create sample Memcached Custom Resource

Let's create the sample Memcached Custom Resource manifest at `memcached-sample.yaml` and define the spec as the following.

```
apiVersion: cache.example.com/v1
kind: Memcached
metadata:
  name: memcached-sample
spec:
  # Add fields here
  size: 1
```

### TODO out of place
Create the Custom Resource:

`kubectl apply -f memcached-sample.yaml`
################3

## Implement the Controller

By now we have the API defined in `Memcached.java`, `MemcachedSpec.java`,
`MemcachedStatus.java`. We also have the CRD and the sample Custom Resource.
This isn't enough, we still need a controller to reconcile these items.

The `create api` command will have scaffolded a skeleton `MemcachedController.java`.
This controller implements the `ResourceController` interface from the
`java-operator-sdk`. This interface has some important and useful methods.

Initially the `MemcachedController.java` will contain the empty stubs for
`createOrUpdateResource` and `deleteResource`. In this section we will fill in
the controller logic in these methods. We will also add a
`createMemcachedDeployment` method that will create the Deployment for our
operator.

The `createOrUpdateResource` and `deleteResource` get called whenever some
update/create/delete event occurs in the cluster. This will allow us to react to
changes to the Deployment.

### createOrUpdateResource

In this section we will focus on implementing the `createOrUpdateResource`
method. In the `MemcachedController.java` you will see a `// TODO: fill in logic`
comment. At this line we will first add code to get the Deployment.

```
        Deployment deployment = client.apps()
                .deployments()
                .inNamespace(resource.getMetadata().getNamespace())
                .withName(resource.getMetadata().getName())
                .get();

```

Once we get the `deployment`, we have a couple of decisions to make. If it is
`null` it does not exist which means we need to create the deployment. In the
`MemachedController.java`, in the `createOrUpdateResource` method just below the
get deployment code we added above, add the following:

```
        if (deployment == null) {
            Deployment newDeployment = createMemcachedDeployment(resource);
            client.apps().deployments().create(newDeployment);
            return UpdateControl.noUpdate();
        }
```

In the above code, we are checking to see if the deployment exists, if not we
will create it by calling the yet to be defined `createMemcachedDeployment`
method.


##############################
Below code will verify that Deployment within the cluster got created or not. If deployment is null then it will create deployment. `createMemcachedDeployment(resource)` creates the Deployment and then it will be applied by using `client.apps().deployments().create(newDeployment);` code. `createMemcachedDeployment(resource)` method explained in the next part.

##############################

Once we create the deployment, we need to decide whether we have to reconcile it or not.
If there is no need of reconciliation then return `UpdateControl.noUpdate()`
else we need to return `UpdateControl.updateStatusSubResource(resource)`

After getting the Deployment, we get the current and required replicas. Add the
following lines below the `if (deployment == null)` block in your
`MemcachedController.java` file.

```
        int currentReplicas = deployment.getSpec().getReplicas();
        int requiredReplicas = resource.getSpec().getSize();
```

Once we get the replicas, we need to determine if they are different so that we
can reconcile. If `currentReplicas` does not match the `requiredReplicas` then
we need to update the `Deployment`. Add the following comparison block to your
controller.

```
        if (currentReplicas != requiredReplicas) {
            deployment.getSpec().setReplicas(requiredReplicas);
            client.apps().deployments().createOrReplace(deployment);
            return UpdateControl.noUpdate();
        }
```

The above sections will cover reconciling any `size` changes to the Spec. In the
next section, we will look at handling the changes to the `nodes` list from the
Status.

Let's get the list of pods and their names. In the `MemcachedController.java`,
add the following code below the `if (currentReplicas != requiredReplicas) {`
block.

```
        List<Pod> pods = client.pods()
            .inNamespace(resource.getMetadata().getNamespace())
            .withLabels(labelsForMemcached(resource))
            .list()
            .getItems();

        List<String> podNames =
            pods.stream().map(p -> p.getMetadata().getName()).collect(Collectors.toList());
```

Now that we have the pods and names. What do we do next? Well, we check whether resources
were created. Then, it verifies podnames with the Memcached resources. If there is a
mismatch in either of these conditions then we need to do a reconciliation.

```
        if (resource.getStatus() == null
                || !CollectionUtils.isEqualCollection(podNames, resource.getStatus().getNodes())) {
            if (resource.getStatus() == null) resource.setStatus(new MemcachedStatus());
            resource.getStatus().setNodes(podNames);
            return UpdateControl.updateStatusSubResource(resource);
        }
```

That's it we have completed the `createOrUpdateResource` method. The method
should now look like the following:

```
    @Override
    public UpdateControl<Memcached> createOrUpdateResource(
        Memcached resource, Context<Memcached> context) {
        // TODO: fill in logic
        Deployment deployment = client.apps()
                .deployments()
                .inNamespace(resource.getMetadata().getNamespace())
                .withName(resource.getMetadata().getName())
                .get();

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

        List<Pod> pods = client.pods()
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

Let's recap what we did.

* Get the Deployment, if the deployment does not exist, we create it.
* If the deployment already exists, we get the current replicas and the desired
  replicas.
* We compare the replicas, if they do not match, we replace the deployment with
  the expected values.
* Next we look at the node list from the pods. If they do not match, we update
  and reconcile.

What's left? If you recall, in the if the deployment is `null`, we call
`createMemcachedDeployment(resource)`. This method still needs to get created.
In the next section, we will walk you through creating this helper method.

### createMemcachedDeployment

Creating Kubernetes objects via APIs can be quite verbose which is why putting
them in helper methods can make the code more readable. The
`MemcachedController.java` needs to create a Deployment if it does not exist. In
the `createOrUpdateResource` we make a call to a helper,
`createMemcachedDeployment`.

Let's create the `createMemcachedDeployment` method. The following code will use
the [`fabric8`](https://fabric8.io/) `DeploymentBuilder` class. Notice the
Deployment specifies the `memcached` image for the pod.

Below your `deleteResource(Memcached resource, Context<Memcached> context) {`
block in the `MemcachedController.java`, add the following method.

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

We have not implemented the `MemcachedController.java`.

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

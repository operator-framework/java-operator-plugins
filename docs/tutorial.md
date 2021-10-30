# Java Operator Tutorial
### An in-depth walkthrough of building and running a Java-based operator.

## Prerequisites

- [Operator SDK](https://sdk.operatorframework.io/docs/installation/) v1.8.0 or newer
- [Java](https://java.com/en/download/help/download_options.html) 11
- [Maven 3.6.3](https://maven.apache.org/install.html) or newer
- User authorized with `cluster-admin` permissions.
- [GNU Make](https://www.gnu.org/software/make/)

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

**Note** Please do not commit this file structure to the `GitHub` immediately after the `init` command. The directory structure does not contain any file, and GitHub will not create an empty directory.

#### A note on dependency management

`operator-sdk init` generates `pom.xml` file. This file contains all the dependencies required to run the operator.


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
├── Makefile
├── PROJECT
├── pom.xml
└── src
    └── main
        ├── java
        │   └── com
        │       └── example
        │           ├── Memcached.java
        │           ├── MemcachedController.java
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

@Version("v1")
@Group("cache.example.com")
@Kind("Memcached")
@Plural("memcacheds")
public class Memcached extends CustomResource<MemcachedSpec, MemcachedStatus>
    implements Namespaced {}
```

You have now created the necessary classes for the API.

### Creating Custom Resource and CRD

There are a couple of ways to create the CRD. You can either create the file
manually. Or let the quarkus extensions defined in `pom.xml` use the annotations
on your Spec/Status classes to create the crd files for you.

#### Via Quarkus extension

Running `mvn clean install` will invoke the CRD generator extension which will analyze
the annotations on the model objects,  `Memcached`, `MemcachedSpec`,
`MemcachedStatus`, and generate the CRD in `target/kubernetes`.

CRD generated in `memcacheds.cache.example.com-v1.yml`.

```
.
├── kubernetes.json
├── kubernetes.yml
└── memcacheds.cache.example.com-v1.yml

0 directories, 3 files
```

The content of the `memcacheds.cache.example.com-v1.yml` file is as shown below.

```
# Generated by Fabric8 CRDGenerator, manual edits might get overwritten!
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: memcacheds.cache.example.com
spec:
  group: cache.example.com
  names:
    kind: Memcached
    plural: memcacheds
    singular: memcached
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              size:
                type: integer
            type: object
          status:
            properties:
              nodes:
                items:
                  type: string
                type: array
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
```

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

## Implement the Controller

By now we have the API defined in `Memcached.java`, `MemcachedSpec.java`,
`MemcachedStatus.java`. We also have the CRD and the sample Custom Resource.
This isn't enough, we still need a controller to reconcile these items.

The `create api` command will have scaffolded a skeleton `MemcachedController.java`.
This controller implements the `ResourceController` interface from the
`java-operator-sdk`. This interface has some important and useful methods.

Initially the `MemcachedController.java` will contain the empty stubs for
`createOrUpdateResource`. In this section we will fill in
the controller logic in these methods. We will also add a
`createMemcachedDeployment` method that will create the Deployment for our
operator and a `labelsForMemcached` method that returns the labels.

The `createOrUpdateResource` get called whenever some
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
As well as the `labelsForMemcached` utility method.

Let's create the utility method first.

### labelsForMemcached

A simple utility method to return a map of the labels we want to attach to some
of the resources. Below the `deleteResource` method add the following
helper:

```
    private Map<String, String> labelsForMemcached(Memcached m) {
        Map<String, String> labels = new HashMap<>();
        labels.put("app", "memcached");
        labels.put("memcached_cr", m.getMetadata().getName());
        return labels;
    }
```

In the next section, we will walk you through creating the
`createMemcachedDeployment` utility method.

### createMemcachedDeployment

Creating Kubernetes objects via APIs can be quite verbose which is why putting
them in helper methods can make the code more readable. The
`MemcachedController.java` needs to create a Deployment if it does not exist. In
the `createOrUpdateResource` we make a call to a helper,
`createMemcachedDeployment`.

Let's create the `createMemcachedDeployment` method. The following code will use
the [`fabric8`](https://fabric8.io/) `DeploymentBuilder` class. Notice the
Deployment specifies the `memcached` image for the pod.

Below your `labelsForMemcached(Memcached m)` block in the
`MemcachedController.java`, add the following method.

```
    private Deployment createMemcachedDeployment(Memcached m) {
        return new DeploymentBuilder()
            .withMetadata(
                new ObjectMetaBuilder()
                    .withName(m.getMetadata().getName())
                    .withNamespace(m.getMetadata().getNamespace())
                    .withOwnerReferences(
                        new OwnerReferenceBuilder()
                            .withApiVersion("v1")
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

Now we have a `createOrUpdateResource` method. It calls
`createMemcachedDeployment` which we have implemented above. In the next section
we will discuss the deletion of the resource.

We have now implemented the `MemcachedController.java`.

## Run the Operator

You can run the operator in a couple of ways. You can run it locally where the
operator runs on your development machine and talks to the cluster. Or it can
build images of your operator and run it directly in the cluster.

In this section we will:

* install the CRD
* create a Custom Resource
* run your operator

If you want to run the operator in the cluster see the [Running the operator in
the cluster](#running-the-operator-in-the-cluster) below or if you'd prefer to
run it locally see the
[Running locally outside the cluster](#running-locally-outside-the-cluster)
section instead.

### Running the operator in the cluster

The following steps will show how to run your operator in the cluster.

1. Build and push your operator's image:

The `java-operator-plugins` project will scaffold out a Makefile to give
Operator SDK users a familiar interface. Using the `docker-*` targets you can
conveniently build your and push your operator's image to registry. In our
example, we are using `quay.io`, but any docker registry should work.

```
make docker-build docker-push IMG=quay.io/YOURUSER/memcached-quarkus-operator:0.0.1
```

This will build the docker image
`quay.io/YOURUSER/memcached-quarkus-operator:0.0.1` and push it to the registry.

You can verify it is in your docker registry:

```
$ docker images | grep memcached
quay.io/YOURUSER/memcached-quarkus-operator                   0.0.1               c84d2616bc1b        29 seconds ago       236MB
```

2. Install the CRD

Next we will install the CRD into the `default` namespace. Using the `crd.yaml`
you created in the [Manually created crd.yaml](#manually-create-crdyaml)
section, apply it to the cluster.

<!--
TODO: Uncomment this when the crd generator works properly.
```
make install
customresourcedefinition.apiextensions.k8s.io/memcacheds.cache.example.com created
```
-->

```
$ kubectl apply -f crd.yaml
customresourcedefinition.apiextensions.k8s.io/memcacheds.cache.example.com created
```

3.  Create rbac.yaml file

The RBAC generated in the `kubernetes.yml` only has [view
permissions](https://quarkus.io/guides/deploying-to-kubernetes#using-the-kubernetes-client)
which is not enough to run the operator. For this example, we will simply grant
cluster-admin to the `memcached-quarkus-operator-operator` service account.

Create a file called `rbac.yaml` with the following contents:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: memcached-operator-admin
subjects:
- kind: ServiceAccount
  name: memcached-quarkus-operator-operator
  namespace: default
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: ""
```

Do not apply this yet. We will do that in a later step.

4. Deploy the operator

Let's deploy your operator to the cluster. The `Makefile` has a convenience
target that can do this for you:

```
make deploy
```

5. Grant `cluster-admin` to service account

Once you've deployed the operator, you will need to grant the
`memcached-quarkus-operator-operator` service account the right privileges.

```
kubectl apply -f rbac.yaml
```

6. Verify the operator is running

Ensure the `memcached-quarkus-operator-operator-XXX` pod is in a `Running`
status.

```
$ kubectl get all -n default
NAME                                                       READY     STATUS    RESTARTS   AGE
pod/memcached-quarkus-operator-operator-7db86ccf58-k4mlm   0/1       Running   0          18s
...
```

7. Apply the memcached-sample

Apply the memcached-sample to see the operator create the memcached-sample pod.

```
$ kubectl apply -f memcached-sample.yaml
memcached.cache.example.com/memcached-sample created
```

8. Verify the sample

Now check the cluster to see if the pod has started. Keep watching until the
`memcached-sample-XXX` pod reaches a `Running` status.

```
$ kubectl get all
NAME                                                       READY   STATUS    RESTARTS   AGE
pod/memcached-quarkus-operator-operator-7b766f4896-kxnzt   1/1     Running   1          79s
pod/memcached-sample-6c765df685-mfqnz                      1/1     Running   0          18s
...
```

9. Trigger a reconcile

If you modify the size field of the `memcached-sample.yaml` and re-apply it. The
operator will trigger a reconcile and adjust the sample pods to the size given.

### Running locally outside the cluster

For development purposes, you may want to run your operator locally for faster
iteration. In the following steps, we will show how to run  your operator
locally.

1. Compile your operator with the below command

```
mvn clean install
```

You should see a nice `BUILD SUCCESS` method like the one below:

```
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  11.193 s
[INFO] Finished at: 2021-05-26T12:16:54-04:00
[INFO] ------------------------------------------------------------------------
```

2. Install the CRD

Next we will install the CRD into the `default` namespace. Using the `crd.yaml`
you created in the [Manually created crd.yaml](#manually-create-crdyaml)
section, apply it to the cluster.

```
$ kubectl apply -f crd.yaml
customresourcedefinition.apiextensions.k8s.io/memcacheds.cache.example.com created
```

3.  Create and apply rbac.yaml file

The RBAC generated in the `kubernetes.yml` only has [view
permissions](https://quarkus.io/guides/deploying-to-kubernetes#using-the-kubernetes-client)
which is not enough to run the operator. For this example, we will simply grant
cluster-admin to the `memcached-quarkus-operator-operator` service account.

Create a file called `rbac.yaml` with the following contents:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: memcached-operator-admin
subjects:
- kind: ServiceAccount
  name: memcached-quarkus-operator-operator
  namespace: default
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: ""
```

Let's apply the rbac role to the cluster:

```
kubectl apply -f rbac.yaml
```

4. Run the operator

After running this command, notice there is now a `target` directory. That
directory may be a bit overwhelming, but the key thing to know for running locally
is the `target/memcached-quarkus-operator-0.0.1.jar` file created for your
operator and the `quarkus-run.jar` in `target/quarkus-app` directory.

Now, run the `jar` file using the below command. This command will run your
opeator locally.

```
java -jar target/quarkus-app/quarkus-run.jar
```

**Note** the above will run the operator and remain running until you kill it.
You will need another terminal to complete the rest of these commands.

5. Apply the memcached-sample

Apply the memcached-sample to see the operator create the memcached-sample pod.

```
$ kubectl apply -f memcached-sample.yaml
memcached.cache.example.com/memcached-sample created
```

6. Verify the sample

Now check the cluster to see if the pod has started. Keep watching until the
`memcached-sample-XXX` pod reaches a `Running` status.

```
$ kubectl get all
NAME                                                       READY   STATUS    RESTARTS   AGE
pod/memcached-sample-6c765df685-mfqnz                      1/1     Running   0          18s
...
```

7. Trigger a reconcile

If you modify the size field of the `memcached-sample.yaml` and re-apply it. The
operator will trigger a reconcile and adjust the sample pods to the size given.

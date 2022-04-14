package com.example;

import io.fabric8.kubernetes.api.model.ContainerBuilder;
import io.fabric8.kubernetes.api.model.ContainerPortBuilder;
import io.fabric8.kubernetes.api.model.LabelSelectorBuilder;
import io.fabric8.kubernetes.api.model.ObjectMetaBuilder;
import io.fabric8.kubernetes.api.model.OwnerReferenceBuilder;
import io.fabric8.kubernetes.api.model.Pod;
import io.fabric8.kubernetes.api.model.PodSpecBuilder;
import io.fabric8.kubernetes.api.model.PodTemplateSpecBuilder;
import io.fabric8.kubernetes.api.model.apps.Deployment;
import io.fabric8.kubernetes.api.model.apps.DeploymentBuilder;
import io.fabric8.kubernetes.api.model.apps.DeploymentSpecBuilder;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.javaoperatorsdk.operator.api.reconciler.Context;
import io.javaoperatorsdk.operator.api.reconciler.Reconciler;
import io.javaoperatorsdk.operator.api.reconciler.UpdateControl;
import org.apache.commons.collections.CollectionUtils;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class MemcachedReconciler implements Reconciler<Memcached> { 
  private final KubernetesClient client;

  public MemcachedReconciler(KubernetesClient client) {
    this.client = client;
  }


  @Override
  public UpdateControl<Memcached> reconcile(Memcached resource, Context context) {

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
      return UpdateControl.updateStatus(resource);
    }

    return UpdateControl.noUpdate();
  }

  private Map<String, String> labelsForMemcached(Memcached m) {
    Map<String, String> labels = new HashMap<>();
    labels.put("app", "memcached");
    labels.put("memcached_cr", m.getMetadata().getName());
    return labels;
  }

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
}


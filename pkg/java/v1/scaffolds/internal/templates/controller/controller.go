package controller

import (
	"fmt"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds/internal/templates/util"
)

var _ machinery.Template = &Controller{}

type Controller struct {
	machinery.TemplateMixin

	// Package is the source files package
	Package string

	// Name of the operator used for the main file.
	ClassName string
}

func (f *Controller) SetTemplateDefaults() error {
	if f.ClassName == "" {
		return fmt.Errorf("invalid model name")
	}

	if f.Path == "" {
		f.Path = util.PrependJavaPath(f.ClassName+"Controller.java", util.AsPath(f.Package))
	}

	f.TemplateBody = controllerTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const controllerTemplate = `package {{ .Package }};

import io.fabric8.kubernetes.api.model.*;
import io.fabric8.kubernetes.api.model.apps.Deployment;
import io.fabric8.kubernetes.api.model.apps.DeploymentBuilder;
import io.fabric8.kubernetes.api.model.apps.DeploymentSpecBuilder;
import io.fabric8.kubernetes.client.KubernetesClient;
import io.javaoperatorsdk.operator.api.*;
import io.javaoperatorsdk.operator.api.Context;
import io.javaoperatorsdk.operator.processing.event.EventSourceManager;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import org.apache.commons.collections.CollectionUtils;

@Controller
public class {{ .ClassName }}Controller implements ResourceController<{{ .ClassName }}> {

    private final KubernetesClient client;

    public {{ .ClassName }}Controller(KubernetesClient client) {
        this.client = client;
    }

    // TODO Fill in the rest of the controller

    @Override
    public void init(EventSourceManager eventSourceManager) {
        // TODO: fill in init
    }

    @Override
    public UpdateControl<{{ .ClassName }}> createOrUpdateResource(
        {{ .ClassName }} resource, Context<{{ .ClassName }}> context) {
        // TODO: fill in logic

        return UpdateControl.noUpdate();
    }

    @Override
    public DeleteControl deleteResource({{ .ClassName }} resource, Context<{{ .ClassName }}> context) {
        // nothing to do here...
        // framework takes care of deleting the resource object
        // k8s takes care of deleting deployment and pods because of ownerreference set
        return DeleteControl.DEFAULT_DELETE;
    }
}

`

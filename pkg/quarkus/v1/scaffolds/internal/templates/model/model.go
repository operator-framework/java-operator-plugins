package model

import (
	"fmt"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/quarkus/v1/scaffolds/internal/templates/util"
)

var _ machinery.Template = &Model{}

type Model struct {
	machinery.TemplateMixin

	// Package is the source files package
	Package string

	// Name of the operator used for the main file.
	ClassName string

	Version string
	Group   string
}

func (f *Model) SetTemplateDefaults() error {
	if f.ClassName == "" {
		return fmt.Errorf("invalid model name")
	}

	if f.Path == "" {
		f.Path = util.PrependJavaPath(f.ClassName+".java", util.AsPath(f.Package))
	}

	f.TemplateBody = modelTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const modelTemplate = `package {{ .Package }};

import io.fabric8.kubernetes.api.model.Namespaced;
import io.fabric8.kubernetes.client.CustomResource;
import io.fabric8.kubernetes.model.annotation.Group;
import io.fabric8.kubernetes.model.annotation.Version;

@Version("{{ .Version }}")
@Group("{{ .Group }}")
public class {{ .ClassName }} extends CustomResource<{{ .ClassName }}Spec, {{ .ClassName }}Status>
    implements Namespaced {}

`

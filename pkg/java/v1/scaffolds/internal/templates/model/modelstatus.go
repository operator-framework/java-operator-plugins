package model

import (
	"fmt"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds/internal/templates/util"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &ModelStatus{}

type ModelStatus struct {
	machinery.TemplateMixin

	// Package is the source files package
	Package string

	// Name of the operator used for the main file.
	ClassName string
}

func (f *ModelStatus) SetTemplateDefaults() error {
	if f.ClassName == "" {
		return fmt.Errorf("invalid operator name")
	}

	if f.Path == "" {
		f.Path = util.PrependJavaPath(f.ClassName+"Status.java", util.AsPath(f.Package))
	}

	f.TemplateBody = modelStatusTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const modelStatusTemplate = `package {{ .Package }};

public class {{ .ClassName }}Status {

    // Add Status information here
}
`

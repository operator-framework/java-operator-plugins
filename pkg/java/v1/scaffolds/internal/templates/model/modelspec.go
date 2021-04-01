package model

import (
	"fmt"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds/internal/templates/util"
)

var _ machinery.Template = &ModelSpec{}

type ModelSpec struct {
	machinery.TemplateMixin

	// Package is the source files package
	Package string

	// Name of the operator used for the main file.
	ClassName string
}

func (f *ModelSpec) SetTemplateDefaults() error {
	if f.ClassName == "" {
		return fmt.Errorf("invalid operator name")
	}

	if f.Path == "" {
		f.Path = util.PrependJavaPath(f.ClassName+"Spec.java", util.AsPath(f.Package))
	}

	f.TemplateBody = modelSpecTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const modelSpecTemplate = `package {{ .Package }};

public class {{ .ClassName }}Spec {

	// Add Spec information here
}
`

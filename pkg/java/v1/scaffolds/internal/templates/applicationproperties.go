package templates

import (
	"fmt"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kubebuilder/v3/pkg/model/file"
)

var _ file.Template = &ApplicationPropertiesFile{}

type ApplicationPropertiesFile struct {
	file.TemplateMixin
	OrgName         string
	ProjectName     string
}

const (
	FilePathSeparator = string(filepath.Separator)
	javaPaths = "src" + FilePathSeparator + "main" + FilePathSeparator + "resources"
)

func prependJavaPathResources(filename string) string {
	return javaPaths + FilePathSeparator + filename
}

func asResourcesPath(s string) string {
	return strings.ReplaceAll(s, ".", FilePathSeparator)
}

func (f *ApplicationPropertiesFile) SetTemplateDefaults() error {
	if f.ProjectName == "" {
		return fmt.Errorf("invalid Application Properties name")
	}

	if f.Path == "" {
		f.Path = prependJavaPathResources("application.properties")
	}

	f.TemplateBody = ApplicationPropertiesTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const ApplicationPropertiesTemplate = `quarkus.container-image.build=true
#quarkus.container-image.group=
quarkus.container-image.name={{ .ProjectName }}-service
`

package scaffolds

import (
	"fmt"

	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds/internal/templates/controller"
	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds/internal/templates/model"
	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/util"
)

type apiScaffolder struct {
	fs machinery.Filesystem

	config   config.Config
	resource resource.Resource
}

// NewCreateAPIScaffolder returns a new plugins.Scaffolder for project initialization operations
func NewCreateAPIScaffolder(cfg config.Config, res resource.Resource) plugins.Scaffolder {
	fmt.Println("NewCreateAPIScaffolder called")
	return &apiScaffolder{
		config:   cfg,
		resource: res,
	}
}

func (s *apiScaffolder) InjectFS(fs machinery.Filesystem) {
	fmt.Println("InjectFS called")
	s.fs = fs
}

func (s *apiScaffolder) Scaffold() error {
	fmt.Println("api.Scaffold()")

	if err := s.config.UpdateResource(s.resource); err != nil {
		return err
	}

	// Initialize the machinery.Scaffold that will write the files to disk
	scaffold := machinery.NewScaffold(s.fs,
		// NOTE: kubebuilder's default permissions are only for root users
		machinery.WithDirectoryPermissions(0755),
		machinery.WithFilePermissions(0644),
		machinery.WithConfig(s.config),
		machinery.WithResource(&s.resource),
	)

	var createAPITemplates []machinery.Builder
	createAPITemplates = append(createAPITemplates,
		&model.Model{
			Package:   util.ReverseDomain(s.config.GetDomain()),
			ClassName: util.ToClassname(s.resource.Kind),
			Version:   s.resource.Version,
			Group:     s.resource.Group,
		},
		&model.ModelSpec{
			Package:   util.ReverseDomain(s.config.GetDomain()),
			ClassName: util.ToClassname(s.resource.Kind),
		},
		&model.ModelStatus{
			Package:   util.ReverseDomain(s.config.GetDomain()),
			ClassName: util.ToClassname(s.resource.Kind),
		},
		&controller.Controller{
			Package:   util.ReverseDomain(s.config.GetDomain()),
			ClassName: util.ToClassname(s.resource.Kind),
		},
	)

	return scaffold.Execute(createAPITemplates...)
}

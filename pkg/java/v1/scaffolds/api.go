package scaffolds

import (
	"fmt"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/internal/kubebuilder/cmdutil"
)

type apiScaffolder struct {
	fs machinery.Filesystem
}

func NewCreateAPIScaffolder() cmdutil.Scaffolder {
	return &apiScaffolder{}
}

func (s *apiScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

func (s *apiScaffolder) Scaffold() error {
	fmt.Println("api.Scaffold()")
	return nil
}

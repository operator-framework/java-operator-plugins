// Copyright 2021 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/model/resource"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
	pluginutil "sigs.k8s.io/kubebuilder/v3/pkg/plugin/util"

	"github.com/operator-framework/java-operator-plugins/pkg/quarkus/v1alpha/scaffolds"
)

const filePath = "Makefile"

type createAPIOptions struct {
	CRDVersion string
	Namespaced bool
}

type createAPISubcommand struct {
	config   config.Config
	resource *resource.Resource
	options  createAPIOptions
}

func (opts createAPIOptions) UpdateResource(res *resource.Resource) {

	res.API = &resource.API{
		CRDVersion: opts.CRDVersion,
		Namespaced: opts.Namespaced,
	}

	// Ensure that Path is empty and Controller false as this is not a Go project
	res.Path = ""
	res.Controller = false
}

var (
	_ plugin.CreateAPISubcommand = &createAPISubcommand{}
)

func (p *createAPISubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false
	fs.StringVar(&p.options.CRDVersion, "crd-version", "v1", "crd version to generate")
	fs.BoolVar(&p.options.Namespaced, "namespaced", true, "resource is namespaced")
}

func (p *createAPISubcommand) InjectConfig(c config.Config) error {
	p.config = c

	return nil
}

func (p *createAPISubcommand) Run(fs machinery.Filesystem) error {
	return nil
}

func (p *createAPISubcommand) Validate() error {
	return nil
}

func (p *createAPISubcommand) PostScaffold() error {
	return nil
}

func (p *createAPISubcommand) Scaffold(fs machinery.Filesystem) error {
	scaffolder := scaffolds.NewCreateAPIScaffolder(p.config, *p.resource)

	makefileBytes, err := afero.ReadFile(fs.FS, filePath)
	if err != nil {
		return err
	}

	projectName := p.config.GetProjectName()
	if projectName == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %w", err)
		}
		projectName = strings.ToLower(filepath.Base(dir))
	}

	makefileBytes = append(makefileBytes, []byte(fmt.Sprintf(makefileBundleVarFragment, p.resource.Plural, p.resource.QualifiedGroup(), p.resource.Version, projectName))...)

	makefileBytes = append([]byte(fmt.Sprintf(makefileBundleImageFragement, p.config.GetDomain(), projectName)), makefileBytes...)

	var mode os.FileMode = 0644
	if info, err := fs.FS.Stat(filePath); err == nil {
		mode = info.Mode()
	}
	if err := afero.WriteFile(fs.FS, filePath, makefileBytes, mode); err != nil {
		return fmt.Errorf("error updating Makefile: %w", err)
	}

	scaffolder.InjectFS(fs)
	if err := scaffolder.Scaffold(); err != nil {
		return err
	}

	return nil
}

func (p *createAPISubcommand) InjectResource(res *resource.Resource) error {
	p.resource = res

	// RESOURCE: &{{cache zeusville.com v1 Joke} jokes  0xc00082a640 false 0xc00082a680}
	p.options.UpdateResource(p.resource)

	if err := p.resource.Validate(); err != nil {
		return err
	}

	// Check that resource doesn't have the API scaffolded
	if res, err := p.config.GetResource(p.resource.GVK); err == nil && res.HasAPI() {
		return errors.New("the API resource already exists")
	}

	// Check that the provided group can be added to the project
	if !p.config.IsMultiGroup() && p.config.ResourcesLength() != 0 && !p.config.HasGroup(p.resource.Group) {
		return fmt.Errorf("multiple groups are not allowed by default, to enable multi-group set 'multigroup: true' in your PROJECT file")
	}

	// Selected CRD version must match existing CRD versions.
	if pluginutil.HasDifferentCRDVersion(p.config, p.resource.API.CRDVersion) {
		return fmt.Errorf("only one CRD version can be used for all resources, cannot add %q", p.resource.API.CRDVersion)
	}

	return nil
}

const (
	makefileBundleVarFragment = `
##@Bundle
.PHONY: bundle
bundle:  ## Generate bundle manifests and metadata, then validate generated files.
	cat target/kubernetes/%[1]s.%[2]s-%[3]s.yml target/kubernetes/kubernetes.yml | operator-sdk generate bundle -q --overwrite --version 0.1.1 --default-channel=stable --channels=stable --package=%[4]s
	operator-sdk bundle validate ./bundle
	
.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .
	
.PHONY bundle-push
bundle-push: ## Push the bundle image.
	docker push $(BUNDLE_IMG)
`
)

const (
	makefileBundleImageFragement = `
IMAGE_TAG_BASE ?= %[1]s/%[2]s
BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)
`
)

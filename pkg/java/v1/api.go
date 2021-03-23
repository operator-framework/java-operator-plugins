/*
 * Copyright 2021 The Java Operator SDK Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

import (
	"fmt"

	"github.com/spf13/pflag"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/internal/kubebuilder/cmdutil"
	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds"
)

type createAPISubcommand struct {
	config config.Config
}

var (
	_ plugin.CreateAPISubcommand = &createAPISubcommand{}
)

func (p createAPISubcommand) UpdateContext(ctx *plugin.Context) {
	ctx.Description = `Scaffold a Kubernetes API by creating a Resource definition and / or a Controller.

create resource will prompt the user for if it should scaffold the Resource and / or Controller.  To only
scaffold a Controller for an existing Resource, select "n" for Resource.  To only define
the schema for a Resource without writing a Controller, select "n" for Controller.
`
}

func (p *createAPISubcommand) BindFlags(fs *pflag.FlagSet) {
}

func (p *createAPISubcommand) InjectConfig(c config.Config) {
	p.config = c
}

func (p *createAPISubcommand) Run(fs machinery.Filesystem) error {
	fmt.Println("create called")
	return nil
}

func (p *createAPISubcommand) Validate() error {
	return nil
}

func (p *createAPISubcommand) GetScaffolder() (cmdutil.Scaffolder, error) {
	return scaffolds.NewCreateAPIScaffolder(), nil
}

func (p *createAPISubcommand) PostScaffold() error {
	return nil
}

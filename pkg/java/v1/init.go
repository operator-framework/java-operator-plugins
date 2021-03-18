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

	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/internal/kubebuilder/cmdutil"
	"github.com/java-operator-sdk/kubebuilder-plugin/pkg/java/v1/scaffolds"
	"github.com/spf13/pflag"

	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
)

// This file represents the CLI for this plugin.

type initSubcommand struct {
	config config.Config

	// For help text.
	commandName string

	// Flags
	domain      string
	projectName string
}

var (
	_ plugin.InitSubcommand = &initSubcommand{}
	_ cmdutil.RunOptions    = &initSubcommand{}
)

func (p *initSubcommand) UpdateContext(ctx *plugin.Context) {
	ctx.Description = `Initialize a new project based on the java-operator-sdk project.

Writes the following files:
- a basic, Quarkus-based operator set-up
- a pom.xml file to build the project with Maven
`
	p.commandName = ctx.CommandName
}

func (p *initSubcommand) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&p.domain, "domain", "my.domain", "domain for groups")
	fs.StringVar(&p.projectName, "project-name", "", "name of this project")
	// TODO: include flags required for this plugin
}

func (p *initSubcommand) InjectConfig(c config.Config) {
	p.config = c
}

func (p *initSubcommand) Run(fs machinery.Filesystem) error {
	fmt.Println("init called")

	// TODO: any configuration that needs to be updated happens here
	if err := p.config.SetProjectName(p.projectName); err != nil {
		return err
	}

	if err := p.config.SetDomain(p.domain); err != nil {
		return err
	}

	// TODO: this will run the scaffolders
	if err := cmdutil.Run(p, fs); err != nil {
		return err
	}

	return nil
}

func (p *initSubcommand) Validate() error {
	// TODO: validate the conditions you expect before running the plugin
	return nil
}

func (p *initSubcommand) GetScaffolder() (cmdutil.Scaffolder, error) {
	return scaffolds.NewInitScaffolder(p.config), nil
}

func (p *initSubcommand) PostScaffold() error {
	// TODO: add anything you want to do AFTER the scaffolding has happened.
	return nil
}

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
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
)

type initSubcommand struct {
	config config.Config
	// For help text.
	commandName string
}

var (
	_ plugin.InitSubcommand = &initSubcommand{}
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
}

func (p *initSubcommand) InjectConfig(c config.Config) {
	p.config = c
}

func (p *initSubcommand) Run() error {
	fmt.Println("init called")
	return nil
}

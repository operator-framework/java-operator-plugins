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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	v3 "sigs.k8s.io/kubebuilder/v4/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/stage"
	"sigs.k8s.io/kubebuilder/v4/pkg/plugin"
)

var _ = Describe("v1", func() {
	testPlugin := &Plugin{}

	Describe("Name", func() {
		It("should return the plugin name", func() {
			Expect(testPlugin.Name(), "quarkus.javaoperatorsdk.io")
		})
	})

	Describe("Version", func() {
		It("should return the plugin version", func() {
			Expect(testPlugin.Version(), plugin.Version{Number: 1, Stage: stage.Beta})
		})
	})

	Describe("SupportedProjectVersions", func() {
		It("should return the support project versions", func() {
			Expect(testPlugin.Version(), []config.Version{v3.Version})
		})
	})

	Describe("GetInitSubcommand", func() {
		It("should return the plugin initSubcommand", func() {
			Expect(testPlugin.GetInitSubcommand(), &testPlugin.initSubcommand)
		})
	})

	Describe("GetCreateAPISubcommand", func() {
		It("should return the plugin createAPISubcommand", func() {
			Expect(testPlugin.GetCreateAPISubcommand(), &testPlugin.createAPISubcommand)
		})
	})
})

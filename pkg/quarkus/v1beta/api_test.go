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
	"github.com/spf13/pflag"
	"sigs.k8s.io/kubebuilder/v4/pkg/config"
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
)

var _ = Describe("v1", func() {
	var (
		testAPISubcommand createAPISubcommand
	)

	BeforeEach(func() {
		testAPISubcommand = createAPISubcommand{}
	})

	Describe("UpdateResource", func() {
		It("verify that resource fields were set", func() {
			testAPIOptions := &createAPIOptions{
				CRDVersion: "testVersion",
				Namespaced: true,
			}
			updateTestResource := resource.Resource{}
			testAPIOptions.UpdateResource(&updateTestResource)
			Expect(updateTestResource.API.CRDVersion).To(Equal(testAPIOptions.CRDVersion))
			Expect(updateTestResource.API.Namespaced).To(Equal(testAPIOptions.Namespaced))
			Expect(updateTestResource.Path).To(Equal(""))
			Expect(updateTestResource.Controller).To(BeFalse())
		})
	})

	Describe("BindFlags", func() {
		It("should set SortFlags to false", func() {
			flagTest := pflag.NewFlagSet("testFlag", -1)
			testAPISubcommand.BindFlags(flagTest)
			Expect(flagTest.SortFlags).To(BeFalse())
			Expect(testAPISubcommand.options.CRDVersion).To(Equal("v1"))
			Expect(testAPISubcommand.options.Namespaced).To(BeTrue())
		})
	})

	Describe("InjectConfig", func() {
		It("should set config", func() {
			testConfig, _ := config.New(config.Version{Number: 3})
			err := testAPISubcommand.InjectConfig(testConfig)
			Expect(testAPISubcommand.config).To(Equal(testConfig))
			Expect(err).To(BeNil())
		})
	})

	Describe("Run", func() {
		It("should return nil", func() {
			Expect(testAPISubcommand.Run(machinery.Filesystem{})).To(BeNil())
		})
	})

	Describe("Validate", func() {
		It("should return nil", func() {
			Expect(testAPISubcommand.Validate()).To(BeNil())
		})
	})

	Describe("PostScaffold", func() {
		It("should return nil", func() {
			Expect(testAPISubcommand.PostScaffold()).To(BeNil())
		})
	})

	Describe("InjectResource", func() {
		It("verify that wrong GVKs fail", func() {
			failAPISubcommand := &createAPISubcommand{}
			failResource := resource.Resource{
				GVK: resource.GVK{
					Group:   "Fail-Test-Group",
					Version: "test-version",
					Kind:    "test-kind",
				},
				Plural: "test-plural",
			}

			groupErr := failAPISubcommand.InjectResource(&failResource)
			Expect(failAPISubcommand.resource, failResource)
			Expect(groupErr).To(HaveOccurred())

			failResource.GVK.Group = "test-group"
			versionErr := failAPISubcommand.InjectResource(&failResource)
			Expect(versionErr).To(HaveOccurred())

			failResource.GVK.Version = "v1"
			kindError := failAPISubcommand.InjectResource(&failResource)
			Expect(kindError).To(HaveOccurred())
		})

		It("verify that a correct GVK succeeds", func() {
			testResource := resource.Resource{
				GVK: resource.GVK{
					Group:   "test-group",
					Version: "v1",
					Kind:    "Test-Kind",
				},
				Plural: "test-plural",
			}

			testConfig, _ := config.New(config.Version{Number: 3})
			testAPISubcommand.InjectConfig(testConfig)
			noErr := testAPISubcommand.InjectResource(&testResource)
			Expect(testAPISubcommand.resource, testResource)
			Expect(noErr).To(BeNil())
		})
	})
})

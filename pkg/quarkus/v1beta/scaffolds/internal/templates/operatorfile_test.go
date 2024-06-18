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

package templates

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"text/template"
)

var _ = Describe("operatorfile", func() {

	Describe("SetTemplateDefaults", func() {
		It("Should not stutter operator", func() {

			of := OperatorFile{
				OperatorName: "MemcachedQuarkusOperator",
			}

			err := of.SetTemplateDefaults()
			Expect(err).ToNot(HaveOccurred())
			Expect(of.Path).To(HaveSuffix("MemcachedQuarkusOperator.java"))

			tmpl, err := template.New("operatorfile").Parse(of.TemplateBody)
			Expect(err).ToNot(HaveOccurred())
			buf := new(bytes.Buffer)
			err = tmpl.Execute(buf, of)
			Expect(err).ToNot(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("public class MemcachedQuarkusOperator "))
		})
	})
})

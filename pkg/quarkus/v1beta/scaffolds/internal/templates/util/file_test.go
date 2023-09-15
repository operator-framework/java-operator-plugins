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

package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("file", func() {

	Describe("PrependJavaPath", func() {
		It("should prepend the configured java path to the given file", func() {
			Expect("src/main/java/com/example/MyOperator.java",
				PrependJavaPath("MyOperator.java", "com/example"))
		})
		It("prepend whatever you give it", func() {
			Expect("src/main/java/com.example/MyOperator.java",
				PrependJavaPath("MyOperator.java", "com.example"))
		})
	})

	Describe("PrependResourcePath", func() {
		It("should prepend the configured resource path to the given file", func() {
			Expect("src/main/resources/application.properties",
				PrependResourcePath("application.properties"))
		})
	})

	Describe("AsPath", func() {
		It("should convert a package name to a path", func() {
			Expect("com/example", AsPath("com.example"))
			Expect("/com/example", AsPath("/com.example"))
			Expect("org_unchanged", AsPath("org_unchanged"))
		})
	})
})

package templates

import (
	"bytes"

	. "github.com/onsi/ginkgo"
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

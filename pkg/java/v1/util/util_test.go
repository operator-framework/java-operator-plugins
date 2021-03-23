package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util", func() {

	Describe("ReverseDomain", func() {
		It("should reverse the string", func() {
			Expect("io.java", ReverseDomain("java.io"))
			Expect("com.example", ReverseDomain("example.com"))
			Expect("no/dots/leave/alone", ReverseDomain("no/dots/leave/alone"))
		})
		It("should leave the string alone", func() {
			Expect("no/dots/leave/alone", ReverseDomain("no/dots/leave/alone"))
			Expect("random string/No|dots;at,all",
				ReverseDomain("random string/No|dots;at,all"))
		})
	})

})

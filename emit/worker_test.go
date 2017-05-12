package emit

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	It("Should work", func() {

		var err error = nil
		switch err := err.(type){
		default:
			fmt.Println("coucou %v", err)
			Fail("osef")
		case nil:
			fmt.Println("nil")
			Expect(err).To(BeNil())
		}
	})
})

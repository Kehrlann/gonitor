package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"runtime"
)

var _ = Describe("exec : ", func() {
	Describe("ExecCommand : ", func() {
		resource := Resource{"http://www.example.com", 60, 2, 10, 3, ""}
		message := RecoveryMessage(resource, []int{200, 200, 200})

		It("Not run any command if not defined", func() {
			ret := ExecCommand("", message)

			Expect(ret).To(BeNil())
		})

		It("Should run the global command if defined", func() {
			// TODO : fix me :)
			// Note : don't run on windows because 'echo' is weird on that platform
			if runtime.GOOS != "windows" {
				ret := ExecCommand("echo", message)
				Expect(ret.String()).To(Equal("coucou"))
			}
		})
	})
})

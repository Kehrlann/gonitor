package main

import (
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("exec : ", func() {
	Describe("ExecCommand : ", func() {
		resourceWithoutCommand := Resource{"http://www.example.com", 60, 2, 10, 3, ""}

		message := RecoveryMessage(resourceWithoutCommand, []int{200, 200, 200})

		It("Not run any command if not defined", func() {
			ret := ExecCommand(message, "")
			Expect(ret).To(BeNil())
		})

		// Note : don't run on windows because 'echo' is weird on that platform
		if runtime.GOOS != "windows" {
			It("Should run the global command if defined", func() {
				ret := ExecCommand(message, "echo")
				Expect(ret.String()).To(ContainSubstring("RECOVERY"))
			})

			It("Should run the resource command if defined", func() {
				resourceWithCommand := Resource{"http://www.example.com", 60, 2, 10, 3, "echo"}
				messageWithCommand := RecoveryMessage(resourceWithCommand, []int{200, 200, 200})

				ret := ExecCommand(messageWithCommand, "")
				Expect(ret.String()).To(ContainSubstring("RECOVERY"))
			})

			It("Should override the global command if both are defined", func() {
				resourceWithCommand := Resource{"http://www.example.com", 60, 2, 10, 3, "go"}
				messageWithCommand := RecoveryMessage(resourceWithCommand, []int{200, 200, 200})

				ret := ExecCommand(messageWithCommand, "")
				Expect(ret.String()).To(ContainSubstring("unknown subcommand"))
			})
		}
	})
})

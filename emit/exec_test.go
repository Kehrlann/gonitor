package emit

import (
	"runtime"

	"github.com/kehrlann/gonitor/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("exec : ", func() {
	Describe("execCommand : ", func() {
		resourceWithoutCommand := config.Resource{"http://www.example.com", 60, 2, 10, 3, ""}

		message := RecoveryMessage(resourceWithoutCommand, []int{200, 200, 200})

		It("Not run any command if not defined", func() {
			ret := execCommand(message, "")
			Expect(ret).To(BeNil())
		})

		// Note : don't run on windows because 'echo' is weird on that platform
		if runtime.GOOS != "windows" {
			It("Should run the global command if defined", func() {
				ret := execCommand(message, "echo")
				Expect(ret.String()).To(ContainSubstring("RECOVERY"))
			})

			It("Should run the resource command if defined", func() {
				resourceWithCommand := config.Resource{"http://www.example.com", 60, 2, 10, 3, "echo"}
				messageWithCommand := RecoveryMessage(resourceWithCommand, []int{200, 200, 200})

				ret := execCommand(messageWithCommand, "")
				Expect(ret.String()).To(ContainSubstring("RECOVERY"))
			})

			It("Should override the global command if both are defined", func() {
				resourceWithCommand := config.Resource{"http://www.example.com", 60, 2, 10, 3, "go"}
				messageWithCommand := RecoveryMessage(resourceWithCommand, []int{200, 200, 200})

				ret := execCommand(messageWithCommand, "")
				Expect(ret.String()).To(ContainSubstring("unknown subcommand"))
			})
		}
	})
})

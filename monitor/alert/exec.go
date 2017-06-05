package alert

import (
	"os/exec"
	"bytes"
	"fmt"
	"time"
	"github.com/kehrlann/gonitor/monitor"
)

// commandEmitter emits message by running a command, and passing the fields of the StateChange message as the arguments
// for that command. It tries to run the command associated with the resource, and, if this is not set, it tries to run
// the default global command passed.
// Be careful ! This function waits for the command to return, so it might cause massive problems, leaks ...
type commandEmitter struct {
	defaultCommand string
}

func (c *commandEmitter) Emit(message *monitor.StateChangeMessage) {
	execCommand(message, c.defaultCommand)
}

// Exec command execute the command for the resource. It returns a byte buffer, with both StdErr and StdOut in it,
// mostly for testing purposes.
func execCommand(message *monitor.StateChangeMessage, defaultCommand string) *bytes.Buffer {
	var command string
	if message.Resource.Command != "" {
		command = message.Resource.Command
	} else if defaultCommand != "" {
		command = defaultCommand
	} else {
		return nil
	}

	var recoveryOrFailure string
	if message.IsOk {
		recoveryOrFailure = "RECOVERY"
	} else {
		recoveryOrFailure = "FAILURE"
	}

	cmd := exec.Command(command, recoveryOrFailure, message.Datetime.Format(time.RFC3339), fmt.Sprint(message.Codes), message.Resource.Url)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	return &out
}

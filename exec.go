package main

import (
	"os/exec"
	"bytes"
	"fmt"
	"time"
)


// Exec command execute the command for the resource. If it does not exist, ExecuteCommand tries to run the default
// global command passed. It returns a byte buffer, with both StdErr and StdOut in it, mostly for testing purposes.
// Be careful ! This function waits for the command to return, so it might cause massive problems, leaks ...
func ExecCommand(message *StateChangeMessage, defaultCommand string) *bytes.Buffer {
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

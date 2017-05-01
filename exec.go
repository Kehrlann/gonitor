package main

import (
	"os/exec"
	"bytes"
	"fmt"
)

func ExecCommand(globalCommand string, message *StateChangeMessage) *bytes.Buffer {
	var command string
	if message.Resource.Command != "" {
		command = message.Resource.Command
	} else if globalCommand != "" {
		command = globalCommand
	} else {
		return nil
	}

	var recoveryOrFailure string
	if message.IsOk {
		recoveryOrFailure = "RECOVERY"
	} else {
		recoveryOrFailure = "FAILURE"
	}

	cmd := exec.Command(command, recoveryOrFailure, message.Datetime.String(), fmt.Sprint(message.Codes), message.Resource.Url)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return &out
}

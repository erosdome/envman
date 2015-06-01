package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type EnvironmentKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CommandModel struct {
	Command          string                `json:"command"`
	WorkingDirectory string                `json:"working_directory"`
	LogFilePath      string                `json:"log_file_path"`
	Environments     []EnvironmentKeyValue `json:"environments"`
}

func RunCommandInDirWithArgsEnvsAndWriters(dirPath string, command string, cmdArgs []string, cmdEnvs []string) (int, error) {
	c := exec.Command(command, cmdArgs...)
	c.Env = append(os.Environ(), cmdEnvs...)
	// c.Env = cmdEnvs // only the supported envs, no inherited ones
	if dirPath != "" {
		c.Dir = dirPath
	}

	cmdExitCode := 0
	if err := c.Run(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		var waitStatus syscall.WaitStatus
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			cmdExitCode = waitStatus.ExitStatus()
		}
		return cmdExitCode, err
	}
	return 0, nil
}

func ExecuteCommand(cmdToRun CommandModel) (int, error) {
	fmt.Println("Command start")

	cmdExec := "/bin/bash"
	cmdArgs := []string{
		"--login",
		"-c",
		cmdToRun.Command,
	}
	cmdEnvs := []string{}
	envLength := len(cmdToRun.Environments)
	if envLength > 0 {
		cmdEnvs = make([]string, envLength, envLength)
		for idx, aEnvPair := range cmdToRun.Environments {
			cmdEnvs[idx] = aEnvPair.Key + "=" + aEnvPair.Value
		}
	}

	//
	cmdExitCode, commandErr := RunCommandInDirWithArgsEnvsAndWriters(cmdToRun.WorkingDirectory, cmdExec, cmdArgs, cmdEnvs)

	if commandErr != nil {
		fmt.Println("Command failed: %s", commandErr)
	}

	fmt.Println("Command finished: %s exit code: %d", commandErr, cmdExitCode)
	return cmdExitCode, commandErr
}
package main

import (	
	"fmt"
	"errors"
	"os"
	//"os/exec"
	"io/ioutil"
	"strings"

	"github.com/gkiki90/envman/pathutil"
	"github.com/gkiki90/envman/envutil"
	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
)

var stdinValue string
var configCommandEnvPrefix = "_CMDENV__"

func loadEnvlist() (envutil.EnvListYMLStruct, error) {
	path := pathutil.DefaultEnvlistPath
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		fmt.Println("Failed to check path, err: %s", err)
		return envutil.EnvListYMLStruct{}, err
	}
	if isExists {
		list, err := envutil.ReadEnvListFromFile(pathutil.DefaultEnvlistPath)
		if err != nil {
			fmt.Println("Failed to read envlist, err: %s", err)
			return envutil.EnvListYMLStruct{}, err
		}

		return list, nil
	} else {
		return envutil.EnvListYMLStruct{}, errors.New("No environemt variable list found")
	} 
}

func addCommand(c *cli.Context) {
	envKey := c.String("key")
	envValue := c.String("value")
	if stdinValue != "" {
		envValue = stdinValue
	}

	// Validate input
	if envKey == "" {
		fmt.Println("Invalid environment variable key")
		return
	}
	if envValue == "" {
		fmt.Println("Invalid environment variable value")
		return
	}

	// Load envlist, or create if not exist
	envlist, err := loadEnvlist()
	if err != nil {
		err := pathutil.CreateEnvmanDir()
		if err != nil {
			fmt.Println("Failed to create envman dir, err: %s", err)
			return
		}
	}

	// Add to or update envlist
	alreadyUsedKey := false
	newEnvStruct := envutil.EnvYMLStruct{ envKey, envValue }
	var newEnvList []envutil.EnvYMLStruct
	for i := range envlist.Envlist {
		oldEnvStruct := envlist.Envlist[i]
		if oldEnvStruct.Key ==  newEnvStruct.Key {
			alreadyUsedKey = true
			newEnvList = append(newEnvList, newEnvStruct)
		} else {
			newEnvList = append(newEnvList, oldEnvStruct)
		}
	}
	if alreadyUsedKey == false {
		newEnvList = append(newEnvList, newEnvStruct)
	}
	envlist.Envlist = newEnvList
	err = envutil.WriteEnvListToFile(pathutil.DefaultEnvlistPath, envlist)
	if err != nil {
		fmt.Println("Failed to create store envlist, err: %s", err)
		return
	}
	fmt.Println("New env list: ", newEnvList)

	return
}

func exportCommand(c *cli.Context) {
	envlist, err := loadEnvlist()
	if err != nil {
		fmt.Println("Failed to export environemt variable list, err: %s", err)
		return
	}
	if len(envlist.Envlist) == 0 {
		fmt.Println("Empty environemt variable list")
		return
	}

	for i := range envlist.Envlist {
		env := envlist.Envlist[i]
		os.Setenv(env.Key, env.Value)
		//fmt.Println(env.Key, os.Getenv(env.Key))
	}

	return
}

func runCommand(c *cli.Context) {
	exportCommand(c)

	/*
	doCmdEnvs := getCommandEnvironments()
	doCommand := c.String("run")
	flagCmdWorkDir := c.String("workdir")
	cmdToSend := CommandModel{
		Command:          doCommand,
		Environments:     doCmdEnvs,
		WorkingDirectory: flagCmdWorkDir,
	}

	cmdExitCode, err := ExecuteCommand(cmdToSend)

	fmt.Println(cmdToSend, cmdExitCode, err)
	*/

	executeCmd()

	return
}

func getCommandEnvironments() []EnvironmentKeyValue {
	cmdEnvs := []EnvironmentKeyValue{}

	for _, anEnv := range os.Environ() {
		splits := strings.Split(anEnv, "=")
		keyWithPrefix := splits[0]
		if strings.HasPrefix(keyWithPrefix, configCommandEnvPrefix) {
			cmdEnvItem := EnvironmentKeyValue{
				Key:   keyWithPrefix[len(configCommandEnvPrefix):],
				Value: os.Getenv(keyWithPrefix),
			}
			cmdEnvs = append(cmdEnvs, cmdEnvItem)
		}
	}

	//fmt.Println("cmdEnvs: %#v\n", cmdEnvs)

	return cmdEnvs
}

func main() {
	// Read piped data
	if ! terminal.IsTerminal(0) {
        bytes, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
        	fmt.Print("Failed to read stdin, err: %s", err)
        }
        stdinValue = string(bytes)
    } 

    // Parse cl 
	app := cli.NewApp()
	app.Name = "envman"
	app.Usage = "Environment varaibale manager."
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "run",
			Value: "",
		},
		cli.StringFlag {
			Name: "workdir",
			Value: "",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("run") != "" {
			runCommand(c)
		}
	}
	app.Commands = []cli.Command {
		{
			Name: "add",
			Flags: []cli.Flag {
				cli.StringFlag {
			    Name: "key",
			    Value: "",
			  },
			  cli.StringFlag {
			    Name: "value",
			    Value: "",
			  },
			},
			Action: addCommand,
		},
		{
			Name: "print",
			SkipFlagParsing: true,
			Action: exportCommand,
		},
		{
			Name: "env",
			SkipFlagParsing: true,
			Action: exportCommand,
		},
		{
			Name: "run",
			SkipFlagParsing: true,
			Action: runCommand,
		},
	}

	app.Run(os.Args)
}

package main

import (	
	"fmt"
	"errors"
	"os"
	"io/ioutil"

	"github.com/gkiki90/envman/pathutil"
	"github.com/gkiki90/envman/envutil"
	"github.com/alecthomas/kingpin"
	"code.google.com/p/go.crypto/ssh/terminal"
)

var (
	add = kingpin.Command("add", "Add new environment variable.")
	key = add.Flag("key", "Key for new/exist environment variable.").Required().String()
	value = add.Flag("value", "Value for new/exist environment variable.").String()

	print = kingpin.Command("print", "Load environment variables.")
)

func loadEnvlist() (envutil.EnvListJSONStruct, error) {
	path := pathutil.DefaultEnvlistPath
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		fmt.Println("Failed to check path, err: %s", err)
		return envutil.EnvListJSONStruct{}, err
	}
	if isExists {
		list, err := envutil.ReadEnvListFromFile(pathutil.DefaultEnvlistPath)
		if err != nil {
			fmt.Println("Failed to read envlist, err: %s", err)
			return envutil.EnvListJSONStruct{}, err
		}

		return list, nil
	} else {
		return envutil.EnvListJSONStruct{}, errors.New("No environemt variable list found")
	} 
}

func addEnv(envKey, envValue string) error {
	// Validate input
	if envKey == "" {
		return errors.New("Invalid environment variable key")
	}
	if envValue == "" {
		return errors.New("Invalid environment variable value")
	}

	// Load envlist, or create if not exist
	envlist, err := loadEnvlist()
	if err != nil {
		err := pathutil.CreateEnvmanDir()
		if err != nil {
			fmt.Println("Failed to create envman dir, err: %s", err)
			return err
		}
	}

	// Add to or update envlist
	alreadyUsedKey := false
	newEnvStruct := envutil.EnvJSONStruct{ *key, *value }
	var newEnvList []envutil.EnvJSONStruct
	for i := range envlist.Inputs {
		oldEnvStruct := envlist.Inputs[i]
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
	envlist.Inputs = newEnvList
	err = envutil.WriteEnvListToFile(pathutil.DefaultEnvlistPath, envlist)
	if err != nil {
		fmt.Println("Failed to create store envlist, err: %s", err)
		return err
	}

	fmt.Println("New env list: ", newEnvList)

	return nil
}

func printEnvlist() error {
	envlist, err := loadEnvlist()
	if err != nil {
		fmt.Println("Failed to read environment variable list, err: %s", err)
		return err
	}
	fmt.Println(envlist)
	return nil;
}


func main() {
	stdinValue := ""
	if ! terminal.IsTerminal(0) {
        bytes, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
        	fmt.Print("Failed to read stdin, err: %s", err)
        }
        stdinValue = string(bytes)
    } else {
        fmt.Println("no piped data")
    }

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Bitrise")
	kingpin.CommandLine.Help = "Environment variable manger."
	switch kingpin.Parse() {
		case add.FullCommand(): {
			if stdinValue != "" {
				*value = stdinValue
			}
			kingpin.FatalIfError(addEnv(*key, *value), "Add failed")
		}
		case print.FullCommand(): {
			kingpin.FatalIfError(printEnvlist(), "Print failed")
		}
	}
}

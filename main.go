package main

import (	
	"fmt"

	"github.com/gkiki90/envman/envutil"
	"github.com/gkiki90/envman/pathutil"
	"github.com/alecthomas/kingpin"
)

var (
	add     = kingpin.Command("add", "Add new environment variable.")
	key 	= add.Flag("key", "Key.").Required().String()
	value 	= add.Flag("value", "Value.").Required().String()
)

func versionComand() error {
	fmt.Println("version: ", VersionString)

	return nil
}

func helpComand() error {
	fmt.Println("help")

	return nil;
}

func addComand() error {
	fmt.Println("Add command")

	return nil;
}

func addEnv(envKey, envValue string) error {
	var envlist envutil.EnvListJSONStruct

	isExists, err := pathutil.IsPathExists(DefaultPath)
	if err != nil {
		fmt.Println("err!: %s", err)
		return err
	}
	if isExists {
		list, err := envutil.ReadEnvListFromFile(DefaultPath)
		envlist = list
		if err != nil {
			fmt.Println("err!: %s", err)
		}
	}
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
	envutil.WriteEnvListToFile(DefaultPath, envlist)

	fmt.Println("New env list: ", newEnvList)

	return nil
}

func removeCommand() error {
	fmt.Println("Remove command")
	return nil;
}

func doCommand() error {
	fmt.Println("Do command")
	return nil;
}

func printCommand() error {
	fmt.Println("Print command")
	return nil;
}


func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Alec Thomas")
	kingpin.CommandLine.Help = "Environment manger."
	switch kingpin.Parse() {
	case add.FullCommand():
		kingpin.FatalIfError(addEnv(*key, *value), "Add failed")
	}
}

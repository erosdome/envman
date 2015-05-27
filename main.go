package main

import (	
	"fmt"

	"github.com/gkiki90/envman/pathutil"
	"github.com/gkiki90/envman/envutil"
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

	// Load envlist, or create if not exist
	path := pathutil.DefaultEnvlistPath
	isExists, err := pathutil.IsPathExists(path)
	if err != nil {
		fmt.Println("Failed to check path, err!: %s", err)
		return err
	}
	if isExists {
		list, err := envutil.ReadEnvListFromFile(pathutil.DefaultEnvlistPath)
		envlist = list
		if err != nil {
			fmt.Println("Failed to read envlist, err!: %s", err)
			return err
		}
	} else {
		err := pathutil.CreateEnvmanDir()
		if err != nil {
			fmt.Println("Failed to create envman dir, err!: %s", err)
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
		fmt.Println("Failed to create store envlist, err!: %s", err)
		return err
	}

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

package pathutil

import (
    "os"
    "runtime"
)

var DefaultEnvmanDir string = UserHomeDir() + "/.envman/"
var DefaultEnvlistName string = "environment_variables.yml"
var DefaultEnvlistPath string = DefaultEnvmanDir + DefaultEnvlistName

func IsPathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func UserHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    }
    return os.Getenv("HOME")
}

func CreateEnvmanDir() error {
    path := DefaultEnvmanDir
    exist, _ := IsPathExists(path)
    if exist {
        return nil
    } 
    return createDir(path)
}

func createDir(path string) error {
    err := os.MkdirAll(path, os.ModePerm|os.ModeDir)
    return err
}
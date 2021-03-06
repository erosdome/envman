package envutil

import (
	"io/ioutil"
	"os"
	"errors"

	"gopkg.in/yaml.v2"
)

type EnvYMLStruct struct {
	Key string `yml:"key"`
	Value string `yml:"value"`
}

type EnvListYMLStruct struct {
	Envlist []EnvYMLStruct `yml:"environment_variables"`
}

func ReadEnvListFromFile(fpath string) (EnvListYMLStruct, error) {
	bytes, err := ioutil.ReadFile(fpath)
    if err != nil {
        return EnvListYMLStruct{}, err
    }

	var envlist EnvListYMLStruct
	err = yaml.Unmarshal(bytes, &envlist)
	if err != nil {
		return EnvListYMLStruct{}, err
	}

	return envlist, nil
}

func generateFormattedYMLForEnvList(envlist EnvListYMLStruct) ([]byte, error) {
	bytes, err := yaml.Marshal(envlist)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func WriteEnvListToFile(fpath string, envlist EnvListYMLStruct) error {
	if fpath == "" {
		return errors.New("No path provided")
	}

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonContBytes, err := generateFormattedYMLForEnvList(envlist)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonContBytes)
	if err != nil {
		return err
	}

	return nil
}



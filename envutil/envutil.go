package envutil

import (
	"io"
	"encoding/json"
	"os"
	"errors"

	"github.com/gkiki90/envman/pathutil"
)

type EnvJSONStruct struct {
	Key 			string 		`json:"key"`
	Value 			string		`json:"value"`
}

type EnvListJSONStruct struct {
	Inputs []EnvJSONStruct `json:"environment_variables"`
}

func ReadEnvListFromReader(reader io.Reader) (EnvListJSONStruct, error) {
	var envlist EnvListJSONStruct
	jsonParser := json.NewDecoder(reader)
	if err := jsonParser.Decode(&envlist); err != nil {
		return EnvListJSONStruct{}, err
	}

	return envlist, nil
}

func ReadEnvListFromFile(fpath string) (EnvListJSONStruct, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return EnvListJSONStruct{}, err
	}
	defer file.Close()

	return ReadEnvListFromReader(file)
}

func generateFormattedJSONForEnvList(envlist EnvListJSONStruct) ([]byte, error) {
	jsonContBytes, err := json.MarshalIndent(envlist, "", "\t")
	if err != nil {
		return []byte{}, err
	}
	return jsonContBytes, nil
}

func WriteEnvListToFile(fpath string, envlist EnvListJSONStruct) error {
	if fpath == "" {
		return errors.New("No path provided")
	}

	isExists, err := pathutil.IsPathExists(fpath)
	if err != nil {
		return err
	}
	if isExists {
		// return errors.New("Inputlist file already exists!")
	}

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonContBytes, err := generateFormattedJSONForEnvList(envlist)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonContBytes)
	if err != nil {
		return err
	}

	return nil
}
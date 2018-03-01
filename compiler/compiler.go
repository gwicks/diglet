package compiler

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gwicks/diglet/utils"
)

// CompileFile - Reads the JSON file at the given path and compiles it
func CompileFile(filePath string) (string, error) {
	var rootJSON interface{}

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	json.Unmarshal(fd, &rootJSON)

	resultJSON, rerr := utils.ParseFileRefs(filePath, rootJSON)

	if rerr != nil {
		return "", rerr
	}

	if resultJSONObj, ok := resultJSON.(map[string]interface{}); ok {
		finalJSON, perr := utils.ParseFileParent(filePath, resultJSONObj)
		if perr != nil {
			return "", perr
		}

		validatedJSON, verr := utils.ParseFileSchema(filePath, finalJSON)
		if verr != nil {
			return "", verr
		}

		resultStr, _ := json.MarshalIndent(validatedJSON, "", "    ")

		return string(resultStr), nil
	}
	return "", nil
}

package compiler

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gwicks/diglet/utils"
	"github.com/labstack/gommon/log"
)

// CompileFile - Reads the JSON file at the given path and compiles it
func CompileFile(filePath string) (string, error) {
	var rootJSON interface{}

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return "", err
	}

	json.Unmarshal(fd, &rootJSON)

	resultJSON, rerr := utils.ParseFileRefs(filePath, rootJSON)

	if rerr != nil {
		log.Error(rerr)
		return "", rerr
	}

	if resultJSONObj, ok := resultJSON.(map[string]interface{}); ok {
		finalJSON, perr := utils.ParseFileParent(filePath, resultJSONObj)
		if perr != nil {
			log.Error(perr)
			return "", perr
		}

		validatedJSON, verr := utils.ParseFileSchema(filePath, finalJSON)
		if verr != nil {
			log.Error(verr)
			return "", verr
		}

		resultStr, _ := json.MarshalIndent(validatedJSON, "", "    ")

		return string(resultStr), nil
	}
	return "", nil
}

package compiler

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gwicks/diglet/utils"
)

// BuildOptions - Compiler params
type BuildOptions struct {
	SkipResolve   bool
	SkipParenting bool
	SkipValidate  bool
}

// CompileFile - Reads the JSON file at the given path and compiles it
func CompileFile(filePath string, opts BuildOptions) (string, error) {
	var rootJSON map[string]interface{}

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	json.Unmarshal(fd, &rootJSON)

	var resolveRes interface{}
	var rerr, perr, verr error
	if !opts.SkipResolve {
		resolveRes, rerr = utils.ParseFileRefs(filePath, rootJSON)
		if rerr != nil {
			return "", rerr
		}
	} else {
		resolveRes = rootJSON
	}

	if resultJSONObj, ok := resolveRes.(map[string]interface{}); ok {
		if !opts.SkipParenting {
			resultJSONObj, perr = utils.ParseFileParent(filePath, resultJSONObj)
			if perr != nil {
				return "", perr
			}
		}

		if !opts.SkipValidate {
			resultJSONObj, verr = utils.ParseFileSchema(filePath, resultJSONObj)
			if verr != nil {
				return "", verr
			}
		}

		resultStr, _ := json.MarshalIndent(resultJSONObj, "", "    ")

		return string(resultStr), nil
	}
	return "", nil
}

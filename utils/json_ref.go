package utils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	jsonref "github.com/xeipuuv/gojsonreference"
)

var parentJSON map[string]interface{}

var currentPath string

func copyJSON(basePath string, filePath string, targetJSON *map[string]interface{}) string {
	vsj := filepath.Join(filepath.Dir(basePath), filePath)

	rfd, rerr := ioutil.ReadFile(vsj)
	if rerr != nil {
		return ""
	}

	json.Unmarshal(rfd, &targetJSON)
	delete(*targetJSON, "$ref")

	return vsj
}

func isRef(refJSON map[string]interface{}) (string, bool) {
	if len(refJSON) == 1 {
		if rk, rok := refJSON["$ref"].(string); rok {
			return rk, true
		}
	}
	return "", false
}

func resolveRefs(basePath string, rawJSON map[string]interface{}) error {
	for k, v := range rawJSON {
		switch v.(type) {
		case map[string]interface{}:
			if rv, ok := v.(map[string]interface{}); ok {
				refPath, ir := isRef(rv)
				if ir {
					if string(refPath[0]) == "#" {
						jref, _ := jsonref.NewJsonReference(refPath)
						refVal, _, _ := jref.GetPointer().Get(parentJSON)

						rawJSON[k] = refVal
					} else {
						currentPaths := copyJSON(basePath, refPath, &rv)
						resolveRefs(currentPaths, rv)
					}
				} else {
					resolveRefs(basePath, rv)
				}
			}
		case []interface{}:
			if rv, ok := v.([]interface{}); ok {
				for _, itm := range rv {
					if itmv, iok := itm.(map[string]interface{}); iok {
						resolveRefs(basePath, itmv)
					}
				}
			}
		case string:
			refPath, ir := isRef(rawJSON)
			if ir {
				currentPaths := copyJSON(basePath, refPath, &rawJSON)
				resolveRefs(currentPaths, rawJSON)
			}
		}
	}
	return nil
}

// ParseFileRefs Dummy
func ParseFileRefs(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	parentJSON = inJSON

	resolveRefs(filePath, parentJSON)

	return parentJSON, nil
}

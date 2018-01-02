package utils

import (
	"encoding/json"
	"io/ioutil"
)

var parentParseJSON map[string]interface{}
var outJSON map[string]interface{}

var currentPathParent string

func hasParent(jsonData map[string]interface{}) bool {
	if pa, ok := jsonData["@parent"].([]interface{}); ok {
		if len(pa) > 0 {
			return true
		}
	}
	return false
}

func resolveParents(basePath string, rawJSON map[string]interface{}) (map[string]interface{}, error) {
	var outData map[string]interface{}

	outData = make(map[string]interface{})
	for k, v := range rawJSON {
		if k == "@parent" {
			if vp, ok := v.([]interface{}); ok {
				for _, it := range vp {
					if itm, mok := it.(map[string]interface{}); mok {
						for ck, cv := range itm {
							if ck == "$ref" {
								if cvs, cok := cv.(string); cok {
									copyJSON(basePath, cvs, &outData)
								}
							}
						}
					}
				}
			}
		}
	}

	// delete(parentParseJSON, "@parent")
	// for k, v := range parentParseJSON {
	// 	outData[k] = v
	// }

	return outData, nil
}

// ParseFileParent Dummy
func ParseFileParent(filePath string) (map[string]interface{}, error) {
	parentParseJSON = make(map[string]interface{})

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fd, &parentParseJSON)

	if hasParent(parentParseJSON) {
		outJSON, _ = resolveParents(filePath, parentParseJSON)

		return outJSON, nil
	}

	return parentParseJSON, nil
}

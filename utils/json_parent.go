package utils

import "fmt"

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

	return outData, nil
}

// ParseFileParent Dummy
func ParseFileParent(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	parentParseJSON = inJSON

	if hasParent(parentParseJSON) {
		fmt.Println("HAS PARENT")
		fmt.Println(parentParseJSON)
		outJSON, _ = resolveParents(filePath, parentParseJSON)

		return outJSON, nil
	}

	return parentParseJSON, nil
}

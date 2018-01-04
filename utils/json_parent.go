package utils

import "fmt"

var parentParseJSON map[string]interface{}
var outJSON map[string]interface{}

var lastParent map[string]interface{}
var currentPathParent string

func hasParent(jsonData map[string]interface{}) bool {
	if pa, ok := jsonData["@parent"].([]interface{}); ok {
		if len(pa) > 0 {
			return true
		}
	} else {
		if _, sok := jsonData["@parent"].(map[string]interface{}); sok {
			return true
		}
	}
	return false
}

func resolveParents(basePath string, rawJSON map[string]interface{}) {
	fmt.Println("PARENT RESOLVE FOR")
	fmt.Println(rawJSON)
	// fmt.Println("WRITING TO OUTPUT OBJECT")
	// fmt.Println(outJSON)
	for k, v := range rawJSON {
		if k == "@parent" {
			if vp, ok := v.([]interface{}); ok {
				fmt.Println("FOUND ARRAY OF PARENTS")
				for _, it := range vp {
					if itm, mok := it.(map[string]interface{}); mok {
						if hasParent(itm) {
							lastParent = rawJSON
							resolveParents(basePath, itm)
						} else {
							fmt.Println("---------")
							fmt.Println(itm)
							fmt.Println("COPY INTO")
							fmt.Println(lastParent)
							delete(rawJSON, "@parent")
							if lastParent != nil {
								for vk, vv := range itm {
									lastParent[vk] = vv
								}
								delete(lastParent, "@parent")
							} else {
								for vk, vv := range itm {
									rawJSON[vk] = vv
								}
							}

						}
					}
				}
			} else {
				if sp, sok := v.(map[string]interface{}); sok {
					fmt.Println("FOUND SINGLE PARENT")
					if hasParent(sp) {
						lastParent = rawJSON
						resolveParents(basePath, sp)
					} else {
						delete(rawJSON, "@parent")
						if lastParent != nil {
							for vk, vv := range sp {
								lastParent[vk] = vv
							}
							delete(lastParent, "@parent")
						} else {
							for vk, vv := range sp {
								rawJSON[vk] = vv
							}
						}
					}
				}
			}
		} else {
			fmt.Println("NO PARENT")
			if targetObj, ok := v.(map[string]interface{}); ok {
				resolveParents(basePath, targetObj)
			}
		}
	}
}

// ParseFileParent Dummy
func ParseFileParent(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	outJSON = make(map[string]interface{})
	parentParseJSON = inJSON

	resolveParents(filePath, parentParseJSON)

	return parentParseJSON, nil
}

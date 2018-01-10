package utils

import (
	"sync"

	"github.com/imdario/mergo"
)

var parentParseJSON map[string]interface{}
var outJSON map[string]interface{}

var lastParent map[string]interface{}
var currentPathParent string

var parentMutex sync.Mutex

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

func checkIfLocked(k string, lockedNames []interface{}) bool {
	for _, i := range lockedNames {
		if is, ok := i.(string); ok {
			if k == is {
				return true
			}
		}
	}
	return false
}

func getLockedNames(rawJSON map[string]interface{}) []interface{} {
	if rawJSON["@lock_names"] != nil {
		if vl, lok := rawJSON["@lock_names"].([]interface{}); lok {
			return vl
		}
	}
	return nil
}

func resolveParents(basePath string, inputJSON interface{}, lastObject map[string]interface{}) {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		for k, v := range rawJSON {
			if k == "@parent" {
				if vp, ok := v.([]interface{}); ok {
					for _, it := range vp {
						if itm, mok := it.(map[string]interface{}); mok {
							if hasParent(itm) {
								resolveParents(basePath, itm, rawJSON)
							} else {
								lnames := getLockedNames(itm)
								delete(rawJSON, "@parent")
								for vk, vv := range itm {
									if rawJSON[vk] == nil {
										parentMutex.Lock()
										rawJSON[vk] = vv
										parentMutex.Unlock()
									} else {
										if checkIfLocked(vk, lnames) {
											parentMutex.Lock()
											rawJSON[vk] = vv
											parentMutex.Unlock()
										} else {
											if baseKeys, bkok := vv.(map[string]interface{}); bkok {
												if newKeys, nkok := rawJSON[vk].(map[string]interface{}); nkok {
													parentMutex.Lock()
													mergo.Merge(&newKeys, baseKeys)
													parentMutex.Unlock()
													resolveParents(basePath, newKeys, nil)
												}
											}
										}
									}
								}
								resolveParents(basePath, lastObject, nil)
							}
						}
					}
				} else {
					if sp, sok := v.(map[string]interface{}); sok {
						if hasParent(sp) {
							resolveParents(basePath, sp, rawJSON)
						} else {
							lnames := getLockedNames(sp)
							delete(rawJSON, "@parent")
							for vk, vv := range sp {
								if rawJSON[vk] == nil {
									rawJSON[vk] = vv
								} else {
									if checkIfLocked(vk, lnames) {
										rawJSON[vk] = vv
									}
								}
							}
							resolveParents(basePath, lastObject, nil)
						}
					}
				}
			} else {
				resolveParents(basePath, v, rawJSON)
			}
		}
	} else {
		if rawJSON, rok := inputJSON.([]interface{}); rok {
			for _, itm := range rawJSON {
				resolveParents(basePath, itm, nil)
			}
		}
	}
}

// ParseFileParent Dummy
func ParseFileParent(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	outJSON = make(map[string]interface{})
	parentParseJSON = inJSON

	resolveParents(filePath, parentParseJSON, nil)
	return parentParseJSON, nil
}

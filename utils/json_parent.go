package utils

import (
	"sync"

	"github.com/gwicks/mergo"
)

var parentParseJSON map[string]interface{}

var lastParent map[string]interface{}
var currentPathParent string

var parentMutex sync.Mutex

func hasParent(jsonData map[string]interface{}) (bool, interface{}) {
	if pa, ok := jsonData["@parent"].([]interface{}); ok {
		if len(pa) > 0 {
			return true, pa
		}
	} else {
		if pm, sok := jsonData["@parent"].(map[string]interface{}); sok {
			return true, pm
		}
	}
	return false, nil
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

func saveLockObj(names []interface{}, target map[string]interface{}) map[string]interface{} {
	var retObj map[string]interface{}

	retObj = make(map[string]interface{})

	for _, name := range names {
		if ns, lok := name.(string); lok {
			retObj[ns] = target[ns]
		}
	}

	return retObj
}

func mergeObjects(dest map[string]interface{}, src map[string]interface{}) {
	lnames := getLockedNames(src)
	restoreObject := saveLockObj(lnames, src)
	if len(lnames) > 0 {
		mergo.Merge(&dest, src)
		mergo.Merge(&dest, restoreObject, mergo.WithOverride)
	} else {
		mergo.Merge(&dest, src)
	}
}

func resolveParents(basePath string, inputJSON interface{}, lastObject map[string]interface{}) {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		rhp, pd := hasParent(rawJSON)
		if rhp {
			switch pd.(type) {
			case []interface{}:
				if pdarr, ok := pd.([]interface{}); ok {
					for i := len(pdarr)/2 - 1; i >= 0; i-- {
						opp := len(pdarr) - 1 - i
						pdarr[i], pdarr[opp] = pdarr[opp], pdarr[i]
					}
					for _, parentItem := range pdarr {
						delete(rawJSON, "@parent")
						if pdmap, pok := parentItem.(map[string]interface{}); pok {
							php, _ := hasParent(pdmap)
							if php {
								resolveParents(basePath, pdmap, nil)
								mergeObjects(rawJSON, pdmap)
							} else {
								mergeObjects(rawJSON, pdmap)
								resolveParents(basePath, pdmap, nil)
							}
						}
					}
				}
			case map[string]interface{}:
				if pdobj, ok := pd.(map[string]interface{}); ok {
					delete(rawJSON, "@parent")
					php, _ := hasParent(pdobj)
					if php {
						resolveParents(basePath, pdobj, nil)
						mergeObjects(rawJSON, pdobj)
					} else {
						mergeObjects(rawJSON, pdobj)
						resolveParents(basePath, pdobj, nil)
					}
				}
			}
			for _, v := range rawJSON {
				resolveParents(basePath, v, nil)
			}
		} else {
			for _, v := range rawJSON {
				resolveParents(basePath, v, nil)
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

// ParseFileParent Parses out parenting for a given file, with children overriding values in a parent, unless they are explicitly locked.
func ParseFileParent(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	parentParseJSON = inJSON

	resolveParents(filePath, parentParseJSON, nil)

	return parentParseJSON, nil
}

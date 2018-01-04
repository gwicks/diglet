package utils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	jsonref "github.com/xeipuuv/gojsonreference"
)

var parentJSON map[string]interface{}
var parentArr []interface{}

var parentIdx int
var tempParentHolder map[string]interface{}
var tempParentKey string

var currentPath string

func copyJSON(basePath string, filePath string, targetJSON *map[string]interface{}) (string, error) {
	vsj := filepath.Join(filepath.Dir(basePath), filePath)

	rfd, rerr := ioutil.ReadFile(vsj)
	if rerr != nil {
		return "", rerr
	}

	json.Unmarshal(rfd, &targetJSON)
	delete(*targetJSON, "$ref")

	return vsj, nil
}

func targetType(basePath string, filePath string) (interface{}, string, error) {
	var tempJSON interface{}
	vsj := filepath.Join(filepath.Dir(basePath), filePath)

	rfd, rerr := ioutil.ReadFile(vsj)
	if rerr != nil {
		return nil, vsj, rerr
	}

	json.Unmarshal(rfd, &tempJSON)

	return tempJSON, vsj, nil
}

func copyExtRef(basePath string, filePath string) (interface{}, error) {
	var tempJSON interface{}

	vsj := filepath.Join(filepath.Dir(basePath), filePath)
	pathSplit := strings.Split(vsj, "#/")

	rfd, rerr := ioutil.ReadFile(pathSplit[0])

	if rerr != nil {
		return nil, rerr
	}

	qk := "#/" + pathSplit[1]

	json.Unmarshal(rfd, &tempJSON)

	jref, _ := jsonref.NewJsonReference(qk)
	refVal, _, _ := jref.GetPointer().Get(tempJSON)

	return refVal, nil
}

func isRef(refJSON map[string]interface{}) (string, bool) {
	if len(refJSON) == 1 {
		if rk, rok := refJSON["$ref"].(string); rok {
			return rk, true
		}
	}
	return "", false
}

func queryFile(filePath string, query string) interface{} {
	var queryBase map[string]interface{}

	rfd, _ := ioutil.ReadFile(filePath)

	json.Unmarshal(rfd, &queryBase)

	jref, _ := jsonref.NewJsonReference(query)
	refVal, _, _ := jref.GetPointer().Get(queryBase)

	return refVal
}

func resolveRefs(basePath string, inputJSON interface{}) error {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		for k, v := range rawJSON {
			switch v.(type) {
			case map[string]interface{}:
				if rv, ok := v.(map[string]interface{}); ok {
					refPath, ir := isRef(rv)
					if ir {
						if string(refPath[0]) == "#" {
							rawJSON[k] = queryFile(basePath, refPath)

							resolveRefs(basePath, rawJSON[k])
						} else {
							pathSplit := strings.Split(filepath.Join(filepath.Dir(basePath), refPath), "#/")
							if len(pathSplit) == 2 {
								nv, _ := copyExtRef(basePath, refPath)
								rawJSON[k] = nv
							} else {
								tobjt, tp, _ := targetType(basePath, refPath)
								switch tobjt.(type) {
								case map[string]interface{}:
									currentPaths, copyErr := copyJSON(basePath, refPath, &rv)
									if copyErr != nil {
										log.Error(copyErr)
										return copyErr
									}
									resolveRefs(currentPaths, rv)
								case []interface{}:
									rawJSON[k] = tobjt
									resolveRefs(tp, tobjt)
								}

							}
						}
					} else {
						resolveRefs(basePath, rv)
					}
				}
			case []interface{}:
				if rv, ok := v.([]interface{}); ok {
					for pidx, itm := range rv {
						if itmv, iok := itm.(map[string]interface{}); iok {
							tempParentKey = k
							tempParentHolder = rawJSON
							parentIdx = pidx

							resolveRefs(basePath, itmv)
						} else {
							if atmv, aok := itm.([]interface{}); aok {
								for _, btm := range atmv {
									if btmv, bok := btm.(map[string]interface{}); bok {
										resolveRefs(basePath, btmv)
									}
								}
							}
						}
					}
				}
			case string:
				refPath, ir := isRef(rawJSON)
				if ir {
					tobjt, _, _ := targetType(basePath, refPath)
					switch tobjt.(type) {
					case []interface{}:
						if tobjarr, tok := tobjt.([]interface{}); tok {
							for idx, child := range tobjarr {
								if childMap, mok := child.(map[string]interface{}); mok {
									refPathc, isr := isRef(childMap)
									if isr {
										currentPath, copyErr := copyJSON(basePath, refPathc, &childMap)

										if copyErr != nil {
											log.Error(copyErr)
											return copyErr
										}
										tobjarr[idx] = childMap
										if pd, pok := tempParentHolder[tempParentKey].([]interface{}); pok {
											pd[parentIdx] = tobjarr
											resolveRefs(currentPath, childMap)
										}
									} else {
										if pd, pok := tempParentHolder[tempParentKey].([]interface{}); pok {
											pd[parentIdx] = tobjarr
										} else {
											tempParentHolder[tempParentKey] = tobjarr
										}
									}
								}
							}
						}

					case map[string]interface{}:
						currentPaths, copyErr := copyJSON(basePath, refPath, &rawJSON)
						if copyErr != nil {
							log.Error(copyErr)
							return copyErr
						}
						resolveRefs(currentPaths, rawJSON)
					}

				}
			}
		}
	} else {
		if rawJSON, rok := inputJSON.([]interface{}); rok {
			for cidx, child := range rawJSON {
				if childObj, bok := child.(map[string]interface{}); bok {
					refPath, ir := isRef(childObj)
					if ir {
						vd, vp, _ := targetType(basePath, refPath)
						rawJSON[cidx] = vd
						resolveRefs(vp, vd)
					}
				}
			}
		}
	}

	return nil
}

// ParseFileRefs Dummy
func ParseFileRefs(filePath string, inJSON interface{}) (interface{}, error) {

	switch inJSON.(type) {
	case map[string]interface{}:
		if inObj, ok := inJSON.(map[string]interface{}); ok {
			parentJSON = inObj

			resolveRefs(filePath, parentJSON)

			return parentJSON, nil
		}
	case []interface{}:
		if inObj, ok := inJSON.([]interface{}); ok {
			parentArr = inObj
			for _, child := range inObj {
				if childObj, bok := child.(map[string]interface{}); bok {
					parentJSON = childObj

					resolveRefs(filePath, parentJSON)

					return parentArr, nil
				}

			}

		}
	}

	return nil, nil
}

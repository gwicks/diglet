package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	jsonref "github.com/xeipuuv/gojsonreference"
)

var parentJSON map[string]interface{}
var parentArr []interface{}

var parentIdx int

var currentPath string

var refMutex sync.Mutex

func copyJSON(basePath string, filePath string, targetJSON *map[string]interface{}) (string, error) {
	vsj := filepath.Join(filepath.Dir(basePath), filePath)

	rfd, rerr := ioutil.ReadFile(vsj)
	if rerr != nil {
		return vsj, rerr
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

func logRefCall(bp string, ij interface{}) {
	fmt.Print("RESOLVE WITH PATH: ")
	fmt.Print(bp)
	fmt.Print(" AND OBJ: ")
	fmt.Print(ij)
	fmt.Print(" (")
	fmt.Print(reflect.TypeOf(ij))
	fmt.Println(")")
}

func logKV(k string, v interface{}) {
	fmt.Print("K: ")
	fmt.Print(k)
	fmt.Print(" V: ")
	fmt.Print(v)
	fmt.Println("")
}

func resolveRefs(basePath string, inputJSON interface{}, parentHolder map[string]interface{}, parentKey string) error {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		for k, v := range rawJSON {
			switch v.(type) {
			case map[string]interface{}:
				if rv, ok := v.(map[string]interface{}); ok {
					refPath, ir := isRef(rv)
					if ir {
						if string(refPath[0]) == "#" {
							refMutex.Lock()
							rawJSON[k] = queryFile(basePath, refPath)
							refMutex.Unlock()

							resolveRefs(basePath, rawJSON[k], nil, "")
						} else {
							pathSplit := strings.Split(filepath.Join(filepath.Dir(basePath), refPath), "#/")
							if len(pathSplit) == 2 {
								refMutex.Lock()
								nv, _ := copyExtRef(basePath, refPath)
								rawJSON[k] = nv
								refMutex.Unlock()
							} else {
								tobjt, tp, _ := targetType(basePath, refPath)
								switch tobjt.(type) {
								case map[string]interface{}:
									refMutex.Lock()
									currentPaths, copyErr := copyJSON(basePath, refPath, &rv)
									refMutex.Unlock()
									if copyErr != nil {
										log.Error(copyErr)
										return copyErr
									}
									resolveRefs(currentPaths, rv, nil, "")
								case []interface{}:
									rawJSON[k] = tobjt
									resolveRefs(tp, tobjt, nil, "")
								}

							}
						}
					} else {
						resolveRefs(basePath, rv, nil, "")
					}
				}
			case []interface{}:
				if rv, ok := v.([]interface{}); ok {

					resolveRefs(basePath, rv, rawJSON, k)
				}
			case string:
				refPath, ir := isRef(rawJSON)
				if ir {
					tobjt, rp, _ := targetType(basePath, refPath)
					switch tobjt.(type) {
					case []interface{}:
						if tobjarr, tok := tobjt.([]interface{}); tok {
							for idx, child := range tobjarr {
								if childMap, mok := child.(map[string]interface{}); mok {
									refPathc, isr := isRef(childMap)
									if isr {
										currentPath, copyErr := copyJSON(rp, refPathc, &childMap)
										if copyErr != nil {
											log.Error(copyErr)
											return copyErr
										}
										tobjarr[idx] = childMap
										if parentHolder != nil {
											if pd, pok := parentHolder[parentKey].([]interface{}); pok {
												pd[parentIdx] = tobjarr
												resolveRefs(currentPath, childMap, nil, "")
											}
										}

									} else {
										if parentHolder != nil {
											if pd, pok := parentHolder[parentKey].([]interface{}); pok {
												pd[parentIdx] = tobjarr
												resolveRefs(currentPath, childMap, nil, "")
											} else {
												parentHolder[parentKey] = tobjarr
											}
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
						resolveRefs(currentPaths, rawJSON, nil, "")
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
						resolveRefs(vp, vd, nil, "")
					} else {
						resolveRefs(basePath, childObj, nil, "")
					}
				} else {
					if childObj, bok := child.([]interface{}); bok {
						resolveRefs(basePath, childObj, nil, "")
					}
				}
			}
		}
	}
	return nil
}

// ParseFileRefs Dummy
func ParseFileRefs(filePath string, inJSON interface{}) (interface{}, error) {

	if inObj, ok := inJSON.(map[string]interface{}); ok {
		parentJSON = inObj

		resolveRefs(filePath, parentJSON, nil, "")

		return parentJSON, nil
	}

	return nil, nil
}

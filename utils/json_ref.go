package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/gwicks/mergo"
	jsonref "github.com/xeipuuv/gojsonreference"
)

var parentJSON map[string]interface{}
var parentArr []interface{}

var parentIdx int

var currentPath string

var refMutex sync.Mutex

func getRef(refJSON map[string]interface{}) (interface{}, bool) {
	if len(refJSON) == 1 {
		if _, rok := refJSON["$ref"].(string); rok {
			return refJSON, true
		}
	}
	return "", false
}

func refType(refPath string) int {
	if refPath[0] == '#' {
		return 2
	} else if strings.Index(refPath, "#/") != -1 {
		return 1
	} else {
		return 0
	}
}

func getResolve(refObj interface{}) bool {
	if refObjMap, refOk := refObj.(map[string]interface{}); refOk {
		if docDir, dok := refObjMap["@doc"].(map[string]interface{}); dok {
			if docResolve, rok := docDir["resolve"].(bool); rok {
				return docResolve
			}
		}
	}
	return true
}

func fetchRef(basePath string, refJSON map[string]interface{}, parentHolder interface{}, parentKey string, parentIdx int) string {
	if refStr, ok := refJSON["$ref"].(string); ok {
		var holderJSON interface{}

		refT := refType(refStr)

		switch refT {
		case 0:
			vsj := filepath.Join(filepath.Dir(basePath), refStr)

			rfd, _ := ioutil.ReadFile(vsj)

			json.Unmarshal(rfd, &holderJSON)

			if getResolve(holderJSON) {
				resolveRefs(vsj, holderJSON, nil, "", -1)
			}

			delete(refJSON, "$ref")

			if parentHolderObj, ok := parentHolder.(map[string]interface{}); ok {
				parentHolderObj[parentKey] = holderJSON
			} else {
				if parentHolderArr, ok := parentHolder.([]interface{}); ok {
					parentHolderArr[parentIdx] = holderJSON
				} else {
					mergo.Merge(&refJSON, holderJSON, mergo.WithOverride)
				}
			}

			return vsj
		case 1:
			pathSplit := strings.Split(filepath.Join(filepath.Dir(basePath), refStr), "#/")
			if len(pathSplit) == 2 {
				qk := "#/" + pathSplit[1]

				holderJSON = queryFile(pathSplit[0], qk)

				if getResolve(holderJSON) {
					resolveRefs(pathSplit[0], holderJSON, nil, "", -1)
				}

				if parentHolderObj, ok := parentHolder.(map[string]interface{}); ok {
					parentHolderObj[parentKey] = holderJSON
				} else {
					if parentHolderArr, ok := parentHolder.([]interface{}); ok {
						parentHolderArr[parentIdx] = holderJSON
					} else {
						mergo.Merge(&refJSON, holderJSON, mergo.WithOverride)
					}
				}
			}

		case 2:
			holderJSON = queryFile(basePath, refStr)

			if getResolve(holderJSON) {
				resolveRefs(basePath, holderJSON, nil, "", -1)
			}

			if parentHolderObj, ok := parentHolder.(map[string]interface{}); ok {
				parentHolderObj[parentKey] = holderJSON
			} else {
				if parentHolderArr, ok := parentHolder.([]interface{}); ok {
					parentHolderArr[parentIdx] = holderJSON
				} else {
					mergo.Merge(&refJSON, holderJSON, mergo.WithOverride)
				}
			}
		}
	}
	return basePath
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

func resolveRefs(basePath string, inputJSON interface{}, parentHolder interface{}, parentKey string, parentIdx int) error {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		refP, isr := getRef(rawJSON)
		if isr {
			if refPMap, pok := refP.(map[string]interface{}); pok {
				fetchRef(basePath, refPMap, parentHolder, parentKey, parentIdx)
			}
		} else {
			for k, v := range rawJSON {
				if getResolve(rawJSON) {
					resolveRefs(basePath, v, rawJSON, k, parentIdx)
				}

			}
		}
	} else {
		if rawJSON, arok := inputJSON.([]interface{}); arok {
			for pidx, obj := range rawJSON {
				resolveRefs(basePath, obj, rawJSON, "", pidx)
			}
		}
	}
	return nil
}

// ParseFileRefs Dummy
func ParseFileRefs(filePath string, inJSON interface{}) (interface{}, error) {

	if inObj, ok := inJSON.(map[string]interface{}); ok {
		parentJSON = inObj

		resolveRefs(filePath, parentJSON, nil, "", -1)

		return parentJSON, nil
	}

	return nil, nil
}

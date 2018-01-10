package utils

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema"
)

var rootJSON map[string]interface{}
var validErr error

var schemaMutex sync.Mutex

func schemaForObject(targetObj map[string]interface{}) (string, interface{}) {
	var retURI string
	var retSchema map[string]interface{}
	if schemaData, sok := targetObj["@schemas"].(map[string]interface{}); sok {
		for vk, vv := range schemaData {
			if childData, ok := vv.(map[string]interface{}); ok {
				retURI = vk
				retSchema = childData
			}
		}
		delete(targetObj, "@schemas")
		return retURI, retSchema
	}
	return "", nil
}

func doValidation(scmURI string, scmDat map[string]interface{}, targetObj interface{}) error {
	url := scmURI + ".json"
	schemaMutex.Lock()
	marshaledSchema, _ := json.Marshal(scmDat)
	marshaledObj, _ := json.Marshal(targetObj)
	schemaMutex.Unlock()

	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft6

	if err := compiler.AddResource(url, strings.NewReader(string(marshaledSchema))); err != nil {
		return err
	}

	schema, serr := compiler.Compile(url)
	if serr != nil {
		return serr
	}

	if verr := schema.Validate(strings.NewReader(string(marshaledObj))); verr != nil {
		return verr
	}
	return nil
}

func validateSchema(basePath string, inputJSON interface{}) error {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		for k, v := range rawJSON {
			if k == "@schemas" {
				schmURI, currentSchema := schemaForObject(rawJSON)
				if currentSchemaMap, sok := currentSchema.(map[string]interface{}); sok {
					tmperr := doValidation(schmURI, currentSchemaMap, rawJSON)
					if tmperr != nil {
						return tmperr
					}
				}
			} else {
				switch v.(type) {
				case map[string]interface{}:
					if objData, ok := v.(map[string]interface{}); ok {
						schmURI, currentSchema := schemaForObject(objData)
						if currentSchemaMap, sok := currentSchema.(map[string]interface{}); sok {
							tmperr := doValidation(schmURI, currentSchemaMap, v)
							if tmperr != nil {
								return tmperr
							}
						}
						validateSchema(basePath, objData)
					}
				case []interface{}:
					if rv, ok := v.([]interface{}); ok {
						validateSchema(basePath, rv)
					}
				}
			}
		}
	} else {
		if rawJSON, rok := inputJSON.([]interface{}); rok {
			for _, child := range rawJSON {
				validateSchema(basePath, child)
			}
		}
	}
	return nil
}

// ParseFileSchema Dummy
func ParseFileSchema(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	validErr = nil
	rootJSON = inJSON

	validErr := validateSchema(filePath, rootJSON)
	if validErr != nil {
		return rootJSON, validErr
	}

	return rootJSON, nil
}

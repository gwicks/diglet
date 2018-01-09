package utils

import (
	"encoding/json"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

var rootJSON map[string]interface{}
var validErr error

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
	marshaledSchema, _ := json.Marshal(scmDat)

	marshaledObj, _ := json.Marshal(targetObj)

	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft6

	if err := compiler.AddResource(url, strings.NewReader(string(marshaledSchema))); err != nil {
		return err
	}

	schema, err := compiler.Compile(url)
	if err != nil {
		return err
	}

	if err = schema.Validate(strings.NewReader(string(marshaledObj))); err != nil {
		return err
	}
	return nil
}

func validateSchema(basePath string, inputJSON interface{}) {
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		for k, v := range rawJSON {
			if k == "@schemas" {
				schmURI, currentSchema := schemaForObject(rawJSON)
				if currentSchemaMap, sok := currentSchema.(map[string]interface{}); sok {
					validErr = doValidation(schmURI, currentSchemaMap, rawJSON)
				}
			} else {
				switch v.(type) {
				case map[string]interface{}:
					if objData, ok := v.(map[string]interface{}); ok {
						schmURI, currentSchema := schemaForObject(objData)
						if currentSchemaMap, sok := currentSchema.(map[string]interface{}); sok {
							validErr = doValidation(schmURI, currentSchemaMap, v)
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
}

// ParseFileSchema Dummy
func ParseFileSchema(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	rootJSON = inJSON

	validateSchema(filePath, rootJSON)
	if validErr != nil {
		return rootJSON, validErr
	}

	return rootJSON, nil
}

package utils

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema"
)

var rootJSON map[string]interface{}

var schemaMutex sync.Mutex

func hasSchema(jsonData map[string]interface{}) (bool, interface{}) {
	if pm, sok := jsonData["@schemas"].(map[string]interface{}); sok {
		if len(pm) > 0 {
			return true, pm
		}
	}
	return false, nil
}

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

	return schema.Validate(strings.NewReader(string(marshaledObj)))
}

func validateSchema(basePath string, inputJSON interface{}) error {
	var subErr error
	if rawJSON, rok := inputJSON.(map[string]interface{}); rok {
		rhs, _ := hasSchema(rawJSON)
		if rhs {
			schmURI, currentSchema := schemaForObject(rawJSON)
			if currentSchemaMap, sok := currentSchema.(map[string]interface{}); sok {
				validationErr := doValidation(schmURI, currentSchemaMap, rawJSON)
				if validationErr != nil {
					return validationErr
				}
				subErr = validateSchema(basePath, rawJSON)
				if subErr != nil {
					return subErr
				}
			}
		} else {
			for _, v := range rawJSON {
				subErr = validateSchema(basePath, v)
				if subErr != nil {
					return subErr
				}
			}
		}
	} else {
		if rawJSON, rok := inputJSON.([]interface{}); rok {
			for _, child := range rawJSON {
				subErr = validateSchema(basePath, child)
				if subErr != nil {
					return subErr
				}
			}
		}
	}
	return nil
}

// ParseFileSchema Handles schema validation to the JSON Schema Draft 6 specification, returns any validations errors.
func ParseFileSchema(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	rootJSON = inJSON

	validErr := validateSchema(filePath, rootJSON)
	if validErr != nil {
		return rootJSON, validErr
	}

	return rootJSON, nil
}

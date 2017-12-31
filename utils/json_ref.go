package utils

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

var parentJSON map[string]interface{}

var currentPath string

func copyJSON(basePath string, filePath string, targetJSON *map[string]interface{}) string {
	vsj := filepath.Join(filepath.Dir(basePath), filePath)

	rfd, rerr := ioutil.ReadFile(vsj)
	if rerr != nil {
		return ""
	}

	json.Unmarshal(rfd, &targetJSON)
	delete(*targetJSON, "$ref")

	return vsj
}

func resolveRefs(basePath string, rawJSON map[string]interface{}) error {
	for _, v := range rawJSON {
		switch v.(type) {
		case map[string]interface{}:
			if rv, ok := v.(map[string]interface{}); ok {
				for ak, av := range rv {
					if ak == "$ref" {
						if avk, aok := av.(string); aok {
							currentPaths := copyJSON(basePath, avk, &rv)
							resolveRefs(currentPaths, rv)
						}
					}
				}
			}
		}
	}
	return nil
}

// ParseFileRefs Dummy
func ParseFileRefs(filePath string) (map[string]interface{}, error) {
	parentJSON = make(map[string]interface{})

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fd, &parentJSON)

	resolveRefs(filePath, parentJSON)

	return parentJSON, nil
}

package utils

import (
	"encoding/json"
	"io/ioutil"
)

// ParseFileParent Dummy
func ParseFileParent(filePath string) (map[string]interface{}, error) {
	parentJSON = make(map[string]interface{})

	fd, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fd, &parentJSON)

	resolveRefs(filePath, parentJSON)

	return parentJSON, nil
}

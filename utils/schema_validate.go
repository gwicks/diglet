package utils

var rootJSON map[string]interface{}

func validateSchema(basePath string, rawJSON map[string]interface{}) error {
	// fmt.Println(basePath)

	return nil
}

// ParseFileSchema Dummy
func ParseFileSchema(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	rootJSON = inJSON

	validateSchema(filePath, rootJSON)

	return rootJSON, nil
}

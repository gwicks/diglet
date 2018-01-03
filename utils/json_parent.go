package utils

var parentParseJSON map[string]interface{}
var outJSON map[string]interface{}

var currentPathParent string

func hasParent(jsonData map[string]interface{}) bool {
	if pa, ok := jsonData["@parent"].([]interface{}); ok {
		if len(pa) > 0 {
			return true
		}
	} else {
		if _, sok := jsonData["@parent"].(map[string]interface{}); sok {
			return true
		}
	}
	return false
}

func resolveParents(basePath string, rawJSON map[string]interface{}) {
	for k, v := range rawJSON {
		if k == "@parent" {
			if vp, ok := v.([]interface{}); ok {
				for _, it := range vp {
					if itm, mok := it.(map[string]interface{}); mok {
						if hasParent(itm) {
							resolveParents(basePath, itm)
						} else {
							for vk, vv := range itm {
								outJSON[vk] = vv
							}
						}
					}
				}
			} else {
				if sp, sok := v.(map[string]interface{}); sok {
					if hasParent(sp) {
						resolveParents(basePath, sp)
					} else {
						for vk, vv := range sp {
							outJSON[vk] = vv
						}
					}
				}
			}
		} else {
			outJSON[k] = v
		}
	}
}

// ParseFileParent Dummy
func ParseFileParent(filePath string, inJSON map[string]interface{}) (map[string]interface{}, error) {
	outJSON = make(map[string]interface{})
	parentParseJSON = inJSON

	resolveParents(filePath, parentParseJSON)

	return outJSON, nil
}

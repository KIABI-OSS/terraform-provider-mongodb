package provider

import "errors"

// Convert an index type declared in terraform as string into a type and value expected by Mongo's client.
func convertToMongoIndexType(indexType string) interface{} {
	switch indexType {
	case "asc":
		return 1
	case "desc":
		return -1
	default:
		return indexType
	}
}

// Convert an index type returned by Mongo's client into a string understood by terraform.
func convertToTfIndexType(indexType interface{}) (string, error) {
	intValue, isInt := indexType.(int32)
	if isInt {
		switch intValue {
		case 1:
			return "asc", nil
		case -1:
			return "desc", nil
		default:
			return "", errors.New("if typeIndex is int, it MUST have value 1 ou -1")
		}
	}

	strValue, isStr := indexType.(string)
	if isStr {
		return strValue, nil
	}

	return "", errors.New("typeIndex MUST be int32 or string")
}

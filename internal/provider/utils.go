package provider

import (
	"errors"
	"strings"
)

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

type indexId struct {
	database   string
	collection string
	indexName  string
}

func parseIndexId(path string) (*indexId, error) {
	splitPath := strings.Split(path, ".")
	if len(splitPath) != 3 {
		return nil, errors.New("Index id's format must be <database>.<collection>.<index_name>")
	}

	return &indexId{database: splitPath[0], collection: splitPath[1], indexName: splitPath[2]}, nil
}

package maputil

import (
	"fmt"

	"github.com/alex-slynko/manifest_splitter/types"
)

func ExtractOperations(first, second map[string]interface{}) ([]types.Operation, error) {
	var operations []types.Operation

	for k, v := range first {
		value, ok := second[k]
		if !ok {
			op := types.Operation{
				Path:  "/" + k + "?",
				Type:  "replace",
				Value: v,
			}
			operations = append(operations, op)
		} else {
			newOps, err := generateOperationForValues("/"+k, v, value)
			if err != nil {
				return operations, err
			}
			operations = append(operations, newOps...)
		}
	}

	for k, _ := range second {
		_, ok := first[k]
		if !ok {
			op := types.Operation{
				Path: "/" + k,
				Type: "remove",
			}
			operations = append(operations, op)
		}
	}

	return operations, nil
}

func generateOperationForValues(key string, newValue, oldValue interface{}) ([]types.Operation, error) {
	var operations []types.Operation

	if oldMap, ok := oldValue.(map[interface{}]interface{}); ok {
		if newMap, ok := newValue.(map[interface{}]interface{}); ok {
			for k, v := range newMap {
				value, ok := oldMap[k]
				if !ok {
					op := types.Operation{
						Path:  key + "/" + k.(string) + "?",
						Type:  "replace",
						Value: v,
					}
					operations = append(operations, op)
				} else {
					newOps, _ := generateOperationForValues(key+"/"+k.(string), v, value)
					operations = append(operations, newOps...)
				}
			}
			for k := range oldMap {
				_, ok := newMap[k]
				if !ok {
					op := types.Operation{
						Path: key + "/" + k.(string),
						Type: "remove",
					}
					operations = append(operations, op)
				}
			}
		} else {
			return operations, fmt.Errorf("Can not replace %s. New value is not a map, but old value is (%#v and %#v)", key, newValue, oldValue)
		}
	} else if oldSlice, ok := oldValue.([]interface{}); ok {
		if newSlice, ok := newValue.([]interface{}); ok {
			for _, v := range newSlice {
				if !contains(v, oldSlice) {
					op := types.Operation{
						Path:  key + "/-",
						Type:  "replace",
						Value: v,
					}
					operations = append(operations, op)

				}
			}
		}
	} else if newValue != oldValue {
		if _, ok := newValue.(map[interface{}]interface{}); ok {
			return operations, fmt.Errorf("Can not replace %s. New value is a map, but old value is not (%#v and %#v)", key, newValue, oldValue)
		}
		op := types.Operation{
			Path:  key,
			Type:  "replace",
			Value: newValue,
		}
		operations = append(operations, op)
	}
	return operations, nil
}

func contains(value interface{}, slice []interface{}) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

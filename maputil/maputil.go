package maputil

import (
	"fmt"

	"github.com/alex-slynko/manifest_splitter/types"
)

// ExtractOperations ...
func ExtractOperations(newValue, oldValue map[string]interface{}) ([]types.Operation, error) {
	var operations []types.Operation

	for k, v := range newValue {
		value, ok := oldValue[k]
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

	for k := range oldValue {
		_, ok := newValue[k]
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

func generateOperationForMaps(key string, oldMap map[interface{}]interface{}, newValue interface{}) (operations []types.Operation, err error) {

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
		return operations, fmt.Errorf("Can not replace %s. New value is not a map, but old value is (%#v and %#v)", key, newValue, oldMap)
	}
	return
}

func generateOperationForSlices(key string, oldSlice []interface{}, newValue interface{}) (operations []types.Operation, err error) {
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

		for i, v := range oldSlice {
			if !contains(v, newSlice) {
				op := types.Operation{
					Path:  fmt.Sprintf("%s/%d", key, i),
					Type:  "remove",
					Value: nil,
				}
				operations = append(operations, op)
			}
		}
	} else {
		return operations, fmt.Errorf("Can not replace %s. New value is not a slice, but old value is (%#v)", key, oldSlice)
	}
	return
}

func generateOperationForValues(key string, newValue, oldValue interface{}) ([]types.Operation, error) {
	var operations []types.Operation

	if oldMap, ok := oldValue.(map[interface{}]interface{}); ok {
		ops, err := generateOperationForMaps(key, oldMap, newValue)
		if err != nil {
			return operations, err
		}
		operations = append(operations, ops...)
	} else if oldSlice, ok := oldValue.([]interface{}); ok {
		ops, err := generateOperationForSlices(key, oldSlice, newValue)
		if err != nil {
			return operations, err
		}
		operations = append(operations, ops...)
	} else if newValue != oldValue {
		if _, ok := newValue.(map[interface{}]interface{}); ok {
			return operations, fmt.Errorf("Can not replace %s. New value is a map, but old value is not (%#v and %#v)", key, newValue, oldValue)
		}
		if _, ok := newValue.([]interface{}); ok {
			return operations, fmt.Errorf("Can not replace %s. New value is a slice, but old value is not (%#v and %#v)", key, newValue, oldValue)
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

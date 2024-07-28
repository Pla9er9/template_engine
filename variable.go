package templateEngine

import (
	"errors"
	"reflect"
	"strings"
)

func getVariable(variableName string, variables *map[string]any) (any, error) {
	properties := strings.Split(variableName, ".")
	if len(properties) > 1 {
		obj := (*variables)[properties[0]]
		return getPropertyFromObject(obj, properties[1:])
	} else {
		return (*variables)[variableName], nil
	}
}

func getPropertyFromObject(struct_ any, fields []string) (any, error) {
	if struct_ == nil {
		return nil, errors.New("nil passed")
	}

	v, err := getField(&struct_, fields[0])
	if err != nil {
		return nil, err
	}

	if len(fields) == 1 {
		return v, nil
	}

	return getPropertyFromObject(v, fields[1:])
}

func getField(v *any, field string) (any, error) {
	r := reflect.ValueOf(*v)
	f := reflect.Indirect(r).FieldByName(field)

	if f.Kind().String() == "invalid" || !f.CanInterface() {
		return nil, errors.New("field doesnt exist")
	}

	return f.Interface(), nil
}

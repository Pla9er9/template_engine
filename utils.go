package templateEngine

import (
	"errors"
	"reflect"
)

func convertAnyToSlice(data any) ([]any, error) {
	value := reflect.ValueOf(data)
	kind := value.Kind()

	if kind == reflect.Slice {
		sliceValue := value
		newSlice := make([]any, 0, sliceValue.Len())

		for i := 0; i < sliceValue.Len(); i++ {
			elementValue := sliceValue.Index(i)
			newSlice = append(newSlice, elementValue.Interface())
		}

		return newSlice, nil
	} else {
		return nil, errors.New("passed data is not slice")
	}
}

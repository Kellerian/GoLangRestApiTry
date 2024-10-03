package main

import (
	"errors"
	"fmt"
	"reflect"
)

func ValidateStruct(s interface{}) (err error) {
	structType := reflect.TypeOf(s)
	if structType.Kind() != reflect.Struct {
		return errors.New("input param should be a struct")
	}

	structVal := reflect.ValueOf(s)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name
		var isSet bool
		if field.Kind() == reflect.Pointer {
			if !field.IsNil() {
				isSet = field.IsValid() && !field.IsZero()
			} else {
				isSet = true
			}
		} else {
			isSet = field.IsValid() && !field.IsZero()
		}
		if !isSet {
			err = fmt.Errorf("%v%s in not set; ", err, fieldName)
		}
	}

	return err
}

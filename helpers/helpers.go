package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	fmt.Println(typ)

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsNil() {
			continue
		}
		fieldName := typ.Field(i).Tag.Get("db")
		fieldValueKind := val.Field(i).Kind()
		var fieldValue interface{}
		if fieldValueKind == reflect.Struct {
			fieldValue = StructToMap(val.Field(i).Interface())
		} else {
			fieldValue = val.Field(i).Interface()
		}
		result[fieldName] = fieldValue
	}

	return result
}

func GetPageSize(r *http.Request) (int, error) {
	strPageSize := r.URL.Query().Get("size")
	pageSize := 100
	if strPageSize != "" {
		pageSize, err := strconv.Atoi(strPageSize)
		if err != nil || pageSize < 0 {
			return 0, errors.New("incorrect page size")
		}
	}
	return pageSize, nil
}

func GetPageNumber(r *http.Request) (int, error) {
	strPageNumber := r.URL.Query().Get("number")
	pageNumber := 0
	if strPageNumber != "" {
		pageNumber, err := strconv.Atoi(strPageNumber)
		if err != nil || pageNumber < 0 {
			return 0, errors.New("incorrect page number")
		}
	}
	return pageNumber, nil
}

func CalcOffset(s, n int) uint {
	return uint(s * n)
}

package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
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
	log.Println("Размер страницы в запросе", strPageSize)
	if strPageSize != "" {
		pageSize, err := strconv.Atoi(strPageSize)
		log.Println("Размер страницы в запросе", strPageSize)
		if err != nil || pageSize < 0 {
			return 0, errors.New("incorrect page size")
		}
		return pageSize, nil
	}
	log.Println("Итоговый размер страницы в запросе", pageSize)
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

func CalcOffset(page_size, page_number int) uint {
	return uint(page_size * page_number)
}

func BuildWhereClause(q *CodeFilterParams) ([]goqu.Expression, error) {
	var conditions []goqu.Expression
	v := reflect.ValueOf(q).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		fieldValue := v.Field(i)

		// Обработка тегов goqu
		goquTags := parseGoquTags(field.Tag.Get("goqu"))
		if shouldOmit(fieldValue, goquTags) {
			continue
		}

		// Для указателей: получаем значение
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		// Обработка разных типов данных
		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			if fieldValue.Len() == 0 {
				continue
			}
			values := make([]interface{}, fieldValue.Len())
			for j := 0; j < fieldValue.Len(); j++ {
				values[j] = fieldValue.Index(j).Interface()
			}
			conditions = append(conditions, goqu.I(dbTag).In(values))

		default:
			conditions = append(conditions, goqu.I(dbTag).Eq(fieldValue.Interface()))
		}
	}

	return conditions, nil
}

func parseGoquTags(tag string) map[string]bool {
	tags := strings.Split(tag, ",")
	result := make(map[string]bool)
	for _, t := range tags {
		result[t] = true
	}
	return result
}

func shouldOmit(v reflect.Value, tags map[string]bool) bool {
	// Для указателей проверяем omitnil
	if v.Kind() == reflect.Ptr {
		if tags["omitnil"] && v.IsNil() {
			return true
		}
		if !v.IsNil() {
			v = v.Elem()
		}
	}

	// Проверяем omitempty для значений
	if tags["omitempty"] {
		switch v.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return v.Len() == 0
		case reflect.String:
			return v.String() == ""
		case reflect.Ptr:
			return v.IsNil()
		}
	}
	return false
}

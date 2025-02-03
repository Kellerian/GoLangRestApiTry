package helpers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// Преобразуем CustomTime в time.Time
	t := time.Time(ct)

	// Форматируем время в строку
	formatted := t.Format("2006-01-02 15:04:05.999999")

	// Возвращаем строку в кавычках (как JSON-строку)
	return json.Marshal(formatted)
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// Убираем кавычки из JSON-строки
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return fmt.Errorf("invalid time format: %s", str)
	}
	str = str[1 : len(str)-1]

	// Парсим строку в time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", str)
	if err != nil {
		return err
	}

	// Сохраняем результат в CustomTime
	*ct = CustomTime(parsedTime)
	return nil
}

func (ct CustomTime) ToTime() time.Time {
	return time.Time(ct)
}

func (ct *CustomTime) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", string(text))
	if err != nil {
		return err
	}

	// Сохраняем результат в CustomTime
	*ct = CustomTime(parsedTime)
	return nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime(v)
	case []byte:
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", string(v))
		if err != nil {
			return err
		}
		*ct = CustomTime(parsedTime)
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", v)
		if err != nil {
			return err
		}
		*ct = CustomTime(parsedTime)
	default:
		return fmt.Errorf("unsupported type for CustomTime: %T", value)
	}

	return nil
}

func (ct CustomTime) Value() (driver.Value, error) {
	return time.Time(ct), nil
}

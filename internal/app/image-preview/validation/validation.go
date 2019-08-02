package validation

import (
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v9"
)

var _validate *validator.Validate

func init() {
	_validate = validator.New()
}

func Validate(i interface{}) error {
	itemType := reflect.TypeOf(i)
	switch itemType.Kind() {
	case reflect.Struct:
		return _validate.Struct(i)
	case reflect.Slice:
		itemValue := reflect.ValueOf(i)
		for i := 0; i < itemValue.Len(); i++ {
			if err := _validate.Struct(itemValue.Index(i)); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported item type: %s", itemType.Name())
	}
}

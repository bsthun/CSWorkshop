package text

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func Validate(schema any, tag ...string) error {
	s := reflect.ValueOf(schema)

	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	if s.Kind() == reflect.Struct {
		var excepts []string
		for i := 0; i < s.NumField(); i++ {
			field := s.Field(i)
			kind := field.Kind()

			if kind == reflect.Ptr {
				field = field.Elem()
				kind = field.Kind()
			}

			if kind == reflect.Struct {
				excepts = append(excepts, s.Type().Field(i).Name)
				if err := Validate(field.Interface()); err != nil {
					return err
				}
			}

			if kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice {
				excepts = append(excepts, s.Type().Field(i).Name)
				if err := Validate(field.Interface(), s.Type().Field(i).Tag.Get("validate")); err != nil {
					return err
				}
			}
		}
		return Validator.StructExcept(schema, excepts...)
	}

	if s.Kind() == reflect.Map {
		kind := s.Type().Elem().Kind()
		if kind == reflect.Ptr {
			kind = s.Type().Elem().Elem().Kind()
		}

		if kind == reflect.Struct {
			for _, e := range s.MapKeys() {
				val := s.MapIndex(e)
				if err := Validate(val.Interface()); err != nil {
					return err
				}
			}
		} else {
			for _, e := range s.MapKeys() {
				val := s.MapIndex(e)
				if err := Validator.Var(val, tag[0]); err != nil {
					return err
				}
			}
		}
	}

	if s.Kind() == reflect.Array || s.Kind() == reflect.Slice {
		if len(tag) != 1 {
			return errors.New("no tag inserted for an array")
		}

		kind := s.Type().Elem().Kind()
		if kind == reflect.Ptr {
			kind = s.Type().Elem().Elem().Kind()
		}

		if kind == reflect.Struct {
			for i := 0; i < s.Len(); i++ {
				if err := Validate(s.Index(i).Interface()); err != nil {
					return err
				}
			}
		} else {
			for i := 0; i < s.Len(); i++ {
				if err := Validator.Var(s.Index(i).Interface(), tag[0]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

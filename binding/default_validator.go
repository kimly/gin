// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"sync"

	"gopkg.in/go-playground/validator.v9"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"fmt"
	"errors"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
	trans ut.Translator
}

var _ StructValidator = &defaultValidator{}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			//return err
			errs := err.(validator.ValidationErrors)
			errstring := ""
			for _, e := range errs {
				// can translate each error one at a time.
				errstring += fmt.Sprintf("%s", e.Translate(v.trans))
			}
			return errors.New(errstring)
		}
	}

	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		zh := zh.New()
		uni := ut.New(zh, zh)
		v.trans ,_ = uni.GetTranslator("zh")

		v.validate = validator.New()
		v.validate.SetTagName("binding")

		zh_translations.RegisterDefaultTranslations(v.validate, v.trans)
	})
}


func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
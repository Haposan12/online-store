package validator

import (
	ut "github.com/go-playground/universal-translator"
	validatorGo "github.com/go-playground/validator/v10"
	"log"
	"strings"
)

func registerCustomIndonesianTranslator(v *validatorGo.Validate, trans ut.Translator) {
	if err := v.RegisterTranslation("name", trans, func(ut ut.Translator) error {
		if err := ut.Add("name", "{0} tidak diperbolehkan mengandung angka atau tanda baca selain '.' dan ','.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}
	if err := v.RegisterTranslation("email_address", trans, func(ut ut.Translator) error {
		if err := ut.Add("email_address", "{0} format tidak valid.", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}
	if err := v.RegisterTranslation("address", trans, func(ut ut.Translator) error {
		if err := ut.Add("address", "{0} tidak diperbolehkan mengandung tanda baca selain ,.()-'/ (koma, titik, tanda kurung, strip, petik, garis miring).", false); err != nil {
			return err
		}
		return nil
	}, func(ut ut.Translator, fe validatorGo.FieldError) string {
		// first, clean/remove the comma
		cleaned := strings.Replace(fe.Param(), "-", " ", -1)

		// convert 'cleaned' comma separated string to slice
		strSlice := strings.Fields(cleaned)

		t, err := ut.T(fe.Tag(), fe.Field(), strings.Join(strSlice, ","))
		if err != nil {
			log.Printf("warning: error translating FieldError: %#v", fe)
			return fe.(error).Error()
		}
		return t
	}); err != nil {
		panic(err)
	}
}

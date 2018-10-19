package validation

import (
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	// Use single instance of Validate, it caches struct info
	validate *validator.Validate
)

// getValidator : Singleton returning a validator object
func getValidator() *validator.Validate {

	if validate == nil {
		// Init a new Validator
		validate = validator.New()
	}

	return validate
}

// ValidateStruct : Validate Struct passed as param. If validation fails, returns an array of fields that failed validations.
func ValidateStruct(s interface{}) ([]string, error) {

	// Fields that failed validation
	var fields []string

	// Validate Struct's exposed fields (Starting with a capitalized letter)
	err := getValidator().Struct(s)

	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fields = append(fields, err.Field())
		}

		return fields, err
	}

	return nil, nil
}

// ValidateStructExcept : Validate Struct passed as param except specified fields. If validation fails, returns an array of fields that failed validations.
func ValidateStructExcept(s interface{}, excludedFields ...string) ([]string, error) {

	// Fields that failed validation
	var fields []string

	// Validate Struct's exposed fields (Starting with a capitalized letter), except excluded fields
	err := getValidator().StructExcept(s, excludedFields...)

	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {

			fields = append(fields, err.Field())
		}

		return fields, err
	}

	return nil, nil

}

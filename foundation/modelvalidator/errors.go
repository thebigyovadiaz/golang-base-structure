package modelvalidator

import (
	"encoding/json"
	"errors"
)

// FieldError is used to indicate an error with a specific field.
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// FieldErrors represents a collection of fields that have errors.
type FieldErrors []FieldError

// NewFieldsError creates a FieldError.
func NewFieldsError(field string, err error) error {
	return FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Fields returns the fields that failed the validation.
func (fe FieldErrors) Fields() map[string]any {
	m := make(map[string]any)
	for _, fld := range fe {
		m[fld.Field] = fld.Err
	}
	return m
}

// IsFieldErrors checks if the given error contains a FieldErrors.
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

// GetFieldErrors returns a copy of the FieldErrors pointer.
func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}

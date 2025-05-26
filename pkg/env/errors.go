package env

import (
	"fmt"
)

// ErrMissingRequiredValue returned when a field with required=true contains no value or default
type ErrMissingRequiredValue struct {
	Value string
}

func (e ErrMissingRequiredValue) Error() string {
	return fmt.Sprintf("value for this env is missing and it's set as required [%s]", e.Value)
}

type ErrInvalidValue struct {
	Value string
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("value for this env is invalid [%s]", e.Value)
}

type ErrUnexportedType struct {
	Value string
}

func (e ErrUnexportedType) Error() string {
	return fmt.Sprintf("value for this env is unexported [%s]", e.Value)
}

type ErrUnsupportedField struct {
	Value string
}

func (e ErrUnsupportedField) Error() string {
	return fmt.Sprintf("value for this env is unsupported [%s]", e.Value)
}

type ErrInvalidEnvSet struct {
	Value string
}

func (e ErrInvalidEnvSet) Error() string {
	return fmt.Sprintf("items in environ must have format key=value [%s]", e.Value)
}

package secret

import (
	"fmt"
)

type ErrProjectIdBlank struct {
	Value string
}

func (e ErrProjectIdBlank) Error() string {
	return fmt.Sprintf("project id can't be blank [%s]", e.Value)
}

type ErrFailedToCreateClient struct {
	Value string
}

func (e ErrFailedToCreateClient) Error() string {
	return fmt.Sprintf("failed to create a new secret client [%s]", e.Value)
}

type ErrInvalidSecretName struct {
	Value string
}

func (e ErrInvalidSecretName) Error() string {
	return fmt.Sprintf("invalid secret name, secret name can't be blank [%s]", e.Value)
}

type ErrInvalidSecretVersion struct {
	Value string
}

func (e ErrInvalidSecretVersion) Error() string {
	return fmt.Sprintf("invalid secret version, secret version can't be blank [%s]", e.Value)
}

type ErrInvalidPattern struct {
	Value string
}

func (e ErrInvalidPattern) Error() string {
	return fmt.Sprintf("invalid pattern [%s]", e.Value)
}

type ErrFailedToListSecrets struct {
	Value string
}

func (e ErrFailedToListSecrets) Error() string {
	return fmt.Sprintf("failed to list secrets [%s]", e.Value)
}

type ErrFailedToCreateSecret struct {
	Value string
}

func (e ErrFailedToCreateSecret) Error() string {
	return fmt.Sprintf("failed to create secret [%s]", e.Value)
}

type ErrInvalidSecretData struct {
	Value string
}

func (e ErrInvalidSecretData) Error() string {
	return fmt.Sprintf("invalid secret data [%s]", e.Value)
}

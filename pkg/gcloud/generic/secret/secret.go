package secret

import (
	"encoding/json"
	"fmt"
	"regexp"

	sm "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"google.golang.org/api/iterator"
)

// New creates a new Secret client
// Parameters:
// opt: Option [Optional configuration for the client]
// Returns:
//   - SecretInterface[T]: The secret client
//   - error: An error if one occurs.
func New[T any](opts ...Option) (SecretInterface[T], error) {
	c := &secret[T]{conf: defaultConfig()}
	for _, opt := range opts {
		opt(c.conf)
	}
	client, err := sm.NewClient(c.conf.Context)
	if err != nil {
		return nil, ErrFailedToCreateClient{
			Value: fmt.Sprintf("failed to create a new secret client: %v", err),
		}
	}
	return &secret[T]{
		conf:   c.conf,
		client: client,
	}, nil
}

type SecretData struct {
	Data []byte
	Name string
}

// GetBytes gets a secret from Secret Manager
// Parameters:
//   - name: string [The secret name]
//
// Returns:
//   - []byte: The secret
//   - error: An error if one occurs.
func (s *secret[T]) GetBytes(name string) ([]byte, error) {
	if name == "" {
		return nil, ErrInvalidSecretName{Value: "invalid secret name"}
	}
	return s.GetVersion(name, "latest")
}

// Get gets a secret from Secret Manager
// Parameters:
//   - name: string [The secret name]
//
// Returns:
//   - T: The secret
//   - error: An error if one occurs.
func (s *secret[T]) Get(name string) (T, error) {
	var t T
	if name == "" {
		return t, ErrInvalidSecretName{Value: "invalid secret name"}
	}
	sec, err := s.GetVersion(name, "latest")
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(sec, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// GetVersion fetches secret based on name and version
// Parameters:
//   - name: string [The secret name]
//   - version: string [The secret version]
//
// Returns:
//   - []byte: The secret
//   - error: An error if one occurs.
func (s *secret[T]) GetVersion(name string, version string) ([]byte, error) {
	if name == "" {
		return nil, ErrInvalidSecretName{Value: "invalid secret name"}
	}
	if version == "" {
		return nil, ErrInvalidSecretVersion{Value: "invalid secret version"}
	}
	secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", s.conf.ProjectId, name, version)
	if s.client == nil {
		return nil, ErrFailedToCreateClient{
			Value: "secret manager client is not initialized",
		}
	}
	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	// Call the API.
	result, err := s.client.AccessSecretVersion(s.conf.Context, req)
	if err != nil {
		err = fmt.Errorf("failed to access secret version: %v", err)
		return nil, err
	}
	return result.Payload.Data, nil
}

// GetSecrets from projectId using secretsRegexp
// Parameters:
//   - secretsRegexp: *regexp.Regexp [The regular expression to match secrets]
//
// Returns:
//   - []SecretData: The secret data
//   - error: An error if one occurs.
func (s *secret[T]) GetSecrets(secretsRegexp *regexp.Regexp) ([]SecretData, error) {
	var secretsData []SecretData

	if secretsRegexp == nil {
		return nil, ErrInvalidPattern{Value: "nil pattern"}
	}

	if s.client == nil {
		return nil, ErrFailedToCreateClient{Value: "secret manager client is not initialized"}
	}

	parent := fmt.Sprintf("projects/%s", s.conf.ProjectId)
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}
	it := s.client.ListSecrets(s.conf.Context, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return secretsData, ErrFailedToListSecrets{Value: fmt.Sprintf("failed to fetch next secret: %v", err)}
		}

		if secretsRegexp.MatchString(resp.Name) {
			// Build the request.
			req := &secretmanagerpb.AccessSecretVersionRequest{
				Name: resp.Name + "/versions/latest",
			}

			// Call the API.
			result, err := s.client.AccessSecretVersion(s.conf.Context, req)
			if err != nil {
				return secretsData, fmt.Errorf("failed to access latest secret version %s: %v", resp.Name, err)
			}
			secretData := SecretData{
				Data: result.Payload.Data,
				Name: resp.Name,
			}
			secretsData = append(secretsData, secretData)
		}
	}
	return secretsData, nil
}

// AddSecretVersion adds a secret version to Secret Manager
// Parameters:
//   - secretName: string [The secret name]
//   - payload: []byte [The secret payload]
//
// Returns:
//   - error: An error if one occurs.
func (s *secret[T]) AddSecretVersion(secretName string, payload []byte) error {
	if secretName == "" {
		return ErrInvalidSecretName{Value: "invalid secret name"}
	}

	if payload == nil {
		return ErrInvalidSecretData{Value: "nil data"}
	}

	if len(payload) == 0 {
		return ErrInvalidSecretData{Value: "empty data"}
	}

	parent := "projects/" + s.conf.ProjectId + "/secrets/" + secretName

	if s.client == nil {
		return ErrFailedToCreateClient{Value: "secret manager client is not initialized"}
	}

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	_, err := s.client.AddSecretVersion(s.conf.Context, req)
	if err != nil {
		return ErrFailedToCreateSecret{Value: fmt.Sprintf("failed to add secret version: %v", err)}
	}
	return nil
}

// CreateSecret creates a secret in Secret Manager
// Parameters:
//   - secretName: string [The secret name]
//
// Returns:
//   - error: An error if one occurs.
func (s *secret[T]) CreateSecret(secretName string) error {
	if secretName == "" {
		return ErrInvalidSecretName{Value: "invalid secret name"}
	}

	if s.conf.ProjectId == "" {
		return ErrProjectIdBlank{Value: "project ID is required"}
	}

	parent := "projects/" + s.conf.ProjectId

	if s.client == nil {
		return ErrFailedToCreateClient{Value: "secret manager client is not initialized"}
	}

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretName,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API.
	_, err := s.client.CreateSecret(s.conf.Context, req)
	if err != nil {
		return ErrFailedToCreateSecret{Value: fmt.Sprintf("failed to create secret: %v", err)}
	}
	return nil
}

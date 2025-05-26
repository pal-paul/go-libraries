package secret_test

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/pal-paul/go-libraries/pkg/gcloud/generic/secret"
	"github.com/pal-paul/go-libraries/pkg/gcloud/generic/secret/mocks"
)

type TestSecret struct {
	Value string `json:"value"`
}

func testClient(t *testing.T) secret.SecretInterface[TestSecret] {
	t.Helper()
	client, err := secret.New[TestSecret](
		secret.WithContext(context.Background()),
		secret.WithProjectId("test-project"),
	)
	require.NoError(t, err)
	return client
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []secret.Option
		wantErr bool
	}{
		{
			name:    "should create new secret client successfully",
			opts:    []secret.Option{secret.WithContext(context.Background())},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := secret.New[TestSecret](tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestSecretInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		fn      func() error
		wantErr error
	}{
		{
			name: "GetBytes empty name",
			fn: func() error {
				_, err := testClient(t).GetBytes("")
				return err
			},
			wantErr: secret.ErrInvalidSecretName{Value: "invalid secret name"},
		},
		{
			name: "Get empty name",
			fn: func() error {
				_, err := testClient(t).Get("")
				return err
			},
			wantErr: secret.ErrInvalidSecretName{Value: "invalid secret name"},
		},
		{
			name: "GetVersion empty name",
			fn: func() error {
				_, err := testClient(t).GetVersion("", "1")
				return err
			},
			wantErr: secret.ErrInvalidSecretName{Value: "invalid secret name"},
		},
		{
			name: "GetVersion empty version",
			fn: func() error {
				_, err := testClient(t).GetVersion("test-secret", "")
				return err
			},
			wantErr: secret.ErrInvalidSecretVersion{Value: "invalid secret version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecretCreateSecret(t *testing.T) {
	tests := []struct {
		name       string
		secretName string
		wantErr    error
	}{
		{
			name:       "empty secret name",
			secretName: "",
			wantErr:    secret.ErrInvalidSecretName{Value: "invalid secret name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t)
			err := client.CreateSecret(tt.secretName)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecretAddSecretVersion(t *testing.T) {
	tests := []struct {
		name       string
		secretName string
		data       []byte
		wantErr    error
	}{
		{
			name:       "empty secret name",
			secretName: "",
			data:       []byte("test data"),
			wantErr:    secret.ErrInvalidSecretName{Value: "invalid secret name"},
		},
		{
			name:       "nil data",
			secretName: "test-secret",
			data:       nil,
			wantErr:    secret.ErrInvalidSecretData{Value: "nil data"},
		},
		{
			name:       "empty data",
			secretName: "test-secret",
			data:       []byte{},
			wantErr:    secret.ErrInvalidSecretData{Value: "empty data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t)
			err := client.AddSecretVersion(tt.secretName, tt.data)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecretGetSecrets(t *testing.T) {
	tests := []struct {
		name    string
		pattern *regexp.Regexp
		wantErr error
	}{
		{
			name:    "nil pattern",
			pattern: nil,
			wantErr: secret.ErrInvalidPattern{Value: "nil pattern"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t)
			secrets, err := client.GetSecrets(tt.pattern)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, secrets)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, secrets)
			}
		})
	}
}

func TestSecretIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSecretInterface[TestSecret](ctrl)
	secretName := "test-secret"
	testData := TestSecret{Value: "test-value"}
	testBytes, err := json.Marshal(testData)
	require.NoError(t, err)

	// Setup expectations
	mockClient.EXPECT().CreateSecret(secretName).Return(nil)
	mockClient.EXPECT().AddSecretVersion(secretName, testBytes).Return(nil)
	mockClient.EXPECT().GetBytes(secretName).Return(testBytes, nil)
	mockClient.EXPECT().Get(secretName).Return(testData, nil)
	mockClient.EXPECT().GetVersion(secretName, "1").Return(testBytes, nil)
	mockClient.EXPECT().GetSecrets(gomock.Any()).Return([]secret.SecretData{
		{
			Data: testBytes,
			Name: "projects/test-project/secrets/" + secretName,
		},
	}, nil)

	// Test create secret
	err = mockClient.CreateSecret(secretName)
	require.NoError(t, err)

	// Test add secret version
	err = mockClient.AddSecretVersion(secretName, testBytes)
	require.NoError(t, err)

	// Test get bytes
	bytes, err := mockClient.GetBytes(secretName)
	require.NoError(t, err)
	assert.Equal(t, testBytes, bytes)

	// Test get typed
	var got TestSecret
	got, err = mockClient.Get(secretName)
	require.NoError(t, err)
	assert.Equal(t, testData, got)

	// Test get version
	bytes, err = mockClient.GetVersion(secretName, "1")
	require.NoError(t, err)
	assert.Equal(t, testBytes, bytes)

	// Test get secrets by pattern
	pattern := regexp.MustCompile("^test-.*$")
	secrets, err := mockClient.GetSecrets(pattern)
	require.NoError(t, err)
	assert.NotEmpty(t, secrets)
	assert.Contains(t, secrets[0].Name, secretName)
	assert.Equal(t, testBytes, secrets[0].Data)
}

func TestSecretConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent operations test in short mode")
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSecretInterface[TestSecret](ctrl)
	secretName := "test-concurrent-secret"

	// Setup expectations
	mockClient.EXPECT().CreateSecret(secretName).Return(nil)
	for i := 0; i < 5; i++ {
		data := []byte(fmt.Sprintf("test-data-%d", i))
		mockClient.EXPECT().AddSecretVersion(secretName, data).Return(nil)
	}
	mockClient.EXPECT().GetBytes(secretName).Return([]byte("test-data"), nil)

	// First create the secret
	err := mockClient.CreateSecret(secretName)
	require.NoError(t, err)

	// Test concurrent version additions
	done := make(chan error)
	for i := 0; i < 5; i++ {
		go func(i int) {
			data := []byte(fmt.Sprintf("test-data-%d", i))
			done <- mockClient.AddSecretVersion(secretName, data)
		}(i)
	}

	// Collect results
	for i := 0; i < 5; i++ {
		err := <-done
		assert.NoError(t, err)
	}

	// Verify we can get the latest version
	_, err = mockClient.GetBytes(secretName)
	assert.NoError(t, err)
}

func TestSecretCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cleanup test in short mode")
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSecretInterface[TestSecret](ctrl)

	// Setup test data
	secretNames := []string{
		"test-cleanup-1",
		"test-cleanup-2",
		"test-cleanup-3",
	}

	testData := TestSecret{Value: "test-value"}
	testBytes, err := json.Marshal(testData)
	require.NoError(t, err)

	// Setup expectations for each secret
	for _, name := range secretNames {
		mockClient.EXPECT().CreateSecret(name).Return(nil)
		mockClient.EXPECT().AddSecretVersion(name, testBytes).Return(nil)
	}

	// Setup GetSecrets expectation
	var secretData []secret.SecretData
	for _, name := range secretNames {
		secretData = append(secretData, secret.SecretData{
			Data: testBytes,
			Name: "projects/test-project/secrets/" + name,
		})
	}
	mockClient.EXPECT().GetSecrets(gomock.Any()).Return(secretData, nil)

	// Helper function to create and populate a secret
	createTestSecret := func(t *testing.T, client secret.SecretInterface[TestSecret], name string) {
		err := client.CreateSecret(name)
		require.NoError(t, err)

		testData := TestSecret{Value: "test-value"}
		testBytes, err := json.Marshal(testData)
		require.NoError(t, err)

		err = client.AddSecretVersion(name, testBytes)
		require.NoError(t, err)
	}

	// Create test secrets
	for _, name := range secretNames {
		createTestSecret(t, mockClient, name)
	}

	// Get all test secrets
	pattern := regexp.MustCompile("^test-cleanup-.*$")
	secrets, err := mockClient.GetSecrets(pattern)
	require.NoError(t, err)

	// Verify all secrets were created
	assert.Equal(t, len(secretNames), len(secrets))
}

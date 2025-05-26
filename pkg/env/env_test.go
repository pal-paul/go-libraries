// Package env provides functionality for parsing environment variables into Go structs
package env

import (
	"os"
	"testing"
	"time"
)

type testConfig struct {
	String      string        `env:"TEST_STRING"`
	StringPtr   *string       `env:"TEST_STRING_PTR"`
	Bool        bool          `env:"TEST_BOOL"`
	Int         int           `env:"TEST_INT"`
	Float64     float64       `env:"TEST_FLOAT"`
	Duration    time.Duration `env:"TEST_DURATION"`
	Required    string        `env:"TEST_REQUIRED,required"`
	WithDefault string        `env:"TEST_DEFAULT,default=defaultValue"`
	Multiple    string        `env:"TEST_MULTIPLE_1,TEST_MULTIPLE_2"`
}

type customUnmarshaler struct {
	value string
}

func (c *customUnmarshaler) UnmarshalEnvironmentValue(data string) error {
	c.value = "custom_" + data
	return nil
}

type testConfigWithCustom struct {
	Custom customUnmarshaler `env:"TEST_CUSTOM"`
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		cfg     interface{}
		wantErr bool
	}{
		{
			name: "basic types",
			envs: map[string]string{
				"TEST_STRING":   "test",
				"TEST_BOOL":     "true",
				"TEST_INT":      "42",
				"TEST_FLOAT":    "3.14",
				"TEST_DURATION": "1h",
				"TEST_REQUIRED": "value",
			},
			cfg: &testConfig{},
		},
		{
			name: "with required field",
			envs: map[string]string{
				"TEST_REQUIRED": "value",
			},
			cfg: &testConfig{},
		},
		{
			name:    "missing required field",
			envs:    map[string]string{}, // Empty map to test required field validation
			cfg:     &testConfig{},
			wantErr: true,
		},
		{
			name: "with default value",
			envs: map[string]string{
				"TEST_REQUIRED": "value",
			},
			cfg: &testConfig{},
		},
		{
			name: "custom unmarshaler",
			envs: map[string]string{
				"TEST_CUSTOM": "value",
			},
			cfg: &testConfigWithCustom{},
		},
		{
			name: "multiple env names",
			envs: map[string]string{
				"TEST_MULTIPLE_2": "second",
				"TEST_REQUIRED":   "value",
			},
			cfg: &testConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set up environment variables
			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("Failed to set environment variable %s: %v", k, err)
				}
			}

			// Create envSet directly from test case
			es := make(envSet)
			for k, v := range tt.envs {
				es[k] = v
			}

			err := unmarshal(es, tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if _, ok := err.(*ErrMissingRequiredValue); !ok && tt.name == "missing required field" {
					t.Errorf("Expected ErrMissingRequiredValue but got %v", err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				switch c := tt.cfg.(type) {
				case *testConfig:
					if tt.envs["TEST_STRING"] != "" && c.String != tt.envs["TEST_STRING"] {
						t.Errorf("String = %v, want %v", c.String, tt.envs["TEST_STRING"])
					}
					if tt.envs["TEST_BOOL"] == "true" && !c.Bool {
						t.Errorf("Bool = %v, want true", c.Bool)
					}
					if c.WithDefault != "defaultValue" && tt.envs["TEST_DEFAULT"] == "" {
						t.Errorf("WithDefault = %v, want defaultValue", c.WithDefault)
					}
					if tt.envs["TEST_MULTIPLE_2"] != "" && c.Multiple != tt.envs["TEST_MULTIPLE_2"] {
						t.Errorf("Multiple = %v, want %v", c.Multiple, tt.envs["TEST_MULTIPLE_2"])
					}
				case *testConfigWithCustom:
					if tt.envs["TEST_CUSTOM"] != "" && c.Custom.value != "custom_"+tt.envs["TEST_CUSTOM"] {
						t.Errorf("Custom = %v, want %v", c.Custom.value, "custom_"+tt.envs["TEST_CUSTOM"])
					}
				}
			}
		})
	}
}

func TestInvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr error
	}{
		{
			name:    "nil pointer",
			input:   nil,
			wantErr: &ErrInvalidValue{Value: "<nil>"},
		},
		{
			name:    "non-pointer",
			input:   testConfig{},
			wantErr: &ErrInvalidValue{Value: "env.testConfig"},
		},
		{
			name:    "pointer to non-struct",
			input:   new(string),
			wantErr: &ErrInvalidValue{Value: "*string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := make(envSet)
			err := unmarshal(es, tt.input)
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}

			if err.Error() != tt.wantErr.Error() {
				t.Errorf("Got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		cfg     interface{}
		wantErr bool
	}{
		{
			name: "invalid bool",
			envs: map[string]string{
				"TEST_BOOL": "not-a-bool",
			},
			cfg:     &testConfig{},
			wantErr: true,
		},
		{
			name: "invalid int",
			envs: map[string]string{
				"TEST_INT": "not-an-int",
			},
			cfg:     &testConfig{},
			wantErr: true,
		},
		{
			name: "invalid float",
			envs: map[string]string{
				"TEST_FLOAT": "not-a-float",
			},
			cfg:     &testConfig{},
			wantErr: true,
		},
		{
			name: "invalid duration",
			envs: map[string]string{
				"TEST_DURATION": "not-a-duration",
			},
			cfg:     &testConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("Failed to set environment variable %s: %v", k, err)
				}
			}

			es, err := envToEnvSet(os.Environ())
			if err != nil {
				t.Fatalf("Failed to create envSet: %v", err)
			}

			err = unmarshal(es, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

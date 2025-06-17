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

type testConfigWithRequiredFalse struct {
	Optional       string `env:"TEST_OPTIONAL,required=false"`
	OptionalWith0  string `env:"TEST_OPTIONAL_0,required=0"`
	OptionalWithNo string `env:"TEST_OPTIONAL_NO,required=no"`
	Required       string `env:"TEST_REQUIRED,required"`
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
		{
			name: "required=false functionality",
			envs: map[string]string{
				"TEST_OPTIONAL":   "optional",
				"TEST_OPTIONAL_0": "optional_with_0",
				"TEST_REQUIRED":   "value",
			},
			cfg: &testConfigWithRequiredFalse{},
		},
		{
			name: "required=false should not require values",
			envs: map[string]string{
				"TEST_REQUIRED": "value", // Only provide required field
				// Deliberately omit TEST_OPTIONAL, TEST_OPTIONAL_0, TEST_OPTIONAL_NO
			},
			cfg: &testConfigWithRequiredFalse{},
		},
		{
			name: "required=false with provided values",
			envs: map[string]string{
				"TEST_REQUIRED":    "value",
				"TEST_OPTIONAL":    "optional_value",
				"TEST_OPTIONAL_0":  "zero_value",
				"TEST_OPTIONAL_NO": "no_value",
			},
			cfg: &testConfigWithRequiredFalse{},
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
				case *testConfigWithRequiredFalse:
					if tt.envs["TEST_OPTIONAL"] != "" && c.Optional != tt.envs["TEST_OPTIONAL"] {
						t.Errorf("Optional = %v, want %v", c.Optional, tt.envs["TEST_OPTIONAL"])
					}
					if tt.envs["TEST_OPTIONAL_0"] != "" && c.OptionalWith0 != tt.envs["TEST_OPTIONAL_0"] {
						t.Errorf("OptionalWith0 = %v, want %v", c.OptionalWith0, tt.envs["TEST_OPTIONAL_0"])
					}
					if tt.envs["TEST_OPTIONAL_NO"] != "" && c.OptionalWithNo != tt.envs["TEST_OPTIONAL_NO"] {
						t.Errorf("OptionalWithNo = %v, want %v", c.OptionalWithNo, tt.envs["TEST_OPTIONAL_NO"])
					}
					if c.Required != tt.envs["TEST_REQUIRED"] {
						t.Errorf("Required = %v, want %v", c.Required, tt.envs["TEST_REQUIRED"])
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

func TestRequiredFalseFeature(t *testing.T) {
	type config struct {
		RequiredField string `env:"REQUIRED_FIELD,required"`
		OptionalField string `env:"OPTIONAL_FIELD,required=false"`
	}

	// Test that we can unmarshal successfully with only the required field
	t.Run("missing optional field should not cause error", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("REQUIRED_FIELD", "test")
		// Deliberately do not set OPTIONAL_FIELD

		es, err := envToEnvSet(os.Environ())
		if err != nil {
			t.Fatalf("Failed to create envSet: %v", err)
		}

		cfg := &config{}
		err = unmarshal(es, cfg)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if cfg.RequiredField != "test" {
			t.Errorf("RequiredField = %v, want test", cfg.RequiredField)
		}
		if cfg.OptionalField != "" {
			t.Errorf("OptionalField = %v, want empty string", cfg.OptionalField)
		}
	})

	// Test that both fields work when provided
	t.Run("both fields provided", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("REQUIRED_FIELD", "test")
		os.Setenv("OPTIONAL_FIELD", "optional")

		es, err := envToEnvSet(os.Environ())
		if err != nil {
			t.Fatalf("Failed to create envSet: %v", err)
		}

		cfg := &config{}
		err = unmarshal(es, cfg)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if cfg.RequiredField != "test" {
			t.Errorf("RequiredField = %v, want test", cfg.RequiredField)
		}
		if cfg.OptionalField != "optional" {
			t.Errorf("OptionalField = %v, want optional", cfg.OptionalField)
		}
	})
}

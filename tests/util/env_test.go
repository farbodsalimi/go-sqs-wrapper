package util

import (
	"go-sqs-wrapper/src/util"
	"os"
	"strconv"
	"testing"
)

func TestGetEnvStrDefaultValue(t *testing.T) {
	d := "default value"
	env, _ := util.GetEnvStr(util.GetEnvParams{Key: "mock_key", DefaultValue: d})
	if env != "default value" {
		t.Errorf("TestGetEnvStrDefaultValue FAILED, expected %v but got value %v", d, env)
	}
}

func TestGetEnvStr(t *testing.T) {
	mockKey := "mock_key"
	mockValue := "mock_value"
	os.Setenv(mockKey, mockValue)

	env, _ := util.GetEnvStr(util.GetEnvParams{Key: "mock_key", DefaultValue: "default"})
	if env != mockValue {
		t.Errorf("TestGetEnvStr FAILED, expected %v but got value %v", mockValue, env)

	}

	os.Unsetenv(mockKey)
}

func TestGetEnvStrError(t *testing.T) {
	mockKey := "mock_key"

	_, err := util.GetEnvStr(util.GetEnvParams{Key: mockKey})
	if err == nil {
		t.Errorf("TestGetEnvStrError FAILED, expected %v but got value %v", util.ErrEnvVarEmpty(mockKey), err)
	}
}

func TestGetEnvIntDefaultValue(t *testing.T) {
	d := "123"
	env, _ := util.GetEnvInt(util.GetEnvParams{Key: "mock_key", DefaultValue: d})
	if env != 123 {
		t.Errorf("TestGetEnvIntDefaultValue FAILED, expected %v but got value %v", d, env)
	}
}

func TestGetEnvInt(t *testing.T) {
	mockKey := "mock_key"
	mockValue := 123

	os.Setenv(mockKey, strconv.Itoa(mockValue))

	env, err := util.GetEnvInt(util.GetEnvParams{Key: "mock_key", DefaultValue: "321"})
	if env != mockValue {
		t.Log(err)
		t.Errorf("TesGetEnvInt FAILED, expected %v but got value %v", mockValue, env)
	}

	os.Unsetenv(mockKey)
}

func TestGetEnvIntInvalidIntiger(t *testing.T) {
	mockKey := "mock_key"
	os.Setenv(mockKey, "abc")

	env, _ := util.GetEnvInt(util.GetEnvParams{Key: "mock_key", DefaultValue: "321"})
	if env != 0 {
		t.Errorf("TesGetEnvInt FAILED, expected %v but got value %v", 0, env)
	}

	os.Unsetenv(mockKey)
}

func TestGetEnvIntError(t *testing.T) {
	mockKey := "mock_key"

	_, err := util.GetEnvInt(util.GetEnvParams{Key: mockKey})
	if err == nil {
		t.Errorf("TestGetEnvIntError FAILED, expected %v but got value %v", util.ErrEnvVarEmpty(mockKey), err)
	}
}

func TestGetEnvBoolDefaultValue(t *testing.T) {
	d := "true"
	env, _ := util.GetEnvBool(util.GetEnvParams{Key: "mock_key", DefaultValue: d})
	if env != true {
		t.Errorf("TestGetEnvBoolDefaultValue FAILED, expected %v but got value %v", d, env)
	}
}

func TestGetEnvBool(t *testing.T) {
	mockKey := "mock_key"
	mockValue := "true"
	mockValueBool, _ := strconv.ParseBool(mockValue)
	os.Setenv(mockKey, mockValue)

	env, _ := util.GetEnvBool(util.GetEnvParams{Key: "mock_key", DefaultValue: "false"})
	if env != mockValueBool {
		t.Errorf("TestGetEnvBool FAILED, expected %v but got value %v", mockValueBool, env)

	}

	os.Unsetenv(mockKey)
}

func TestGetEnvInvalidBool(t *testing.T) {
	mockKey := "mock_key"
	mockValue := "not_true"
	os.Setenv(mockKey, mockValue)

	_, err := util.GetEnvBool(util.GetEnvParams{Key: "mock_key"})
	if err == nil {
		t.Errorf("TestGetEnvInvalidBool FAILED, expected %v but got value %v", "invalid syntax err", err)
	}

	os.Unsetenv(mockKey)
}

func TestGetEnvBoolError(t *testing.T) {
	mockKey := "mock_key"

	_, err := util.GetEnvBool(util.GetEnvParams{Key: mockKey})
	if err == nil {
		t.Errorf("TestGetEnvBoolError FAILED, expected %v but got value %v", util.ErrEnvVarEmpty(mockKey), err)
	}
}

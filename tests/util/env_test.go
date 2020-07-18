package util

import (
	"os"
	"strconv"
	"testing"

	"go-sqs-wrapper/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvStrDefaultValue(t *testing.T) {
	d := "default value"
	env, _ := util.GetEnvStr(util.GetEnvParams{Key: "mock_key", DefaultValue: d})
	assert.Equal(t, env, d)
}

func TestGetEnvStr(t *testing.T) {
	mockKey := "mock_key"
	mockValue := "mock_value"
	os.Setenv(mockKey, mockValue)

	env, _ := util.GetEnvStr(util.GetEnvParams{Key: "mock_key", DefaultValue: "default"})
	assert.Equal(t, env, mockValue)

	os.Unsetenv(mockKey)
}

func TestGetEnvStrError(t *testing.T) {
	_, err := util.GetEnvStr(util.GetEnvParams{Key: "mock_key"})
	assert.NotNil(t, err)
}

func TestGetEnvIntDefaultValue(t *testing.T) {
	env, _ := util.GetEnvInt(util.GetEnvParams{Key: "mock_key", DefaultValue: "123"})
	assert.Equal(t, env, 123)
}

func TestGetEnvInt(t *testing.T) {
	mockKey := "mock_key"
	mockValue := 123

	os.Setenv(mockKey, strconv.Itoa(mockValue))

	env, _ := util.GetEnvInt(util.GetEnvParams{Key: "mock_key", DefaultValue: "321"})
	assert.Equal(t, env, mockValue)

	os.Unsetenv(mockKey)
}

func TestGetEnvIntInvalidIntiger(t *testing.T) {
	mockKey := "mock_key"
	os.Setenv(mockKey, "abc")

	env, _ := util.GetEnvInt(util.GetEnvParams{Key: "mock_key", DefaultValue: "321"})
	assert.Equal(t, env, 0)

	os.Unsetenv(mockKey)
}

func TestGetEnvIntError(t *testing.T) {
	_, err := util.GetEnvInt(util.GetEnvParams{Key: "mock_key"})
	assert.NotNil(t, err)
}

func TestGetEnvInt64DefaultValue(t *testing.T) {
	env, _ := util.GetEnvInt64(util.GetEnvParams{Key: "mock_key", DefaultValue: "123"})
	assert.Equal(t, env, int64(123))
}

func TestGetEnvInt64(t *testing.T) {
	mockKey := "mock_key"
	mockValue := int64(123)

	os.Setenv(mockKey, strconv.Itoa(int(mockValue)))

	env, _ := util.GetEnvInt64(util.GetEnvParams{Key: "mock_key", DefaultValue: "321"})
	assert.Equal(t, env, mockValue)

	os.Unsetenv(mockKey)
}

func TestGetEnvInt64InvalidIntiger(t *testing.T) {
	mockKey := "mock_key"
	os.Setenv(mockKey, "abc")

	env, _ := util.GetEnvInt64(util.GetEnvParams{Key: "mock_key", DefaultValue: "321"})
	assert.Equal(t, env, int64(0))

	os.Unsetenv(mockKey)
}

func TestGetEnvInt64Error(t *testing.T) {
	_, err := util.GetEnvInt64(util.GetEnvParams{Key: "mock_key"})
	assert.NotNil(t, err)
}

func TestGetEnvBoolDefaultValue(t *testing.T) {
	env, _ := util.GetEnvBool(util.GetEnvParams{Key: "mock_key", DefaultValue: "true"})
	assert.True(t, env)
}

func TestGetEnvBool(t *testing.T) {
	mockKey := "mock_key"
	mockValue := "true"
	os.Setenv(mockKey, mockValue)

	env, _ := util.GetEnvBool(util.GetEnvParams{Key: "mock_key", DefaultValue: "false"})
	assert.True(t, env)

	os.Unsetenv(mockKey)
}

func TestGetEnvInvalidBool(t *testing.T) {
	mockKey := "mock_key"
	mockValue := "not_true"
	os.Setenv(mockKey, mockValue)

	_, err := util.GetEnvBool(util.GetEnvParams{Key: "mock_key"})
	assert.NotNil(t, err)

	os.Unsetenv(mockKey)
}

func TestGetEnvBoolError(t *testing.T) {
	_, err := util.GetEnvBool(util.GetEnvParams{Key: "mock_key"})
	assert.NotNil(t, err)
}

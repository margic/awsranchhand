package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitEnv(t *testing.T) {
	key, value := splitEnv("MYKEY=MYVAL")
	assert.Equal(t, "MYKEY", key, "Testing key is correct")
	assert.Equal(t, "MYVAL", value, "Testing value is correct")
}

func TestUpdateEnvironment(t *testing.T) {
	var envs []string
	envs = append(envs, "MYKEY=MYVAL")
	envs = append(envs, "KEY2=123")
	envs = append(envs, "KEY3=321")

	envMap := make(map[string]interface{})

	envMap = updateEnvironment(envMap, envs, false)

	assert.Equal(t, "MYVAL", envMap["MYKEY"], "Environment variable not added to map")
	assert.Equal(t, "123", envMap["KEY2"], "Environment variable not added to map")
	assert.Equal(t, "321", envMap["KEY3"], "Environment variable not added to map")
}

func TestUpdateEnvironmentWithClear(t *testing.T) {
	var envs []string
	envs = append(envs, "MYKEY=MYVAL")
	envs = append(envs, "KEY2=123")
	envs = append(envs, "KEY3=321")

	envMap := make(map[string]interface{})
	envMap["existing"] = "myval"
	envMap = updateEnvironment(envMap, envs, true)

	assert.Nil(t, envMap["existing"], "Should remove existing value")
	assert.Equal(t, "MYVAL", envMap["MYKEY"], "Environment variable not added to map")
	assert.Equal(t, "123", envMap["KEY2"], "Environment variable not added to map")
	assert.Equal(t, "321", envMap["KEY3"], "Environment variable not added to map")
}

func TestUpdateEnvironmentUpdateValue(t *testing.T) {
	var envs []string
	envs = append(envs, "MYKEY=MYVAL")
	envs = append(envs, "KEY2=123")
	envs = append(envs, "KEY3=321")

	envMap := make(map[string]interface{})
	// put in existing value
	envMap["MYKEY"] = "NOTMYVAL"

	envMap = updateEnvironment(envMap, envs, false)

	assert.Equal(t, "MYVAL", envMap["MYKEY"], "Environment variable not updated")
	assert.Equal(t, "123", envMap["KEY2"], "Environment variable not added to map")
	assert.Equal(t, "321", envMap["KEY3"], "Environment variable not added to map")
}

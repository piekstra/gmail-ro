package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigCommand(t *testing.T) {
	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "config", configCmd.Use)
	})

	t.Run("has short description", func(t *testing.T) {
		assert.NotEmpty(t, configCmd.Short)
	})

	t.Run("has subcommands", func(t *testing.T) {
		subcommands := configCmd.Commands()
		assert.GreaterOrEqual(t, len(subcommands), 3)

		var names []string
		for _, cmd := range subcommands {
			names = append(names, cmd.Name())
		}
		assert.Contains(t, names, "show")
		assert.Contains(t, names, "test")
		assert.Contains(t, names, "clear")
	})
}

func TestConfigShowCommand(t *testing.T) {
	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "show", configShowCmd.Use)
	})

	t.Run("requires no arguments", func(t *testing.T) {
		err := configShowCmd.Args(configShowCmd, []string{})
		assert.NoError(t, err)

		err = configShowCmd.Args(configShowCmd, []string{"extra"})
		assert.Error(t, err)
	})

	t.Run("has short description", func(t *testing.T) {
		assert.NotEmpty(t, configShowCmd.Short)
	})

	t.Run("has long description", func(t *testing.T) {
		assert.NotEmpty(t, configShowCmd.Long)
	})
}

func TestConfigTestCommand(t *testing.T) {
	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "test", configTestCmd.Use)
	})

	t.Run("requires no arguments", func(t *testing.T) {
		err := configTestCmd.Args(configTestCmd, []string{})
		assert.NoError(t, err)

		err = configTestCmd.Args(configTestCmd, []string{"extra"})
		assert.Error(t, err)
	})

	t.Run("has short description", func(t *testing.T) {
		assert.NotEmpty(t, configTestCmd.Short)
	})

	t.Run("has long description", func(t *testing.T) {
		assert.NotEmpty(t, configTestCmd.Long)
	})
}

func TestConfigClearCommand(t *testing.T) {
	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "clear", configClearCmd.Use)
	})

	t.Run("requires no arguments", func(t *testing.T) {
		err := configClearCmd.Args(configClearCmd, []string{})
		assert.NoError(t, err)

		err = configClearCmd.Args(configClearCmd, []string{"extra"})
		assert.Error(t, err)
	})

	t.Run("has short description", func(t *testing.T) {
		assert.NotEmpty(t, configClearCmd.Short)
	})

	t.Run("has long description", func(t *testing.T) {
		assert.NotEmpty(t, configClearCmd.Long)
		assert.Contains(t, configClearCmd.Long, "token")
	})
}

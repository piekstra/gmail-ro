package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "gmail-ro", rootCmd.Use)
	})

	t.Run("has short description", func(t *testing.T) {
		assert.NotEmpty(t, rootCmd.Short)
		assert.Contains(t, rootCmd.Short, "read-only")
	})

	t.Run("has long description", func(t *testing.T) {
		assert.NotEmpty(t, rootCmd.Long)
		assert.Contains(t, rootCmd.Long, "Gmail")
		assert.Contains(t, rootCmd.Long, "read-only")
	})

	t.Run("has subcommands", func(t *testing.T) {
		commands := rootCmd.Commands()
		assert.NotEmpty(t, commands)

		// Check for expected commands
		var names []string
		for _, cmd := range commands {
			names = append(names, cmd.Name())
		}
		assert.Contains(t, names, "search")
		assert.Contains(t, names, "read")
		assert.Contains(t, names, "thread")
		assert.Contains(t, names, "version")
	})
}

func TestVersionCommand(t *testing.T) {
	t.Run("outputs version", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)

		// Set a known version for testing
		oldVersion := Version
		Version = "test-version"
		defer func() { Version = oldVersion }()

		err := versionCmd.RunE(versionCmd, []string{})
		assert.NoError(t, err)

		// Note: output goes to stdout, not the buffer set on rootCmd
		// This test verifies the command doesn't panic
	})

	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "version", versionCmd.Use)
	})
}

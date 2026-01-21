package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadCommand(t *testing.T) {
	t.Run("has correct use", func(t *testing.T) {
		assert.Equal(t, "thread <id>", threadCmd.Use)
	})

	t.Run("requires exactly one argument", func(t *testing.T) {
		err := threadCmd.Args(threadCmd, []string{})
		assert.Error(t, err)

		err = threadCmd.Args(threadCmd, []string{"thread123"})
		assert.NoError(t, err)

		err = threadCmd.Args(threadCmd, []string{"thread1", "thread2"})
		assert.Error(t, err)
	})

	t.Run("has json flag", func(t *testing.T) {
		flag := threadCmd.Flags().Lookup("json")
		assert.NotNil(t, flag)
		assert.Equal(t, "j", flag.Shorthand)
		assert.Equal(t, "false", flag.DefValue)
	})

	t.Run("has short description", func(t *testing.T) {
		assert.NotEmpty(t, threadCmd.Short)
		assert.Contains(t, threadCmd.Short, "thread")
	})

	t.Run("long description explains thread ID", func(t *testing.T) {
		assert.Contains(t, threadCmd.Long, "thread ID")
		assert.Contains(t, threadCmd.Long, "message ID")
	})
}

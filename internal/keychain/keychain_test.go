package keychain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestConfigFile_TokenRoundTrip(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Create test token
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour).Truncate(time.Second),
	}

	// Store token
	err := setInConfigFile(token)
	require.NoError(t, err)

	// Retrieve token
	retrieved, err := getFromConfigFile()
	require.NoError(t, err)

	assert.Equal(t, token.AccessToken, retrieved.AccessToken)
	assert.Equal(t, token.RefreshToken, retrieved.RefreshToken)
	assert.Equal(t, token.TokenType, retrieved.TokenType)
	// Compare times with tolerance for JSON marshaling
	assert.WithinDuration(t, token.Expiry, retrieved.Expiry, time.Second)
}

func TestConfigFile_Permissions(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	token := &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}

	err := setInConfigFile(token)
	require.NoError(t, err)

	// Check file permissions
	path := filepath.Join(tmpDir, serviceName, tokenFile)
	info, err := os.Stat(path)
	require.NoError(t, err)

	// Verify 0600 permissions (read/write for owner only)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestConfigFile_DirectoryPermissions(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	token := &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}

	err := setInConfigFile(token)
	require.NoError(t, err)

	// Check directory permissions
	dir := filepath.Join(tmpDir, serviceName)
	info, err := os.Stat(dir)
	require.NoError(t, err)

	// Verify 0700 permissions (read/write/execute for owner only)
	assert.Equal(t, os.FileMode(0700), info.Mode().Perm())
}

func TestConfigFile_NotFound(t *testing.T) {
	// Create temp directory (empty)
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	_, err := getFromConfigFile()
	assert.ErrorIs(t, err, ErrTokenNotFound)
}

func TestConfigFile_InvalidJSON(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Create config directory
	configDir := filepath.Join(tmpDir, serviceName)
	err := os.MkdirAll(configDir, 0700)
	require.NoError(t, err)

	// Write invalid JSON
	path := filepath.Join(configDir, tokenFile)
	err = os.WriteFile(path, []byte("invalid json"), 0600)
	require.NoError(t, err)

	_, err = getFromConfigFile()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token file")
}

func TestConfigFile_Overwrite(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Store first token
	token1 := &oauth2.Token{
		AccessToken: "first-token",
		TokenType:   "Bearer",
	}
	err := setInConfigFile(token1)
	require.NoError(t, err)

	// Store second token
	token2 := &oauth2.Token{
		AccessToken: "second-token",
		TokenType:   "Bearer",
	}
	err = setInConfigFile(token2)
	require.NoError(t, err)

	// Retrieve should return second token
	retrieved, err := getFromConfigFile()
	require.NoError(t, err)
	assert.Equal(t, "second-token", retrieved.AccessToken)
}

func TestConfigFile_Delete(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Store token
	token := &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}
	err := setInConfigFile(token)
	require.NoError(t, err)

	// Delete token
	err = deleteFromConfigFile()
	require.NoError(t, err)

	// Should be gone
	_, err = getFromConfigFile()
	assert.ErrorIs(t, err, ErrTokenNotFound)
}

func TestConfigFile_DeleteNonExistent(t *testing.T) {
	// Create temp directory (empty)
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Delete should not error on non-existent file
	err := deleteFromConfigFile()
	assert.NoError(t, err)
}

func TestMigrateFromFile_NoFile(t *testing.T) {
	// Migration should succeed (no-op) when file doesn't exist
	err := MigrateFromFile("/nonexistent/path/token.json")
	assert.NoError(t, err)
}

func TestMigrateFromFile_InvalidJSON(t *testing.T) {
	// This test verifies that invalid JSON in a token file causes an error
	// when migration is attempted without existing secure storage.
	//
	// Note: On macOS, if there's already a token in the system keychain,
	// MigrateFromFile will skip reading the file (already migrated).
	// This is expected behavior - the test validates the JSON parsing path
	// when no secure storage token exists.

	tmpDir := t.TempDir()
	configDir := t.TempDir()

	// Override config directory so we don't find existing tokens
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", configDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Create temp file with invalid JSON
	tokenPath := filepath.Join(tmpDir, "token.json")
	err := os.WriteFile(tokenPath, []byte("invalid json"), 0600)
	require.NoError(t, err)

	// If secure storage has a token (e.g., from real keychain), migration is skipped
	// In that case, we test the direct file parsing instead
	if IsSecureStorage() && HasStoredToken() {
		t.Skip("Secure storage already has a token, migration would be skipped")
	}

	err = MigrateFromFile(tokenPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token file")
}

func TestMigrateFromFile_Success(t *testing.T) {
	// This test verifies that migration from file to storage works correctly.
	// On macOS, if there's already a token in the system keychain,
	// migration is skipped (considered already migrated).

	// Skip if secure storage already has a token
	if IsSecureStorage() && HasStoredToken() {
		t.Skip("Secure storage already has a token, migration would be skipped")
	}

	// Create temp directories
	tmpDir := t.TempDir()
	configDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", configDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Create valid token file
	token := &oauth2.Token{
		AccessToken:  "migrated-token",
		RefreshToken: "migrated-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}
	tokenPath := filepath.Join(tmpDir, "token.json")
	data, err := json.Marshal(token)
	require.NoError(t, err)
	err = os.WriteFile(tokenPath, data, 0600)
	require.NoError(t, err)

	// Run migration
	err = MigrateFromFile(tokenPath)
	require.NoError(t, err)

	// Verify token was stored (in config file since we use temp config dir)
	retrieved, err := getFromConfigFile()
	require.NoError(t, err)
	assert.Equal(t, "migrated-token", retrieved.AccessToken)

	// Verify backup was created
	backupPath := tokenPath + ".backup"
	_, err = os.Stat(backupPath)
	assert.NoError(t, err)

	// Verify original was removed
	_, err = os.Stat(tokenPath)
	assert.True(t, os.IsNotExist(err))
}

func TestHasStoredToken_ConfigFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Override config directory for test
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Test file-based storage specifically
	// Note: HasStoredToken() may find tokens in system keychain on macOS,
	// so we test the underlying file functions directly

	// Should return error when no token file
	_, err := getFromConfigFile()
	assert.ErrorIs(t, err, ErrTokenNotFound)

	// Store a token in config file
	token := &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}
	err = setInConfigFile(token)
	require.NoError(t, err)

	// Should successfully retrieve from config file
	retrieved, err := getFromConfigFile()
	require.NoError(t, err)
	assert.Equal(t, "test-token", retrieved.AccessToken)
}

func TestGetStorageBackend(t *testing.T) {
	// Just verify it returns a valid backend
	backend := GetStorageBackend()
	assert.Contains(t, []StorageBackend{BackendKeychain, BackendSecretTool, BackendFile}, backend)
}

func TestIsSecureStorage(t *testing.T) {
	// This will vary by platform - just verify it returns a bool
	result := IsSecureStorage()
	assert.IsType(t, true, result)
}

func TestTokenFilePath(t *testing.T) {
	// Test with XDG_CONFIG_HOME set
	tmpDir := t.TempDir()
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	path, err := tokenFilePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(tmpDir, serviceName, tokenFile), path)
}

func TestConfigDir(t *testing.T) {
	// Test with XDG_CONFIG_HOME set
	tmpDir := t.TempDir()
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	dir, err := configDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(tmpDir, serviceName), dir)
}

func TestConfigDir_NoXDG(t *testing.T) {
	// Test without XDG_CONFIG_HOME
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	dir, err := configDir()
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".config", serviceName), dir)
}

package host

import (
	"os"
	"path/filepath"

	"github.com/tetratelabs/wazero"
)

// NewPersistentCompilationCache returns a compilation cache that persists on disk.
// It uses the user's cache directory (e.g., ~/.cache/<appName>/wasm on Linux).
// If disk access fails, it falls back to an in-memory cache.
func NewPersistentCompilationCache(appName string) wazero.CompilationCache {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to home if UserCacheDir fails
		home, err := os.UserHomeDir()
		if err != nil {
			return wazero.NewCompilationCache()
		}
		cacheDir = filepath.Join(home, ".cache")
	}

	path := filepath.Join(cacheDir, appName, "wasm")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return wazero.NewCompilationCache()
	}

	cache, err := wazero.NewCompilationCacheWithDir(path)
	if err != nil {
		return wazero.NewCompilationCache()
	}

	return cache
}

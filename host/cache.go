package host

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tetratelabs/wazero"
)

// CompilationCache defines the behavior for storing compiled Wasm modules.
// This interface decouples the SDK from the specific runtime (wazero).
type CompilationCache interface {
	Close(ctx context.Context) error
}

// wazeroCacheAdapter wraps wazero.CompilationCache to satisfy our local interface.
type wazeroCacheAdapter struct {
	inner wazero.CompilationCache
}

func (w *wazeroCacheAdapter) Close(ctx context.Context) error {
	return w.inner.Close(ctx)
}

// NewPersistentCompilationCache returns a compilation cache that persists on disk.
// It uses the user's cache directory (e.g., ~/.cache/<appName>/wasm on Linux).
// If disk access fails, it falls back to an in-memory cache.
func NewPersistentCompilationCache(appName string) CompilationCache {
	var cache wazero.CompilationCache
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to home if UserCacheDir fails
		home, err := os.UserHomeDir()
		if err != nil {
			return &wazeroCacheAdapter{inner: wazero.NewCompilationCache()}
		}
		cacheDir = filepath.Join(home, ".cache")
	}

	path := filepath.Join(cacheDir, appName, "wasm")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return &wazeroCacheAdapter{inner: wazero.NewCompilationCache()}
	}

	cache, err = wazero.NewCompilationCacheWithDir(path)
	if err != nil {
		return &wazeroCacheAdapter{inner: wazero.NewCompilationCache()}
	}

	return &wazeroCacheAdapter{inner: cache}
}

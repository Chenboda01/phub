package discovery

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoadOptions_usesEnvironmentOverrides_whenValuesAreConfigured(t *testing.T) {
	// Given
	roots := []string{
		filepath.Join(t.TempDir(), "projects"),
		filepath.Join(t.TempDir(), "code"),
	}
	t.Setenv(RootsEnvironment, roots[0]+string(os.PathListSeparator)+roots[1])
	t.Setenv(DepthEnvironment, "2")

	// When
	got, err := LoadOptions()

	// Then
	if err != nil {
		t.Fatalf("LoadOptions() error = %v", err)
	}
	if !slices.Equal(got.Roots(), roots) {
		t.Fatalf("Roots() = %q, want %q", got.Roots(), roots)
	}
	if got.MaxDepth() != 2 {
		t.Fatalf("MaxDepth() = %d, want 2", got.MaxDepth())
	}
}

func TestLoadOptions_returnsInvalidDepthError_whenDepthIsNotAnInteger(t *testing.T) {
	// Given
	t.Setenv(DepthEnvironment, "many")

	// When
	_, err := LoadOptions()

	// Then
	if !errors.Is(err, ErrInvalidDepth) {
		t.Fatalf("LoadOptions() error = %v, want ErrInvalidDepth", err)
	}
}

func TestLoadOptions_returnsInvalidDepthError_whenDepthIsNegative(t *testing.T) {
	// Given
	t.Setenv(DepthEnvironment, "-1")

	// When
	_, err := LoadOptions()

	// Then
	if !errors.Is(err, ErrInvalidDepth) {
		t.Fatalf("LoadOptions() error = %v, want ErrInvalidDepth", err)
	}
}

func TestLoadOptions_allowsNoRoots_whenRootsOverrideIsEmpty(t *testing.T) {
	// Given
	t.Setenv(RootsEnvironment, "")

	// When
	got, err := LoadOptions()

	// Then
	if err != nil {
		t.Fatalf("LoadOptions() error = %v", err)
	}
	if len(got.Roots()) != 0 {
		t.Fatalf("Roots() = %q, want no roots", got.Roots())
	}
}

func TestLoadOptions_usesHomeDirectoryAsDefaultRoot_whenRootsAreNotConfigured(t *testing.T) {
	// Given
	previousRoots, wasConfigured := os.LookupEnv(RootsEnvironment)
	if err := os.Unsetenv(RootsEnvironment); err != nil {
		t.Fatalf("unset %s: %v", RootsEnvironment, err)
	}
	t.Cleanup(func() {
		if wasConfigured {
			_ = os.Setenv(RootsEnvironment, previousRoots)
			return
		}
		_ = os.Unsetenv(RootsEnvironment)
	})
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("resolve home directory: %v", err)
	}

	// When
	got, err := LoadOptions()

	// Then
	if err != nil {
		t.Fatalf("LoadOptions() error = %v", err)
	}
	if !slices.Equal(got.Roots(), []string{home}) {
		t.Fatalf("Roots() = %q, want %q", got.Roots(), []string{home})
	}
}

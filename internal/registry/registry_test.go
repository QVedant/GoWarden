package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempRegistry(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "languages.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp registry: %v", err)
	}
	return path
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempRegistry(t, `
languages:
  python:
    display_name: "Python 3"
    extension: ".py"
    compile: []
    run: ["python3", "{file}"]
    timeout_seconds: 5
    memory_mb: 128
    image_toolchain: "python3.11"
`)

	reg, err := Load(path)
	if err != nil {
		t.Fatalf("expected valid registry to load, got error: %v", err)
	}
	if _, ok := reg.Get("python"); !ok {
		t.Fatal("expected 'python' language to be present")
	}
}

func TestLoad_MissingFilePlaceholder(t *testing.T) {
	path := writeTempRegistry(t, `
languages:
  broken:
    display_name: "Broken"
    extension: ".x"
    compile: []
    run: ["echo", "hi"]
    timeout_seconds: 5
    memory_mb: 128
    image_toolchain: "none"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for run command missing {file} placeholder")
	}
}

func TestLoad_CompileWithoutBinaryPlaceholder(t *testing.T) {
	path := writeTempRegistry(t, `
languages:
  c:
    display_name: "C"
    extension: ".c"
    compile: ["gcc", "-o", "out", "{file}"]
    run: ["out"]
    timeout_seconds: 5
    memory_mb: 128
    image_toolchain: "gcc"
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when compile step exists but run has no {binary} placeholder")
	}
}

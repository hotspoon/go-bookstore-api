package logging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesSeparateLogFiles(t *testing.T) {
	dir := t.TempDir()
	appPath := filepath.Join(dir, "app.log")
	accessPath := filepath.Join(dir, "access.log")

	files, err := Init(appPath, accessPath)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	files.Close()

	for _, path := range []string{appPath, accessPath} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("log file %q was not created: %v", path, err)
		}
	}
}

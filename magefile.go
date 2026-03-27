//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Check runs the primary local verification suite.
func Check() error {
	if err := VerifyBootstrap(); err != nil {
		return err
	}
	if err := FmtCheck(); err != nil {
		return err
	}
	if err := Test(); err != nil {
		return err
	}
	return nil
}

// CI runs the CI-equivalent verification suite.
func CI() error {
	return Check()
}

// VerifyBootstrap ensures the expected repository seed files exist.
func VerifyBootstrap() error {
	required := []string{
		"AGENTS.md",
		"PLAN.md",
		"README.md",
		"magefile.go",
		"go.mod",
		".github/workflows/ci.yml",
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("verify bootstrap %q: %w", path, err)
		}
	}
	return nil
}

// Fmt formats Go files with gofmt.
func Fmt() error {
	files, err := goFiles(".")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"-w"}, files...)
	return run("gofmt", args...)
}

// FmtCheck reports unformatted Go files.
func FmtCheck() error {
	files, err := goFiles(".")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"-l"}, files...)
	out, err := output("gofmt", args...)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("gofmt required for:\n%s", strings.TrimSpace(out))
	}
	return nil
}

// Test runs the Go test suite.
func Test() error {
	return run("go", "test", "./...")
}

// Build compiles the demo command when it exists.
func Build() error {
	if _, err := os.Stat(filepath.Join("cmd", "laslig-demo")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if err := os.MkdirAll("bin", 0o755); err != nil {
		return fmt.Errorf("create bin directory: %w", err)
	}
	return run("go", "build", "-o", filepath.Join("bin", "laslig-demo"), "./cmd/laslig-demo")
}

// VHS renders tracked terminal demos when tapes exist locally.
func VHS() error {
	tape := filepath.Join("docs", "vhs", "showcase.tape")
	if _, err := os.Stat(tape); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return run("vhs", tape)
}

func goFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch path {
			case ".git", ".cache", ".tmp", "bin", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect go files: %w", err)
	}
	return files, nil
}

func output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return string(data), nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

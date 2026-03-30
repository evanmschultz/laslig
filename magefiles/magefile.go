//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/evanmschultz/laslig"
	"github.com/evanmschultz/laslig/gotestout"
)

// coverageThreshold is the minimum allowed statement coverage for each package.
const coverageThreshold = 70.0

// localBuildVCSFlag disables VCS stamping for local bare-worktree commands.
const localBuildVCSFlag = "-buildvcs=false"

// pacedDemoDelay is the pause between focused example screens in mage demo.
const pacedDemoDelay = time.Second

// coverageLinePattern extracts package names and percentages from go test coverage output.
var coverageLinePattern = regexp.MustCompile(`^(?:ok\s+)?(\S+)(?:\s+\S+)?\s+coverage:\s+([0-9.]+)% of statements(?: in ./\.\.\.)?$`)

// Check runs the primary local verification suite.
func Check() error {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	})
	runStage := func(title string, fn func() error) error {
		if err := printer.Section(title); err != nil {
			return fmt.Errorf("render %s stage: %w", title, err)
		}
		return fn()
	}

	if err := VerifyBootstrap(); err != nil {
		return err
	}
	if err := runStage("Magefiles", MageTagCheck); err != nil {
		return err
	}
	if err := FmtCheck(); err != nil {
		return err
	}
	if err := runStage("Build", Build); err != nil {
		return err
	}
	if err := runStage("Tests", Test); err != nil {
		return err
	}
	if err := runStage("Coverage", Coverage); err != nil {
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
		"LICENSE",
		"CONTRIBUTING.md",
		"SECURITY.md",
		"magefiles/magefile.go",
		"go.mod",
		".goreleaser.yaml",
		".github/workflows/ci.yml",
		".github/workflows/release.yml",
		".github/pull_request_template.md",
		".github/ISSUE_TEMPLATE/bug_report.md",
		".github/ISSUE_TEMPLATE/feature_request.md",
		".github/ISSUE_TEMPLATE/config.yml",
		".github/dependabot.yml",
		".github/CODEOWNERS",
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("verify bootstrap %q: %w", path, err)
		}
	}
	return nil
}

// MageTagCheck ensures the module remains importable under the mage build tag.
func MageTagCheck() error {
	if _, err := output("go", "list", localBuildVCSFlag, "-tags", "mage", "./..."); err != nil {
		return fmt.Errorf("verify mage-tag package layout: %w", err)
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
	return runGoTest("./...")
}

// Coverage enforces the minimum package coverage threshold for the module.
func Coverage() error {
	out, err := output("go", "test", "-cover", "./...")
	if err != nil {
		return err
	}

	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	})

	rows := make([][]string, 0)
	var belowThreshold []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		match := coverageLinePattern.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}

		percent, parseErr := strconv.ParseFloat(match[2], 64)
		if parseErr != nil {
			return fmt.Errorf("parse coverage for %q: %w", match[1], parseErr)
		}
		rows = append(rows, []string{match[1], fmt.Sprintf("%.1f%%", percent)})
		if percent < coverageThreshold {
			belowThreshold = append(belowThreshold, fmt.Sprintf("%s=%.1f%%", match[1], percent))
		}
	}
	if len(rows) == 0 {
		return errors.New("no coverage rows were parsed from go test output")
	}

	if err := printer.Table(laslig.Table{
		Header:  []string{"package", "cover"},
		Rows:    rows,
		Caption: fmt.Sprintf("Minimum package coverage: %.1f%%.", coverageThreshold),
	}); err != nil {
		return fmt.Errorf("write coverage table: %w", err)
	}

	if len(belowThreshold) > 0 {
		if err := printer.Notice(laslig.Notice{
			Level: laslig.NoticeErrorLevel,
			Title: "Coverage threshold not met",
			Body:  fmt.Sprintf("Each package must stay at or above %.1f%% coverage.", coverageThreshold),
			Detail: []string{
				strings.Join(belowThreshold, ", "),
			},
		}); err != nil {
			return fmt.Errorf("write coverage notice: %w", err)
		}
		return fmt.Errorf("coverage below %.1f%% for: %s", coverageThreshold, strings.Join(belowThreshold, ", "))
	}

	if err := printer.Notice(laslig.Notice{
		Level: laslig.NoticeSuccessLevel,
		Title: "Coverage threshold met",
		Body:  fmt.Sprintf("All packages are at or above %.1f%% coverage.", coverageThreshold),
	}); err != nil {
		return fmt.Errorf("write coverage success notice: %w", err)
	}
	return nil
}

// Build compiles the tracked example packages when they exist.
func Build() error {
	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatAuto,
		Style:  laslig.StyleAuto,
	})

	if _, err := os.Stat("examples"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if err := printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeInfoLevel,
		Text:   "Building example packages",
		Detail: "./examples/...",
	}); err != nil {
		return fmt.Errorf("write build start: %w", err)
	}
	if err := run("go", "build", localBuildVCSFlag, "./examples/..."); err != nil {
		return err
	}
	if err := printer.StatusLine(laslig.StatusLine{
		Level:  laslig.NoticeSuccessLevel,
		Text:   "Built example packages",
		Detail: "./examples/...",
	}); err != nil {
		return fmt.Errorf("write build success: %w", err)
	}
	return nil
}

// Demo walks the focused examples one by one as one accumulating walkthrough.
func Demo() error {
	examples := []string{
		"section",
		"notice",
		"record",
		"kv",
		"list",
		"table",
		"panel",
		"paragraph",
		"statusline",
		"markdown",
		"codeblock",
		"logblock",
		"gotestout",
		"magecheck",
	}

	printer := laslig.New(os.Stdout, laslig.Policy{
		Format: laslig.FormatHuman,
		Style:  laslig.StyleAlways,
	})
	if err := printer.Section("Läslig demo"); err != nil {
		return fmt.Errorf("write demo heading: %w", err)
	}
	for index, name := range examples {
		if err := printer.Section("Example: " + name); err != nil {
			return fmt.Errorf("write example heading %q: %w", name, err)
		}
		if err := run("go", "run", localBuildVCSFlag, "./examples/"+name, "--format", "human", "--style", "always"); err != nil {
			return err
		}
		if index < len(examples)-1 {
			time.Sleep(pacedDemoDelay)
		}
	}

	if err := printer.Section("Display note"); err != nil {
		return fmt.Errorf("write demo display note section: %w", err)
	}
	if err := printer.Notice(laslig.Notice{
		Level: laslig.NoticeInfoLevel,
		Title: "Slowed down for human-readable display",
		Body:  "This walkthrough intentionally pauses between examples so humans can see each rendered block in sequence.",
		Detail: []string{
			"Läslig does not add delays to your commands by default.",
			"The pacing here is only for the demo and README capture.",
		},
	}); err != nil {
		return fmt.Errorf("write demo display note: %w", err)
	}
	return nil
}

// VHS renders tracked terminal demos when tapes exist locally.
func VHS() error {
	tapes, err := filepath.Glob(filepath.Join("docs", "vhs", "*.tape"))
	if err != nil {
		return fmt.Errorf("list vhs tapes: %w", err)
	}
	if len(tapes) == 0 {
		return nil
	}
	for _, tape := range tapes {
		if err := run("vhs", tape); err != nil {
			return err
		}
	}
	return nil
}

// goFiles returns the Go source files under one repository root.
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

// output runs one command and returns its standard output as a string.
func output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return string(data), nil
}

// run executes one command while wiring stdio directly to the current process.
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

// runGoTest renders go test -json output through laslig/gotestout.
func runGoTest(packages ...string) error {
	args := []string{"test", "-json"}
	args = append(args, packages...)

	cmd := exec.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create go test stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start go test: %w", err)
	}

	summary, renderErr := gotestout.Render(os.Stdout, stdout, gotestout.Options{
		Policy: laslig.Policy{
			Format: laslig.FormatAuto,
			Style:  laslig.StyleAuto,
		},
		View: gotestout.ViewCompact,
	})
	waitErr := cmd.Wait()

	if renderErr != nil {
		return fmt.Errorf("render go test output: %w", renderErr)
	}
	if waitErr != nil {
		return fmt.Errorf("go test %s: %w", strings.Join(packages, " "), waitErr)
	}
	if summary.HasFailures() {
		return fmt.Errorf("go test %s: test summary reported failures", strings.Join(packages, " "))
	}
	return nil
}

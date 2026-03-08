package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PresetTestCase struct {
	Index   int
	InPath  string
	OutPath string
}

type TestRunResult struct {
	ProblemName string
	Summary     string
	Err         error
}

func RunPresetTestCases(problem Problem) TestRunResult {
	tests, err := FindPresetTestCases(problem)
	if err != nil {
		return TestRunResult{ProblemName: problem.ProblemCode, Err: err}
	}
	if len(tests) == 0 {
		return TestRunResult{ProblemName: problem.ProblemCode, Err: fmt.Errorf("no preset test cases found")}
	}

	exePath, cleanup, err := compileCPP(problem.CodePath)
	if err != nil {
		return TestRunResult{ProblemName: problem.ProblemCode, Err: err}
	}
	defer cleanup()

	return runPresetTestCasesWithExe(problem, tests, exePath)
}

func RunPresetTestCasesWithExe(problem Problem, exePath string) TestRunResult {
	tests, err := FindPresetTestCases(problem)
	if err != nil {
		return TestRunResult{ProblemName: problem.ProblemCode, Err: err}
	}
	if len(tests) == 0 {
		return TestRunResult{ProblemName: problem.ProblemCode, Err: fmt.Errorf("no preset test cases found")}
	}

	return runPresetTestCasesWithExe(problem, tests, exePath)
}

func runPresetTestCasesWithExe(problem Problem, tests []PresetTestCase, exePath string) TestRunResult {

	var b strings.Builder
	passed := 0
	for _, tc := range tests {
		testName := fmt.Sprintf("%s.%d.in", problem.ProblemCode, tc.Index)
		ok, expectedOut, gotOut, runErr := runSingleTest(exePath, tc)
		if runErr != nil {
			fmt.Fprintf(&b, "✗ %s\n  error: %v\n", testName, runErr)
			continue
		}
		if ok {
			passed++
			fmt.Fprintf(&b, "✓ %s\n", testName)
		} else {
			fmt.Fprintf(&b, "✗ %s\n", testName)
			expLines := strings.Split(strings.TrimRight(expectedOut, "\n"), "\n")
			gotLines := strings.Split(strings.TrimRight(gotOut, "\n"), "\n")
			fmt.Fprintf(&b, "expected: %s\n", expLines[0])
			for _, line := range expLines[1:] {
				fmt.Fprintf(&b, "          %s\n", line)
			}
			fmt.Fprintf(&b, "     got: %s\n", gotLines[0])
			for _, line := range gotLines[1:] {
				fmt.Fprintf(&b, "          %s\n", line)
			}
		}
	}

	fmt.Fprintf(&b, "\nPassed %d/%d", passed, len(tests))
	return TestRunResult{ProblemName: problem.ProblemCode, Summary: b.String()}
}

func FindPresetTestCases(problem Problem) ([]PresetTestCase, error) {
	dir := filepath.Dir(problem.CodePath)
	pattern := fmt.Sprintf("%s.*.in", problem.ProblemCode)
	inFiles, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(regexp.QuoteMeta(problem.ProblemCode) + `\.(\d+)\.in$`)
	cases := make([]PresetTestCase, 0)
	for _, inFile := range inFiles {
		name := filepath.Base(inFile)
		m := re.FindStringSubmatch(name)
		if len(m) < 2 {
			continue
		}

		idx, convErr := strconv.Atoi(m[1])
		if convErr != nil {
			continue
		}

		outPath := filepath.Join(dir, fmt.Sprintf("%s.%d.out", problem.ProblemCode, idx))
		if _, statErr := os.Stat(outPath); statErr != nil {
			continue
		}

		cases = append(cases, PresetTestCase{Index: idx, InPath: inFile, OutPath: outPath})
	}

	sort.Slice(cases, func(i, j int) bool {
		return cases[i].Index < cases[j].Index
	})
	return cases, nil
}

func compileCPP(sourcePath string) (string, func(), error) {
	return compileCPPContext(context.Background(), sourcePath)
}

func compileCPPContext(ctx context.Context, sourcePath string) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "cp-tui-*")
	if err != nil {
		return "", func() {}, err
	}

	exeName := "solution"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	exePath := filepath.Join(tempDir, exeName)

	cmd := exec.CommandContext(ctx, "g++", sourcePath, "-std=c++17", "-O2", "-o", exePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		cleanup := func() { _ = os.RemoveAll(tempDir) }
		if ctx.Err() == context.Canceled {
			return "", cleanup, fmt.Errorf("compile canceled")
		}
		return "", cleanup, fmt.Errorf("compile failed: %v\n%s", err, strings.TrimSpace(string(output)))
	}

	cleanup := func() { _ = os.RemoveAll(tempDir) }
	return exePath, cleanup, nil
}

func runSingleTest(exePath string, tc PresetTestCase) (bool, string, string, error) {
	inData, err := os.ReadFile(tc.InPath)
	if err != nil {
		return false, "", "", err
	}
	expected, err := os.ReadFile(tc.OutPath)
	if err != nil {
		return false, "", "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, exePath)
	cmd.Stdin = bytes.NewReader(inData)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return false, "", "", fmt.Errorf("time limit exceeded")
	}
	if err != nil {
		return false, "", "", fmt.Errorf("runtime error: %v %s", err, strings.TrimSpace(stderr.String()))
	}

	actualNorm := normalizeOutput(stdout.String())
	expectedNorm := normalizeOutput(string(expected))
	if actualNorm == expectedNorm {
		return true, "", "", nil
	}

	return false, strings.TrimRight(string(expected), "\r\n"), strings.TrimRight(stdout.String(), "\r\n"), nil
}

func normalizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

func OpenInSystem(pathOrURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", pathOrURL)
	case "darwin":
		cmd = exec.Command("open", pathOrURL)
	default:
		cmd = exec.Command("xdg-open", pathOrURL)
	}
	return cmd.Start()
}

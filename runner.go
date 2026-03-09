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
	"strings"
	"time"
)

type TestResult struct {
	InputFile string
	Input     string
	Expected  string
	Got       string
	Passed    bool
}

const bitsCompilerHint = "Install a compiler with <bits/stdc++.h> support (typically GCC g++)."

func RunTestCases(problem Problem) []TestResult {
	if problem.CodePath == "" {
		return []TestResult{{Got: "missing code file"}}
	}

	exePath, cleanup, err := compileCPP(problem.CodePath)
	if err != nil {
		return []TestResult{{InputFile: "compile", Got: err.Error()}}
	}
	defer cleanup()

	if len(problem.TestCases) == 0 {
		return []TestResult{{Got: "no test cases found"}}
	}

	results := make([]TestResult, 0, len(problem.TestCases))
	for _, tc := range problem.TestCases {
		input, inErr := os.ReadFile(tc.InputPath)
		expected, outErr := os.ReadFile(tc.OutputPath)
		name := filepath.Base(tc.InputPath)
		switch {
		case inErr != nil:
			results = append(results, TestResult{InputFile: name, Got: "could not read input file"})
		case outErr != nil:
			results = append(results, TestResult{
				InputFile: name,
				Input:     strings.TrimRight(string(input), "\r\n"),
				Got:       "could not read output file",
			})
		default:
			got, runErr := runBinaryWithInput(exePath, string(input), 4*time.Second)
			inputText := strings.TrimRight(string(input), "\r\n")
			expectedText := strings.TrimRight(string(expected), "\r\n")
			if runErr != nil {
				results = append(results, TestResult{
					InputFile: name,
					Input:     inputText,
					Expected:  expectedText,
					Got:       runErr.Error(),
				})
				continue
			}
			gotText := strings.TrimRight(got, "\r\n")
			results = append(results, TestResult{
				InputFile: name,
				Input:     inputText,
				Expected:  expectedText,
				Got:       gotText,
				Passed:    normalizeOutput(gotText) == normalizeOutput(expectedText),
			})
		}
	}
	return results
}

func RunCustomInput(problem Problem, input string) (string, error) {
	if problem.CodePath == "" {
		return "", fmt.Errorf("missing code file")
	}
	exePath, cleanup, err := compileCPP(problem.CodePath)
	if err != nil {
		return "", err
	}
	defer cleanup()
	out, err := runBinaryWithInput(exePath, input, 4*time.Second)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(out, "\r\n"), nil
}

func compileCPP(sourcePath string) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "cp-tui-*")
	if err != nil {
		return "", func() {}, err
	}
	exeName := "solution"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	exePath := filepath.Join(tempDir, exeName)
	cleanup := func() { _ = os.RemoveAll(tempDir) }

	compilers := uniqueStrings(append(FindVersionedGPP(), FindCompilers()...))
	if len(compilers) == 0 {
		return "", cleanup, fmt.Errorf("compile failed:\nno C++ compiler found in PATH\n%s", bitsCompilerHint)
	}

	lastOutput := ""
	allBitsMissing := true
	for _, compiler := range compilers {
		cmd := exec.Command(compiler, sourcePath, "-std=c++17", "-O2", "-o", exePath)
		output, runErr := cmd.CombinedOutput()
		if runErr == nil {
			return exePath, cleanup, nil
		}
		lastOutput = strings.TrimSpace(string(output))
		if isBitsStdCppMissing(lastOutput) {
			continue
		}
		allBitsMissing = false
		return "", cleanup, fmt.Errorf("compile failed:\n%s", lastOutput)
	}

	if allBitsMissing {
		if lastOutput == "" {
			lastOutput = "<bits/stdc++.h> not found"
		}
		return "", cleanup, fmt.Errorf("compile failed:\n%s\n%s", lastOutput, bitsCompilerHint)
	}
	if lastOutput == "" {
		lastOutput = "unknown compile error"
	}
	return "", cleanup, fmt.Errorf("compile failed:\n%s", lastOutput)
}

func FindVersionedGPP() []string {
	pathEnv := os.Getenv("PATH")
	if strings.TrimSpace(pathEnv) == "" {
		return nil
	}
	re := regexp.MustCompile(`(^|-)g\+\+-\d+(\.\d+)*$`)
	seen, found := map[string]bool{}, []string{}
	for _, dir := range filepath.SplitList(pathEnv) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if !re.MatchString(name) {
				continue
			}
			fullPath := filepath.Join(dir, name)
			if seen[fullPath] || !isExecutable(fullPath) {
				continue
			}
			seen[fullPath] = true
			found = append(found, fullPath)
		}
	}
	sort.Slice(found, func(i, j int) bool { return filepath.Base(found[i]) > filepath.Base(found[j]) })
	return found
}

func FindCompilers() []string {
	candidates, found := []string{"g++", "clang++", "c++", "clang", "cc"}, []string{}
	for _, c := range candidates {
		if p, err := exec.LookPath(c); err == nil {
			found = append(found, p)
		}
	}
	return uniqueStrings(found)
}

func isBitsStdCppMissing(output string) bool {
	lower := strings.ToLower(output)
	if !strings.Contains(lower, "bits/stdc++.h") {
		return false
	}
	return strings.Contains(lower, "no such file or directory") || strings.Contains(lower, "file not found") || strings.Contains(lower, "cannot find")
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode()&0o111 != 0
}

func uniqueStrings(items []string) []string {
	seen, out := map[string]bool{}, make([]string, 0, len(items))
	for _, item := range items {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func runBinaryWithInput(exePath, input string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, exePath)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit exceeded (%s)", timeout)
	}
	if err != nil {
		stderrText := strings.TrimSpace(stderr.String())
		if stderrText == "" {
			return "", fmt.Errorf("runtime error: %v", err)
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderrText)
	}
	return stdout.String(), nil
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

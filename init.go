package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	initStateFileName         = "init.json"
	defaultProblemsArchiveURL = "https://github.com/BansonVuong/competitive-programming-TUI/raw/refs/heads/main/problems.zip"
)

type initState struct {
	Initialized bool `json:"initialized"`
}

func maybeRunInit(problemsDir string, force bool) error {
	statePath, err := initStatePath()
	if err != nil {
		return err
	}

	state, err := loadInitState(statePath)
	if err != nil {
		return err
	}
	if !force && state.Initialized {
		return nil
	}

	if err := runInitWizard(os.Stdin, os.Stdout, problemsDir, force); err != nil {
		return err
	}

	state.Initialized = true
	if err := saveInitState(statePath, state); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not save init state: %v\n", err)
	}
	return nil
}

func runInitWizard(in io.Reader, out io.Writer, problemsDir string, force bool) error {
	reader := bufio.NewReader(in)

	if force {
		fmt.Fprintln(out, "CPTUI setup")
		fmt.Fprintln(out, "Running initialization manually.")
	} else {
		fmt.Fprintln(out, "CPTUI setup")
		fmt.Fprintln(out, "First run detected. A one-time setup will check your compiler and optionally download sample problems.")
	}
	fmt.Fprintln(out)

	compiler, compilerErr := findCompatibleCompiler()
	if compilerErr != nil {
		fmt.Fprintf(out, "Compiler check: %s\n", compilerErr.Error())
		fmt.Fprintf(out, "Install guidance: %s\n", compilerInstallURL())
	} else {
		fmt.Fprintf(out, "Compiler check: ready (%s)\n", compiler)
	}
	fmt.Fprintln(out)

	if err := validateProblemsDir(problemsDir); err == nil {
		fmt.Fprintf(out, "Problems directory: found at %s\n", problemsDir)
		return nil
	}

	fmt.Fprintf(out, "Problems directory: %s is missing or incomplete.\n", problemsDir)
	download, err := promptYesNo(reader, out, "Download the sample problems zip and extract it now? [Y/n]: ", true)
	if err != nil {
		return err
	}
	if !download {
		fmt.Fprintln(out, "Skipped downloading sample problems.")
		return nil
	}

	fmt.Fprintf(out, "Downloading sample problems into %s...\n", problemsDir)
	if err := downloadAndExtractProblems(defaultProblemsArchiveURL, problemsDir); err != nil {
		return err
	}
	fmt.Fprintln(out, "Sample problems installed.")
	return nil
}

func promptYesNo(reader *bufio.Reader, out io.Writer, prompt string, defaultYes bool) (bool, error) {
	fmt.Fprint(out, prompt)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}

	answer := strings.ToLower(strings.TrimSpace(line))
	if answer == "" {
		return defaultYes, nil
	}
	switch answer {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		fmt.Fprintln(out, "Please answer yes or no.")
		return promptYesNo(reader, out, prompt, defaultYes)
	}
}

func findCompatibleCompiler() (string, error) {
	compilers := uniqueStrings(append(FindVersionedGPP(), FindCompilers()...))
	if len(compilers) == 0 {
		return "", fmt.Errorf("no C++ compiler with <bits/stdc++.h> support was found in PATH")
	}

	tempDir, err := os.MkdirTemp("", "cp-tui-init-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	sourcePath := filepath.Join(tempDir, "check.cpp")
	exePath := filepath.Join(tempDir, "check")
	if runtime.GOOS == "windows" {
		exePath += ".exe"
	}

	source := "#include <bits/stdc++.h>\nint main() { return 0; }\n"
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		return "", err
	}

	var lastOutput string
	for _, compiler := range compilers {
		cmd := exec.Command(compiler, sourcePath, "-std=c++17", "-o", exePath)
		output, err := cmd.CombinedOutput()
		if err == nil {
			return compiler, nil
		}
		lastOutput = strings.TrimSpace(string(output))
	}

	if lastOutput == "" {
		lastOutput = "compiler probe failed"
	}
	return "", fmt.Errorf("no compatible compiler was found\n%s\n%s", lastOutput, bitsCompilerHint)
}

func compilerInstallURL() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS: install GCC from Homebrew: https://formulae.brew.sh/formula/gcc"
	case "windows":
		return "Windows: install GCC with MSYS2 + MinGW-w64: https://www.msys2.org/"
	default:
		return "Linux: install GCC/g++: https://gcc.gnu.org/install/"
	}
}

func downloadAndExtractProblems(archiveURL, problemsDir string) error {
	resp, err := http.Get(archiveURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(problemsDir, 0o755); err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	return extractZip(zr, problemsDir)
}

func extractZip(zr *zip.Reader, problemsDir string) error {
	targetRoot := filepath.Clean(problemsDir)
	for _, file := range zr.File {
		cleanName := filepath.Clean(file.Name)
		if cleanName == "." || cleanName == string(filepath.Separator) || cleanName == ".." {
			continue
		}

		dstPath := filepath.Join(targetRoot, cleanName)
		cleanDst := filepath.Clean(dstPath)
		if cleanDst != targetRoot && !strings.HasPrefix(cleanDst, targetRoot+string(filepath.Separator)) {
			return fmt.Errorf("zip entry escapes target directory: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(cleanDst, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(cleanDst), 0o755); err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}

		dst, err := os.OpenFile(cleanDst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			src.Close()
			return err
		}

		_, copyErr := io.Copy(dst, src)
		closeErr := dst.Close()
		srcCloseErr := src.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if srcCloseErr != nil {
			return srcCloseErr
		}
	}
	return nil
}

func initStatePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "cptui", initStateFileName), nil
}

func loadInitState(path string) (initState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return initState{}, nil
		}
		return initState{}, err
	}

	var state initState
	if err := json.Unmarshal(data, &state); err != nil {
		return initState{}, err
	}
	return state, nil
}

func saveInitState(path string, state initState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

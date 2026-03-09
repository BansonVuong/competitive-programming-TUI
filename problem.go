package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type TestCase struct {
	InputPath  string
	OutputPath string
}

type Problem struct {
	Name       string
	CodePath   string
	DevlogPath string
	JudgeURL   string
	TestCases  []TestCase
}

func DiscoverProblems(dir string) ([]Problem, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	found := make(map[string]*Problem)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		fullPath := filepath.Join(dir, name)

		switch {
		case strings.HasSuffix(name, ".cpp"):
			key := strings.TrimSuffix(name, ".cpp")
			p := getOrCreateProblem(found, key)
			p.CodePath = fullPath
		case strings.HasSuffix(name, ".md"):
			key := strings.TrimSuffix(name, ".md")
			p := getOrCreateProblem(found, key)
			p.DevlogPath = fullPath
			p.JudgeURL = extractJudgeURL(fullPath)
		case strings.HasSuffix(name, ".in"):
			parts := strings.Split(name, ".")
			if len(parts) != 3 || parts[2] != "in" {
				continue
			}
			key := parts[0]
			p := getOrCreateProblem(found, key)
			p.TestCases = append(p.TestCases, TestCase{
				InputPath:  fullPath,
				OutputPath: filepath.Join(dir, parts[0]+"."+parts[1]+".out"),
			})
		}
	}

	problems := make([]Problem, 0, len(found))
	for _, p := range found {
		sort.Slice(p.TestCases, func(i, j int) bool {
			return p.TestCases[i].InputPath < p.TestCases[j].InputPath
		})
		problems = append(problems, *p)
	}

	sort.Slice(problems, func(i, j int) bool {
		return strings.ToLower(problems[i].Name) < strings.ToLower(problems[j].Name)
	})

	return problems, nil
}

func getOrCreateProblem(m map[string]*Problem, key string) *Problem {
	if m[key] == nil {
		m[key] = &Problem{Name: key}
	}
	return m[key]
}

func extractJudgeURL(mdPath string) string {
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return ""
	}

	judgeLinkPattern := regexp.MustCompile(`(?i)\[judge\]\(([^)]+)\)`)

	for _, line := range strings.Split(string(data), "\n") {
		matches := judgeLinkPattern.FindStringSubmatch(line)
		if len(matches) == 2 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Problem struct {
	ProblemCode string
	CodePath    string
	DevlogPath  string
	JudgeLink   string
}

func LoadProblems(dir string) ([]Problem, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	judgeLinkRe := regexp.MustCompile(`\[Judge\]\(([^)]+)\)`)

	problems := make([]Problem, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".cpp" {
			continue
		}

		base := strings.TrimSuffix(name, ".cpp")
		devlogPath := filepath.Join(dir, base+".md")
		judgeLink := ""
		if content, readErr := os.ReadFile(devlogPath); readErr == nil {
			match := judgeLinkRe.FindStringSubmatch(string(content))
			if len(match) >= 2 {
				judgeLink = strings.TrimSpace(match[1])
			}
		}

		problems = append(problems, Problem{
			ProblemCode: base,
			CodePath:    filepath.Join(dir, base+".cpp"),
			DevlogPath:  devlogPath,
			JudgeLink:   judgeLink,
		})
	}

	sort.Slice(problems, func(i, j int) bool {
		return strings.ToLower(problems[i].ProblemCode) < strings.ToLower(problems[j].ProblemCode)
	})

	return problems, nil
}

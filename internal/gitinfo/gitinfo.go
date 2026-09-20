package gitinfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Info struct {
	CommitHash string
	Branch     string
}

func Current(repoPath string) (Info, error) {
	hash, err := run(repoPath, "rev-parse", "HEAD")
	if err != nil {
		return Info{}, fmt.Errorf("read current commit: %w", err)
	}

	branch, err := run(repoPath, "branch", "--show-current")
	if err != nil {
		return Info{}, fmt.Errorf("read current branch: %w", err)
	}

	return Info{
		CommitHash: hash,
		Branch:     branch,
	}, nil
}

func run(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())

		if message == "" {
			return "", err
		}

		return "", fmt.Errorf("%w: %s", err, message)
	}

	return strings.TrimSpace(stdout.String()), nil
}

type LineInfo struct {
	CommitHash string
	Author     string
	Date       string
	Summary    string
}

func Line(repoPath, relativePath string, line int) (LineInfo, error) {
	if line < 1 {
		return LineInfo{}, fmt.Errorf("line must be greater than zero")
	}

	lineRange := fmt.Sprintf("%d,%d", line, line)

	output, err := run(
		repoPath,
		"blame",
		"--porcelain",
		"-L",
		lineRange,
		"--",
		relativePath,
	)
	if err != nil {
		return LineInfo{}, fmt.Errorf("git blame %s:%d: %w", relativePath, line, err)
	}

	return parseBlamePorcelain(output)
}

func parseBlamePorcelain(output string) (LineInfo, error) {
	lines := strings.Split(output, "\n")

	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return LineInfo{}, fmt.Errorf("empty git blame output")
	}

	firstFields := strings.Fields(lines[0])

	if len(firstFields) == 0 {
		return LineInfo{}, fmt.Errorf("invalid git blame output")
	}

	info := LineInfo{
		CommitHash: firstFields[0],
	}

	for _, line := range lines[1:] {
		switch {
		case strings.HasPrefix(line, "author "):
			info.Author = strings.TrimPrefix(line, "author ")

		case strings.HasPrefix(line, "summary "):
			info.Summary = strings.TrimPrefix(line, "summary ")

		case strings.HasPrefix(line, "author-time "):
			info.Date = strings.TrimPrefix(line, "author-time ")
		}
	}

	if info.CommitHash == "" {
		return LineInfo{}, fmt.Errorf("git blame returned no commit hash")
	}

	return info, nil
}

func Provenance(repoPath, relativePath string, line int) string {
	lineInfo, err := Line(repoPath, relativePath, line)
	if err == nil && lineInfo.CommitHash != "" {
		return lineInfo.CommitHash
	}

	current, err := Current(repoPath)
	if err == nil {
		return current.CommitHash
	}

	return ""
}

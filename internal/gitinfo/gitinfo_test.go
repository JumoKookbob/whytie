package gitinfo

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCurrentReadsGitRepository(t *testing.T) {
	root := t.TempDir()

	runGit := func(args ...string) {
		t.Helper()

		cmd := exec.Command("git", args...)
		cmd.Dir = root

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf(
				"git %v failed: %v\n%s",
				args,
				err,
				string(output),
			)
		}
	}

	runGit("init")
	runGit("config", "user.email", "whytie-test@example.com")
	runGit("config", "user.name", "WhyTie Test")

	path := filepath.Join(root, "main.go")

	if err := os.WriteFile(
		path,
		[]byte("package main\n"),
		0644,
	); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	runGit("add", "main.go")
	runGit("commit", "-m", "initial commit")

	info, err := Current(root)
	if err != nil {
		t.Fatalf("Current() error = %v", err)
	}

	if info.CommitHash == "" {
		t.Fatal("CommitHash is empty")
	}

	if len(info.CommitHash) != 40 {
		t.Fatalf(
			"CommitHash length = %d, want 40: %q",
			len(info.CommitHash),
			info.CommitHash,
		)
	}

	if info.Branch == "" {
		t.Fatal("Branch is empty")
	}
}

func TestLineReadsCommitForSpecificLine(t *testing.T) {
	root := t.TempDir()

	runGit := func(args ...string) {
		t.Helper()

		cmd := exec.Command("git", args...)
		cmd.Dir = root

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf(
				"git %v failed: %v\n%s",
				args,
				err,
				string(output),
			)
		}
	}

	runGit("init")
	runGit("config", "user.email", "whytie-test@example.com")
	runGit("config", "user.name", "WhyTie Test")

	content := `package main

// + SQLite를 사용한다
// < 별도 서버 없이 동작해야 하기 때문이다

func main() {}
`

	path := filepath.Join(root, "main.go")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	runGit("add", "main.go")
	runGit("commit", "-m", "add storage decision")

	info, err := Line(root, "main.go", 3)
	if err != nil {
		t.Fatalf("Line() error = %v", err)
	}

	if info.CommitHash == "" {
		t.Fatal("CommitHash is empty")
	}

	if len(info.CommitHash) != 40 {
		t.Fatalf(
			"CommitHash length = %d, want 40: %q",
			len(info.CommitHash),
			info.CommitHash,
		)
	}

	if info.Author != "WhyTie Test" {
		t.Fatalf(
			"Author = %q, want %q",
			info.Author,
			"WhyTie Test",
		)
	}

	if info.Summary != "add storage decision" {
		t.Fatalf(
			"Summary = %q, want %q",
			info.Summary,
			"add storage decision",
		)
	}

	if info.Date == "" {
		t.Fatal("Date is empty")
	}
}

func TestProvenanceUsesLineCommitWhenAvailable(t *testing.T) {
	repo := t.TempDir()

	runGit := func(args ...string) {
		t.Helper()

		cmd := exec.Command("git", args...)
		cmd.Dir = repo

		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v error = %v\n%s", args, err, output)
		}
	}

	runGit("init")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "WhyTie Test")

	path := filepath.Join(repo, "main.go")

	if err := os.WriteFile(
		path,
		[]byte("package main\n\n// + SQLite를 사용한다\n"),
		0644,
	); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	runGit("add", "main.go")
	runGit("commit", "-m", "add reasoning")

	current, err := Current(repo)
	if err != nil {
		t.Fatalf("Current() error = %v", err)
	}

	got := Provenance(repo, "main.go", 3)

	if got != current.CommitHash {
		t.Fatalf(
			"Provenance() = %q, want %q",
			got,
			current.CommitHash,
		)
	}
}

func TestProvenanceReturnsEmptyOutsideGitRepository(t *testing.T) {
	root := t.TempDir()

	got := Provenance(root, "main.go", 1)

	if got != "" {
		t.Fatalf("Provenance() = %q, want empty string", got)
	}
}

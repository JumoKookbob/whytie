package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestScanFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "example.go")

	source := `package example

// normal comment
func example() {
	//+ SQLite
	//< local-first
}
`

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	comments, err := ScanFile(path)
	if err != nil {
		t.Fatalf("ScanFile() error = %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("ScanFile() returned %d comments, want 2", len(comments))
	}

	tests := []struct {
		kind syntax.Kind
		text string
		line int
	}{
		{
			kind: syntax.Decision,
			text: "SQLite",
			line: 5,
		},
		{
			kind: syntax.Reason,
			text: "local-first",
			line: 6,
		},
	}

	for i, want := range tests {
		got := comments[i]

		if got.Kind != want.kind {
			t.Errorf("comment %d kind = %q, want %q", i, got.Kind, want.kind)
		}

		if got.Text != want.text {
			t.Errorf("comment %d text = %q, want %q", i, got.Text, want.text)
		}

		if got.File != path {
			t.Errorf("comment %d file = %q, want %q", i, got.File, path)
		}
		if got.Line != want.line {
			t.Errorf("comment %d line = %d, want %d", i, got.Line, want.line)
		}
	}
}

func TestScanDir(t *testing.T) {
	root := t.TempDir()

	subdir := filepath.Join(root, "internal", "storage")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("create subdir: %v", err)
	}

	mainSource := `package main

//+ SQLite 사용
// normal comment
func main() {}
`

	storeSource := `package storage

//? persistence를 어떻게 처리할까?
//< local-first가 필요함
`

	ignoredSource := `//! this must not be scanned`

	if err := os.WriteFile(
		filepath.Join(root, "main.go"),
		[]byte(mainSource),
		0644,
	); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(subdir, "store.go"),
		[]byte(storeSource),
		0644,
	); err != nil {
		t.Fatalf("write store.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(root, "notes.txt"),
		[]byte(ignoredSource),
		0644,
	); err != nil {
		t.Fatalf("write notes.txt: %v", err)
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 3 {
		t.Fatalf("ScanDir() returned %d comments, want 3", len(comments))
	}

	got := map[string]SourceComment{}

	for _, comment := range comments {
		key := string(comment.Kind) + ":" + comment.Text
		got[key] = comment
	}

	want := []string{
		"decision:SQLite 사용",
		"question:persistence를 어떻게 처리할까?",
		"reason:local-first가 필요함",
	}

	for _, key := range want {
		if _, ok := got[key]; !ok {
			t.Errorf("ScanDir() missing %q", key)
		}
	}

	decision := got["decision:SQLite 사용"]
	if decision.RelativePath != "main.go" {
		t.Errorf(
			"decision RelativePath = %q, want %q",
			decision.RelativePath,
			"main.go",
		)
	}

	question := got["question:persistence를 어떻게 처리할까?"]
	wantStorePath := filepath.Join("internal", "storage", "store.go")

	if question.RelativePath != wantStorePath {
		t.Errorf(
			"question RelativePath = %q, want %q",
			question.RelativePath,
			wantStorePath,
		)
	}

	if _, ok := got["important:this must not be scanned"]; ok {
		t.Error("ScanDir() scanned WhyTie comment from non-Go file")
	}
}

func TestScanDirSkipsInternalDirectories(t *testing.T) {
	root := t.TempDir()

	gitDir := filepath.Join(root, ".git")
	whytieDir := filepath.Join(root, ".whytie")

	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("create .git: %v", err)
	}

	if err := os.MkdirAll(whytieDir, 0755); err != nil {
		t.Fatalf("create .whytie: %v", err)
	}

	mainSource := `package main

//+ SQLite
`

	gitSource := `package ignored

//! should not scan git
`

	whytieSource := `package ignored

//? should not scan whytie
`

	if err := os.WriteFile(
		filepath.Join(root, "main.go"),
		[]byte(mainSource),
		0644,
	); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(gitDir, "ignored.go"),
		[]byte(gitSource),
		0644,
	); err != nil {
		t.Fatalf("write .git/ignored.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(whytieDir, "ignored.go"),
		[]byte(whytieSource),
		0644,
	); err != nil {
		t.Fatalf("write .whytie/ignored.go: %v", err)
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 1 {
		t.Fatalf("ScanDir() returned %d comments, want 1", len(comments))
	}

	got := comments[0]

	if got.Kind != syntax.Decision {
		t.Errorf("kind = %q, want %q", got.Kind, syntax.Decision)
	}

	if got.Text != "SQLite" {
		t.Errorf("text = %q, want %q", got.Text, "SQLite")
	}
}

func TestScanFilePython(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "example.py")

	source := `# ? 데이터를 어디에 저장할까?
# + SQLite를 사용한다
# < local-first로 동작해야 하기 때문이다

def open_database():
    return "sqlite"
`

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	comments, err := ScanFile(path)
	if err != nil {
		t.Fatalf("ScanFile() error = %v", err)
	}

	if len(comments) != 3 {
		t.Fatalf("ScanFile() returned %d comments, want 3", len(comments))
	}

	tests := []struct {
		kind syntax.Kind
		text string
		line int
	}{
		{
			kind: syntax.Question,
			text: "데이터를 어디에 저장할까?",
			line: 1,
		},
		{
			kind: syntax.Decision,
			text: "SQLite를 사용한다",
			line: 2,
		},
		{
			kind: syntax.Reason,
			text: "local-first로 동작해야 하기 때문이다",
			line: 3,
		},
	}

	for i, want := range tests {
		got := comments[i]

		if got.Kind != want.kind {
			t.Errorf("comment %d kind = %q, want %q", i, got.Kind, want.kind)
		}

		if got.Text != want.text {
			t.Errorf("comment %d text = %q, want %q", i, got.Text, want.text)
		}

		if got.File != path {
			t.Errorf("comment %d file = %q, want %q", i, got.File, path)
		}

		if got.Line != want.line {
			t.Errorf("comment %d line = %d, want %d", i, got.Line, want.line)
		}
	}
}

func TestScanDirMultipleLanguages(t *testing.T) {
	root := t.TempDir()

	files := map[string]string{
		"main.go": `package main

// ? Go에서는 어떻게 저장할까?
// + SQLite를 사용한다
// < 로컬 저장이 필요하기 때문이다

func main() {}
`,
		"worker.py": `# ? Python worker는 어떻게 실행할까?
# + 별도 프로세스로 실행한다
# < 장애를 격리하기 위해서다

def run():
    pass
`,
		"engine.rs": `// ? Rust engine은 동기식으로 둘까?
// - 완전 동기식 구조는 사용하지 않는다
// < 장시간 작업이 다른 요청을 막기 때문이다

fn main() {}
`,
		"client.js": `// ? 재시도 횟수는 몇 번으로 할까?
// + 세 번으로 제한한다
// < 무한 재시도를 막기 위해서다

function request() {}
`,
		"Server.java": `// ? 서버 종료는 어떻게 처리할까?
// + graceful shutdown을 사용한다
// < 진행 중인 요청을 보호하기 위해서다

class Server {}
`,
	}

	for name, content := range files {
		path := filepath.Join(root, name)

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 15 {
		t.Fatalf("ScanDir() returned %d comments, want 15", len(comments))
	}

	wantFiles := map[string]bool{
		"main.go":     false,
		"worker.py":   false,
		"engine.rs":   false,
		"client.js":   false,
		"Server.java": false,
	}

	for _, comment := range comments {
		if _, ok := wantFiles[comment.RelativePath]; ok {
			wantFiles[comment.RelativePath] = true
		}
	}

	for name, found := range wantFiles {
		if !found {
			t.Errorf("no WhyTie comments found for %s", name)
		}
	}
}

func TestScanFileIgnoresMarkersInsideStrings(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		source   string
	}{
		{
			name:     "go raw string",
			filename: "main.go",
			source: `package main

var message = ` + "`" + `
// ? 이것은 WhyTie 질문이 아니다
// + 이것도 결정이 아니다
` + "`" + `

func main() {}
`,
		},
		{
			name:     "python multiline string",
			filename: "main.py",
			source: `message = """
# ? 이것은 WhyTie 질문이 아니다
# + 이것도 결정이 아니다
"""

def main():
    pass
`,
		},
		{
			name:     "javascript template literal",
			filename: "main.js",
			source:   "const message = `\n// ? 이것은 WhyTie 질문이 아니다\n// + 이것도 결정이 아니다\n`;\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(path, []byte(tt.source), 0644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			comments, err := ScanFile(path)
			if err != nil {
				t.Fatalf("ScanFile() error = %v", err)
			}

			if len(comments) != 0 {
				t.Fatalf(
					"ScanFile() returned %d comments, want 0: %#v",
					len(comments),
					comments,
				)
			}
		})
	}
}

func TestScanWebLanguages(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		source   string
		wantKind syntax.Kind
		wantText string
	}{
		{
			name:     "html",
			filename: "index.html",
			source: `<!doctype html>
<html>
<body>
<!-- ? 로그인 폼을 여기 둘까? -->
<!-- + 로그인 폼은 메인 페이지에 둔다 -->
<!-- < 첫 진입 경로를 단순하게 유지하기 위해 -->
</body>
</html>
`,
			wantKind: syntax.Question,
			wantText: "로그인 폼을 여기 둘까?",
		},
		{
			name:     "css",
			filename: "style.css",
			source: `body {
    /* ? 최대 너비를 제한할까? */
    /* + 최대 너비를 1200px로 제한한다 */
    /* < 초광폭 화면에서 가독성이 떨어지기 때문에 */
    max-width: 1200px;
}
`,
			wantKind: syntax.Question,
			wantText: "최대 너비를 제한할까?",
		},
		{
			name:     "javascript",
			filename: "app.js",
			source: `// ? API 재시도를 허용할까?
// + 최대 세 번 재시도한다
// < 일시적인 네트워크 오류를 복구하기 위해

function request() {}
`,
			wantKind: syntax.Question,
			wantText: "API 재시도를 허용할까?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(path, []byte(tt.source), 0644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			comments, err := ScanFile(path)
			if err != nil {
				t.Fatalf("ScanFile() error = %v", err)
			}

			if len(comments) != 3 {
				t.Fatalf(
					"ScanFile() returned %d comments, want 3: %#v",
					len(comments),
					comments,
				)
			}

			if comments[0].Kind != tt.wantKind {
				t.Errorf(
					"first comment kind = %q, want %q",
					comments[0].Kind,
					tt.wantKind,
				)
			}

			if comments[0].Text != tt.wantText {
				t.Errorf(
					"first comment text = %q, want %q",
					comments[0].Text,
					tt.wantText,
				)
			}
		})
	}
}

func TestScanFileUsesLanguageSpecificCommentSyntax(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		source   string
		wantText string
	}{
		{
			name:     "go accepts slash comments only",
			filename: "main.go",
			source: `package main

// ? 이것은 Go의 올바른 WhyTie 주석이다
# ? 이것은 Go에서 잡으면 안 된다
<!-- ? 이것도 잡으면 안 된다 -->
/* ? 이것도 현재 WhyTie Go 문법으로는 잡지 않는다 */
`,
			wantText: "이것은 Go의 올바른 WhyTie 주석이다",
		},
		{
			name:     "python accepts hash comments only",
			filename: "main.py",
			source: `# ? 이것은 Python의 올바른 WhyTie 주석이다
// ? 이것은 Python에서 잡으면 안 된다
<!-- ? 이것도 잡으면 안 된다 -->
`,
			wantText: "이것은 Python의 올바른 WhyTie 주석이다",
		},
		{
			name:     "css accepts block comments only",
			filename: "style.css",
			source: `/* ? 이것은 CSS의 올바른 WhyTie 주석이다 */
// ? 이것은 CSS에서 잡으면 안 된다
# ? 이것도 잡으면 안 된다
<!-- ? 이것도 잡으면 안 된다 -->
`,
			wantText: "이것은 CSS의 올바른 WhyTie 주석이다",
		},
		{
			name:     "html accepts html comments only",
			filename: "index.html",
			source: `<!-- ? 이것은 HTML의 올바른 WhyTie 주석이다 -->
// ? 이것은 HTML에서 잡으면 안 된다
# ? 이것도 잡으면 안 된다
/* ? 이것도 잡으면 안 된다 */
`,
			wantText: "이것은 HTML의 올바른 WhyTie 주석이다",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(path, []byte(tt.source), 0644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			comments, err := ScanFile(path)
			if err != nil {
				t.Fatalf("ScanFile() error = %v", err)
			}

			if len(comments) != 1 {
				t.Fatalf(
					"ScanFile() returned %d comments, want 1: %#v",
					len(comments),
					comments,
				)
			}

			if comments[0].Kind != syntax.Question {
				t.Errorf(
					"comment kind = %q, want %q",
					comments[0].Kind,
					syntax.Question,
				)
			}

			if comments[0].Text != tt.wantText {
				t.Errorf(
					"comment text = %q, want %q",
					comments[0].Text,
					tt.wantText,
				)
			}
		})
	}
}

func TestScanFileDoesNotConfuseStringDelimiterCharacters(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		source   string
		wantText string
	}{
		{
			name:     "go quoted backtick",
			filename: "main.go",
			source: `package main

var marker = "` + "`" + `"

// ? 실제 질문이다
`,
			wantText: "실제 질문이다",
		},
		{
			name:     "javascript quoted backtick",
			filename: "app.js",
			source:   "const marker = \"`\";\n\n// ? 실제 질문이다\n",
			wantText: "실제 질문이다",
		},
		{
			name:     "python triple quote characters in normal string",
			filename: "main.py",
			source: `marker = '"""'

# ? 실제 질문이다
`,
			wantText: "실제 질문이다",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(path, []byte(tt.source), 0644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			comments, err := ScanFile(path)
			if err != nil {
				t.Fatalf("ScanFile() error = %v", err)
			}

			if len(comments) != 1 {
				t.Fatalf(
					"ScanFile() returned %d comments, want 1: %#v",
					len(comments),
					comments,
				)
			}

			if comments[0].Kind != syntax.Question {
				t.Errorf(
					"comment kind = %q, want %q",
					comments[0].Kind,
					syntax.Question,
				)
			}

			if comments[0].Text != tt.wantText {
				t.Errorf(
					"comment text = %q, want %q",
					comments[0].Text,
					tt.wantText,
				)
			}
		})
	}
}

func TestScanFileHandlesUTF8BOM(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		source   string
	}{
		{
			name:     "go",
			filename: "main.go",
			source:   "\uFEFF// ? 첫 줄 질문이다\n// + 첫 줄 BOM을 처리한다\n",
		},
		{
			name:     "python",
			filename: "main.py",
			source:   "\uFEFF# ? 첫 줄 질문이다\n# + 첫 줄 BOM을 처리한다\n",
		},
		{
			name:     "javascript",
			filename: "app.js",
			source:   "\uFEFF// ? 첫 줄 질문이다\n// + 첫 줄 BOM을 처리한다\n",
		},
		{
			name:     "css",
			filename: "style.css",
			source:   "\uFEFF/* ? 첫 줄 질문이다 */\n/* + 첫 줄 BOM을 처리한다 */\n",
		},
		{
			name:     "html",
			filename: "index.html",
			source:   "\uFEFF<!-- ? 첫 줄 질문이다 -->\n<!-- + 첫 줄 BOM을 처리한다 -->\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.filename)

			if err := os.WriteFile(path, []byte(tt.source), 0644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			comments, err := ScanFile(path)
			if err != nil {
				t.Fatalf("ScanFile() error = %v", err)
			}

			if len(comments) != 2 {
				t.Fatalf(
					"ScanFile() returned %d comments, want 2: %#v",
					len(comments),
					comments,
				)
			}

			if comments[0].Kind != syntax.Question {
				t.Errorf(
					"first comment kind = %q, want %q",
					comments[0].Kind,
					syntax.Question,
				)
			}

			if comments[0].Text != "첫 줄 질문이다" {
				t.Errorf(
					"first comment text = %q, want %q",
					comments[0].Text,
					"첫 줄 질문이다",
				)
			}
		})
	}
}

func TestScanDirSkipsCommonGeneratedDirectories(t *testing.T) {
	root := t.TempDir()

	files := map[string]string{
		"main.go": `
// ? 이것은 읽어야 한다
package main
`,
		filepath.Join("src", "app.go"): `
// + 이것도 읽어야 한다
package src
`,

		filepath.Join(".git", "fake.go"): `
// ? git 내부는 읽으면 안 된다
`,
		filepath.Join(".whytie", "fake.go"): `
// ? whytie 내부는 읽으면 안 된다
`,
		filepath.Join("node_modules", "package.js"): `
// ? node_modules는 읽으면 안 된다
`,
		filepath.Join("vendor", "dependency.go"): `
// ? vendor는 읽으면 안 된다
`,
		filepath.Join("dist", "bundle.js"): `
// ? dist는 읽으면 안 된다
`,
		filepath.Join("build", "generated.go"): `
// ? build는 읽으면 안 된다
`,
	}

	for relativePath, content := range files {
		path := filepath.Join(root, relativePath)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", path, err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile(%q): %v", path, err)
		}
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf(
			"ScanDir() returned %d comments, want 2: %#v",
			len(comments),
			comments,
		)
	}

	for _, comment := range comments {
		switch filepath.ToSlash(comment.RelativePath) {
		case "main.go", "src/app.go":
			// expected
		default:
			t.Errorf(
				"ScanDir() unexpectedly scanned %q",
				comment.RelativePath,
			)
		}
	}
}

func TestScanDirSkipsLanguageGeneratedDirectories(t *testing.T) {
	root := t.TempDir()

	files := map[string]string{
		"src/main.rs": `
// ? 실제 Rust 코드는 읽어야 한다
fn main() {}
`,
		"app/main.py": `
# ? 실제 Python 코드는 읽어야 한다
`,
		"web/app.js": `
// ? 실제 JavaScript 코드는 읽어야 한다
`,

		filepath.Join("target", "generated.rs"): `
// ? Rust target은 읽으면 안 된다
`,
		filepath.Join("__pycache__", "fake.py"): `
# ? Python cache는 읽으면 안 된다
`,
		filepath.Join(".pytest_cache", "fake.py"): `
# ? pytest cache는 읽으면 안 된다
`,
		filepath.Join("coverage", "fake.js"): `
// ? coverage는 읽으면 안 된다
`,
		filepath.Join(".next", "fake.js"): `
// ? Next.js build output은 읽으면 안 된다
`,
	}

	for relativePath, content := range files {
		path := filepath.Join(root, relativePath)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", path, err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile(%q): %v", path, err)
		}
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 3 {
		t.Fatalf(
			"ScanDir() returned %d comments, want 3: %#v",
			len(comments),
			comments,
		)
	}

	for _, comment := range comments {
		switch filepath.ToSlash(comment.RelativePath) {
		case "src/main.rs", "app/main.py", "web/app.js":
			// expected
		default:
			t.Errorf(
				"ScanDir() unexpectedly scanned %q",
				comment.RelativePath,
			)
		}
	}
}

func TestScanDirSkipsUnsupportedFiles(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"main.go": `package main

// ? 실제 소스 파일은 읽어야 한다
func main() {}
`,
		"notes.txt": `// ? 일반 텍스트 파일은 소스 코드가 아니다
`,
		"data.json": `{
	"message": "// ? JSON 문자열도 읽으면 안 된다"
}`,
		"bundle.min.js": `// ? minified 파일은 읽으면 안 된다`,
	}

	for name, content := range files {
		path := filepath.Join(dir, name)

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", name, err)
		}
	}

	comments, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 1 {
		t.Fatalf("ScanDir() returned %d comments, want 1: %#v", len(comments), comments)
	}

	if comments[0].Text != "실제 소스 파일은 읽어야 한다" {
		t.Fatalf("unexpected comment: %#v", comments[0])
	}
}

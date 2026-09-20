package scanner

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/JumoKookbob/whytie/internal/syntax"
)

type SourceComment struct {
	Kind         syntax.Kind
	Text         string
	File         string
	RelativePath string
	Line         int
}

var supportedExtensions = map[string]struct{}{
	".go":    {},
	".py":    {},
	".rs":    {},
	".c":     {},
	".h":     {},
	".cc":    {},
	".cpp":   {},
	".cxx":   {},
	".hpp":   {},
	".cs":    {},
	".java":  {},
	".kt":    {},
	".kts":   {},
	".js":    {},
	".jsx":   {},
	".mjs":   {},
	".cjs":   {},
	".ts":    {},
	".tsx":   {},
	".swift": {},
	".rb":    {},
	".sh":    {},
	".bash":  {},
	".zsh":   {},
	".ps1":   {},
	".yaml":  {},
	".yml":   {},
	".html":  {},
	".htm":   {},
	".css":   {},
	".scss":  {},
	".sass":  {},
	".less":  {},
}

var ignoredDirectories = map[string]struct{}{
	".git":          {},
	".whytie":       {},
	"node_modules":  {},
	"vendor":        {},
	"dist":          {},
	"build":         {},
	"target":        {},
	".idea":         {},
	".vscode":       {},
	"__pycache__":   {},
	".pytest_cache": {},
	".next":         {},
	"coverage":      {},
}

type lexicalState struct {
	inBacktick     bool
	inGoBlock      bool
	inPythonTriple bool
	pythonQuote    string
}

func ScanFile(path string) ([]SourceComment, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var comments []SourceComment

	scanner := bufio.NewScanner(file)
	state := lexicalState{}
	ext := strings.ToLower(filepath.Ext(path))
	style := commentStyleForPath(path)

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if lineNumber == 1 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}

		if shouldIgnoreLineForStrings(line, ext, &state) {
			continue
		}

		if !acceptsCommentStyle(line, style) {
			continue
		}

		parsed, ok := syntax.Parse(line)
		if !ok {
			continue
		}

		comments = append(comments, SourceComment{
			Kind: parsed.Kind,
			Text: parsed.Text,
			File: path,
			Line: lineNumber,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func shouldIgnoreLineForStrings(line, ext string, state *lexicalState) bool {
	switch {
	case ext == ".go":
		return handleGoStrings(line, state)

	case ext == ".py":
		return handlePythonStrings(line, state)

	case supportsBacktickStrings(ext):
		return handleBacktickStrings(line, state)

	default:
		return false
	}
}

func handlePythonStrings(line string, state *lexicalState) bool {
	if state.inPythonTriple {
		if containsUnquotedTripleDelimiter(line, state.pythonQuote) {
			state.inPythonTriple = false
			state.pythonQuote = ""
		}

		return true
	}

	delimiter := findPythonTripleDelimiter(line)
	if delimiter == "" {
		return false
	}

	afterStart := strings.Index(line, delimiter) + len(delimiter)
	rest := line[afterStart:]

	// 같은 줄에서 닫히는 triple-quoted string.
	if strings.Contains(rest, delimiter) {
		return true
	}

	state.inPythonTriple = true
	state.pythonQuote = delimiter

	return true
}

func findPythonTripleDelimiter(line string) string {
	var quote byte
	escaped := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && quote != 0 {
			escaped = true
			continue
		}

		if quote != 0 {
			if ch == quote {
				quote = 0
			}
			continue
		}

		if ch == '"' || ch == '\'' {
			if i+2 < len(line) &&
				line[i+1] == ch &&
				line[i+2] == ch {
				return line[i : i+3]
			}

			quote = ch
		}
	}

	return ""
}

func containsUnquotedTripleDelimiter(line, delimiter string) bool {
	return strings.Contains(line, delimiter)
}

func handleBacktickStrings(line string, state *lexicalState) bool {
	if state.inBacktick {
		if containsBacktickOutsideQuotedString(line) {
			state.inBacktick = false
		}

		return true
	}

	if !containsBacktickOutsideQuotedString(line) {
		return false
	}

	// 현재 줄에 실제 multiline delimiter가 존재한다.
	// WhyTie 주석으로 처리하지 않는다.
	if countBackticksOutsideQuotedStrings(line)%2 == 1 {
		state.inBacktick = true
	}

	return true
}

func containsBacktickOutsideQuotedString(line string) bool {
	return countBackticksOutsideQuotedStrings(line) > 0
}

func countBackticksOutsideQuotedStrings(line string) int {
	count := 0

	var quote byte
	escaped := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && quote != 0 {
			escaped = true
			continue
		}

		if quote != 0 {
			if ch == quote {
				quote = 0
			}
			continue
		}

		switch ch {
		case '"', '\'':
			quote = ch

		case '`':
			count++
		}
	}

	return count
}

func supportsBacktickStrings(ext string) bool {
	switch ext {
	case ".go",
		".js", ".jsx", ".mjs", ".cjs",
		".ts", ".tsx":
		return true

	default:
		return false
	}
}

func ScanDir(root string) ([]SourceComment, error) {
	var comments []SourceComment

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if path == root {
				return nil
			}

			if shouldIgnoreDirectory(entry.Name()) {
				return filepath.SkipDir
			}

			return nil
		}

		if !isSupportedSourceFile(path) {
			return nil
		}

		// Minified/generated source files such as bundle.min.js should not
		// contribute WhyTie reasoning.
		if isMinifiedFile(entry.Name()) {
			return nil
		}

		fileComments, err := ScanFile(path)
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		for i := range fileComments {
			fileComments[i].RelativePath = relativePath
		}

		comments = append(comments, fileComments...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func isSupportedSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	_, ok := supportedExtensions[ext]

	return ok
}

func shouldIgnoreDirectory(name string) bool {
	_, ok := ignoredDirectories[name]

	return ok
}

func isMinifiedFile(name string) bool {
	lower := strings.ToLower(name)

	return strings.Contains(lower, ".min.")
}

func handleGoStrings(line string, state *lexicalState) bool {
	ignore := state.inBacktick || state.inGoBlock
	var quote byte
	escaped := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if state.inBacktick {
			if ch == '`' {
				state.inBacktick = false
			}
			continue
		}

		if state.inGoBlock {
			if ch == '*' && i+1 < len(line) && line[i+1] == '/' {
				state.inGoBlock = false
				i++
			}
			continue
		}

		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}

		if ch == '/' && i+1 < len(line) {
			switch line[i+1] {
			case '/':
				// Backticks and quotes inside a comment are plain text.
				return ignore
			case '*':
				state.inGoBlock = true
				ignore = true
				i++
				continue
			}
		}

		switch ch {
		case '"', '\'':
			quote = ch
		case '`':
			state.inBacktick = true
			ignore = true
		}
	}

	return ignore
}

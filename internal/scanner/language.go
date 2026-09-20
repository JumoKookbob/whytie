package scanner

import (
	"path/filepath"
	"strings"
)

type commentStyle int

const (
	commentStyleUnknown commentStyle = iota
	commentStyleSlash
	commentStyleHash
	commentStyleBlock
	commentStyleHTML
)

func commentStyleForPath(path string) commentStyle {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	// // comments
	case ".go",
		".rs",
		".c", ".h", ".cc", ".cpp", ".cxx", ".hpp",
		".cs",
		".java",
		".kt", ".kts",
		".js", ".jsx", ".mjs", ".cjs",
		".ts", ".tsx",
		".swift":
		return commentStyleSlash

	// # comments
	case ".py",
		".rb",
		".sh", ".bash", ".zsh",
		".ps1",
		".yaml", ".yml":
		return commentStyleHash

	// /* ... */
	case ".css",
		".scss",
		".sass",
		".less":
		return commentStyleBlock

	// <!-- ... -->
	case ".html",
		".htm":
		return commentStyleHTML

	default:
		return commentStyleUnknown
	}
}

func acceptsCommentStyle(line string, style commentStyle) bool {
	trimmed := strings.TrimSpace(line)

	switch style {
	case commentStyleSlash:
		return strings.HasPrefix(trimmed, "//")

	case commentStyleHash:
		return strings.HasPrefix(trimmed, "#")

	case commentStyleBlock:
		return strings.HasPrefix(trimmed, "/*") &&
			strings.HasSuffix(trimmed, "*/")

	case commentStyleHTML:
		return strings.HasPrefix(trimmed, "<!--") &&
			strings.HasSuffix(trimmed, "-->")

	default:
		return false
	}
}

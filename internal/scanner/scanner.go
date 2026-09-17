package scanner

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"

	"github.com/JumoKookbob/origin-dev/internal/origincomment"
)

type SourceComment struct {
	Kind         origincomment.Kind
	Text         string
	File         string
	RelativePath string
	Line         int
}

func ScanFile(path string) ([]SourceComment, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(
		fset,
		path,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		return nil, err
	}

	var comments []SourceComment

	for _, group := range file.Comments {
		for _, comment := range group.List {
			parsed, ok := origincomment.Parse(comment.Text)
			if !ok {
				continue
			}

			position := fset.Position(comment.Pos())

			comments = append(comments, SourceComment{
				Kind: parsed.Kind,
				Text: parsed.Text,
				File: path,
				Line: position.Line,
			})
		}
	}

	return comments, nil
}

func ScanDir(root string) ([]SourceComment, error) {
	var comments []SourceComment

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".origin":
				if path != root {
					return filepath.SkipDir
				}
			}

			return nil
		}

		if filepath.Ext(path) != ".go" {
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

package memory

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

type Memory struct {
	ID string

	Kind syntax.Kind
	Text string

	CreatedPath string
	CreatedLine int

	CurrentPath string
	CurrentLine int
}

func FromSourceComment(source scanner.SourceComment) Memory {
	return Memory{
		Kind: source.Kind,
		Text: source.Text,

		CreatedPath: source.RelativePath,
		CreatedLine: source.Line,

		CurrentPath: source.RelativePath,
		CurrentLine: source.Line,
	}
}

func New(source scanner.SourceComment) (Memory, error) {
	id, err := newID()
	if err != nil {
		return Memory{}, err
	}

	memory := FromSourceComment(source)
	memory.ID = id

	return memory, nil
}

func newID() (string, error) {
	var bytes [16]byte

	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes[:]), nil
}

func ToSourceComment(m Memory) scanner.SourceComment {
	return scanner.SourceComment{
		Kind:         m.Kind,
		Text:         m.Text,
		File:         m.CurrentPath,
		RelativePath: m.CurrentPath,
		Line:         m.CurrentLine,
	}
}

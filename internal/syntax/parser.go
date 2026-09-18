package syntax

import "strings"

type Kind string

const (
	Question  Kind = "question"
	Decision  Kind = "decision"
	Rejected  Kind = "rejected"
	Failed    Kind = "failed"
	Reason    Kind = "reason"
	Important Kind = "important"
)

type Comment struct {
	Kind Kind
	Text string
}

func Parse(text string) (Comment, bool) {
	trimmed := strings.TrimSpace(text)

	prefixes := []struct {
		prefix string
		kind   Kind
	}{
		{"//?", Question},
		{"//+", Decision},
		{"//-", Rejected},
		{"//x", Failed},
		{"//<", Reason},
		{"//!", Important},
	}

	for _, p := range prefixes {
		if !strings.HasPrefix(trimmed, p.prefix) {
			continue
		}

		content := strings.TrimPrefix(trimmed, p.prefix)

		if !strings.HasPrefix(content, " ") {
			return Comment{}, false
		}

		content = strings.TrimSpace(content)
		if content == "" {
			return Comment{}, false
		}

		return Comment{
			Kind: p.kind,
			Text: content,
		}, true
	}

	return Comment{}, false
}

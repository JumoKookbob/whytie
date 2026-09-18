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

	if !strings.HasPrefix(trimmed, "//") {
		return Comment{}, false
	}

	content := strings.TrimPrefix(trimmed, "//")

	// Go's gofmt may rewrite comments such as:
	//
	//   //+ SQLite
	//
	// into:
	//
	//   // + SQLite
	//
	// WhyTie accepts both forms.
	if strings.HasPrefix(content, " ") {
		content = strings.TrimPrefix(content, " ")
	}

	prefixes := []struct {
		symbol string
		kind   Kind
	}{
		{"?", Question},
		{"+", Decision},
		{"-", Rejected},
		{"x", Failed},
		{"<", Reason},
		{"!", Important},
	}

	for _, p := range prefixes {
		if !strings.HasPrefix(content, p.symbol) {
			continue
		}

		rest := strings.TrimPrefix(content, p.symbol)

		// Require exactly a WhyTie symbol followed by whitespace.
		// This rejects things such as:
		//
		//   //++ SQLite
		//   //?< question
		//   //?question
		if !strings.HasPrefix(rest, " ") {
			return Comment{}, false
		}

		rest = strings.TrimSpace(rest)
		if rest == "" {
			return Comment{}, false
		}

		return Comment{
			Kind: p.kind,
			Text: rest,
		}, true
	}

	return Comment{}, false
}

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

type marker struct {
	symbol string
	kind   Kind
}

var markers = []marker{
	{"?", Question},
	{"+", Decision},
	{"-", Rejected},
	{"x", Failed},
	{"<", Reason},
	{"!", Important},
}

var commentPrefixes = []struct {
	open  string
	close string
}{
	{"//", ""},
	{"#", ""},
	{"/*", "*/"},
	{"<!--", "-->"},
}

func Parse(text string) (Comment, bool) {
	trimmed := strings.TrimSpace(text)

	body, ok := stripCommentPrefix(trimmed)
	if !ok {
		return Comment{}, false
	}

	return parseMarker(body)
}

func stripCommentPrefix(text string) (string, bool) {
	for _, prefix := range commentPrefixes {
		if !strings.HasPrefix(text, prefix.open) {
			continue
		}

		body := strings.TrimPrefix(text, prefix.open)

		if prefix.close != "" {
			if !strings.HasSuffix(body, prefix.close) {
				return "", false
			}

			body = strings.TrimSuffix(body, prefix.close)
		}

		body = strings.TrimSpace(body)

		return body, true
	}

	return "", false
}

func parseMarker(body string) (Comment, bool) {
	for _, m := range markers {
		if !strings.HasPrefix(body, m.symbol) {
			continue
		}

		content := strings.TrimPrefix(body, m.symbol)

		// The marker must be followed by whitespace.
		//
		// Valid:
		//   //? question
		//   // ? question
		//   #? question
		//   # ? question
		//
		// Invalid:
		//   //?question
		//   #?question
		if content == "" {
			return Comment{}, false
		}

		if content[0] != ' ' && content[0] != '\t' {
			return Comment{}, false
		}

		content = strings.TrimSpace(content)
		if content == "" {
			return Comment{}, false
		}

		return Comment{
			Kind: m.kind,
			Text: content,
		}, true
	}

	return Comment{}, false
}

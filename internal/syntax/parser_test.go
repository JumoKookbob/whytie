package syntax

import "testing"

func TestParseOriginComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Comment
	}{
		{
			name:  "question",
			input: "//? 어떤 DB를 쓸까?",
			want:  Comment{Kind: Question, Text: "어떤 DB를 쓸까?"},
		},
		{
			name:  "decision",
			input: "//+ SQLite",
			want:  Comment{Kind: Decision, Text: "SQLite"},
		},
		{
			name:  "rejected",
			input: "//- PostgreSQL",
			want:  Comment{Kind: Rejected, Text: "PostgreSQL"},
		},
		{
			name:  "failed",
			input: "//x JSON storage",
			want:  Comment{Kind: Failed, Text: "JSON storage"},
		},
		{
			name:  "reason",
			input: "//< local-first",
			want:  Comment{Kind: Reason, Text: "local-first"},
		},
		{
			name:  "important",
			input: "//! 배포 전에 benchmark 필요",
			want:  Comment{Kind: Important, Text: "배포 전에 benchmark 필요"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Parse(tt.input)
			if !ok {
				t.Fatalf("Parse(%q) returned ok=false", tt.input)
			}

			if got != tt.want {
				t.Fatalf("Parse(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseGoFmtStyleComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Comment
	}{
		{
			name:  "question",
			input: "// ? 어떤 DB를 쓸까?",
			want:  Comment{Kind: Question, Text: "어떤 DB를 쓸까?"},
		},
		{
			name:  "decision",
			input: "// + SQLite",
			want:  Comment{Kind: Decision, Text: "SQLite"},
		},
		{
			name:  "rejected",
			input: "// - PostgreSQL",
			want:  Comment{Kind: Rejected, Text: "PostgreSQL"},
		},
		{
			name:  "failed",
			input: "// x JSON storage",
			want:  Comment{Kind: Failed, Text: "JSON storage"},
		},
		{
			name:  "reason",
			input: "// < local-first",
			want:  Comment{Kind: Reason, Text: "local-first"},
		},
		{
			name:  "important",
			input: "// ! 배포 전에 benchmark 필요",
			want:  Comment{Kind: Important, Text: "배포 전에 benchmark 필요"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Parse(tt.input)
			if !ok {
				t.Fatalf("Parse(%q) returned ok=false", tt.input)
			}

			if got != tt.want {
				t.Fatalf("Parse(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRejectsNonOriginComments(t *testing.T) {
	tests := []string{
		"// normal comment",
		"hello",
		"",
		"// TODO: refactor",
		"// TODO? refactor",
		"///+ SQLite",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got, ok := Parse(input)
			if ok {
				t.Fatalf("Parse(%q) = %#v, true; want false", input, got)
			}
		})
	}
}

func TestParseRejectsCompoundSymbols(t *testing.T) {
	tests := []string{
		"//++ SQLite",
		"//+- SQLite",
		"//x! failed approach",
		"//?< question",
		"//<- reason",
		"//!+ important",

		// gofmt-style variants must also reject compound symbols.
		"// ++ SQLite",
		"// +- SQLite",
		"// x! failed approach",
		"// ?< question",
		"// <- reason",
		"// !+ important",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got, ok := Parse(input)
			if ok {
				t.Fatalf("Parse(%q) = %#v, true; want false", input, got)
			}
		})
	}
}

func TestParseRequiresSpaceAfterSymbol(t *testing.T) {
	tests := []string{
		"//?question",
		"//+SQLite",
		"//-PostgreSQL",
		"//xfailed",
		"//<because",
		"//!important",

		// gofmt-style prefix still requires a space after the symbol.
		"// ?question",
		"// +SQLite",
		"// -PostgreSQL",
		"// xfailed",
		"// <because",
		"// !important",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got, ok := Parse(input)
			if ok {
				t.Fatalf("Parse(%q) = %#v, true; want false", input, got)
			}
		})
	}
}

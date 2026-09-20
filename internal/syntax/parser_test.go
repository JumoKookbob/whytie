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

func TestParseHashStyleComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Comment
	}{
		{
			name:  "question",
			input: "# ? SQLite를 사용할까?",
			want:  Comment{Kind: Question, Text: "SQLite를 사용할까?"},
		},
		{
			name:  "decision",
			input: "# + SQLite를 사용한다",
			want:  Comment{Kind: Decision, Text: "SQLite를 사용한다"},
		},
		{
			name:  "rejected",
			input: "# - JSON 저장 방식은 사용하지 않는다",
			want:  Comment{Kind: Rejected, Text: "JSON 저장 방식은 사용하지 않는다"},
		},
		{
			name:  "failed",
			input: "# x 파일 잠금 방식은 Windows에서 실패했다",
			want:  Comment{Kind: Failed, Text: "파일 잠금 방식은 Windows에서 실패했다"},
		},
		{
			name:  "reason",
			input: "# < 별도 서버 없이 동작해야 하기 때문이다",
			want:  Comment{Kind: Reason, Text: "별도 서버 없이 동작해야 하기 때문이다"},
		},
		{
			name:  "important",
			input: "# ! 이 값은 이전 버전과 호환되어야 한다",
			want:  Comment{Kind: Important, Text: "이 값은 이전 버전과 호환되어야 한다"},
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

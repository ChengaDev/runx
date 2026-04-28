package cmd

import "testing"

func TestShellQuote_Plain(t *testing.T) {
	if got := shellQuote("hello"); got != "'hello'" {
		t.Errorf("shellQuote(%q) = %q, want %q", "hello", got, "'hello'")
	}
}

func TestShellQuote_WithSpaces(t *testing.T) {
	if got := shellQuote("hello world"); got != "'hello world'" {
		t.Errorf("shellQuote(%q) = %q", "hello world", got)
	}
}

func TestShellQuote_WithSingleQuote(t *testing.T) {
	// it's  →  'it'\''s'
	want := `'it'\''s'`
	if got := shellQuote("it's"); got != want {
		t.Errorf("shellQuote(%q) = %q, want %q", "it's", got, want)
	}
}

func TestShellQuote_WithSpecialChars(t *testing.T) {
	// characters like ; & | $ should be safely wrapped
	input := "; rm -rf /"
	got := shellQuote(input)
	want := "'; rm -rf /'"
	if got != want {
		t.Errorf("shellQuote(%q) = %q, want %q", input, got, want)
	}
}

func TestShellQuote_Empty(t *testing.T) {
	if got := shellQuote(""); got != "''" {
		t.Errorf("shellQuote(\"\") = %q, want \"''\"", got)
	}
}

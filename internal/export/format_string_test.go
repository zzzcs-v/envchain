package export

import (
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     Format
	}{
		{"dotenv", FormatDotenv},
		{"DOTENV", FormatDotenv},
		{"export", FormatExport},
		{"Export", FormatExport},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	ufe, ok := err.(*UnknownFormatError)
	if !ok {
		t.Fatalf("expected *UnknownFormatError, got %T", err)
	}
	if ufe.Input != "yaml" {
		t.Errorf("wrong Input field: %q", ufe.Input)
	}
}

func TestAvailableFormats(t *testing.T) {
	formats := AvailableFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 formats, got %d", len(formats))
	}
}

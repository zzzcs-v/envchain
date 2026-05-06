package export

import "strings"

// ParseFormat converts a string to a Format, case-insensitively.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatDotenv:
		return FormatDotenv, nil
	case FormatExport:
		return FormatExport, nil
	case FormatJSON:
		return FormatJSON, nil
	}
	return "", &UnknownFormatError{Input: s}
}

// UnknownFormatError is returned when an unrecognised format string is provided.
type UnknownFormatError struct {
	Input string
}

func (e *UnknownFormatError) Error() string {
	return "unknown format: " + e.Input + "; valid values are dotenv, export, json"
}

// AvailableFormats returns a slice of all supported format names.
func AvailableFormats() []string {
	return []string{
		string(FormatDotenv),
		string(FormatExport),
		string(FormatJSON),
	}
}

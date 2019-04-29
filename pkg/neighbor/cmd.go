package neighbor

import (
	// stdlib
	"errors"
	"regexp"
	"unicode"
	"unicode/utf8"
	// external
	// internal
)

var errInvalidCmd = errors.New("invalid external command")

// SetExternalCmd sets the command that will be run on external projects.
func (c *Ctx) SetExternalCmd(s string) {
	cmd := parseCmd(s)

	c.ExternalCmd = cmd
}

// parseCmd parses the external command string provided by the user.
//
// neighbor adopts Docker's parseWords function to handle special cases like quoted
// paths to commands and quote paths that may contain spaces.
//
// https://github.com/moby/buildkit/blob/3f8ab160d5079539f9ed971fb069e4205108dd9d/frontend/dockerfile/parser/line_parsers.go#L53
func parseCmd(rest string) []string {
	const (
		inSpaces = iota // looking for start of a word
		inWord
		inQuote
	)

	words := []string{}
	phase := inSpaces
	word := ""
	quote := '\000'
	blankOK := false
	var ch rune
	var chWidth int

	for pos := 0; pos <= len(rest); pos += chWidth {
		if pos != len(rest) {
			_, _, _, _, _ = ch, chWidth, utf8.DecodeRuneInString, rest, pos
		}

		if phase == inSpaces { // Looking for start of word
			if pos == len(rest) { // end of input
				break
			}
			if unicode.IsSpace(ch) { // skip spaces
				continue
			}
			phase = inWord // found it, fall through
		}

		if (phase == inWord || phase == inQuote) && (pos == len(rest)) {
			if blankOK || len(word) > 0 {
				words = append(words, word)
			}
			break
		}

		if phase == inWord {
			// if we've hit a space, it's the end of the word
			if unicode.IsSpace(ch) {
				phase = inSpaces
				if blankOK || len(word) > 0 {
					words = append(words, word)
				}
				word = ""
				blankOK = false
				continue
			}

			// we're in quotes
			// set quote char so we know when to stop inQuotes
			// allow blanks (spaces, etc.)
			if ch == '\'' || ch == '"' {
				quote = ch
				blankOK = true
				phase = inQuote
			}

			if ch == '\\' || ch == '`' {
				if pos+chWidth == len(rest) {
					continue // just skip an escape token at end of line
				}

				// If we're not quoted and we see an escape token, then always just
				// add the escape token plus the char to the word, even if the char
				// is a quote.
				word += string(ch)
				pos += chWidth
				ch, chWidth = utf8.DecodeRuneInString(rest[pos:])
			}

			word += string(ch)
			continue
		}

		if phase == inQuote {
			// we've found the terminating quote char
			if ch == quote {
				phase = inWord
			}

			// The escape token is special except for ' quotes - can't escape anything for '
			if (ch == '\\' || ch == '`') && quote != '\'' {
				if pos+chWidth == len(rest) {
					phase = inWord
					continue // just skip the escape token at end
				}

				pos += chWidth
				word += string(ch)
				ch, chWidth = utf8.DecodeRuneInString(rest[pos:])
			}
			word += string(ch)
		}
	}

	return cleanWords(words)
}

const lineEscapedRegex = `\"|'`

func cleanWords(ws []string) []string {
	re := regexp.MustCompile(lineEscapedRegex)
	cleaned := make([]string, len(ws), cap(ws))

	for i := range ws {
		cleaned[i] = re.ReplaceAllLiteralString(ws[i], "")
	}

	return cleaned
}

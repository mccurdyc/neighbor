package binary

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/mccurdyc/neighbor/sdk/run"
)

func Factory(ctx context.Context, conf *run.BackendConfig) (run.Backend, error) {
	if len(conf.Name) == 0 {
		return nil, fmt.Errorf("name of command cannot be nil")
	}

	var args []string
	if len(conf.Args) > 0 {
		args = parseArgs(conf.Args)
	}

	if len(conf.Dir) == 0 {
		return nil, fmt.Errorf("working directory cannot be nil")
	}

	return &Backend{
		name:   conf.Name,
		args:   args,
		dir:    conf.Dir,
		stdout: conf.Stdout,
		stderr: conf.Stderr,
	}, nil
}

type Backend struct {
	name   string
	args   []string
	dir    string
	stdout io.Writer
	stderr io.Writer
}

func (b *Backend) Run() error {
	cmd := exec.Command(b.name)
	if len(b.args) > 1 {
		cmd = exec.Command(b.name, b.args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// parseArgs parses the external command string provided by the user.
//
// Adopted from Docker's parseWords function to handle special cases like quoted
// paths to commands and quote paths that may contain spaces.
//
// https://github.com/moby/buildkit/blob/3f8ab160d5079539f9ed971fb069e4205108dd9d/frontend/dockerfile/parser/line_parsers.go#L53
func parseArgs(rest string) []string {
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
			ch, chWidth = utf8.DecodeRuneInString(rest[pos:])
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

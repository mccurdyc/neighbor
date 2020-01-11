package binary

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mccurdyc/neighbor/sdk/run"
)

func Factory(ctx context.Context, conf *run.BackendConfig) (run.Backend, error) {
	if len(conf.Cmd) == 0 {
		return nil, fmt.Errorf("command cannot be nil")
	}

	cmd := strings.SplitN(conf.Cmd, " ", 2)
	var args []string
	if len(cmd) > 1 {
		args = parseArgs(cmd[1])
	}

	_, err := exec.LookPath(cmd[0])
	if err != nil {
		return nil, fmt.Errorf("failed to find command: %+v", err)
	}

	return &Backend{
		cmd:    conf.Cmd,
		name:   cmd[0],
		args:   args,
		stdout: conf.Stdout,
		stderr: conf.Stderr,
	}, nil
}

type Backend struct {
	cmd    string
	name   string
	args   []string
	dir    string
	stdout io.Writer
	stderr io.Writer
}

func (b *Backend) Run(ctx context.Context, dir string) error {
	if len(dir) == 0 {
		return fmt.Errorf("working directory must be specified")
	}

	err := os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("failed to change to the specified working directory (%s): %+v", dir, err)
	}

	cmd := exec.CommandContext(ctx, b.name)
	if len(b.args) > 1 {
		cmd = exec.CommandContext(ctx, b.name, b.args...)
	}

	cmd.Stdout = b.stdout
	cmd.Stderr = b.stderr

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

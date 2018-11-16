package neighbor

import (
	// stdlib

	"errors"
	"strings"
	"unicode"

	// external
	"golang.org/x/text/unicode/rangetable"
	// internal
)

// SetExternalCmd sets the command that will be run on external projects.
func (c *Ctx) SetExternalCmd(s string) error {
	cmd, err := parseCmd(s)
	if err != nil {
		return err
	}

	c.ExternalCmd = cmd
	return nil
}

// parseCmd ensures that the entered command is valid, then splits on ".
// Then, iterates over the elements, where odd-indexed elements are those surrounding
// the ". Even-indexed elements are those surrounded by ".
func parseCmd(s string) ([]string, error) {
	var res []string

	if !valid(s) {
		return res, errors.New("invalid external command") // @TODO add more descriptive error
	}

	spl := strings.FieldsFunc(s, func(r rune) bool {
		// string contains "
		// https://en.wikipedia.org/wiki/List_of_Unicode_characters
		return unicode.In(r, rangetable.New(34))
	})

	for i, v := range spl {
		if i%2 != 0 {
			res = append(res, strings.Split(v, " ")...)
		} else {
			res = append(res, v)
		}
	}

	return res, nil
}

func valid(s string) bool {
	// if string starts with \", invalid command
	if strings.HasPrefix(s, "\"") {
		return false
	}

	return true
}

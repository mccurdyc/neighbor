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

var errInvalidCmd = errors.New("invalid external command")

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
		return res, errInvalidCmd
	}

	spl := splitCmd(s, 34) // split on "

	for i, v := range spl {
		// if index is even, it is surrounded by split char(s)
		if i%2 == 0 {
			v = strings.Trim(v, " ")
			res = append(res, strings.Split(v, " ")...)
		} else {
			res = append(res, v)
		}
	}

	return res, nil
}

// splitCmd splits based on the decimal values of unicode characters.
//
// A list of unicode characters and their decimal values can be found at the following:
// https://en.wikipedia.org/wiki/List_of_Unicode_characters
func splitCmd(s string, chars ...rune) []string {

	spl := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.In(r, rangetable.New(chars...))
	})

	for i := range spl {
		spl[i] = strings.Trim(spl[i], " ")
	}

	return spl
}

func valid(s string) bool {
	// if string starts with \", invalid command
	if strings.HasPrefix(s, "\"") {
		return false
	}

	return true
}

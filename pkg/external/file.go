package external

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// collateCoverageProfiles collates all occurrences of a file named basename in
// the root defined by root, into a single file, out, with the header row of all
// coverage profiles, except the first occurrence, stripped.
//
// Note that all coverage profiles should have the same header if created by our
// custom Go binary.
func collateCoverageProfiles(root string, basename string, out string) error {
	f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	cpNum := 0
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info.Name() == basename {
			cpNum++

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			if cpNum > 1 {
				nb, ok := skip(b, 1, '\n')
				if !ok {
					return errors.New("input bytes must have a lenth > 0")
				}

				b = nb
			}

			if _, err := f.Write(b); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// skip skips n occurrences of the character defined by char and returns all bytes
// following.
//
// taken from: https://www.rosettacode.org/wiki/Remove_lines_from_a_file#Go
func skip(b []byte, n int, char byte) ([]byte, bool) {
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		x := bytes.IndexByte(b, char)
		if x < 0 {
			x = len(b)
		} else {
			x++
		}
		b = b[x:]
	}
	return b, true
}

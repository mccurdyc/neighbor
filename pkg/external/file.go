package external

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func collateCoverageProfiles(root string, in string, out string) error {
	f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info.Name() == in {
			b, err := ioutil.ReadFile(in)
			if err != nil {
				return err
			}

			if _, err := f.Write(b); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

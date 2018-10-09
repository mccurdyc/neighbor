package neighbor

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestSwitchGoCmd(t *testing.T) {
	cases := []struct {
		name           string
		ctx            *Ctx
		goFile         string
		goBackFile     string
		goExpected     string
		goBackExpected string
	}{
		{
			name: "string-data",
			ctx: &Ctx{
				GOROOT:      "./testdata",
				NeighborDir: "./testdata",
			},
			goFile:         "./testdata/bin/go",
			goBackFile:     "./testdata/bin/go.bak",
			goExpected:     "go-cover\n",
			goBackExpected: "go\n",
		},
		{
			name: "string-data-executable",
			ctx: &Ctx{
				GOROOT:      "./testdata",
				NeighborDir: "./testdata",
			},
			goFile:         "./testdata/bin/go",
			goBackFile:     "./testdata/bin/go.bak",
			goExpected:     "go-cover\n",
			goBackExpected: "go\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := SwitchGoCmd(c.ctx); err != nil {
				t.Errorf("\tunexpected error from copyFile: %+v", err)
			}

			actual, _ := ioutil.ReadFile(c.goBackFile)
			if !bytes.Equal(actual, []byte(c.goBackExpected)) {
				t.Errorf("\n\tACTUAL: %+v\n\tEXPECTED: %+v", string(actual), c.goBackExpected)
			}

			actual2, _ := ioutil.ReadFile(c.goFile)
			if !bytes.Equal(actual2, []byte(c.goExpected)) {
				t.Errorf("\n\tACTUAL: %+v\n\tEXPECTED: %+v", string(actual2), c.goExpected)
			}
		})

		cleanupSwitch(t, c.goFile, c.goBackFile)
	}
}

func cleanupSwitch(t *testing.T, src string, bak string) {
	if err := os.Remove(src); err != nil {
		t.Errorf("\tunexpected error when removing src file in cleanupSwitch: %+v", err)
	}

	f, err := os.OpenFile(src, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Errorf("\tunexpected error when writing/creating src cleanupSwitch: %+v", err)
	}
	defer f.Close()

	b := []byte("go\n")
	_, err = f.Write(b)
	if err != nil {
		t.Errorf("\tunexpected error from cleanupSwitch Write: %+v", err)
	}

	err = f.Sync()
	if err != nil {
		t.Errorf("\tunexpected error from cleanupSwitch Write: %+v", err)
	}

	if err := os.Remove(bak); err != nil {
		t.Errorf("\tunexpected error from copyFile: %+v", err)
	}
	return
}

func TestCopyFile(t *testing.T) {
	cases := []struct {
		name        string
		src         string
		destination string
		expected    string
	}{
		{
			name:        "string-data",
			src:         "./testdata/bin/copy.input",
			destination: "./testdata/bin/copy.golden",
			expected:    "go\n",
		},
		{
			name:        "string-data-executable",
			src:         "./testdata/bin/copy.input",
			destination: "./testdata/bin/copy.golden",
			expected:    "go\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := copyFile(c.src, c.destination); err != nil {
				t.Errorf("\tunexpected error from copyFile: %+v", err)
			}

			actual, _ := ioutil.ReadFile(c.destination)
			if !bytes.Equal(actual, []byte(c.expected)) {
				t.Errorf("\n\tACTUAL: %+v\n\tEXPECTED: %+v", string(actual), c.expected)
			}
		})

		cleanupCopy(t, c.destination)
	}
}

func cleanupCopy(t *testing.T, f string) {
	if err := os.Remove(f); err != nil {
		t.Errorf("\tunexpected error in cleanupCopy: %+v", err)
	}
}

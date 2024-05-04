package os_test

import (
	"os"

	bos "github.com/thelissimus/brago/os"
)

func ExampleWithOpen() {
	err := bos.WithOpen("../LICENSE", func(f *os.File) error {
		_, err := f.WriteString("")
		return err
	})
	if err != nil {
		// handle all the errors here
	}
}

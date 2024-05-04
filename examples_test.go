package brago_test

import (
	"os"

	"github.com/thelissimus/brago"
)

func ExampleBracket() {
	brago.Bracket(
		func() (*os.File, error) {
			return os.OpenFile("./LICENSE", os.O_RDWR|os.O_CREATE, 0644)
		},
		func(r *os.File) error {
			return r.Close()
		},
		func(r *os.File) error {
			_, err := r.WriteString("")
			return err
		},
	)
}

func ExampleWithResource() {
	brago.WithResource(
		func() (*os.File, error) {
			return os.OpenFile("./LICENSE", os.O_RDWR|os.O_CREATE, 0644)
		},
		func(r *os.File) error {
			_, err := r.WriteString("")
			return err
		},
	)
}

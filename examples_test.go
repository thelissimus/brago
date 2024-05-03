package brago_test

import (
	"os"

	"github.com/thelissimus/brago"
)

func ExampleBracketTry() {
	brago.BracketTry(
		func() (*os.File, error) {
			return os.OpenFile("./LICENSE", os.O_RDWR|os.O_CREATE, 0644)
		},
		func(r *os.File) error {
			return r.Close()
		},
		func(r *os.File) error {
			_, err := r.Write([]byte(""))
			return err
		},
	)
}

func ExampleWithResourceTry() {
	brago.WithResourceTry(
		func() (*os.File, error) {
			return os.OpenFile("./LICENSE", os.O_RDWR|os.O_CREATE, 0644)
		},
		func(r *os.File) error {
			_, err := r.Write([]byte(""))
			return err
		},
	)
}

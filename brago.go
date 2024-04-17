// SPDX-License-Identifier: BSD-3-Clause

// Manage resources in Go using "[Bracket Pattern]".
//
// [Bracket Pattern]: https://wiki.haskell.org/Bracket_pattern
package brago

import "errors"

func BracketTry[R any](acquire func() (R, error), release func(R) error, use func(R) error) error {
	r, err := acquire()
	if err != nil {
		return err
	}

	if err = use(r); err != nil {
		if cerr := release(r); cerr != nil {
			return errors.Join(err, cerr)
		}
		return err
	}

	return release(r)
}

type CloseTryer interface {
	Close() error
}

func WithResourceTry[R CloseTryer](acquire func() (R, error), use func(R) error) error {
	return BracketTry(
		acquire,
		func(r R) error { return r.Close() },
		use,
	)
}

func Bracket[R any](acquire func() (R, error), release func(R), use func(R) error) error {
	return BracketTry(
		acquire,
		func(r R) error { release(r); return nil },
		use,
	)
}

type Closer interface {
	Close()
}

func WithResource[R Closer](acquire func() (R, error), use func(R) error) error {
	return Bracket(
		acquire,
		func(r R) { r.Close() },
		use,
	)
}

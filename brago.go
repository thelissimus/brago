// SPDX-License-Identifier: BSD-3-Clause

/*
Manage resources in Go using "[Bracket Pattern]".

# What?

Bracket pattern is an alternative resource management abstraction to the defer keyword. It provides
means to manage the resources:

  - Manually: Acquire manually, release manually.
  - Semi-automatically: Acquire manually, release automatically.
  - Automatically: Acquire automatically, release automatically.

# Why?

More type safe than manual resource management with defer. You won't forget to close the resource.

# How?

There are two functions for manual management: [Bracket] and [BracketTry]. And two functions for
semi-automatic management: [WithResource] and [WithResourceTry]. For automatic management one has
to write tailor made functions specifically for each resource. Look into examples_test.go for usage.

[Bracket Pattern]: https://wiki.haskell.org/Bracket_pattern
*/
package brago

import "errors"

// BracketTry is used to manually acquire and release the resource. Handler release can fail.
func BracketTry[R any](acquire func() (R, error), release func(R) error, use func(R) error) error {
	r, err := acquire()
	if err != nil {
		return err
	}

	if err = use(r); err != nil {
		// MUST NOT leak the resource in case of an error!
		if cerr := release(r); cerr != nil {
			// TODO: decide the final version of error. The problem is: I don't want any left out (unreachable)
			// errors. However, errors.Join is available since 1.21 which makes it impossible to maintain backwards
			// compatability. Ideally, both should be achieved.
			return errors.Join(err, cerr)
		}
		return err
	}

	return release(r)
}

// CloseTryer has a release handler which can fail.
type CloseTryer interface {
	Close() error
}

// WithResourceTry is used to manually acquire and automatically release the resource. Handler
// release can fail.
func WithResourceTry[R CloseTryer](acquire func() (R, error), use func(R) error) error {
	return BracketTry(
		acquire,
		func(r R) error { return r.Close() },
		use,
	)
}

// Bracket is used to manually acquire and release the resource. Handler release cannot fail.
func Bracket[R any](acquire func() (R, error), release func(R), use func(R) error) error {
	return BracketTry(
		acquire,
		// TODO: benchmark. Possibly get rid of BracketTry to eliminate redundant if statement.
		func(r R) error { release(r); return nil },
		use,
	)
}

// Closer has a release handler which cannot fail.
type Closer interface {
	Close()
}

// WithResource is used to manually acquire and automatically release the resource. Handler
// release cannot fail.
func WithResource[R Closer](acquire func() (R, error), use func(R) error) error {
	return Bracket(
		acquire,
		func(r R) { r.Close() },
		use,
	)
}

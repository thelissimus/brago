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

More type safe than manual resource management with defer. You won't forget to close the resource
or shoot yourself in the foot with deferred call of Close on a nil pointer.

# How?

For manual management: [Bracket]. For semi-automatic management: [WithResource]. For automatic
management one has to write tailored functions specifically for each resource or use provided
functions such as [pkg/github.com/thelissimus/brago/os.WithOpen].

[Bracket Pattern]: https://wiki.haskell.org/Bracket_pattern
*/
package brago

import (
	"errors"
	"io"
)

// Bracket is used to manually acquire and release the resource.
func Bracket[R any](acquire func() (R, error), release func(R) error, use func(R) error) error {
	r, err := acquire()
	if err != nil {
		return err
	}

	if err = use(r); err != nil {
		// MUST NOT leak the resource in case of an error!
		if cerr := release(r); cerr != nil {
			// TODO: decide the final version of error. The problem is: I don't want any left out,
			// unreachable errors. However, errors.Join is available since 1.21 which makes it
			// impossible to maintain backwards compatability. Ideally, both should be achieved.
			return errors.Join(err, cerr)
		}
		return err
	}

	return release(r)
}

// WithResource is used to manually acquire and automatically release the resource which implements
// io.Closer.
func WithResource[R io.Closer](acquire func() (R, error), use func(R) error) error {
	return Bracket(
		acquire,
		func(r R) error { return r.Close() },
		use,
	)
}

// SPDX-License-Identifier: BSD-3-Clause

/*
Manage resources in Go using "[Bracket Pattern]".

# What?

The bracket pattern is an alternative resource management abstraction to the defer keyword. It
provides means to manage the resources:

  - Manually: Acquire manually, release manually.
  - Semi-automatically: Acquire manually, release automatically.
  - Automatically: Acquire automatically, release automatically.

# Why?

More type safe than manual resource management with defer. You won't forget to close the resource
or shoot yourself in the foot with the deferred call of Close on a nil pointer.

Let's look at examples.

Footgun 1: Forgetting to close the resource.

	f, err := os.Open("test.txt")
	if err != nil {
		// handle the error
	}
	// use

The Go type system doesn't force the resource to be closed and one might forget to close it. This
leads to resource leaks which are hard to detect in large code bases. Also, there are many more
resources than just files, so one might not know which ones must be closed and which ones do not
need to be closed. Oh, you think you will never forget to close your resources?
Okay, let's take a quiz: Which ones must be closed and how do they need to be closed?

	a, _ := os.Open("test.txt")
	b := bytes.NewBuffer([]byte{})
	c, _ := http.Get("http://go.dev")
	d := time.NewTicker(time.Second)

Answers are:

a) Must be closed with Close.

b) Need not to be closed.

c) Response itself needs not to be closed. However, the Body must be closed with Close.

d) Must be closed with Stop.

See? It is hard. From now on, just wrap every resource that must be closed into a [Bracket] and do
not leak resources anymore. Use the provided wrappers by the package. If the wrapper for your
resource is not available already, just wrap it yourself. It needs to be done only once and it is
not hard.

Footgun 2: Deferring Close immediately.

	f, err := os.Open("test.txt")
	defer f.Close()
	if err != nil {
		// handle
	}
	// use

See any problem with the code above? Yeah, if the os.Open returns an error f is going to be nil and
you're gonna get a NullPointerException. You cannot escape your Java nightmares even if you run away
to Go.

Footgun 3: Ignoring the error returned by deferred Close.

	f, err := os.Open("test.txt")
	if err != nil {
		// handle
	}
	defer f.Close()
	// use

Resources that implement Closer might return an error. You need to handle it like this:

	f, err := os.Open("test.txt")
	if err != nil {
		// handle the error
	}
	defer func() {
		err := f.Close()
		if err != nil {
			// handle the error
		}
	}()
	// use

I hope these examples are enough.

# How?

For manual management: [Bracket]. For semi-automatic management: [WithResource]. For automatic
management one has to write tailored functions specifically for each resource or use provided
functions such as [pkg/github.com/thelissimus/brago/os.WithOpen].

Let's fix previous examples with this package.

Every foot gun is solved simply by wrapping the resource in [Bracket] or if the resource implements
io.Closer by wrapping it in [WithResource].

Footgun 1:

	func WithOpen(name string, use func(*os.File) error) error {
		return WithResource(func() (*os.File, error) { return os.Open(name) }, use)
	}

	func main() {
		err := WithOpen("test.txt", func(f *os.File) error {
			// use
			return nil
		})
		if err != nil {
			// handle the error
		}
	}

You won't forget to close it because it is done automatically by [WithResource]. Also, there are
some wrappers provided by the package for common resources. This foot gun could've been solved
by just using [pkg/github.com/thelissimus/brago/os.WithOpen].

Footgun 1 quiz:

a) Same solution as the above "Footgun 1" solution.

c) Solution:

	func WithHttpResponse(acquire func() (*http.Response, error), use func(r *http.Response) error) error {
		return Bracket(acquire, func(r *http.Response) error { return r.Body.Close() }, use)
	}

	func main() {
		err := WithHttpResponse(func() (*http.Response, error) { return http.Get("http://go.dev") }, func(r *http.Response) error {
			// use
			return nil
		})
		if err != nil {
			// handle the error
		}
	}

d) Solution:

	func WithTicker(t time.Duration, use func(r *time.Ticker) error) error {
		return Bracket(
			func() (*time.Ticker, error) { return time.NewTicker(t), nil },
			func(r *time.Ticker) error { r.Stop(); return nil },
			use,
		)
	}

	func main() {
		err := WithTicker(time.Second, func(r *time.Ticker) error {
			// use
			return nil
		})
		if err != nil {
			// handle the error
		}
	}

Footguns 2 and 3 are solved similarly to Footgot 1.

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

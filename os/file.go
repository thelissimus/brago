// SPDX-License-Identifier: BSD-3-Clause

/* Wrappers of brago for stdlib os package. */
package os

import (
	"os"

	"github.com/thelissimus/brago"
)

// WithOpen is a wrapper for [pkg/os.Open].
func WithOpen(name string, use func(*os.File) error) error {
	return brago.WithResource(func() (*os.File, error) { return os.Open(name) }, use)
}

// WithCreate is a wrapper for [pkg/os.Create].
func WithCreate(name string, use func(*os.File) error) error {
	return brago.WithResource(func() (*os.File, error) { return os.Create(name) }, use)
}

// WithOpenFile is a wrapper for [pkg/os.OpenFile].
func WithOpenFile(name string, flag int, perm os.FileMode, use func(*os.File) error) error {
	return brago.WithResource(func() (*os.File, error) { return os.OpenFile(name, flag, perm) }, use)
}

package logger

import "github.com/Mrzrb/goerr/core"

type Unit struct{}

// File implements core.Outputer.
func (u *Unit) File() string {
	return "test.gen.go"
}

// Imports implements core.Outputer.
func (u *Unit) Imports() []string {
	return nil
}

// Output implements core.Outputer.
func (u *Unit) Output() []byte {
	return nil
}

// Package implements core.Outputer.
func (u *Unit) Package() string {
	return "main"
}

var _ core.Outputer = (*Unit)(nil)

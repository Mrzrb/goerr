package logger

import (
	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
)

// @Getter
type Unit struct {
	IsMethod  bool
	Receiver  string
	FuncName  string
	Params    []core.Ident
	Param     string
	Returns   []core.Ident
	Return    string
	CallParam string
	HasReturn bool
	Packages  string

	Context core.Ident

	Import []string
	file   string
}

// Valid implements core.Outputer.
func (u *Unit) Valid() bool {
	if len(u.Params) <= 0 {
		return false
	}
	if u.Params[0].Type != "*gin.Context" {
		return false
	}
	u.Context = u.Params[0]
	return true
}

// File implements core.Outputer.
func (u *Unit) File() string {
	return u.file
}

// Imports implements core.Outputer.
func (u *Unit) Imports() []string {
	u.Import = append(u.Import, "time", "git.zuoyebang.cc/pkg/golib/v2/zlog")
	return utils.Uniq(u.Import, func(t string) string { return t })
}

// Output implements core.Outputer.
func (u *Unit) Output() []byte {
	return utils.Must(core.ExecuteTemplate(temp.Lookup(tplName), u))
}

// Package implements core.Outputer.
func (u *Unit) Package() string {
	return u.Packages
}

var _ core.Outputer = (*Unit)(nil)

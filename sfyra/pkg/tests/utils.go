// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"errors"
	"io"
)

var errMain = errors.New("testing: unexpected use of func Main")

type matchStringOnly func(pat, str string) (bool, error)

func (f matchStringOnly) MatchString(pat, str string) (bool, error) { return f(pat, str) }

func (f matchStringOnly) StartCPUProfile(w io.Writer) error { return errMain }

func (f matchStringOnly) StopCPUProfile() {}

func (f matchStringOnly) WriteProfileTo(string, io.Writer, int) error { return errMain }

func (f matchStringOnly) ImportPath() string { return "" }

func (f matchStringOnly) StartTestLog(io.Writer) {}

func (f matchStringOnly) StopTestLog() error { return errMain }

func (f matchStringOnly) SetPanicOnExit0(bool) {}

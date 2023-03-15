// Copyright (c) 2023 ezotaka
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package pecofy

import "strings"

type LineContainer interface {
	String() string
	DisplayString() string
}

type StringLineContainer string

func (s StringLineContainer) String() string {
	return string(s)
}

func (s StringLineContainer) DisplayString() string {
	return string(s)
}

type Input struct {
	lines []LineContainer
}

func newInput() *Input {
	return &Input{}
}

func (i *Input) AddLineContainers(lines ...LineContainer) *Input {
	i.lines = append(i.lines, lines...)
	return i
}

func (i *Input) Get(index uint64) LineContainer {
	return i.lines[index]
}

func (i *Input) String() string {
	var lines []string
	for _, l := range i.lines {
		// line must not contain "\n"
		lines = append(lines, strings.ReplaceAll(l.DisplayString(), "\n", "\\n"))
	}
	return strings.Join(lines, "\n")
}

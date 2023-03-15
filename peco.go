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

import (
	"context"
	"errors"
	"strings"

	"github.com/google/btree"
	"github.com/peco/peco"
	"github.com/peco/peco/line"
)

type Runner interface {
	Run(ctx context.Context, option *Option, input *Input) ([]LineContainer, error)
}

type Peco struct {
	peco   *peco.Peco
	option *Option
}

func newPeco() *Peco {
	return &Peco{peco: peco.New(), option: &Option{}}
}

// Run implements Peco
func (p *Peco) Run(ctx context.Context, option *Option, input *Input) ([]LineContainer, error) {
	p.peco.Stdin = strings.NewReader(input.String())
	p.peco.Argv = *option

	if err := p.peco.Run(ctx); err != nil {
		if err.Error() == "collect results" {
			// peco returns this error when user selected a line
			selection := p.peco.Selection()
			if selection.Len() == 0 {
				lineNumber := p.peco.Location().LineNumber()
				if l, err := p.peco.CurrentLineBuffer().LineAt(lineNumber); err == nil {
					selection.Add(l)
				}
			}

			var selected []LineContainer
			p.peco.Selection().Ascend(func(it btree.Item) bool {
				selected = append(selected, input.Get(it.(line.Line).ID()))
				return true
			})
			return selected, nil
		}
		if err.Error() == "user canceled" {
			return []LineContainer{}, nil
		}
	}
	return []LineContainer{}, errors.New("unknown peco exit status")
}

type Option []string

func (o *Option) addArgs(args ...string) {
	*o = append(*o, args...)
}

// --prompt option for peco
func (o *Option) Prompt(prompt string) *Option {
	o.addArgs("--prompt", prompt)
	return o
}

// --query option for peco
func (o *Option) Query(query string) *Option {
	o.addArgs("--query", query)
	return o
}

// --select-1 option for peco
func (o *Option) SelectOne() *Option {
	o.addArgs("--select-1")
	return o
}

// --selection-prefix option for peco
func (o *Option) SelectionPrefix(prefix string) *Option {
	o.addArgs("--selection-prefix", prefix)
	return o
}

type pecoMock struct {
	selectedIndexes []uint64
}

func newMock(indexes ...uint64) *pecoMock {
	return &pecoMock{selectedIndexes: indexes}
}

// Run implements Peco
func (p *pecoMock) Run(ctx context.Context, option *Option, input *Input,
) (selected []LineContainer, err error) {
	for _, i := range p.selectedIndexes {
		selected = append(selected, input.Get(i))
	}
	return selected, nil
}

var _ Runner = &pecoMock{}

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

// Use peco to let users of the CLI tool select any string
package pecofy

import (
	"context"
)

type pecofy struct {
	option Option
	peco   Runner
}

func New() *pecofy {
	return &pecofy{Option{}, newPeco()}
}

func (p *pecofy) Option() *Option {
	return &p.option
}

func NewMock(indexes ...uint64) *pecofy {
	return &pecofy{Option{}, newMock(indexes...)}
}

func (p *pecofy) Run(ctx context.Context, lines []string) (selected []string, err error) {
	selected = []string{}
	var input Input
	for _, l := range lines {
		input.AddLineContainers(StringLineContainer(l))
	}
	if result, err := p.peco.Run(ctx, &p.option, &input); err == nil {
		for _, r := range result {
			selected = append(selected, string(r.(StringLineContainer)))
		}
	}
	return
}

func (p *pecofy) RunContainers(ctx context.Context, lines []LineContainer) (selected []LineContainer, err error) {
	return p.peco.Run(ctx, &p.option, NewInput().AddLineContainers(lines...))
}

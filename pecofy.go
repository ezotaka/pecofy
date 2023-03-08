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
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/peco/peco"
)

// peco Runner
type Runner struct {
	argv  []string
	input strings.Builder
}

// Instantiate peco runner
func New() *Runner {
	return &Runner{}
}

// --prompt option for peco
func (p *Runner) Prompt(prompt string) *Runner {
	p.argv = append(p.argv, "--prompt", prompt)
	return p
}

// --query option for peco
func (p *Runner) Query(query string) *Runner {
	p.argv = append(p.argv, "--query", query)
	return p
}

// --select-1 option for peco
func (p *Runner) SelectOne() *Runner {
	p.argv = append(p.argv, "--select-1")
	return p
}

// --selection-prefix option for peco
func (p *Runner) SelectionPrefix(prefix string) *Runner {
	p.argv = append(p.argv, "--selection-prefix", "prefix")
	return p
}

// Pass string to peco
func (p *Runner) InputString(input string) *Runner {
	p.input.WriteString(input)
	if !strings.HasSuffix(input, "\n") {
		p.input.WriteString("\n")
	}
	return p
}

// Pass strings to peco
func (p *Runner) InputStrings(input []string) *Runner {
	for _, i := range input {
		p.InputString(i)
	}
	return p
}

// Run peco command
func (p *Runner) Run() (selected []string, ok bool, err error) {
	return p.RunWithContext(context.Background())
}

// Run peco command with Context
func (p *Runner) RunWithContext(ctx context.Context) (selected []string, ok bool, err error) {
	result := bytes.NewBufferString("")

	peco := peco.New()
	peco.Stdin = strings.NewReader(p.input.String())
	peco.Stdout = result
	peco.Argv = p.argv

	err = peco.Run(ctx)
	if err != nil {
		if err.Error() == "collect results" {
			// peco returns this error when user selected a line
			peco.PrintResults()
			return strings.Split(strings.TrimSuffix(result.String(), "\n"), "\n"), true, nil
		}
		if err.Error() == "user canceled" {
			return []string{}, false, nil
		}
	}

	return []string{}, false, errors.New("unknown peco exit status")
}

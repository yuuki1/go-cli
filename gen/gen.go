// Package gen generates github.com/yuuki1/go-cli.Command from
// function docs in files specified.
package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"regexp"
	"strings"
)

var (
	FileFormat = `// auto-generated file

package main

import "github.com/yuuk1/go-cli"

func init() {%s}
`
	CommandFormat = `
	cli.Use(
		&cli.Command{
			Name:   %q,
			Action: %s,
			Short:  %q,
			Long:   %q,
		},
	)
`
)

/*
Generate reads source file for command actions with their usage documentations
and writes Go code that registers the command to cli.

Usage documentation should be like below:

	// +command <name> - <short>
	//
	// <usage line>
	//
	// <long description>...
	func action(flags *flag.FlagSet, args []string) {
	}

Currently, generated files will be like below:

	// auto-generated file

	package main

	import "github.com/yuuki1/go-cli"

	func init() {
	    cli.Use(
	        &cli.Command{
	            Name:   "foo",
	            ...
	        }
	    )
	    ...
	}
*/
func Generate(w io.Writer, path string, src interface{}) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	commandCodes := []string{}

	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Doc == nil {
			continue
		}

		doc := funcDecl.Doc.Text()
		pos := strings.Index(doc, "+command")
		if pos == -1 {
			continue
		}

		doc = doc[pos+len("+command "):]
		re := regexp.MustCompile(`^ *(\S+) +- +(.+)\n((?s).+)`)
		m := re.FindStringSubmatch(doc)
		if m == nil {
			continue
		}

		var (
			name  = m[1]
			short = m[2]
			long  = strings.TrimSpace(m[3])
		)

		commandCodes = append(
			commandCodes,
			fmt.Sprintf(CommandFormat, name, funcDecl.Name.Name, short, long),
		)
	}

	code := fmt.Sprintf(FileFormat, strings.Join(commandCodes, ""))

	_, err = w.Write([]byte(code))
	return err
}

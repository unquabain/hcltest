/* hcltest
 *
 * This is a little program to explore HashiCorp's HCL library in Go
 * for the purpose of seeing whether it would be a good format for configuration
 * of simple programs and microservices written in Go.
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/pflag"
	"github.com/zclconf/go-cty/cty"
)

// Server is an example block structure to be read from a config file.
type Server struct {
	// The label value in the tag makes this one of the string values that come
	// before the first curly brace.
	Name string `hcl:"name,label"`

	Addr string `hcl:"addr"`
}

// Config is an example structure to be read from a config file.
type Config struct {
	// The block value means that this variable takes attributes
	// in curly braces.
	Servers []Server `hcl:"server,block"`
}

type ArgsRunMode int

const (
	DocumentFuncMode ArgsRunMode = iota
	DocumentAllFuncsMode
	SerdeMode
)

// Args just handles CLI argument parsing.
type Args struct {
	Funcs     bool
	Func      string
	Filename  string
	Variables string
	Output    string
}

func (t *Args) Parse() error {
	fs := pflag.NewFlagSet(`hcltest`, pflag.ExitOnError)
	fs.BoolVarP(&t.Funcs, `funcs`, `f`, false, `Display documentation about all template functions`)
	fs.StringVarP(&t.Func, `func`, `F`, ``, `Display documentation about a particular function`)
	fs.StringVarP(&t.Variables, `vars`, `v`, ``, `Filename to read for variables`)
	fs.StringVarP(&t.Output, `output`, `o`, `json`, `Output format. Valid values are json and hcl`)
	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		return err
	}
	if t.Func == `` && !t.Funcs && fs.NArg() == 0 {
		return fmt.Errorf(`you must specify the name of a HCL file to parse`)
	}
	t.Filename = fs.Arg(0)
	return nil
}

func (t *Args) runMode() ArgsRunMode {
	if t.Funcs {
		return DocumentAllFuncsMode
	}
	if t.Func != `` {
		return DocumentFuncMode
	}
	return SerdeMode
}

// Fatalf is an alien who has eaten too many cats.
func (Args) Fatalf(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(-1)
}

func (t *Args) Run() {
	switch m := t.runMode(); m {
	case DocumentAllFuncsMode:
		t.documentAllFuncs()
	case DocumentFuncMode:
		t.documentFunc()
	case SerdeMode:
		t.serde()
	default:
		t.Fatalf(`unknown run mode: %d`, m)
	}
}

func (Args) documentAllFuncs() {
	fds := FuncDescriptions()
	for _, fname := range slices.Sorted(maps.Keys(fds)) {
		fdesc := fds[fname]
		fmt.Printf("%s\n\t%s\n", fname, fdesc)
	}
}

func (t *Args) documentFunc() {
	fname, desc, found := FuncDescription(t.Func)
	if !found {
		t.Fatalf("Unknown function %s\n", t.Func)
	}
	fmt.Printf("%s\n\t%s\n", fname, desc)
}

func (t *Args) make_context() *hcl.EvalContext {
	ctx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
		Functions: funcmap,
	}
	t.get_variables(ctx)
	return ctx

}

func (t *Args) get_variables(ctx *hcl.EvalContext) {
	// Variables to deserialize from a separate file, optionally.
	if t.Variables == `` {
		return
	}
	// Read the variables file into the map
	if err := hclsimple.DecodeFile(t.Variables, nil, &ctx.Variables); err != nil {
		t.Fatalf(`could not read variable file: %s`, err.Error())
	}
}

func (t *Args) serde() {
	// A config file to try to deserialize
	var config Config
	t.deserialize(&config)
	t.serialize(config)

}

func (t *Args) deserialize(config *Config) {
	ctx := t.make_context()
	// Read the config file from the source.
	/*
		// This is the simple implementation if we're not using any extensions.
		if err := hclsimple.DecodeFile(args.Filename, &ctx, &config); err != nil {
			fmt.Fprintf(os.Stderr, "could not decode %s: %s\n", args.Filename, err.Error())
			os.Exit(-1)
		}
	*/

	// This is a more involved example that shows how to use an extension to pre-parse the body.
	if file, err := hclparse.NewParser().ParseHCLFile(t.Filename); err.HasErrors() {
		t.Fatalf(`could not parse file: %s`, err.Error())
	} else if err := gohcl.DecodeBody(dynblock.Expand(file.Body, ctx), ctx, config); err.HasErrors() {
		t.Fatalf(`could not decode body: %s`, err.Error())
	}
}

func (t *Args) serialize(config Config) {
	switch t.Output {
	case `json`:

		// Print this out using the normal Go JSON encoder
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(config)
	case `hcl`:

		// We can also print out the pre-processed HCL with all the
		// variables and functions evaluated.
		hclFile := hclwrite.NewEmptyFile()
		gohcl.EncodeIntoBody(config, hclFile.Body())
		os.Stdout.Write(hclFile.Bytes())
	default:
		t.Fatalf(`unknown output format %s; valid formats are json and hcl`, t.Output)
	}
}

func main() {
	args := new(Args)
	if err := args.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(-1)
	}
	args.Run()
}

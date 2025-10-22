/* hcltest
 *
 * This is a little function to explore HashiCorp's HCL library in Go
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

func main() {
	args := new(Args)
	if err := args.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(-1)
	}

	// Handle some query options to inspect some of the standard library
	// template functions. See function.go for details.
	if args.Func != `` {
		fname, desc, found := FuncDescription(args.Func)
		if !found {
			fmt.Fprintf(os.Stderr, "Unknown function %s\n", args.Func)
			os.Exit(-1)
		}
		fmt.Printf("%s\n\t%s\n", fname, desc)
		return
	} else if args.Funcs {
		fds := FuncDescriptions()
		for _, fname := range slices.Sorted(maps.Keys(fds)) {
			fdesc := fds[fname]
			fmt.Printf("%s\n\t%s\n", fname, fdesc)
		}
		return
	}

	// A config file to try to deserialize
	var config Config

	// Variables to deserialize from a separate file, optionally.
	vars := make(map[string]cty.Value)
	ctx := hcl.EvalContext{
		Variables: vars,
		Functions: funcmap,
	}
	if args.Variables != `` {
		// Read the variables file into the map
		if err := hclsimple.DecodeFile(args.Variables, nil, &vars); err != nil {
			fmt.Fprintf(os.Stderr, `could not read variable file: %s`, err.Error())
			os.Exit(-1)
		}
		ctx.Variables = vars
	}

	// Read the config file from the source.
	/*
		// This is the simple implementation if we're not using any extensions.
		if err := hclsimple.DecodeFile(args.Filename, &ctx, &config); err != nil {
			fmt.Fprintf(os.Stderr, "could not decode %s: %s\n", args.Filename, err.Error())
			os.Exit(-1)
		}
	*/

	// This is a more involved example that shows how to use an extension to pre-parse the body.
	if file, err := hclparse.NewParser().ParseHCLFile(args.Filename); err.HasErrors() {
		fmt.Fprintf(os.Stderr, `could not parse file: %s`, err.Error())
		os.Exit(-1)
	} else if err := gohcl.DecodeBody(dynblock.Expand(file.Body, &ctx), &ctx, &config); err.HasErrors() {
		fmt.Fprintf(os.Stderr, `could not decode body: %s`, err.Error())
		os.Exit(-1)
	}
	switch args.Output {
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
		fmt.Fprintf(os.Stderr, `unknown output format %s; valid formats are json and hcl`, args.Output)
		os.Exit(-1)
	}
}

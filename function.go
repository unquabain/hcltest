package main

/* function.go is mostly generated from a sed script.
 *
 * This file creates a standard library of template functions to make available
 * to the HCL parsers. Most of them are taken from the cty/function/stdlib
 * package, but there is one example of creating a new, custom function.
 */
import (
	"bytes"
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var funcmap = map[string]function.Function{
	// Standard Library Functions
	`absolute`:               stdlib.AbsoluteFunc,
	`add`:                    stdlib.AddFunc,
	`and`:                    stdlib.AndFunc,
	`assertnotnull`:          stdlib.AssertNotNullFunc,
	`byteslen`:               stdlib.BytesLenFunc,
	`bytesslice`:             stdlib.BytesSliceFunc,
	`csvdecode`:              stdlib.CSVDecodeFunc,
	`ceil`:                   stdlib.CeilFunc,
	`chomp`:                  stdlib.ChompFunc,
	`chunklist`:              stdlib.ChunklistFunc,
	`coalesce`:               stdlib.CoalesceFunc,
	`coalescelist`:           stdlib.CoalesceListFunc,
	`compact`:                stdlib.CompactFunc,
	`concat`:                 stdlib.ConcatFunc,
	`contains`:               stdlib.ContainsFunc,
	`distinct`:               stdlib.DistinctFunc,
	`divide`:                 stdlib.DivideFunc,
	`element`:                stdlib.ElementFunc,
	`equal`:                  stdlib.EqualFunc,
	`flatten`:                stdlib.FlattenFunc,
	`floor`:                  stdlib.FloorFunc,
	`formatdate`:             stdlib.FormatDateFunc,
	`format`:                 stdlib.FormatFunc,
	`formatlist`:             stdlib.FormatListFunc,
	`greaterthan`:            stdlib.GreaterThanFunc,
	`greaterthanorequalto`:   stdlib.GreaterThanOrEqualToFunc,
	`hasindex`:               stdlib.HasIndexFunc,
	`indent`:                 stdlib.IndentFunc,
	`index`:                  stdlib.IndexFunc,
	`int`:                    stdlib.IntFunc,
	`jsondecode`:             stdlib.JSONDecodeFunc,
	`jsonencode`:             stdlib.JSONEncodeFunc,
	`join`:                   stdlib.JoinFunc,
	`keys`:                   stdlib.KeysFunc,
	`length`:                 stdlib.LengthFunc,
	`lessthan`:               stdlib.LessThanFunc,
	`lessthanorequalto`:      stdlib.LessThanOrEqualToFunc,
	`log`:                    stdlib.LogFunc,
	`lookup`:                 stdlib.LookupFunc,
	`lower`:                  stdlib.LowerFunc,
	`max`:                    stdlib.MaxFunc,
	`merge`:                  stdlib.MergeFunc,
	`min`:                    stdlib.MinFunc,
	`modulo`:                 stdlib.ModuloFunc,
	`multiply`:               stdlib.MultiplyFunc,
	`negate`:                 stdlib.NegateFunc,
	`notequal`:               stdlib.NotEqualFunc,
	`not`:                    stdlib.NotFunc,
	`or`:                     stdlib.OrFunc,
	`parseint`:               stdlib.ParseIntFunc,
	`pow`:                    stdlib.PowFunc,
	`range`:                  stdlib.RangeFunc,
	`regexall`:               stdlib.RegexAllFunc,
	`regex`:                  stdlib.RegexFunc,
	`regexreplace`:           stdlib.RegexReplaceFunc,
	`replace`:                stdlib.ReplaceFunc,
	`reverse`:                stdlib.ReverseFunc,
	`reverselist`:            stdlib.ReverseListFunc,
	`sethaselement`:          stdlib.SetHasElementFunc,
	`setintersection`:        stdlib.SetIntersectionFunc,
	`setproduct`:             stdlib.SetProductFunc,
	`setsubtract`:            stdlib.SetSubtractFunc,
	`setsymmetricdifference`: stdlib.SetSymmetricDifferenceFunc,
	`setunion`:               stdlib.SetUnionFunc,
	`signum`:                 stdlib.SignumFunc,
	`slice`:                  stdlib.SliceFunc,
	`sort`:                   stdlib.SortFunc,
	`split`:                  stdlib.SplitFunc,
	`strlen`:                 stdlib.StrlenFunc,
	`substr`:                 stdlib.SubstrFunc,
	`subtract`:               stdlib.SubtractFunc,
	`timeadd`:                stdlib.TimeAddFunc,
	`title`:                  stdlib.TitleFunc,
	`trim`:                   stdlib.TrimFunc,
	`trimprefix`:             stdlib.TrimPrefixFunc,
	`trimspace`:              stdlib.TrimSpaceFunc,
	`trimsuffix`:             stdlib.TrimSuffixFunc,
	`upper`:                  stdlib.UpperFunc,
	`values`:                 stdlib.ValuesFunc,
	`zipmap`:                 stdlib.ZipmapFunc,

	// Example custom function
	`urlify`: function.New(&URLify),
}

// FuncDescription return the documentation strings for a function.
func FuncDescription(funcname string) (string, string, bool) {
	f, found := funcmap[funcname]
	if !found {
		return ``, ``, false
	}
	return fsignature(funcname, f), f.Description(), true
}

// fsignature prints the function signature of a function,
// including the names and types of its parameters.
func fsignature(fname string, f function.Function) string {
	buff := bytes.NewBufferString(fname)
	buff.WriteRune('(')
	hasStatParams := false
	for i, param := range f.Params() {
		if i > 0 {
			buff.WriteString(`, `)
		}
		fmt.Fprintf(buff, "%s: %s", param.Name, param.Type.FriendlyName())
		hasStatParams = true
	}
	if vp := f.VarParam(); vp != nil {
		if hasStatParams {
			buff.WriteString(`, `)
		}
		fmt.Fprintf(buff, "%s: ...%s", vp.Name, vp.Type.FriendlyName())
	}
	buff.WriteRune(')')

	return buff.String()
}

// FuncDescriptions returns a map of the function signature and documentation for
// each function in the funcmap, defined above.
func FuncDescriptions() map[string]string {
	descs := make(map[string]string)
	for fname, f := range funcmap {
		descs[fsignature(fname, f)] = f.Description()
	}
	return descs
}

// An example custom template function
var URLify = function.Spec{
	Description: `an example function that takes a domain name and returns an obvious URL`,
	Params: []function.Parameter{
		{
			Name: `domain`,
			Type: cty.String,
		},
	},
	Type: func(_ []cty.Value) (cty.Type, error) {
		return cty.String, nil
	},
	Impl: func(args []cty.Value, _retType cty.Type) (cty.Value, error) {
		if len(args) != 1 {
			return cty.NullVal(cty.String), fmt.Errorf(`unexpected number of arguments`)
		}
		s := fmt.Sprintf(`https://www.%s.com`, args[0].AsString())
		return cty.StringVal(s), nil
	},
}

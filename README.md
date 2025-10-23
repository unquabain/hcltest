# hcltest

This is a little program to explore HashiCorp's HCL library in Go
for the purpose of seeing whether it would be a good format for
configuration of simple programs and microservices written in Go.

## What is HCL?

HCL is the HashiCorp Configuration Language, and it the native language
of Terraform and some of HashiCorp's other projects.

The idea is that JSON is too restrictive: it makes you type too much,
there's no way to write multi-line strings, no comments, etc.

YAML is a bit better, but it suffers from the same problem as Python:
you can write invisible syntax errors. Also, the latest spec of YAML
actually _removes_ some very useful features like anchors that make
parsers more complicated, but releave a lot of repetitive work for
users.

TOML is very good for very flat structures, but gets very messy if your
structures start to nest too deeply.

Some projects, like Taskfile and Helm, use YAML that they first run
through a template engine—often the Go template engine—first to add
some of the missing functionality back. This also has the problem of
having multiple syntaxes that often compete with one another. Not to
mention that if you're not a Go developer, the Go template syntax and
documentation will be foreign to you.

Some projects opt to use a full language like JavaScript or Ruby for
their configuration files. This allows you a lot of freedom, but also
allows you to write side effects and non-deterministic code that should
very much _not_ go in a configuration.

HCL is a solution that combines a structured declarative language that's
JSON-compatible, a simple expression evaluation language, and a template
language into one syntax. It's as simple to read as TOML, but doesn't
get fouled up by deep nesting or suffer from semantically-significant
whitespace like YAML.

HashiCorp's implementation of HCL is in Go, and the libraries are
public. So while there are implementations of HCL in other languages,
the Go version is the most complete and up-to-date.

This is very similar to Apple's PKL language, though that is implemented
in Java, and is much less mature.

## What is this Program?

This program is basically a do-nothing test of the HCL libraries. It
reads in a variables file and an HCL file, deserializes them into
Go structures, and prints the Go structures out using either the Go
standard library's JSON encoding, or HCL's encoding.

So it could be useful to translate HCL into JSON, or to pre-parse
dynamic HCL expressions, functions and variables into simpler HCL.
Mostly it exists as a learning exercize for me.

## How do I run it?

If you have `Taskfile` installed, simply run `task test clean`. That
should build the program, run the tests, and clean itself up. Simple as...

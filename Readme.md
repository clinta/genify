Genify is yet another code generator to simulate generics. It works as a simple
text processor, multiplying blocks of code with replacements specified in
comments.

## QuickStart

```go src.gen
//genify:generic=foo,bar

func (g *Generic) GetThing() {
	return g.genericThing
}
```

```bash
genify -in src.gen
```

```go
func (g *Foo) GetThing() {
	return g.fooThing
}

func (g *Bar) GetThing() {
	return g.barThing
}
```

Inspired by [cheekybits/genny](https://github.com/cheekybits/genny)

### How is this different than genny?:

* Genify does not care if your source is valid
go code. It does not expect you to create special generic types.
In practice using an empty interface often doesn't result in valid go code
anyway, since the empty interface does not implement what you may be using
in your generic code.

* Genify splits the source file into blocks, terminated by
a non-indented line followed by a blank line. For every block, genify looks for
source substrings that need to be replaced, and if they are found a copy of the
block with the source replaced with all replacements is written. This is simple
string replacement, so it works with comments, variables, types, functions and
anything else. And Documentation comments which preceed function and type
definitions are handled in the same block as the function that they document.

* Genify replacement definitions are configured via special comments in the
template file. To configure a replacement just add a comment like this
`//genify:generic=foo,bar`. Any code after this comment which contains the
string `generic` will be written out twice, once with `generic` replaced with
`foo` and again with `generic` replaced with `bar`. Replacement definitions only
affect code after the definition has been defined. And definitions can be
overwritten or redefined to allow different behavior in later parts of the
template.

* Genify operates on unexported and exported types of the same name. If you
setup a replacement of `//genify:generic=foo,bar` all instances of `Generic`
will also be replaced with `Foo`.

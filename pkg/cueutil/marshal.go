package cueutil

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/token"
	"cuelang.org/go/encoding/gocode/gocodec"
)

func Marshal(v interface{}) ([]byte, error) {
	var r cue.Runtime
	var c gocodec.Config

	codec := gocodec.New(&r, &c)
	value, err := codec.Decode(v)
	if err != nil {
		return nil, err
	}

	syn := []cue.Option{
		cue.Final(),
		cue.Docs(true),
		cue.Definitions(true),
		cue.Attributes(true),
		cue.Optional(true),
	}
	opts := []format.Option{
		format.Simplify(),
		format.UseSpaces(4),
		format.TabIndent(false),
	}

	n := value.Syntax(syn...)
	bs, err := format.Node(toFile(n), opts...)
	if err != nil {
		return bs, err
	}

	return bs, nil
}

func toFile(n ast.Node) *ast.File {
	switch x := n.(type) {
	case nil:
		return nil
	case *ast.StructLit:
		return &ast.File{Decls: x.Elts}
	case ast.Expr:
		ast.SetRelPos(x, token.NoSpace)
		return &ast.File{Decls: []ast.Decl{&ast.EmbedDecl{Expr: x}}}
	case *ast.File:
		return x
	default:
		panic(fmt.Sprintf("Unsupported node type %T", x))
	}
}

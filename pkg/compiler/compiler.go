package compiler

import (
	"io"

	"github.com/liamg/bfc/pkg/backend/generators"
	"github.com/liamg/bfc/pkg/frontend/lexer"
	"github.com/liamg/bfc/pkg/frontend/parser"
)

func Compile(r io.Reader, w io.Writer, gen generators.Generator) error {

	tokens, err := lexer.Lex(r)
	if err != nil {
		return err
	}

	statements, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	return gen.Generate(statements, w)
}

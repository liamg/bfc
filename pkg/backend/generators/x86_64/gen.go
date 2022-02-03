package x86_64

import (
	"io"

	"github.com/liamg/bfc/pkg/backend/generators"
	"github.com/liamg/bfc/pkg/frontend/parser"
)

type generator struct{}

func New() generators.Generator {
	return &generator{}
}

func (g *generator) Generate(statements []parser.Statement, w io.Writer) error {
	for _, statement := range statements {
		// TODO: compile statement
		_ = statement
	}
	return nil
}

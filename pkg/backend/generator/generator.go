package generator

import (
	"io"

	"github.com/liamg/bfc/pkg/frontend/parser"
)

type Generator interface {
	Generate(statements []parser.Statement, w io.Writer) error
}

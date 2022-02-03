package parser

import "github.com/liamg/bfc/pkg/frontend/lexer"

type Statement struct {
	Type  lexer.TokenType
	Count int
}

func Parse(tokens []lexer.Token) ([]Statement, error) {
	var statements []Statement
	var statement Statement
	for _, token := range tokens {
		if token.Type == statement.Type {
			statement.Count++
			continue
		}
		if statement.Count > 0 {
			statements = append(statements, statement)
		}
		statement.Type = token.Type
		statement.Count = 1
	}
	if statement.Count > 0 {
		statements = append(statements, statement)
	}
	return statements, nil
}

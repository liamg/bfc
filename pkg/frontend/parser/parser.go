package parser

import (
	"github.com/liamg/bfc/pkg/frontend/lexer"
)

type Statement struct {
	Type  lexer.TokenType
	Count int
	Label int
	Jump  int
}

func Parse(tokens []lexer.Token) ([]Statement, error) {

	// TODO validate stuff here
	// e.g. check for ] without [

	// merge adjacent statements together
	var statements []Statement
	var statement Statement
	for _, token := range tokens {
		if token.Type == lexer.TokenComment {
			continue
		}
		if token.Type == statement.Type && isRepeatable(token.Type) {
			statement.Count++
			continue
		}
		if statement.Count > 0 {
			statements = append(statements, statement)
		}
		statement.Label = -1
		statement.Type = token.Type
		statement.Count = 1
	}
	if statement.Count > 0 {
		statements = append(statements, statement)
	}

	// store metadata about JMP locations
	var opens []int
	for i, statement := range statements {
		switch statement.Type {
		case lexer.TokenJumpForward:
			statements[i].Label = i
			opens = append(opens, i)
		case lexer.TokenJumpBackward:
			lastJump := opens[len(opens)-1]
			opens = opens[:len(opens)-1]
			statements[i].Label = i
			statements[i].Jump = lastJump
			statements[lastJump].Jump = i
		}
	}

	return statements, nil
}

func isRepeatable(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.TokenMoveLeft, lexer.TokenMoveRight, lexer.TokenIncrement, lexer.TokenDecrement:
		return true
	default:
		return false
	}
}

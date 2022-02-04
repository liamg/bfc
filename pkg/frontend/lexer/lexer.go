package lexer

import (
	"bufio"
	"io"
)

func Lex(r io.Reader) ([]Token, error) {
	br := bufio.NewReader(r)
	var tokens []Token
	var id int64
	for {
		r, _, err := br.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch TokenType(r) {
		case TokenUnknown,
			TokenMoveRight,
			TokenMoveLeft,
			TokenIncrement,
			TokenDecrement,
			TokenOutput,
			TokenInput,
			TokenJumpForward,
			TokenJumpBackward:
			tokens = append(tokens, Token{
				ID:   id,
				Type: TokenType(r),
			})
		default:
			tokens = append(tokens, Token{
				ID:   id,
				Type: TokenComment,
			})
		}
	}
	return tokens, nil
}

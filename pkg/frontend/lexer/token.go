package lexer

type TokenType rune

/*
> 	Move the pointer to the right
< 	Move the pointer to the left
+ 	Increment the memory cell at the pointer
- 	Decrement the memory cell at the pointer
. 	Output the character signified by the cell at the pointer
, 	Input a character and store it in the cell at the pointer
[ 	Jump past the matching ] if the cell at the pointer is 0
] 	Jump back to the matching [ if the cell at the pointer is nonzero
*/

const (
	TokenUnknown      TokenType = 0
	TokenMoveRight    TokenType = '>'
	TokenMoveLeft     TokenType = '<'
	TokenIncrement    TokenType = '+'
	TokenDecrement    TokenType = '-'
	TokenOutput       TokenType = '.'
	TokenInput        TokenType = ','
	TokenJumpForward  TokenType = '['
	TokenJumpBackward TokenType = ']'
	TokenComment      TokenType = '/'
)

type Token struct {
	ID   int64
	Type TokenType
}

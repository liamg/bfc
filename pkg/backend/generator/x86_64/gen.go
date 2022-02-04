package x86_64

import (
	"fmt"
	"io"

	"github.com/liamg/bfc/pkg/backend/generator"
	"github.com/liamg/bfc/pkg/frontend/lexer"
	"github.com/liamg/bfc/pkg/frontend/parser"
)

type gen struct {
	w io.Writer
}

func New() generator.Generator {
	return &gen{}
}

func (g *gen) Generate(statements []parser.Statement, w io.Writer) error {

	g.w = w

	if err := g.setup(); err != nil {
		return err
	}

	for _, statement := range statements {
		if err := g.compileStatement(statement); err != nil {
			return err
		}
	}

	_, err := g.w.Write([]byte(`
    ; all done!
    mov rax, 60
    mov rdi, 0
    syscall
`))
	return err
}

func (g *gen) setup() error {
	g.write([]byte(`
section .data
    index: dq 0

section .bss
    table: resq 30000

section .text
    global _start

in:
    ; store pointer to character at index in rcx
    mov rcx, table
    add qword rcx, [index]

    ; read character from stdin
    mov rax, 0
    mov rdi, 0
    mov rsi, rcx
    mov rdx, 1
    syscall

    ret

    ret

out:
    ; store pointer to character at index in rcx
    mov rcx, table
    add qword rcx, [index]

    ; print character to stdout
    mov rax, 1
    mov rdi, 1
    mov rsi, rcx
    mov rdx, 1
    syscall

    ret

; move left: decrement the index by 1
left:
    sub qword [index], 8
    ret

; move right: increment the index by 1
right:
    add qword [index], 8
    ret

; move left: decrement the index by rbx
leftm:
    mov rax, rbx
    mov rbx, 8
    mul rbx
    sub qword [index], rax
    ret

; move right: increment the index by rbx
rightm:
    mov rax, rbx
    mov rbx, 8
    mul rbx
    add qword [index], rax
    ret

; inc cell at [index] by 1
inc:
    mov rcx, table
    add rcx, [index]
    add qword [rcx], 1
    ret

; dec cell at [index] by 1
dec:
    mov rcx, table
    add rcx, [index]
    sub qword [rcx], 1
    ret

; inc cell at [index] by rbx
incm:
    mov rcx, table
    add rcx, [index]
    add [rcx], rbx
    ret

; dec cell at [index] by rbx
decm:
    mov rcx, table
    add rcx, [index]
    sub [rcx], rbx
    ret

; read cell at [index] and store in rcx
readcell: 
    mov rcx, table
    add rcx, [index]
    mov rcx, [rcx]
    ret
   
; write cell at [index] to value stored in rcx
writecell:
    mov rdx, table
    add rdx, [index]
    mov qword [rdx], rcx
    ret

_start:

`)...)
	return nil
}

func (g *gen) write(data ...byte) error {
	_, err := g.w.Write(data)
	return err
}

func (g *gen) writeString(l string) error {
	return g.write([]byte(l)...)
}

func (g *gen) writeLine(l string) error {
	return g.writeString("    " + l + "\n")
}

func (g *gen) compileStatement(stmt parser.Statement) error {
	switch stmt.Type {
	case lexer.TokenMoveLeft:
		if stmt.Count == 1 {
			g.writeLine("call left")
		} else {
			g.writeLine(fmt.Sprintf("mov rbx, %d", stmt.Count))
			g.writeLine("call leftm")
		}
	case lexer.TokenMoveRight:
		if stmt.Count == 1 {
			g.writeLine("call right")
		} else {
			g.writeLine(fmt.Sprintf("mov rbx, %d", stmt.Count))
			g.writeLine("call rightm")
		}
	case lexer.TokenIncrement:
		if stmt.Count == 1 {
			g.writeLine("call inc")
		} else {
			g.writeLine(fmt.Sprintf("mov rbx, %d", stmt.Count))
			g.writeLine("call incm")
		}
	case lexer.TokenDecrement:
		if stmt.Count == 1 {
			g.writeLine("call dec")
		} else {
			g.writeLine(fmt.Sprintf("mov rbx, %d", stmt.Count))
			g.writeLine("call decm")
		}
	case lexer.TokenOutput:
		for i := 0; i < stmt.Count; i++ {
			g.writeLine("call out")
		}
	case lexer.TokenInput:
		for i := 0; i < stmt.Count; i++ {
			g.writeLine("call in")
		}
	case lexer.TokenJumpForward:
		g.writeString(fmt.Sprintf("\nlstart%d:\n", stmt.Label))
		g.writeLine("call readcell")
		g.writeLine("add rcx, 0")
		g.writeLine(fmt.Sprintf("jz lend%d\n", stmt.Label))
	case lexer.TokenJumpBackward:
		g.writeString(fmt.Sprintf("\nlend%d:\n", stmt.Jump))
		g.writeLine("call readcell")
		g.writeLine("add rcx, 0")
		g.writeLine(fmt.Sprintf("jnz lstart%d\n", stmt.Jump))
	default:
		return fmt.Errorf("cannot compile %c", stmt.Type)
	}
	return nil
}

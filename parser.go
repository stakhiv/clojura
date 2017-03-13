package main

import (
	"errors"
	"io"
	"strconv"
)

type CoreType uint8

const (
	TypeLiteral CoreType = iota
	TypeNumber
	TypeBoolean
	TypeExpression
	TypeFunction
	TypeMacros
	TypeList
	TypeRecur
)

type Sexpr interface {
	Type() CoreType
	Append(Sexpr) error
	String() string
	Eval(*Context) Sexpr
	Bool() bool
}

type Recur struct {
	Args []Sexpr
}

func (r *Recur) Type() CoreType {
	return TypeRecur
}

func (r *Recur) Bool() bool {
	if len(r.Args) > 0 {
		return true
	}
	return false
}

func (r *Recur) Append(s Sexpr) error {
	return nil
}

func (r *Recur) String() string {
	res := "("
	for _, s := range r.Args {
		res += s.String() + " "
	}
	res += ")"
	return res
}

func (r *Recur) Eval(c *Context) Sexpr {
	return r
}

type Expression struct {
	Elements []Sexpr
}

func (e *Expression) Type() CoreType {
	return TypeExpression
}

func (e *Expression) Bool() bool {
	if len(e.Elements) > 0 {
		return true
	}
	return false
}

func (e *Expression) Append(s Sexpr) error {
	e.Elements = append(e.Elements, s)
	return nil
}

func (e *Expression) String() string {
	res := "("
	for _, s := range e.Elements {
		res += s.String() + " "
	}
	res += ")"
	return res
}

func (e *Expression) Eval(c *Context) Sexpr {
	if len(e.Elements) < 1 {
		return nil
	}
	f := e.Elements[0].Eval(c)
	if f != nil && (f.Type() == TypeFunction || f.Type() == TypeMacros) {
		if f.Type() == TypeMacros {
			return f.(Macros).Create(c, e.Elements)
		}
		args := make([]Sexpr, len(e.Elements)-1)
		for i, arg := range e.Elements[1:] {
			args[i] = arg.Eval(c)
		}
		fun, ok := f.(Function)
		if !ok {
			return nil
		}
		return fun.Call(args)
	}
	return nil
}

type Literal string

func (l Literal) Bool() bool {
	return false
}

func (l Literal) String() string {
	return string(l)
}

func (l Literal) Type() CoreType {
	return TypeLiteral
}

func (l Literal) Append(s Sexpr) error {
	return errors.New("cannot append")
}

func (l Literal) Eval(c *Context) Sexpr {
	val, ok := c.Get(l)
	if ok {
		return val
	}
	return l
}

type Number int

func (n Number) Bool() bool {
	return n != 0
}

func (n Number) String() string {
	return strconv.Itoa(int(n))
}

func (n Number) Type() CoreType {
	return TypeNumber
}

func (n Number) Append(s Sexpr) error {
	return errors.New("cannot append")
}

func (n Number) Eval(c *Context) Sexpr {
	return n
}

type Boolean struct {
	Val bool
}

func (b *Boolean) Bool() bool {
	return b.Val
}

func (b *Boolean) String() string {
	if b.Val {
		return "true"
	}
	return "false"
}

func (b *Boolean) Type() CoreType {
	return TypeBoolean
}

func (b *Boolean) Append(s Sexpr) error {
	return errors.New("cannot append")
}

func (b *Boolean) Eval(c *Context) Sexpr {
	return b
}

type Parser struct {
	lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (p *Parser) Parse() ([]Sexpr, error) {
	resp := make([]Sexpr, 0)
	match := 0

	var s Sexpr
	stack := make([]Sexpr, 0)

	eval := true

	for {
		t, err := p.lexer.ReadToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t {
		case "'":
			log.Debug("EVAL FALSE")
			eval = false
		case "(":
			match += 1
			if s != nil {
				stack = append(stack, s)
			}
			if eval {
				s = &Expression{}
			} else {
				s = NewList()
			}
			eval = true
		case ")":
			if match-1 < 0 {
				return nil, errors.New("unmatched pair")
			}
			match -= 1
			if len(stack) > 0 {
				p := stack[len(stack)-1]
				err := p.Append(s)
				if err != nil {
					return nil, err
				}
				s = p
				stack = stack[:len(stack)-1]
			} else {
				resp = append(resp, s)
				s = nil
			}
		default:
			n, err := strconv.Atoi(t)
			if err == nil {
				s.Append(Number(n))
			} else {
				s.Append(Literal(t))
			}
		}
	}
	if match != 0 {
		return nil, errors.New("unmatched pair")
	}
	return resp, nil
}

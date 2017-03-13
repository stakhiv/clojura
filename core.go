package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	True  = Boolean(true)
	False = Boolean(false)
)

type Context struct {
	vars   map[Literal]Sexpr
	parent *Context
}

func NewContext(parent *Context) *Context {
	return &Context{
		vars:   make(map[Literal]Sexpr),
		parent: parent,
	}
}

func (c *Context) String() string {
	s := ""
	for k, v := range c.vars {
		s += fmt.Sprintf("%s: %s\n", k, v)
	}
	if c.parent != nil {
		s += c.parent.String()
	}
	return s
}

func (c *Context) Get(key Literal) (Sexpr, bool) {
	val, ok := c.vars[key]
	if !ok && c.parent != nil {
		return c.parent.Get(key)
	}
	return val, ok
}

func (c *Context) Set(key Literal, s Sexpr) {
	c.vars[key] = s
}

type Function interface {
	Call([]Sexpr) Sexpr
}

type Macros interface {
	Create(*Context, []Sexpr) Sexpr
}

type function struct {
	name    string
	args    []Literal
	body    Sexpr
	context *Context
}

func NewFunction(name string, args []Literal, body Sexpr, c *Context) *function {
	return &function{
		name:    name,
		args:    args,
		body:    body,
		context: c,
	}
}

func (f function) Call(args []Sexpr) Sexpr {
	context := NewContext(f.context)

	var res Sexpr
	for {
		if len(args) != len(f.args) {
			fmt.Println("Invalid number of arguments for", f.name)
			return nil
		}
		for i, arg := range args {
			context.Set(f.args[i], arg)
		}

		res = f.body.Eval(context)
		if res != nil && res.Type() == TypeRecur {
			args = res.(*Recur).Args
			continue
		}
		break
	}
	return res
}

func (f function) String() string {
	return "func " + f.name
}

func (f function) Type() CoreType {
	return TypeFunction
}

func (f function) Append(s Sexpr) error {
	return errors.New("cannot append")
}

func (f function) Eval(c *Context) Sexpr {
	return f
}

func (f function) Bool() bool {
	return true
}

type coreF func([]Sexpr) Sexpr

func (cf coreF) Bool() bool {
	return true
}

func (cf coreF) String() string {
	return "core func"
}

func (cf coreF) Type() CoreType {
	return TypeFunction
}

func (cf coreF) Append(s Sexpr) error {
	return errors.New("cannot append")
}

func (cf coreF) Eval(c *Context) Sexpr {
	return cf
}

func (cf coreF) Call(args []Sexpr) Sexpr {
	return cf(args)
}

type macros func(*Context, []Sexpr) Sexpr

func (m macros) Bool() bool {
	return true
}
func (m macros) String() string {
	return "core macros"
}

func (m macros) Type() CoreType {
	return TypeMacros
}

func (m macros) Append(s Sexpr) error {
	return errors.New("cannot append")
}

func (m macros) Eval(c *Context) Sexpr {
	return m
}

func (m macros) Create(c *Context, args []Sexpr) Sexpr {
	return m(c, args)
}

var coreContext *Context

func init() {
	coreContext = NewContext(nil)
	coreContext.Set("true", True)
	coreContext.Set("false", False)
	coreContext.Set("+", coreF(coreAdd))
	coreContext.Set("-", coreF(coreSub))
	coreContext.Set("def", macros(coreDef))
	coreContext.Set("println", coreF(corePrintln))
	coreContext.Set("fn", macros(coreFn))
	coreContext.Set("not", coreF(coreNot))
	coreContext.Set("eq", coreF(coreEq))
	coreContext.Set("if", macros(ifMacro))
	coreContext.Set("head", coreF(coreHead))
	coreContext.Set("tail", coreF(coreTail))
	coreContext.Set("conj", coreF(coreConj))
	coreContext.Set("do", macros(doMacro))
	coreContext.Set("time", macros(coreTime))
	coreContext.Set("cons", coreF(coreCons))
	coreContext.Set("recur", coreF(coreRecur))
	coreContext.Set("range", coreF(coreRange))
	coreContext.Set("odd?", coreF(coreOdd))
}

func coreAdd(args []Sexpr) Sexpr {
	var res Number
	for _, arg := range args {
		n, ok := arg.(Number)
		if ok {
			res += n
		}
	}

	return Number(res)
}

func coreSub(args []Sexpr) Sexpr {
	if len(args) < 1 {
		return Number(0)
	}

	res, _ := args[0].(Number)
	for _, arg := range args[1:] {
		r, ok := arg.(Number)
		if ok {
			res -= r
		}
	}

	return res

}

func coreDef(c *Context, args []Sexpr) Sexpr {
	args = args[1:]
	var res Sexpr
	if len(args) > 1 {
		n, ok := args[0].(Literal)
		if ok {
			res = args[1].Eval(c)
			coreContext.Set(n, res)
		}
	}
	return res
}

func coreEq(args []Sexpr) Sexpr {
	if len(args) < 1 {
		return True
	}

	if len(args) > 2 {
		return False
	}

	a, b := args[0], args[1]

	if a.Type() != b.Type() {
		return False
	}

	switch a.Type() {
	case TypeNumber:
		if a.(Number) == b.(Number) {
			return True
		}
		return False
	case TypeBoolean:
		if a.Bool() == b.Bool() {
			return True
		}
		return False
	case TypeMacros:
		fallthrough
	case TypeFunction:
		if a == b {
			return True
		}
		return False
	}
	return False
}

func coreNot(args []Sexpr) Sexpr {
	if len(args) < 1 {
		return True
	}

	if len(args) > 2 {
		return False
	}

	if args[0].Bool() {
		return False
	}
	return True
}

func corePrintln(args []Sexpr) Sexpr {
	res := ""
	for _, arg := range args {
		if arg == nil {
			res += "nil "
		} else {
			res += arg.String() + " "
		}
	}
	fmt.Println("->", res)
	return nil
}

func coreFn(c *Context, args []Sexpr) Sexpr {
	if len(args) < 3 {
		return nil
	}

	name := args[0].String()
	argp, body := args[1], args[2]
	args = argp.(*Expression).Elements
	largs := make([]Literal, len(args))

	for i, arg := range args {
		a, ok := arg.(Literal)
		if !ok {
			fmt.Println("fn argument should be a literal")
			return nil
		}
		largs[i] = a
	}
	return NewFunction(name, largs, body, c)
}

func ifMacro(c *Context, args []Sexpr) Sexpr {
	if len(args) < 3 {
		return nil
	}

	clause := args[1]
	good := args[2]

	var bad Sexpr
	if len(args) == 4 {
		bad = args[3]
	}

	res := clause.Eval(c)
	if res != nil && res.Bool() {
		return good.Eval(c)
	} else if bad != nil {
		return bad.Eval(c)
	}
	return nil
}

func coreHead(args []Sexpr) Sexpr {
	if len(args) != 1 {
		log.Error("bad num of args for head")
		return nil
	}
	l := args[0]

	lst, ok := l.(*List)
	if !ok {
		log.Error("head arg shold be list, got", l.Type())
		return nil
	}

	return lst.Head()
}

func coreTail(args []Sexpr) Sexpr {
	if len(args) != 1 {
		log.Error("bad num of args for tail")
		return nil
	}
	l := args[0]

	lst, ok := l.(*List)
	if !ok {
		log.Error("head arg shold be list, got", l.Type())
		return nil
	}

	return lst.GetTail()
}

func coreConj(args []Sexpr) Sexpr {
	if len(args) != 2 {
		log.Error("bad num of args for conj")
		return nil
	}
	l1 := args[0]
	l2 := args[1]
	if l1 == nil && l2 == nil {
		log.Error("conj args shold be list, got", nil, nil)
		return nil
	}

	if l1 == nil {
		if l, ok := l2.(*List); ok {
			return l.Copy()
		}
	} else if l2 == nil {
		if l, ok := l1.(*List); ok {
			return l.Copy()
		}
	}

	lst1, ok1 := l1.(*List)
	lst2, ok2 := l2.(*List)
	if !ok1 || !ok2 {
		log.Error("conj args shold be list, got", l1.Type(), l2.Type())
		return nil
	}

	return lst2.Conj(lst1)
}

func doMacro(c *Context, args []Sexpr) Sexpr {
	if len(args) < 2 {
		return nil
	}

	var res Sexpr
	for _, arg := range args[1:] {
		res = arg.Eval(c)
	}
	return res
}

func coreTime(c *Context, args []Sexpr) Sexpr {
	if len(args) != 2 {
		return nil
	}

	start := time.Now()
	res := args[1].Eval(c)
	fmt.Printf("Executed %s in %s\n", args[1], time.Since(start))
	return res
}

func coreRecur(args []Sexpr) Sexpr {
	return &Recur{args}
}

func coreCons(args []Sexpr) Sexpr {
	if len(args) != 2 {
		log.Error("bad num of args for cons")
		return nil
	}
	v := args[0]
	l := args[1]
	if l == nil {
		log.Error("cons first arg shold be list, got", nil)
		return nil
	}

	lst, ok := l.(*List)
	if !ok {
		log.Error("cons first arg shold be list, got", l.Type())
		return nil
	}

	return lst.Add(v)
}

func coreRange(args []Sexpr) Sexpr {
	if len(args) != 1 {
		log.Error("bad num of args for range")
		return nil
	}
	n, ok := args[0].(Number)
	if !ok {
		log.Error("range first arg shold be number, got", args[0].Type())
		return nil
	}

	res := NewList()
	var i Number
	for i = n - 1; i >= 0; i-- {
		res.addFast(i)
	}

	return res
}

func coreOdd(args []Sexpr) Sexpr {
	if len(args) != 1 {
		log.Error("bad num of args for range")
		return nil
	}
	n, ok := args[0].(Number)
	if !ok {
		log.Error("range first arg shold be number, got", args[0].Type())
		return nil
	}

	return Boolean((n % 2) != 0)
}

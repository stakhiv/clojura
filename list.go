package main

import (
	"fmt"
)

type Lister interface {
	Add(Sexpr) *List
	Head() Sexpr
	GetTail() *List
	Length() int
	Conj(*List) *List
}

type Node struct {
	Next *Node
	Val  Sexpr
}

type List struct {
	Node *Node
	Len  int
	Tail *Node
}

func NewList() *List {
	return &List{
		Node: nil,
		Len:  0,
		Tail: nil,
	}
}

func (l *List) Copy() *List {
	if l.Node == nil {
		cl := *l
		return &cl
	}
	n := l.Node
	i := *l.Node
	v := &i
	for n.Next != nil {
		n = n.Next
		c := *n
		v.Next = &c
		v = &c
	}
	return &List{
		Len:  l.Len,
		Node: &i,
		Tail: v,
	}
}

func (l *List) Conj(a *List) *List {
	lc := l.Copy()
	ac := a.Copy()

	if lc.Node == nil {
		return ac
	} else if ac.Node == nil {
		return lc
	}
	lc.Tail.Next = ac.Node
	lc.Len += ac.Len
	lc.Tail = ac.Tail
	return lc
}

func (l *List) addFast(s Sexpr) {
	node := &Node{
		Next: l.Node,
		Val:  s,
	}
	l.Node = node
	l.Len += 1
}

func (l *List) Add(s Sexpr) *List {
	node := &Node{
		Next: l.Node,
		Val:  s,
	}
	tail := l.Tail
	if l.Len == 0 {
		tail = node
	}
	return &List{
		Node: node,
		Len:  l.Len + 1,
		Tail: tail,
	}
}

func (l *List) Head() Sexpr {
	if l.Node != nil {
		return l.Node.Val
	}
	return nil
}
func (l *List) GetTail() *List {
	if l.Node != nil {
		return &List{
			Node: l.Node.Next,
			Len:  l.Len - 1,
			Tail: l.Tail,
		}
	}

	return &List{
		Node: nil,
		Len:  0,
		Tail: nil,
	}
}
func (l *List) Length() int {
	return l.Len
}

func (l *List) Type() CoreType {
	return TypeList
}
func (l *List) Append(s Sexpr) error {
	c := *l
	n := (&c).Add(s)
	*l = *n
	return nil
}
func (l *List) String() string {
	res := "'("
	n := l.Node
	if n != nil {
		res += fmt.Sprintf("%s ", n.Val)
		for n.Next != nil {
			n = n.Next
			res += fmt.Sprintf("%s ", n.Val)
		}
	}
	res += ")"
	return res
}

// TODO: MAKE IMMUTABLE
func (l *List) Eval(c *Context) Sexpr {
	l = l.Copy()
	n := l.Node
	if n != nil {
		if n.Val != nil {
			n.Val = n.Val.Eval(c)
		}

		for n.Next != nil {
			n = n.Next
			if n.Val != nil {
				n.Val = n.Val.Eval(c)
			}
		}
	}
	return l
}
func (l *List) Bool() bool {
	if l.Length() > 0 {
		return true
	}
	return false
}

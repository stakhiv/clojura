package main

import (
	"testing"
)

// func TestListAppend(t *testing.T) {
// 	list := &list{val: &number{1}}
// 	list.Append(&number{2})
//
// 	t.Log(list)
// 	t.Log(list.next)
// }
//
// func TestListCopy(t *testing.T) {
// 	list := NewList(&number{1})
// 	list = list.Add(&number{2})
//
// 	list2 := list.Copy()
//
// 	list.next.val = &number{100}
//
// 	t.Log(list)
// 	t.Log(list2)
// }

func BenchmarkCoreType(b *testing.B) {
	var s Sexpr = Number(1)
	for i := 0; i < b.N; i++ {
		switch t := s.(type) {
		case Number:
			_ = t
			s = Literal("")
		case Literal:
			_ = t
			s = Number(1)
		}
	}
}

func BenchmarkCustomType(b *testing.B) {
	var s Sexpr = Number(1)
	for i := 0; i < b.N; i++ {
		switch s.Type() {
		case TypeNumber:
			_ = s.(Number)
			s = Literal("")
		case TypeLiteral:
			_ = s.(Literal)
			s = Number(1)
		}
	}
}

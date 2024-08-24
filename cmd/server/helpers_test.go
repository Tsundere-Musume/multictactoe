package main

import (
	"testing"
)

func TestStack(t *testing.T) {
	t.Run("StringStack", func(t *testing.T) {
		stk := NewStack[string]()
		stk.push("hello")
		stk.push("world")
		stk.pop()
		stk.push("gophers")
		val, _ := stk.pop()
		assert_eq("gophers", val)
		val, _ = stk.pop()
		assert_eq("hello", val)

	})
	t.Run("IntStack", func(t *testing.T) {
		stk := NewStack[int]()
		stk.push(1)
		stk.push(2)
		stk.pop()
		stk.push(3)
		val, _ := stk.pop()
		assert_eq(3, val)
		val, _ = stk.pop()
		assert_eq(1, val)

	})
}

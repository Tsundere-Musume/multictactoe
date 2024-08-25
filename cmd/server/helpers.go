package main

import (
	"errors"
	"fmt"
	"sync"
	message "tictactoe/internal"
)

func any_f[T comparable](arr []T, predicate func(int, T) bool) bool {
	for idx, val := range arr {
		if predicate(idx, val) {
			return true
		}
	}
	return false
}

type Stack[T any] struct {
	data []T
	mu   *sync.Mutex
}

func (s *Stack[T]) pop() (T, error) {
	var result T
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) <= 0 {
		return result, errors.New("empty stack.")
	}
	l := len(s.data)
	result = s.data[l-1]
	s.data = s.data[:l-1]
	return result, nil
}

func (s *Stack[T]) push(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, value)
}

func (s *Stack[T]) length() int {
	return len(s.data)
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		mu: &sync.Mutex{},
	}
}

func assert_eq[T comparable](left, right T) {
	if left != right {
		panic(fmt.Sprintf("not equal, left = %v and right = %v", left, right))
	}
}

func createServerReponse(gameId string, cmdType uint8, body string) message.Message {
	msg := message.Message{GameId: gameId, Type: message.ServerResponse, CommandType: cmdType, Body: body}
	return msg
}

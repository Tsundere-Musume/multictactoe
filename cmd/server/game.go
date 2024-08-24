package main

import (
	"fmt"
	"net"
	"strconv"
)

const (
	movePlayer1 = 5
	movePlayer2 = 7
)

type Game struct {
	board          [9]int
	currPlayerMove int
	players        [2]net.Conn
	curr           int
	ready          bool
}

func newGame() *Game {
	return &Game{
		currPlayerMove: movePlayer1,
	}
}

func (g *Game) printBoard() {
	for r := range 9 {
		if r%3 == 0 {
			fmt.Println()
		}
		fmt.Printf("%v ", g.board[r])
	}
}
func (g *Game) checkWin() bool {
	sums := []int{}
	for i := range 3 {
		row := i * 3
		sums = append(sums, g.board[row]+g.board[row+1]+g.board[row+2])
		sums = append(sums, g.board[i]+g.board[i+3]+g.board[i+6])
	}
	//diagonals check
	sums = append(sums, g.board[0]+g.board[4]+g.board[8])
	sums = append(sums, g.board[2]+g.board[4]+g.board[6])
	return any_f(sums, func(_ int, val int) bool { return val == movePlayer1*3 || val == movePlayer2*3 })
}

func (g *Game) handleCommand(conn net.Conn, cm []byte) {
	for r := range 9 {
		if r%3 == 0 {
			fmt.Println()
		}
		fmt.Printf("%v ", g.board[r])
	}
	fmt.Println()
	cmd := string(cm)
	fmt.Println(cmd)
	cell, _ := strconv.Atoi(cmd)
	if conn != g.players[g.curr] {
		_, err := conn.Write([]byte("notify:not your turn\n"))
		if err != nil {
			//TODO: server error
			return
		}
		return
	}
	fmt.Printf("%v played %v\n", conn.RemoteAddr(), cell)
	if g.curr%2 == 0 {
		g.board[cell] = movePlayer1
	} else {
		g.board[cell] = movePlayer2
	}
	g.broadcast([]byte(fmt.Sprintf("move:%v\n", cell)))
	if g.checkWin() {
		g.printBoard()
		g.ready = false
		g.writeWin()
	}
	g.curr = (g.curr + 1) % 2
}

func (g *Game) writeWin() {
	for idx, conn := range g.players {
		if idx == g.curr {
			conn.Write([]byte("end:Winner congrats!\n"))
		} else {
			conn.Write([]byte("end:The other player won.\n"))
		}
	}
}

func (s *Game) broadcast(msg []byte) {
	for _, conn := range s.players {
		_, err := conn.Write(msg)
		if err != nil {
			continue
		}
	}
}

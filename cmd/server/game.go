package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
}

func newGame() *Game {
	return &Game{
		currPlayerMove: movePlayer1,
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
	return any(sums, func(_ int, val int) bool { return val == movePlayer1*3 || val == movePlayer2*3 })
}

func (g *Game) start() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%#v\n", g.board)
		if g.checkWin() {
			return
		}

		scanner.Scan()
		cell, err := strconv.Atoi(scanner.Text())
		if err != nil || cell >= 9 || cell < 0 {
			continue
		}
		g.board[cell] = g.currPlayerMove

		if g.currPlayerMove == movePlayer1 {
			g.currPlayerMove = movePlayer2
		} else {
			g.currPlayerMove = movePlayer1
		}

	}
}

func (g *Game) handleCommand(conn net.Conn, cmd string) {
	for r := range 9 {
		if r%3 == 0 {
			fmt.Println()
		}
		fmt.Printf("%v ", g.board[r])
	}
	fmt.Println()
	cell, _ := strconv.Atoi(cmd)
	if conn != g.players[g.curr] {
		_, err := conn.Write([]byte("Not your turn\n"))
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
	if g.checkWin() {
		for r := range 9 {
			if r%3 == 0 {
				fmt.Println()
			}
			fmt.Printf("%v ", g.board[r])
		}
		g.writeWin()
	}
	g.curr = (g.curr + 1) % 2
}

func (g *Game) writeWin() {
	for idx, conn := range g.players {
		if idx == g.curr {
			conn.Write([]byte("Winner congrats!"))
		} else {
			conn.Write([]byte("The other player won."))
		}
	}
}

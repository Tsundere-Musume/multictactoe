package main

import (
	"fmt"
	"net"
	"strconv"
	message "tictactoe/internal"

	"github.com/google/uuid"
)

const (
	movePlayer1 = 5
	movePlayer2 = 7
)

var playerMoves = [2]int{movePlayer1, movePlayer2}

type Game struct {
	id             string
	board          [9]int
	currPlayerMove int
	players        [2]net.Conn
	curr           int
	ready          bool
}

func newGame() *Game {
	return &Game{
		currPlayerMove: movePlayer1,
		id:             uuid.NewString(),
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

func (g *Game) handleCommand(conn net.Conn, msg message.Message) {
	if conn != g.players[g.curr] {
		err := message.Send(conn, createServerReponse(g.id, message.Notify, "Not your turn."))
		if err != nil {
			//TODO: server error
			return
		}
		return
	}

	switch msg.CommandType {
	case message.Move:
		fmt.Println(msg.Body)
		cell, err := strconv.Atoi(msg.Body)
		fmt.Println(cell)
		if err != nil || cell < 0 || cell >= 9 {
			message.Send(conn, createServerReponse(g.id, message.Notify, "Invalid Move"))
			break
		}
		fmt.Printf("%v tried playing %v\n", conn.RemoteAddr(), cell)
		g.board[cell] = playerMoves[g.curr]
		g.broadcast(msg)
		if g.checkWin() {
			g.printBoard()
			g.ready = false
			g.writeWin()
		}
		g.curr = (g.curr + 1) % 2
		fmt.Printf("%v played %v\n", conn.RemoteAddr(), cell)
	}
	g.printBoard()
}

func (g *Game) writeWin() {
	for idx, conn := range g.players {
		if idx == g.curr {
			message.Send(conn, createServerReponse(g.id, message.EndGame, "Winner Congrats!!!"))
		} else {
			message.Send(conn, createServerReponse(g.id, message.EndGame, "Other player wins."))
		}
	}
}

func (g *Game) broadcast(msg message.Message) {
	for _, conn := range g.players {
		err := message.Send(conn, createServerReponse(g.id, message.Move, msg.Body))
		if err != nil {
			continue
		}
	}
}

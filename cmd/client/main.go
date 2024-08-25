package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	message "tictactoe/internal"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	movePlayer1 = 5
	movePlayer2 = 7
)

type Game struct {
	gameId         string
	board          [9]int
	cursor         int
	conn           net.Conn
	cmds           chan message.Message
	currPlayerMove int
	done           bool
	message        string
	ready          bool
}

func (m *Game) listen() {
	defer m.conn.Close()
	dec := json.NewDecoder(m.conn)
	for {
		var resp message.Message
		err := dec.Decode(&resp)
		if err != nil {
			if err == io.EOF {
				panic("Connection closed from server\n")
			}
			return
		}
		m.cmds <- resp
	}
}

type responseMsg message.Message
type quitMsg struct{}

func waitForActivity(sub chan message.Message) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func handleCommand(m Game, resp message.Message) (Game, tea.Cmd) {
	if resp.Type != message.ServerResponse {
		return m, waitForActivity(m.cmds)
	}
	if !m.ready && resp.Type == message.GameReady {
		m.gameId = resp.GameId
		m.ready = true
		m.message = "Game is ready."
		return m, waitForActivity(m.cmds)
	}
	if m.ready {
		switch resp.CommandType {
		case message.Move:
			idx, err := strconv.Atoi(resp.Body)
			if err != nil {
				//log this out
				panic("Server error: malformed command.")
			}
			m.board[idx] = m.currPlayerMove
			m.currPlayerMove = -(m.currPlayerMove - movePlayer1 - movePlayer2)
		case message.EndGame:
			m.message = resp.Body
			m.done = true
			return m, tea.Quit

		case message.Notify:
			m.message = resp.Body
		}
	}
	return m, waitForActivity(m.cmds)

}

func (m Game) Init() tea.Cmd {
	message.Send(m.conn, createServerCmd(m.gameId, message.Join, ""))
	return waitForActivity(m.cmds)
}

func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "l":
			m.cursor += 1
		case "h":
			m.cursor -= 1
		case "j":
			m.cursor += 3
		case "k":
			m.cursor -= 3
		case " ":
			err := message.Send(m.conn, createGameCmd(m.gameId, message.Move, fmt.Sprint(m.cursor)))
			if err != nil {
				fmt.Println("Error writing to the server", err)
			}

		}

	case responseMsg:
		m.message = ""
		m, cmd = handleCommand(m, message.Message(msg))
		cmds = append(cmds, cmd)

	case quitMsg:
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

func (m Game) View() string {
	s := m.message + "\n\n"
	for idx := range 9 {
		style := base
		if idx%3 == 0 && idx != 0 {
			s += "\n"
			s += borderStyle.Render("━━━╋━━━╋━━━")
			s += "\n"
		}
		if idx == m.cursor {
			style = style.Background(valid)
		}

		if idx%3 != 0 {
			s += borderStyle.Render("┃")
		}
		if m.board[idx] == 0 {
			s += style.Render(" ")
		} else if m.board[idx] == movePlayer1 {
			s += style.Render("O")
		} else {
			s += style.Render("X")
		}

	}
	s += "\n\n"
	return s
}

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Fatal("Couldn't connect to the game server\n")
	}
	game := Game{cmds: make(chan message.Message), conn: conn, currPlayerMove: movePlayer1}
	go game.listen()
	p := tea.NewProgram(game)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

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
	cmds           chan string
	currPlayerMove int
	done           bool
	message        string
	ready          bool
}

func (m *Game) listen() {
	defer m.conn.Close()
	reader := bufio.NewReader(m.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("Connection closed from the server")
			}
			return
		}
		msg = strings.TrimSpace(msg)
		m.cmds <- msg
	}

}

type responseMsg string
type quitMsg struct{}

func waitForActivity(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func handleCommand(m Game, resp string) (Game, tea.Cmd) {
	if !m.ready && resp[:1] == "0" {
		m.gameId = resp[1:]
		m.ready = true
		m.message = "Game is ready."
		return m, waitForActivity(m.cmds)
	}
	if m.ready {
		out := strings.Split(resp, ":")
		if len(out) < 2 {
			panic(fmt.Sprintf("Couldn't parse command, %s", out[0]))
		}
		switch out[0] {
		case "move":
			idx, err := strconv.Atoi(out[1])
			if err != nil {
				panic("Server error: malformed command.")
			}
			m.board[idx] = m.currPlayerMove
			m.currPlayerMove = -(m.currPlayerMove - movePlayer1 - movePlayer2)
		case "end":
			m.message = out[1]
			m.done = true
			return m, tea.Quit

		case "notify":
			m.message = out[1]
		}
	}
	return m, waitForActivity(m.cmds)

}

func (m Game) Init() tea.Cmd {
	m.conn.Write([]byte("0\n"))
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
			_, err := m.conn.Write([]byte(fmt.Sprintf("1%v%v\n", m.gameId, m.cursor)))
			if err != nil {
				fmt.Println("Error writing to the server", err)
			}

		}

	case responseMsg:
		m.message = ""
		m, cmd = handleCommand(m, string(msg))
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
	game := Game{cmds: make(chan string), conn: conn, currPlayerMove: movePlayer1}
	go game.listen()
	p := tea.NewProgram(game)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	message "tictactoe/internal"
)

type GameServer struct {
	game    map[string]*Game
	clients []net.Conn
	mu      sync.Mutex
	ready   bool
	pool    chan net.Conn
}

func newSever() *GameServer {
	return &GameServer{
		game: make(map[string]*Game),
		pool: make(chan net.Conn, 10),
	}
}

func (s *GameServer) createGame() {
	for {
		player1 := <-s.pool
		player2 := <-s.pool
		game := newGame()
		s.mu.Lock()
		game.players[0] = player1
		game.players[1] = player2
		s.game[game.id] = game
		s.mu.Unlock()
		//TODO: handle errors
		game.broadcast(createServerReponse(game.id, message.GameReady, ""))
	}
}

func (s *GameServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	go s.createGame()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting the connection ", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *GameServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		var m message.Message
		err := dec.Decode(&m)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed from %v\n", conn.RemoteAddr())
			}
			log.Println(err)
			return
		}
		log.Printf("Got %#v from %v", m, conn.RemoteAddr())
		s.handleCommand(conn, m)
	}
}

func (s *GameServer) handleCommand(conn net.Conn, msg message.Message) {
	switch msg.Type {
	case message.ServerCommand:
		s.pool <- conn
	case message.GameCommand:
		game, ok := s.game[msg.GameId]
		if !ok {
			return
		}
		game.handleCommand(conn, msg)
	}
}

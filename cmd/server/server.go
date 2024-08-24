package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
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
		game: make(map[string]*Game, 10),
		pool: make(chan net.Conn, 10),
	}
}

func (s *GameServer) createGame() {
	for {
		player1 := <-s.pool
		player2 := <-s.pool
		gameId := uuid.NewString()
		game := newGame()
		s.mu.Lock()
		game.players[0] = player1
		game.players[1] = player2
		s.game[gameId] = game
		s.mu.Unlock()
		//TODO: handle errors
		game.broadcast([]byte("0" + gameId + "\n"))
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
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("Error reading from the connection, ", err)
			}
			return
		}
		msg = bytes.TrimSpace(msg)
		log.Printf("Got %s from %v", string(msg), conn.RemoteAddr())
		s.handleCommand(conn, msg)
	}
}

func (s *GameServer) handleCommand(conn net.Conn, msg []byte) {
	switch msg[0] {
	case '0':
		s.pool <- conn
	case '1':
		if len(msg) < 38 {
			return
		}
		gameId := msg[1:37]
		game, ok := s.game[string(gameId)]
		if !ok {
			return
		}
		game.handleCommand(conn, msg[37:])
	}
}

package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type GameServer struct {
	game    *Game
	clients []net.Conn
	mu      sync.Mutex
	ready   bool
}

func newSever() *GameServer {
	return &GameServer{
		game: newGame(),
	}
}

func (s *GameServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting the connection ", err)
			continue
		}
		s.mu.Lock()
		s.clients = append(s.clients, conn)
		if len(s.clients) == 2 {
			s.ready = true
			s.game.players[0] = s.clients[0]
			s.game.players[1] = s.clients[1]
		}
		s.mu.Unlock()
		go s.handleConnection(conn)
	}
}

func (s *GameServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	for func() bool { return !s.ready }() {
	}
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("Error reading from the connection, ", err)
			}
			return
		}
		msg = strings.TrimSpace(msg)
		log.Printf("Got %s from %v", msg, conn.RemoteAddr())

		s.game.handleCommand(conn, msg)
	}
}

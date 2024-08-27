package message

import (
	"encoding/json"
	"net"
)

type Message struct {
	GameId      string `json:"game_id"`
	Type        uint8  `json:"type"`
	CommandType uint8  `json:"command_type"`
	Body        string `json:"command"`
}

const (
	ServerCommand uint8 = iota + 1
	GameCommand

	ServerResponse
)

const (
	Notify uint8 = iota + 1
	Move
	GameReady
	EndGame
	Join
	Create
)

func Send(conn net.Conn, msg Message) error {
	data, err := Encode(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

func Encode(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func Decode(msg []byte) (*Message, error) {
	var result Message
	err := json.Unmarshal(msg, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

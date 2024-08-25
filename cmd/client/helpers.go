package main

import message "tictactoe/internal"

func createGameCmd(gameId string, cmdType uint8, body string) message.Message {
	msg := message.Message{GameId: gameId, Type: message.GameCommand, CommandType: cmdType, Body: body}
	return msg
}

func createServerCmd(gameId string, cmdType uint8, body string) message.Message {
	msg := message.Message{GameId: gameId, Type: message.ServerCommand, CommandType: cmdType, Body: body}
	return msg
}

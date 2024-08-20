package main

import "log"

func main() {
	// clients = make(map[net.Conn]struct{})
	// listener, err := net.Listen("tcp", ":8000")
	// if err != nil {
	// 	log.Fatalln("Error starting server:", err)
	// }
	//
	// defer listener.Close()
	//
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Println("error accepting", err)
	// 		continue
	// 	}
	// 	clients[conn] = struct{}{}
	// 	go handleConnection(conn)
	// }

	// game := newGame()
	// game.start()
	server := newSever()
	log.Fatal(server.Start(":8000"))
}

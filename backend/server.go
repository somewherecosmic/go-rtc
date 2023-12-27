package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	// alternatively the websocket remote addr pointing to the connection pointer
	conns map[*websocket.Conn]bool
}

func newServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWebsocketConnection(ws *websocket.Conn) {
	fmt.Println("New incoming connection from address:", ws.RemoteAddr())

	// maps in golang aren't safe - should use some kind of mutex here in future
	// eliminate race conditions
	s.conns[ws] = true

	s.readLoop(ws)

}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf) // returns int n, marker of how far the frame data goes into the buffer
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read Error: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("Thank you for the message"))
	}
}

func main() {
	server := newServer()
	http.Handle("/ws", websocket.Handler(server.handleWebsocketConnection))
	http.ListenAndServe(":80", nil)
	fmt.Println("Websocket Server Running on Port 80")
}

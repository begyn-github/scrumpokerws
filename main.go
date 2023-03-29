package main

import (
	"fmt"
	"lesson1/statemachine"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "WebSocket: /ws")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	root := statemachine.GetRoot()
	var userState statemachine.UserState = statemachine.UserState{User: "", ActualState: &root}

	if ws != nil {
		err = ws.WriteMessage(1, []byte(strings.Join(userState.ActualState.GetMenu(), "\n")))
		if err != nil {
			log.Println(err)
		}
	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws, userState)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn, us statemachine.UserState) {
	toExit := false

	for !toExit {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		fmt.Println("From Client: " + string(p))

		state, err2 := us.ActualState.GoTo(string(p))

		if err2 == nil {
			us.ActualState = &state
			fmt.Println(strings.Join(state.GetMenu(), "\n"))
		}

		if err := conn.WriteMessage(messageType, []byte(strings.Join(state.GetMenu(), "\n"))); err != nil {
			log.Println(err)
			return
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Hello World")
	fmt.Println(port)
	setupRoutes()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

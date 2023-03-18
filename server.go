package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func challenge(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ccDjjFl80ofw71TxqR2KK5Iva-vB7diFWothQXXQYhw.RrcB1hbRA1D23gYXZu1iOPpRTFX-9ovExmpGXFhC9lw")
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
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/.well-known/acme-challenge/ccDjjFl80ofw71TxqR2KK5Iva-vB7diFWothQXXQYhw", challenge)
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	//port := os.Getenv("PORT")

	port := "8080"
	fmt.Println("Hello World")
	fmt.Println(port)
	setupRoutes()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

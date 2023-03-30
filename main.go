package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"scrumpokerws/statemachine"
	"strings"

	"github.com/gorilla/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "WebSocket: /ws")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// transforma conexao em ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	// Recebe o root da state machine
	// e define o estado inicial do usuário
	root := statemachine.GetRoot()
	var userState statemachine.UserState = statemachine.UserState{User: "", ActualState: &root}

	// Escreve o menu para o usuário, logo depois de conectado
	if ws != nil {
		err = ws.WriteMessage(1, []byte(strings.Join(userState.ActualState.GetMenu(), "\n")))
		if err != nil {
			log.Println(err)
		}
	}

	// escuta indefinidamente por novas mensagens pela conexao
	// e passa o state do usuario conectado
	reader(ws, userState)
}

// leitor de words para avanco na state machine
func reader(conn *websocket.Conn, us statemachine.UserState) {
	toExit := false

	for !toExit {
		// le entradas do usuário
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// imprime em tela mensagem do usuário
		fmt.Println("From Client: " + string(p))

		// tenta avançar para o próximo estado
		// de acordo com escolha no menu
		state, err2 := us.ActualState.GoTo(string(p))
		if err2 == nil {
			us.ActualState = &state
			fmt.Println(strings.Join(state.GetMenu(), "\n"))
		}

		// envia o menu do novo estado alcançado ao usuário
		if err := conn.WriteMessage(messageType, []byte(strings.Join(state.GetMenu(), "\n"))); err != nil {
			log.Println(err)
			return
		}
	}
}

// Cadastra as rotas no servidor
func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

// Necessitamos de Read and Write buffers
// sizes para definir o Upgrader
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

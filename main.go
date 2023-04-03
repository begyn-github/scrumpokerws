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
	var userState statemachine.UserState = statemachine.UserState{ActualState: &root}

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

func readOnce(conn *websocket.Conn) (int, string) {
	messageType, p, err := conn.ReadMessage()

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	return messageType, string(p)
}

// leitor de words para avanco na state machine
func reader(conn *websocket.Conn, us statemachine.UserState) {
	toExit := false
	var (
		messageType int
		word        string
	)

	for !toExit {
		// le entradas do usuário
		messageType, word = readOnce(conn)

		// imprime em tela mensagem do usuário
		fmt.Println("From Client: " + word)

		// tenta avançar para o próximo estado
		// de acordo com escolha no menu
		state, err := us.ActualState.GoTo(strings.ToUpper(word))
		if err == nil {
			us.ActualState = &state
		}

		if us.ActualState.Field != "" {
			clientMsg := "Set " + us.ActualState.Field + " (" + us.GetDataValue() + ") => "

			if err := conn.WriteMessage(messageType, []byte(clientMsg)); err != nil {
				log.Println(err)
				return
			}

			messageType, word = readOnce(conn)

			if word != "" {
				us.UpdateDataValue(word)
			}
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

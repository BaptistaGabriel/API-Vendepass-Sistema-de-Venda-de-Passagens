package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type Message struct {
	Type    string `json:"type"`
	Content interface{} `json:"content"`
}

func receiveMessage(connection net.Conn) Message {
	var msg Message
	decoder := json.NewDecoder(connection)
	if err := decoder.Decode(&msg); err != nil {
		fmt.Printf("Erro ao receber mensagem do servidor %v\n", err)
		return Message{Type: "error"}
	}
	fmt.Printf("Mensagem recebida do servidor: %+v\n", msg)
	return msg
}

func sendMessage(connection net.Conn, msg Message) {
	encoder := json.NewEncoder(connection)
	if err := encoder.Encode(msg); err != nil {
		fmt.Printf("Erro ao mandar a mensagem para o servidor %v\n", err)
		return
	}
	fmt.Println("Mensagem enviada para o servidor")
}

func firstMenu(connection net.Conn) {
	var option int
	for {
		fmt.Println("==========================")
		fmt.Printf("\033[34m|     1. Fazer login     |\n|     2. Criar conta     |\033[0m\n")
		fmt.Println("==========================")
		fmt.Scanln(&option)

		sendMessage(connection, Message{Type: "action", Content: option})

		switch option {
		case 1:
			var numberID int
			fmt.Println("Número de identificação do cliente: ")
			fmt.Scanln(&numberID)
			sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(numberID)})

			responseMsg := receiveMessage(connection)
			if responseMsg.Content == "-1" {
				fmt.Println("Usuário não cadastrado.")
			} else {
				fmt.Printf("Olá, %v \n", responseMsg.Content)
				return
			}
		case 2:
			var name string
			fmt.Printf("Nome: ")
			fmt.Scanln(&name)
			sendMessage(connection, Message{Type: "action", Content: name})

			numberIDMsg := receiveMessage(connection)
			fmt.Printf("Número da sua conta: %v \n", numberIDMsg.Content)
		default:
			fmt.Println("Opção inválida!")
		}
	}
}

func secondMenu(connection net.Conn) {
	var option int
	for {
		fmt.Println("==========================")
		fmt.Printf("\033[34m|     1. Comprar passagem     |\n|     2. Cancelar passagem    |\n|     3. Sair                |\033[0m\n")
		fmt.Println("==========================")
		fmt.Scanln(&option)

		sendMessage(connection, Message{Type: "action", Content: option})

		switch option {
		case 1:
			// Comprar passagem
			routesMsg := receiveMessage(connection)
			fmt.Println("Rotas disponíveis: ")
			for i, route := range routesMsg.Content.([]interface{}) {
				fmt.Printf("%d. %v\n", i, route)
			}

			var routeNumber int
			fmt.Println("Escolha o número da rota: ")
			fmt.Scanln(&routeNumber)
			sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(routeNumber)})

			seatsMsg := receiveMessage(connection)
			fmt.Println("Assentos disponíveis: ")
			for i, seat := range seatsMsg.Content.([]interface{}) {
				fmt.Printf("Assento %d - Reservado: %v\n", i, seat)
			}

			var seatNumber int
			fmt.Println("Escolha o número do assento: ")
			fmt.Scanln(&seatNumber)
			sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(seatNumber)})

			responseMsg := receiveMessage(connection)
			fmt.Println(responseMsg.Content)

		case 2:
			// Cancelar passagem
			clientFlightsMsg := receiveMessage(connection)
			fmt.Println("Suas passagens: ")
			for i, flight := range clientFlightsMsg.Content.([]interface{}) {
				fmt.Printf("%d. %v\n", i, flight)
			}

			var flightIndex, seatIndex int
			fmt.Println("Escolha o número do voo para cancelar: ")
			fmt.Scanln(&flightIndex)
			sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(flightIndex)})

			fmt.Println("Escolha o número do assento para cancelar: ")
			fmt.Scanln(&seatIndex)
			sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(seatIndex)})

			responseMsg := receiveMessage(connection)
			fmt.Println(responseMsg.Content)

		case 3:
			// Sair
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
}


func main() {
	connection, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	firstMenu(connection)
	secondMenu(connection)
}
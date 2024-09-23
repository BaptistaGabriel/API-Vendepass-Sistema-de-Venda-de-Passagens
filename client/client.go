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
		fmt.Printf("\033[34m|  1. Comprar passagens  |\n|  2. Cancelar passagens |\n|  3. Sair               |\033[0m\n")
		fmt.Println("==========================")
		fmt.Scanln(&option)

		sendMessage(connection, Message{Type: "action", Content: option})

		switch option {
		case 1:
			exit := true
			for exit {
				listMsg := receiveMessage(connection)
				if listMsg.Type != "list" {
					continue
				}
				list := listMsg.Content.([]interface{})
				for index, item := range list {
					fmt.Printf("%d.........%s\n", index+1, item)
				}
				var routeNumber int
				fmt.Scanln(&routeNumber)
				sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(routeNumber - 1)})

				listSeatsMsg := receiveMessage(connection)
				if listSeatsMsg.Type != "list" {
					continue
				}
				listSeats := listSeatsMsg.Content.([]interface{})
				for index, seat := range listSeats {
					if seat.(string) == "false" {
						fmt.Printf("%d......... Livre\n", index+1)
					} else {
						fmt.Printf("%d......... Ocupado\n", index+1)
					}
				}

				var chooseSeat int
				fmt.Scanln(&chooseSeat)
				sendMessage(connection, Message{Type: "action", Content: strconv.Itoa(chooseSeat - 1)})
				fmt.Println(receiveMessage(connection).Content)

				fmt.Println("===================================")
				fmt.Printf("\033[34m|     1. Comprar outra passagem     |\n|     2. Sair                      |\033[0m\n")
				fmt.Println("===================================")
				fmt.Scanln(&option)
				sendMessage(connection, Message{Type: "action", Content: option})
				if option == 2 {
					exit = false
				}
			}
		case 2:
			fmt.Println("------------Cancelando------------")
		case 3:
			fmt.Println("Obrigada por comprar com a gente!")
			return
		default:
			fmt.Println("Opção inválida")
			return
		}
	}
}

func main() {
	connection, err := net.Dial("tcp", "172.16.103.12:7777")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	firstMenu(connection)
	secondMenu(connection)
}
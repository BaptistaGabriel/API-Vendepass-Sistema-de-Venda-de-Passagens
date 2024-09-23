package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

func receiveMessage(connection net.Conn) string {

	// Recebendo mensagem do servidor
	buffer := make([]byte, 1024)

	// Lendo resposta do servidor
	size_bytes, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber mensagem do servidor %v\n", err)
		return "-1"
	}

	message := string(buffer[:size_bytes])

	// Mostrando a mensagem
	fmt.Printf("Mensagem recebida do servidor: %v\n", message)

	return message
}

func sendMessage(connection net.Conn, message string) {

	// Mandando mensagem para o servidor
	_, err := connection.Write([]byte(message))
	if err != nil {
		fmt.Printf("Erro ao mandar a mensagem para o servidor %v\n", err)
		return
	}
	fmt.Println("Mensagem enviada para o servidor")
}

func receiveJSON(connection net.Conn) [] string{
	buffer := make([]byte, 1024)
	var list [] string

	size_bytes, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber o json: %v\n", err)
		return list
	}

	
	err = json.Unmarshal(buffer[:size_bytes], &list)
	if err != nil {
		fmt.Printf("Erro ao converter o JSON para lista: %v\n", err)
		return list
	}
	return list
}

func firstMenu(connection net.Conn) {

	var option int

	for {
		fmt.Println("==========================")
		fmt.Printf("\033[34m|     1. Fazer login     |\n|     2. Criar conta     |\033\n[0m")
		fmt.Println("==========================")
		fmt.Scanln(&option)
		sendMessage(connection, strconv.Itoa(option))

		switch option {

		case 1:
			// Fazer login

			var numberID int
			fmt.Println("Número de identificação do cliente: ")
			fmt.Scanln(&numberID)
			sendMessage(connection, strconv.Itoa(numberID))
			message := receiveMessage(connection)
			_, err := strconv.Atoi(message)
			if err == nil {
				fmt.Println("Usuário não cadastrado.")
			} else {
				// Retornar o nome do cliente
				fmt.Printf("Olá, %v \n", message)
				return
			}

		case 2:
			// Cadastrar cliente
			var name string
			fmt.Printf("Nome: ")
			fmt.Scanln(&name)
			sendMessage(connection, name)
			numberID := receiveMessage(connection)
			fmt.Printf("Número da sua conta: %v \n", numberID)

		default:
			fmt.Println("Opção inválida!")
		}
	}
}

func secondMenu(connection net.Conn, ) {
	var option int

	for {
		fmt.Println("==========================")
		fmt.Printf("\033[34m|  1. Comprar passagens  |\n|  2. Cancelar passagans |\n|  3. Sair               |\033\n[0m")
		fmt.Println("==========================")
		fmt.Scanln(&option)
		sendMessage(connection, strconv.Itoa(option))

		switch option {
		case 1: 
			exit := true
			for exit {
				// Comprar passagens
				list := receiveJSON(connection)
				for index, item := range list {
					fmt.Printf("%d.........%s\n", index + 1, item)
				}
				var route_number int
				fmt.Scanln(&route_number)
				sendMessage(connection, strconv.Itoa(route_number))
				list_seats := receiveJSON(connection)
				for index, seat := range list_seats {
					if seat == "false" {
						fmt.Printf("%d......... Livre\n", index + 1)
					} else {
						fmt.Printf("%d......... Ocupado\n", index + 1)
					}
				}

				var choose_seat int

				fmt.Scanln(&choose_seat)
				sendMessage(connection, strconv.Itoa(choose_seat - 1))
				fmt.Println(receiveMessage(connection))

				fmt.Println("===================================")
				fmt.Printf("\033[34m|     1. Comprar outra passagem     |\n|     2. Sair                      |\033\n[0m")
				fmt.Println("===================================")
				fmt.Scanln(&option)
				sendMessage(connection, strconv.Itoa(option))
				if option == 2 {
					exit = false
				}
			}

		case 2:
			fmt.Println("------------Cancelando------------")
		case 3:
			// Sair
			fmt.Println("Obrigada por comprar com a gente!")
			return
		default:
			fmt.Println("Opção inválida")
		}
	}
}

func main() {

	// Conectando com o servidor
	connection, err := net.Dial("tcp", "192.168.1.3:8080")
	if err != nil {
		fmt.Printf("Erro ao conectar com o servidor %v\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Conectado ao servidor!")

	firstMenu(connection)
	secondMenu(connection)

}

package main

import (
	"fmt"
	"net"
	 "strconv"
)

func receiveMessage(connection net.Conn) string {

	// Recebendo mensagem do servidor
	buffer := make([]byte, 1024)

	// Lendo resposta do servidor
	tam_bytes, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber mensagem do servidor %v\n", err)
		return "-1"
	}

	message := string(buffer[:tam_bytes])

	// Mostrando a mensagem
	fmt.Printf("Mensagem recebida do servidor: %v\n", message)

	return message
}

func sendMessage(connection net.Conn, mensage string) {

	// Mandando mensagem para o servidor
	message := mensage
	_, err := connection.Write([]byte(message))
	if err != nil {
		fmt.Printf("Erro ao mandar a mensagem para o servidor %v\n", err)
		return
	}
	fmt.Println("Mensagem enviada para o servidor")
}

func firstMenu() (string, int) {

	var option int
	
	for {
		fmt.Println("==========================")
		fmt.Printf("\033[34m|     1. Fazer login     |\n|     2. Criar conta     |\033\n[0m")
		fmt.Println("==========================")
		fmt.Scanln(&option)
		
		switch option {
		
		case 1:
			var numberID int
			fmt.Println("Número de identificação do cliente: ")
			fmt.Scanln(&numberID)

			return "", numberID
		
		case 2:
			var name string
			fmt.Printf("Nome: ")
			fmt.Scanln(&name)
			return name, -1
		default:
			fmt.Println("Opção inválida!")
		}
	}
}

func main() {
	
	// Conectando com o servidor
	connection, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("Erro ao conectar com o servidor %v\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Conectado ao servidor!")

	nameClient, numberID := firstMenu()

	if numberID == -1 {
		// Criar conta
		sendMessage(connection, nameClient)
		
		numberID, _ := strconv.Atoi(receiveMessage(connection))
		nameClient, numberID := firstMenu()

	} else {
		// Fazer login

    numberIDstg := strconv.Itoa(numberID) 

		sendMessage(connection, numberIDstg)
		confirmation := receiveMessage(connection)
		
		// Mostrando a confirmação
		fmt.Println(confirmation)

		// Ir para o menu 2 ////////////////
	}

}

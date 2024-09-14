package main

import (
	"fmt"
	"net"
)

func receiveMessage(connection net.Conn) {

	// Recebendo mensagem do servidor
	buffer := make([]byte, 1024)

	// Lendo resposta do servidor
	tam_bytes, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Erro ao receber mensagem do servidor %v\n", err)
		return
	}

	message := string(buffer[:tam_bytes])

	// Mostrando a mensagem
	fmt.Printf("Mensagem recebida do servidor: %v\n", message)
}

func returnMessage(connection net.Conn) {

	// Mandando mensagem para o servidor
	message := "Mensagem para o servidor: Olá server!"
	_, err := connection.Write([]byte(message))
	if err != nil {
		fmt.Printf("Erro ao mandar a mensagem para o servidor %v\n", err)
		return
	}
	fmt.Println("Mensagem enviada para o servidor")
}

func communication(connection net.Conn) {
	defer connection.Close()

	returnMessage(connection)
	receiveMessage(connection)
	
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

	communication(connection)
}

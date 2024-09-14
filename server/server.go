package main

import (
	"fmt"
	"net"
)

func receiveMessage(connection net.Conn) {
	// Criando um buffer
	buffer := make([]byte, 1024)

	// Lendo dados do cliente
	message, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Erro em receber a mensagem %v\n", err)
		return
	}

	// Mostrando a mensagem
	fmt.Printf("Mensagem recebida: %s\n", buffer[:message])

}

func returnMessage(connection net.Conn) {
	// Mandando mensagem para o cliente
	data := []byte("Servidor respondendo...")
	_, err := connection.Write(data)
	if err != nil {
		fmt.Printf("Erro ao mandar a resposta para o cliente %v\n", err)
		return
	}
	fmt.Println("Resposta devolvida\n\n")
}

func communication(connection net.Conn) {
	defer connection.Close()

	receiveMessage(connection)
	returnMessage(connection)
}

func main() {

	// Criando o servidor na porta 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor funcionando na porta 8080...\n")

	// Aceitando conexões em loop
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}
		fmt.Println("Recebendo mensagen...\n")

		go communication(connection)
	}
}

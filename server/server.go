package main

import (
	"fmt"
	"net"
	"log"
)

func receiveMessage(connection net.Conn) string {
	// Recebendo mensagem do cliente

	// Criando um buffer
	buffer := make([]byte, 1024)

	// Lendo dados do cliente
	size_bytes, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Erro em receber a mensagem do cliente %v\n", err)
		return "-1"
	}

	message := string(buffer[:size_bytes])

	// Mostrando a mensagem
	fmt.Printf("Mensagem recebida do cliente: %s\n", message)
	return message
}

func sendMessage(connection net.Conn) {
	// Mandando mensagem para o cliente

	// Mensagem
	data := []byte("Servidor respondendo...")
	_, err := connection.Write(data)
	if err != nil {
		fmt.Printf("Erro ao mandar a resposta para o cliente %v\n", err)
		return
	}
	fmt.Println("Resposta devolvida para o cliente")
}

func communication(connection net.Conn, mapClients *map[int] string) {
	defer connection.Close()

	exit := true
	// Menu 1
	for exit {		
		option := receiveMessage(connection)
		
		// Fazer login
		if option == "1" {
			//
		}
		
		
		// Cadastrar
	}

	receiveMessage(connection)
	sendMessage(connection)
}

func getLocalIP() net.IP {
	// Fazendo uma conexão não efetiva com o servidor DNS da Google 
	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil{
		log.Fatal(err)
	}
	defer connection.Close()

	localAddress := connection.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func main() {

	// Criando lista de clientes
	mapClients := make(map[int]string)
	
	// Pegando o IP do servidor 
	fmt.Printf("IP do servidor %v\n", getLocalIP())

	// Criando o servidor na porta 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor funcionando na porta 8080...")

	// Aceitando conexões em loop
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}
		fmt.Println("Recebendo mensagen...")

		go communication(connection, &mapClients)
	}
}

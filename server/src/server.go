package main

import (
	"fmt"
	"net"
	"log"
	"strconv"
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

func sendMessage(connection net.Conn, message string) {
	// Mandando mensagem para o cliente

	// Mensagem
	data := []byte(message)
	_, err := connection.Write(data)
	if err != nil {
		fmt.Printf("Erro ao mandar a resposta para o cliente %v\n", err)
		return
	}
	fmt.Println("Resposta devolvida para o cliente")
}

func communication(connection net.Conn, mapClients map[int] string, flights []Flight) {
	defer connection.Close()
	exit := true

	// Menu 1
	for exit {		
		option := receiveMessage(connection)
		
		// Fazer login
		if option == "1" {
			numberID,_ := strconv.Atoi(receiveMessage(connection))
			name, exists := mapClients[numberID]
			fmt.Printf("NOMEEEE DO CLIENTE: %v", name)
			if exists{
				sendMessage(connection, name)
				exit = false
			} else {
				sendMessage(connection, "-1")
			}
		// Cadastrar
		} else if option == "2" {
			name := receiveMessage(connection)
			sendMessage(connection, strconv.Itoa(createClient(name, mapClients)))
		} else {
			// Retornar se o cliente cair no primeiro menu
			if option != "0" {
				fmt.Println("SE O CLIENTE CAIR NO PRIMEIRO MENU")
				return
			}
		}
	}

	// Menu 2
	for {
		option := receiveMessage(connection)
		if option == "1"{
			fmt.Println("Finge que está comprando")
			exit := true
			for exit {
				routes := GetRoutes(flights)
				fmt.Println(routes)
				// Logica da compra aqui dentro
				// Lembrar de colocar tudo no nome do cliente para saber quem comprou
				fmt.Println("Finge que está mostrando as rotas aqui tá ligado")
				exit = false
			}
		} else if option == "2" {
			// Cancelar passagem mostrar tudo que ele comprou, só as passagens ativas
			fmt.Println("Finge que está cancelando")
		} else if option == "3" {
			fmt.Println("Saindooooooo")
			return
		} else {
			// Retornar se o cliente cair no primeiro menu
			if option != "0" {
				fmt.Println("Se o cliente cair!! No segundo menu claro")
				return
			}
		}
	}
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

func createClient(name string, mapClients map[int] string) int{
	number := len(mapClients) + 1
	mapClients[number] = name
	return number
}

func main() {

	// Criando lista de clientes
	mapClients := make(map[int]string)

	// Criando nome do arquivo das rotas
	var flight_file string
	
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

	// Criando e salvando as rotas em um arquivo
	flight := CreateRoutes()
	message := SaveFlightsToFile(flight_file, flight)
	
	// Se der erro
	if message != "" {
		return
	}

	// Aceitando conexões em loop
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}
		fmt.Println("Recebendo mensagen...")

		go communication(connection, mapClients, flight)
	}
}

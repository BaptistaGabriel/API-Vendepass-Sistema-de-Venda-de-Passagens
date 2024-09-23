package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

func sendJSON(connection net.Conn, list []string){	
	// Converter para json
	json_flight, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("Erro ao converter para JSON: %v\n", err)
		return 
	}
	
	// Mandando um json para o cliente
	_, err = connection.Write(json_flight)
	if err != nil {
		fmt.Printf("Erro ao enviar o json para o cliente %v\n", err)
		return
	}
	fmt.Println("Lista enviada ao cliente!")
}

func communication(connection net.Conn, flights []Flight) {
	defer connection.Close()
	exit := true

	var numberID int

	mapClients := getClients()

	// Menu 1
	for exit {		

		option := receiveMessage(connection)
		
		// Fazer login
		if option == "1" {
			numberID,_ = strconv.Atoi(receiveMessage(connection))
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
			saveClient(mapClients)
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
			/////////////////// PROTEÇÃO CONTRA CAIR O CLIENTE ////////////////
			fmt.Println("Finge que está comprando")
			exit := true
			for exit {
				routes := GetRoutes(flights)
				sendJSON(connection, routes)
				
				route_number, _ := strconv.Atoi(receiveMessage(connection))
				client_flight := flights[route_number]
				seats := GetSeats(client_flight)

				var list_seats[] string 

				for _, seat := range seats {
					item := strconv.FormatBool(seat.IsReserved)
					list_seats = append(list_seats, item)
				}

				sendJSON(connection, list_seats)
				seat_number, _ := strconv.Atoi(receiveMessage(connection))
				if (ReserveSeat(flights, route_number, seat_number, strconv.Itoa(numberID))) {
					sendMessage(connection, "Assento comprado com sucesso!")
				} else {
					sendMessage(connection, "Erro ao comprar o assento")
				}

				err := SaveFlightsToFile("flights", flights)
				if err != "" {
					fmt.Println("Erro ao salvar a operação!")
				}

				option = receiveMessage(connection)
				if option == "2" {
					exit = false 
				}

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

func getClients() map[int]string {
	
	// Criando lista de clientes
	mapClients := make(map[int]string)

	file, err := os.Create("data/clients.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo dos clientes")
		return nil
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo dos clientes")
		return nil
	}

	err = json.Unmarshal(bytes, &mapClients)
	if err != nil {
		fmt.Println("Erro listar os clientes")
		return nil
	}
	defer file.Close()

	return mapClients
}

func saveClient(mapClients map[int] string) {
	// Garante que a pasta 'data' exista
	os.MkdirAll("data", os.ModePerm)	

	file, err := json.MarshalIndent(mapClients, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}
	
	err = os.WriteFile("data/clients.json", file, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar arquivo: ", err)
		return
	}

	return
}

func main() {
	
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
	// Nome do arquivo
	flight := CreateRoutes()
	message := SaveFlightsToFile("flights.json", flight)
	
	// Se der erro
	if message != "" {
		return
	}

	flight, err = LoadFlightsFromFile("flights.json")
	if err != nil {
		fmt.Printf("Erro ao carregar rotas do arquivo %v:", err)
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

		go communication(connection, flight)
	}
}

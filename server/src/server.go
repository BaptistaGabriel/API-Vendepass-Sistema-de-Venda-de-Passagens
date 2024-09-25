package main

import (
	"encoding/json"
	"fmt"
	"log"
	"io"
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
		if err == io.EOF {
			fmt.Println("Cliente desconectado.")
			return Message{Type: "disconnected"}
		}
		fmt.Printf("Erro em receber a mensagem do cliente %v\n", err)
		return Message{Type: "error"}
	}
	fmt.Printf("Mensagem recebida do cliente: %+v\n", msg)
	return msg
}

func sendMessage(connection net.Conn, msg Message) {
	encoder := json.NewEncoder(connection)
	if err := encoder.Encode(msg); err != nil {
		fmt.Printf("Erro ao mandar a resposta para o cliente %v\n", err)
		return
	}
	fmt.Println("Resposta devolvida para o cliente")
}

func sendJSON(connection net.Conn, list []string) {
	msg := Message{Type: "list", Content: list}
	sendMessage(connection, msg)
}

func communication(connection net.Conn, mapClients map[int]string, flights []Flight) {
	defer connection.Close()
	exit := true
	var numberID int

	// Menu 1 - Login ou Cadastro
	for exit {
		optionMsg := receiveMessage(connection)
		if optionMsg.Type == "disconnected" {
			return // Saia se o cliente desconectar
		}
		if optionMsg.Type != "action" {
			continue
		}

		option := optionMsg.Content.(float64)
		if int(option) == 1 {
			numberIDMsg := receiveMessage(connection)
			if numberIDMsg.Type == "disconnected" {
				return
			}
			numberID, _ = strconv.Atoi(numberIDMsg.Content.(string))
			name, exists := mapClients[numberID]

			if exists {
				sendMessage(connection, Message{Type: "response", Content: name})
				exit = false
			} else {
				sendMessage(connection, Message{Type: "response", Content: "-1"})
			}
		} else if int(option) == 2 {
			nameMsg := receiveMessage(connection)
			if nameMsg.Type == "disconnected" {
				return
			}
			clientID := createClient(nameMsg.Content.(string), mapClients)
			sendMessage(connection, Message{Type: "response", Content: strconv.Itoa(clientID)})
		} else {
			return
		}
	}

	// Menu 2 - Compras e Cancelamentos
	exit = true
	for exit {
		optionMsg := receiveMessage(connection)
		if optionMsg.Type == "disconnected" {
			return
		}
		if optionMsg.Type != "action" {
			continue
		}

		option := int(optionMsg.Content.(float64))
		if option == 1 {
			// Compra de assentos
			routes := GetRoutes(flights)
			sendJSON(connection, routes)

			routeNumberMsg := receiveMessage(connection)
			if routeNumberMsg.Type == "disconnected" {
				return
			}
			routeNumber, _ := strconv.Atoi(routeNumberMsg.Content.(string))
			clientFlight := flights[routeNumber]
			seats := GetSeats(clientFlight)

			var listSeats []string
			for _, seat := range seats {
				listSeats = append(listSeats, strconv.FormatBool(seat.IsReserved))
			}

			sendJSON(connection, listSeats)

			seatNumberMsg := receiveMessage(connection)
			if seatNumberMsg.Type == "disconnected" {
				return
			}
			seatNumber, _ := strconv.Atoi(seatNumberMsg.Content.(string))
			if ReserveSeat(flights, routeNumber, seatNumber, strconv.Itoa(numberID)) {
				sendMessage(connection, Message{Type: "response", Content: "Assento comprado com sucesso!"})
			} else {
				sendMessage(connection, Message{Type: "response", Content: "Erro ao comprar o assento"})
			}

			err := SaveFlightsToFile("flights.json", flights)
			if err != "" {
				fmt.Println("Erro ao salvar a operação!")
			}

		} else if option == 2 {
			// Cancelamento de passagens
			var clientFlights []string

			// Carregar voos atualizados
			flights, err := LoadFlightsFromFile("flights.json")
			if err != nil {
				fmt.Println("Erro ao carregar vôos")
				sendMessage(connection, Message{Type: "error", Content: "Erro ao carregar vôos"})
				continue
			}

			// Mostrar as reservas do cliente
			for index, flight := range flights {
				for seatIndex, seat := range flight.Seats {
					if seat.CustomerID == strconv.Itoa(numberID) {
						clientFlight := fmt.Sprintf("%d - %s -> %s - Assento: %d", index, flight.Origin, flight.Destination, seatIndex)
						clientFlights = append(clientFlights, clientFlight)
					}
				}
			}

			sendJSON(connection, clientFlights)

			// Selecionar voo e assento para cancelamento
			flightIndexMsg := receiveMessage(connection)
			if flightIndexMsg.Type == "disconnected" {
				return
			}
			flightIndex, _ := strconv.Atoi(flightIndexMsg.Content.(string))

			seatIndexMsg := receiveMessage(connection)
			if seatIndexMsg.Type == "disconnected" {
				return
			}
			seatIndex, _ := strconv.Atoi(seatIndexMsg.Content.(string))

			flights, success := CancelSeat(flights, flightIndex, seatIndex)
			if success {
				sendMessage(connection, Message{Type: "response", Content: "Passagem cancelada com sucesso!"})
			} else {
				sendMessage(connection, Message{Type: "response", Content: "Erro ao cancelar a passagem"})
			}

			if SaveFlightsToFile("flights.json", flights) != "" {
				fmt.Println("Erro ao salvar arquivo")
			}

		} else if option == 3 {
			fmt.Println("Cliente saiu.")
			return
		} else {
			return
		}
	}
}
func getLocalIP() net.IP {
	// Fazendo uma conexão não efetiva com o servidor DNS da Google
	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	localAddress := connection.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func createClient(name string, mapClients map[int]string) int {
	number := len(mapClients) + 1
	mapClients[number] = name
	return number
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

	mapClients := make(map[int]string)

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

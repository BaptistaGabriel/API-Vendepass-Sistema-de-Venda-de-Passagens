package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	// Menu 1
	for exit {
		optionMsg := receiveMessage(connection)
		if optionMsg.Type != "action" {
			continue
		}

		option := optionMsg.Content.(float64) // JSON decodifica números como float64
		if int(option) == 1 {
			numberIDMsg := receiveMessage(connection)
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
			clientID := createClient(nameMsg.Content.(string), mapClients)
			sendMessage(connection, Message{Type: "response", Content: strconv.Itoa(clientID)})
		} else if int(option) == 0 {
			break
		}
	}

	// Menu 2
	for {
	 optionMsg := receiveMessage(connection)
		if optionMsg.Type != "action" {
			continue
		}

		option := int(optionMsg.Content.(float64))
		if option == 1 {
			fmt.Println("Finge que está comprando")
			exit := true
			for exit {
				routes := GetRoutes(flights)
				sendJSON(connection, routes)

				routeNumberMsg := receiveMessage(connection)
				routeNumber, _ := strconv.Atoi(routeNumberMsg.Content.(string))
				clientFlight := flights[routeNumber]
				seats := GetSeats(clientFlight)

				var listSeats []string
				for _, seat := range seats {
					listSeats = append(listSeats, strconv.FormatBool(seat.IsReserved))
				}

				sendJSON(connection, listSeats)

				seatNumberMsg := receiveMessage(connection)
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

				optionExitMsg := receiveMessage(connection)
				if optionExitMsg.Type != "action" {
					continue
				}
				optionExit := int(optionExitMsg.Content.(float64))
				if optionExit == 2 {
					exit = false
				}
			}
		} else if option == 2 {
			fmt.Println("Finge que está cancelando")
		} else if option == 3 {
			fmt.Println("Saindooooooo")
			return
		} else if option != 0 {
			fmt.Println("Se o cliente cair!! No segundo menu claro")
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

/*func createFileClients(){
	_, err := os.Create("data/clients.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo dos clientes")
		return
	}
	return
}

func saveClients(name string, mapClients map[int] string) int{

	numberID := createClient(name, mapClients)

	file, err := json.MarshalIndent(mapClients, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return -1
	}

	err = os.WriteFile("data/clients.json", file, 0644)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo: ", err)
		return -1
	}

	return numberID
}

func loadFromClientsFIle(numberID int, mapClients map[int] string) (string, bool) {
	file, err := os.Open("data/clients.json")
	if err != nil {
		return "-1", false
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "-1", false
	}

	err = json.Unmarshal(bytes, &mapClients)
	if err != nil {
		return "-1", false
	}
	name, exists := mapClients[numberID]

	return name, exists

} */

func main() {

	// Pegando o IP do servidor
	fmt.Printf("IP do servidor %v\n", getLocalIP())

	// Criando o servidor na porta 8080
	listener, err := net.Listen("tcp", ":7777")
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

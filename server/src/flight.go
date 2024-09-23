package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Seat struct {
	IsReserved bool
	CustomerID string
}

type Flight struct {
	Origin      string
	Destination string
	Seats       []Seat
}

// Função para criar os assentos
func CreateSeats(numSeats int) []Seat {
	seats := make([]Seat, numSeats)
	for i := range seats {
		seats[i] = Seat{IsReserved: false, CustomerID: ""}
	}
	return seats
}

// Função para reservar um assento (muda para reservado e armazena o ID do cliente)
func ReserveSeat(flights []Flight, flightIndex int, seatNumber int, customerID string) bool {
	if flightIndex >= 0 && flightIndex < len(flights) {
		flight := &flights[flightIndex]
		if seatNumber >= 0 && seatNumber < len(flight.Seats) && !flight.Seats[seatNumber].IsReserved {
			flight.Seats[seatNumber].IsReserved = true
			flight.Seats[seatNumber].CustomerID = customerID
			SaveFlightsToFile("flights.json", flights) // Salvar as alterações
			return true
		}
	}
	return false
}

// Função para cancelar uma reserva (muda para disponível e remove o ID do cliente)
func CancelSeat(flights []Flight, flightIndex int, seatNumber int) ([]Flight, bool) {
	if flightIndex >= 0 && flightIndex < len(flights) {
		flight := &flights[flightIndex]
		if seatNumber >= 0 && seatNumber < len(flight.Seats) && flight.Seats[seatNumber].IsReserved {
			flight.Seats[seatNumber].IsReserved = false
			flight.Seats[seatNumber].CustomerID = ""
			SaveFlightsToFile("flights.json", flights) // Salvar as alterações
			return flights, true
		}
	}
	return flights, false
}

// Função para listar todos os assentos de um voo
func GetSeats(flight Flight) []Seat {
	return flight.Seats
}

// Função para limpar todos os assentos (marca todos como disponíveis e remove IDs)
func ClearSeats(flights []Flight, flightIndex int) []Flight {
	if flightIndex >= 0 && flightIndex < len(flights) {
		for i := range flights[flightIndex].Seats {
			flights[flightIndex].Seats[i].IsReserved = false
			flights[flightIndex].Seats[i].CustomerID = ""
		}
		SaveFlightsToFile("flights.json", flights) // Salvar as alterações
	}
	return flights
}

// Função para criar rotas (inicializa os voos com assentos disponíveis)
func CreateRoutes() []Flight {
	var flights []Flight

	routes := map[string][]string{
		"Kiev":             {"Xique Xique-BA"},
		"Xique Xique-BA":   {"Recife", "Feira de Santana", "Brasilia", "São Paulo", "Rio de Janeiro", "Curitiba"},
		"Recife":           {"Xique Xique-BA", "Feira de Santana"},
		"Feira de Santana": {"Recife", "Xique Xique-BA"},
		"Salvador":         {"Feira de Santana", "Manaus"},
		"Manaus":           {"Palmas", "Rio Branco"},
		"Palmas":           {"Rio Branco", "Manaus"},
		"Rio Branco":       {"Manaus", "Palmas"},
		"Brasilia":         {"Xique Xique-BA"},
		"São Paulo":        {"Xique Xique-BA"},
		"Rio de Janeiro":   {"Xique Xique-BA"},
		"Curitiba":         {"Xique Xique-BA"},
	}

	for origin, destinations := range routes {
		for _, destination := range destinations {
			flight := Flight{
				Origin:      origin,
				Destination: destination,
				Seats:       CreateSeats(12),
			}
			flights = append(flights, flight)
		}
	}

	return flights
}

func SaveFlightsToFile(filename string, flights []Flight) string {
	// Garante que a pasta 'data' exista
	os.MkdirAll("data", os.ModePerm)

	// Atualiza o caminho para salvar na pasta 'data'
	filepath := "data/" + filename

	file, err := json.MarshalIndent(flights, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return "Erro"
	}

	err = os.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar arquivo:", err)
		return "Erro"
	}
	return ""
}

// Função para carregar os voos e assentos de um arquivo JSON
func LoadFlightsFromFile(filename string) ([]Flight, error) {
	var flights []Flight
	path := "data/" + filename

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &flights)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

// Função para pegar as rotas
func GetRoutes(flights []Flight) []string {
	var routes []string
	for _, flight := range flights {
		route := fmt.Sprintf("%s -> %s", flight.Origin, flight.Destination)
		routes = append(routes, route)
	}
	return routes
}

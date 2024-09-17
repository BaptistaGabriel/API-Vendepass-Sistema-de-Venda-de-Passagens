package main

type Seat struct {
	Number      int
	IsAvailable bool
}

type Flight struct {
	Origin      string
	Destination []string
	Seats       []Seat
}

func CreateRoutes() [12]Flight {
	var flights [12]Flight
	var destinations = [12]string{"Kiev", "Xique Xique-BA", "Recife", "Feira de Santana", "Salvador", "Manaus", "Palmas", "Rio Branco", "Brasilia", "São Paulo", "Rio de Janeiro", "Curitiba"}

	routes := map[string][]string{
		"Kiev":          {"Xique Xique-BA"},
		"Xique Xique-BA": {"Recife", "Feira de Santana", "Brasilia", "São Paulo", "Rio de Janeiro", "Curitiba"},
		"Recife":        {"Xique Xique-BA", "Feira de Santana"},
		"Feira de Santana": {"Recife", "Xique Xique-BA"},
		"Salvador":      {"Feira de Santana", "Manaus"},
		"Manaus":        {"Palmas", "Rio Branco"},
		"Palmas":        {"Rio Branco", "Manaus"},
		"Rio Branco":    {"Manaus", "Palmas"},
		"Brasilia":      {"Xique Xique-BA"},
		"São Paulo":     {"Xique Xique-BA"},
		"Rio de Janeiro": {"Xique Xique-BA"},
		"Curitiba":      {"Xique Xique-BA"},
	}

	for i := 0; i < 12; i++ {
		flights[i].Origin = destinations[i]
        flights[i].Destination = routes[destinations[i]]
		for j := 0; j < 10; j++ {
			flights[i].Seats[j].Number = j
			flights[i].Seats[j].IsAvailable = true
		}
	}

    return flights
}

func ReserveSeat(flights [12]Flight, origin string, destination string, seat int) bool {
	for i := 0; i < 12; i++ {
		if flights[i].Origin == origin {
			for j := 0; j < len(flights[i].Destination); j++ {
				if flights[i].Destination[j] == destination {
					if flights[i].Seats[seat].IsAvailable {
						flights[i].Seats[seat].IsAvailable = false
						return true
					}
				}
			}
		}
	}
	return false
}

func CancelSeat(flights [12]Flight, origin string, destination string, seat int) bool {
	for i := 0; i < 12; i++ {
		if flights[i].Origin == origin {
			for j := 0; j < len(flights[i].Destination); j++ {
				if flights[i].Destination[j] == destination {
					if !flights[i].Seats[seat].IsAvailable {
						flights[i].Seats[seat].IsAvailable = true
						return true
					}
				}
			}
		}
	}
	return false
}

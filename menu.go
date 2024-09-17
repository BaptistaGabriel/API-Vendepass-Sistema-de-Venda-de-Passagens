package main 

import(
	"fmt"
)

func createClient(name string, mapClients *map[int] string) int{
	number := len(*mapClients) + 1
	(*mapClients)[number] = name
	return number
}

func login(mapClients *map[int] string) bool{
	var number int 
	fmt.Println("Número do login: ")
	fmt.Scanln(&number)
	if _, ok := (*mapClients)[number]; !ok {
		return false
	}
	return true
}

func main () {
	var option int
	mapClients := make(map[int]string)
	
	for {
		fmt.Println("==========================")
		fmt.Printf("\033[34m|     1. Fazer login     |\n|     2. Criar conta     |\033\n[0m")
		fmt.Println("==========================")
		fmt.Scanln(&option)
		switch option {
		case 1:
			exit := false
			if login(&mapClients) {
				for !exit {
					fmt.Println("Testando 123...")
					exit = true
				}
			} else {
				fmt.Println("Usuário não cadastrado.")
			}
		case 2:
			var name string
			fmt.Printf("Nome: ")
			fmt.Scanln(&name)
			number := createClient(name, &mapClients)
			fmt.Println(number)
			fmt.Println(mapClients)
		default:
			fmt.Println("Opção inválida!")
		}
	}

}
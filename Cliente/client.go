package main

import (
	"fmt"
	"net"
	"strings"
)

// Função para manipular cada conexão de cliente
func manipularConexao(conexao net.Conn) {
	// Fechar a conexão ao final da execução da função
	defer conexao.Close()

	// Identificando o cliente
	enderecoCliente := conexao.RemoteAddr().String()
	ipCliente, portaCliente := dividirIPPorta(enderecoCliente)
	fmt.Printf("Cliente conectado - IP: %s, Porta: %s\n", ipCliente, portaCliente)

	// Lendo dados do cliente
	buffer := make([]byte, 1024) // Buffer de 1 KB para armazenar os dados recebidos
	bytesLidos, err := conexao.Read(buffer)

	if err != nil {
		fmt.Println("Erro ao ler os dados:", err)
		return
	}

	// Convertendo os dados lidos para string
	mensagem := string(buffer[:bytesLidos])
	fmt.Printf("Mensagem recebida: %s\n", mensagem)

	// Convertendo a mensagem para maiúsculas
	resposta := strings.ToUpper(mensagem)
	_, err = conexao.Write([]byte(resposta)) // Enviando a resposta ao cliente

	if err != nil {
		fmt.Println("Erro ao enviar a mensagem:", err)
		return
	}
	fmt.Println("Resposta enviada com sucesso!")
}

// Função auxiliar para dividir o endereço IP e a porta
func dividirIPPorta(endereco string) (string, string) {
	partes := strings.Split(endereco, ":")
	return partes[0], partes[1]
}

func main() {
	// Criando o servidor TCP na porta 7777
	servidor, err := net.Listen("tcp", ":7777")
	if err != nil {
		fmt.Println("Erro ao criar o servidor:", err)
		return
	}
	defer servidor.Close() // Fechar o servidor quando o programa terminar
	fmt.Println("Servidor funcionando na porta 7777...")

	// Loop para aceitar conexões de clientes
	for {
		conexao, err := servidor.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		// Manipula cada conexão de cliente em uma nova goroutine
		go manipularConexao(conexao)
	}
}

package main

import (
	"fmt"
	"net"
	"strings"
)

// Função responsável por manipular a conexão de um cliente
func manipularConexao(cliente net.Conn) {
	defer cliente.Close()

	ip, porta := identificarCliente(cliente)

	fmt.Printf("Cliente conectado: IP %s na porta %s\n", ip, porta)

	mensagem, erro := lerDados(cliente)
	if erro != nil {
		fmt.Printf("Erro ao ler os dados: %v\n", erro)
		return
	}

	fmt.Printf("Mensagem recebida: %s\n", mensagem)

	mensagemModificada := strings.ToUpper(mensagem)
	if erro := enviarDados(cliente, mensagemModificada); erro != nil {
		fmt.Printf("Erro ao enviar a mensagem: %v\n", erro)
		return
	}

	fmt.Println("Mensagem enviada com sucesso!")
}

// Função que identifica o cliente a partir do endereço remoto
func identificarCliente(cliente net.Conn) (string, string) {
	idPorta := cliente.RemoteAddr().String()
	identificador := strings.Split(idPorta, ":")
	ip := identificador[0]
	porta := identificador[1]
	return ip, porta
}

// Função para ler dados do cliente
func lerDados(cliente net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	tamanho, erro := cliente.Read(buffer)
	if erro != nil {
		return "", erro
	}
	return string(buffer[:tamanho]), nil
}

// Função para enviar dados ao cliente
func enviarDados(cliente net.Conn, mensagem string) error {
	_, erro := cliente.Write([]byte(mensagem))
	return erro
}

func main() {
	// Criando o servidor na porta 7777
	servidor, erro := net.Listen("tcp", ":7777")
	if erro != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", erro)
		return
	}
	defer servidor.Close()

	fmt.Println("Servidor funcionando na porta 7777...")

	// Aceitando conexões em loop
	for {
		conexao, erro := servidor.Accept()
		if erro != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", erro)
			continue
		}
		go manipularConexao(conexao) // Manipulando conexões simultâneas
	}
}

<div align="center">

# API para Vendepass Sistema de Venda de Passagens
**Projeto desenvolvido para a disciplina TEC502 - MI Concorrência e Conectividade**
</div>

### Resumo

Este projeto é um **sistema de venda de passagens aéreas**, desenvolvido em **Go (Golang)**, que utiliza **TCP-IP** para comunicação entre o cliente e o servidor. O sistema possibilita que clientes se cadastrem, visualizem voos disponíveis, comprem passagens e cancelem seus pedidos. O objetivo deste trabalho é simular a arquitetura de um sistema distribuído e explorar conceitos de concorrência e comunicação de rede.

## Índice
- [Introdução](#introdução)
  - [Tecnologias Utilizadas](#tecnologias-utilizadas)
  - [Como Executar](#como-executar)
- [Arquitetura do Projeto](#arquitetura-do-projeto)
  - [Client](#client)
  - [Server](#server)
- [Comunicação](#comunicação)
  - [Paradigma de Comunicação](#paradigma-de-comunicação)
  - [Protocolo de Comunicação](#protocolo-de-comunicação)
  - [Formatação e Tratamento de Dados](#formatação-e-tratamento-de-dados)
  - [Tratamento de Conexões Simultâneas](#tratamento-de-conexões-simultâneas)
  - [Tratamento de Concorrência](#tratamento-de-concorrência)
- [Desempenho e Avaliação](#desempenho-e-avaliação)
- [Testes](#testes)
- [Conclusão](#conclusão)

### Introdução
O desenvolvimento de uma comunicação servidor-cliente é fundamental para as empresas aéreas de baixo custo desenvolverem seus negócios. Então, o objetivo deste trabalho é projetar uma comunicação servidor-cliente usando um protocolo de comunicação. O projeto consiste em um servidor TCP-IP escrito em Go para gerenciar rotas de voos e permitir que clientes realizem ações como cadastro, compra e cancelamento de passagens. Este relatório descreve as tecnologias utilizadas, procedimentos para execução do projeto e a explicação de como utilizá-lo.

## Tecnologias Utilizadas

As principais tecnologias usadas no desenvolvimento deste projeto:

- **Go (Golang)**: Linguagem principal do servidor, escolhida por sua eficiência em lidar com concorrência e comunicação de rede.
- **Docker**: Utilizado para criar um contêiner que facilita a execução, portabilidade e replicação do projeto em diversos ambientes.
- **TCP (Transmission Control Protocol)**: Protocolo de comunicação confiável para interação entre o cliente e o servidor.
- **net**: Pacote nativo do Go usado para gerenciar as conexões TCP.
- **Alpine Linux**: Base do contêiner Docker, fornecendo um ambiente leve para executar o servidor Go.

## Como Executar

Siga os passos abaixo para configurar e executar o servidor localmente usando Docker.

1. **Clone o repositório**:

    Primeiro, faça o download do código fonte:

    ```bash
    git clone https://github.com/BaptistaGabriel/API-Vendepass-Sistema-de-Venda-de-Passagens
    cd API-Vendepass-Sistema-de-Venda-de-Passagens
    ```

2. **Construa a imagem Docker**:

    Dentro do diretório onde o `Dockerfile` está localizado, construa a imagem com:

    ```bash
    docker build -t vendepass_server .
    ```

3. **Execute o contêiner**:

    Após a construção da imagem, inicie o contêiner expondo a porta `8080`:

    ```bash
    docker run -d -p 8080:8080 vendepass_server
    ```

4. **Verifique a execução do servidor**:

    O servidor estará funcionando na porta `8080` do localhost. Use um cliente TCP para verificar a comunicação.

## Arquitetura do Projeto

A arquitetura do sistema foi desenvolvida seguindo o modelo cliente-servidor, com o servidor responsável por gerenciar todas as rotas de voos e os clientes se comunicando com ele para realizar as operações de cadastro, compra e cancelamento de passagens. A escolha por uma arquitetura baseada em **TCP** visa garantir uma comunicação confiável entre cliente e servidor.

### Componentes

- **Client**: 
  - O cliente é responsável por enviar solicitações ao servidor e receber respostas. As interações com o servidor incluem o cadastro de clientes, consulta de voos disponíveis, compra e cancelamento de passagens.
  - O cliente foi desenvolvido com foco na simplicidade, utilizando sockets TCP para a comunicação direta com o servidor.

- **Server**: 
  - O servidor gerencia toda a lógica de negócio, sendo responsável por manter o estado coerente das rotas de voos, os clientes cadastrados, e as passagens compradas. 
  - O servidor é **stateful**, ou seja, ele mantém o estado das operações ao longo da comunicação. 
  - O servidor foi desenvolvido utilizando a linguagem **Go**, com **goroutines** para lidar com múltiplas conexões simultâneas e **channels** para a comunicação entre processos concorrentes. A escolha do Go se deve à sua eficiência em lidar com operações paralelas e concorrentes.

### Estrutura das Pastas

Abaixo está a estrutura de pastas e arquivos do projeto:

```bash
API-Vendepass-Sistema-de-Venda-de-Passagens/
├── client/
│   ├── client.go         # Código do cliente para comunicação com o servidor
├── server/
│   ├── data/
│   │   ├── flight.json   # Arquivo json responsável por armazenar os dados dos vôos
│   ├── src/
│   │   ├── flight.go     # Código responsável pelo gerenciamento de vôos
│   │   ├── server.go     # Código principal do servidor, incluindo a lógica de rotas e passagens
├── Dockerfile            # Arquivo Docker para construir o ambiente do servidor
├── go.mod                # Arquivo de dependências Go
├── Makefile              # Script Makefile para automação de tarefas de build e execução
├── .gitignore            # Arquivo para ignorar arquivos desnecessários no versionamentoo

```

## Comunicação


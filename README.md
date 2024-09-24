<div align="center">

# API para Vendepass Sistema de Venda de Passagens
**Projeto desenvolvido para a disciplina TEC502 - MI Concorrência e Conectividade**
</div>

## Resumo

Este projeto é um **sistema de venda de passagens aéreas**, desenvolvido em **Go (Golang)**, que utiliza **TCP-IP** para comunicação entre o cliente e o servidor. O sistema possibilita que clientes se cadastrem, visualizem voos disponíveis, comprem passagens e cancelem seus pedidos. O objetivo deste trabalho é simular a arquitetura de um sistema distribuído e explorar conceitos de concorrência e comunicação de rede.

## Índice
- [Introdução](#introdução) <!-- Introdução -->
- [Fundamentação Teórica](#fundamentação-teórica)
  - [Conceitos de Concorrência](#conceitos-de-concorrência)<!-- Fundamentação Teórica -->
  - [Tecnologias Utilizadas](#tecnologias-utilizadas)
  - [Como Executar](#como-executar)
- [Arquitetura do Projeto](#arquitetura-do-projeto)<!-- Metodologia/Diagramas e figuras -->
  - [Client](#client) 
  - [Server](#server)
- [Comunicação](#comunicação)
  - [Paradigma de Comunicação](#paradigma-de-comunicação)
  - [Protocolo de Comunicação](#protocolo-de-comunicação)
  - [Formatação e Tratamento de Dados](#formatação-e-tratamento-de-dados)
  - [Tratamento de Conexões Simultâneas](#tratamento-de-conexões-simultâneas)
  - [Tratamento de Concorrência](#tratamento-de-concorrência)
- [Desempenho e Avaliação](#desempenho-e-avaliação) <!-- Resultados e discussões -->
- [Testes](#testes)
- [Conclusão](#conclusão) <!-- Conclusão -->
- [Referências](#referências) <!-- Referências -->

### Introdução
O desenvolvimento de uma comunicação servidor-cliente é fundamental para as empresas aéreas de baixo custo desenvolverem seus negócios. Então, o objetivo deste trabalho é projetar uma comunicação servidor-cliente usando um protocolo de comunicação. O projeto consiste em um servidor TCP-IP escrito em Go para gerenciar rotas de voos e permitir que clientes realizem ações como cadastro, compra e cancelamento de passagens. Este relatório descreve as tecnologias utilizadas, procedimentos para execução do projeto e a explicação de como utilizá-lo.

### Fundamentação Teórica

#### Conceitos de Concorrência

O projeto utiliza os conceitos de concorrência oferecidos pela linguagem Go para permitir que o servidor atenda múltiplas conexões de clientes de maneira eficiente. A principal ferramenta utilizada é a goroutine, que permite a execução simultânea de tarefas sem bloquear o servidor. Cada vez que um cliente se conecta, uma nova goroutine é iniciada, permitindo o processamento de várias requisições de forma concorrente.

Esses conceitos formam a base da solução de concorrência adotada no projeto, garantindo escalabilidade e eficiência no processamento de requisições simultâneas.

#### Tecnologias Utilizadas

As principais tecnologias usadas no desenvolvimento deste projeto:

- **Go (Golang)**: Linguagem principal do servidor, escolhida por sua eficiência em lidar com concorrência e comunicação de rede.
- **Docker**: Utilizado para criar um contêiner que facilita a execução, portabilidade e replicação do projeto em diversos ambientes.
- **TCP (Transmission Control Protocol)**: Protocolo de comunicação confiável para interação entre o cliente e o servidor.
- **net**: Pacote nativo do Go usado para gerenciar as conexões TCP.
- **Alpine Linux**: Base do contêiner Docker, fornecendo um ambiente leve para executar o servidor Go.

#### Como Executar

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

##### Client:

- O cliente é responsável por enviar solicitações ao servidor e receber respostas. As interações com o servidor incluem o cadastro de clientes, consulta de voos disponíveis, compra e cancelamento de passagens.
- O cliente foi desenvolvido com foco na simplicidade, utilizando sockets TCP para a comunicação direta com o servidor.

##### Server:

- O servidor gerencia toda a lógica de negócio, sendo responsável por manter o estado coerente das rotas de voos, os clientes cadastrados, e as passagens compradas.
- O servidor é **stateful**, ou seja, ele mantém o estado das operações ao longo da comunicação.
- O servidor foi desenvolvido utilizando a linguagem **Go**, com **goroutines** para lidar com múltiplas conexões simultâneas e **channels** para a comunicação entre processos concorrentes. A escolha do Go se deve à sua eficiência e facilidade em lidar com operações paralelas e concorrentes.

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

### Paradigma de Comunicação

O paradigma de comunicação adotado foi o **stateful**, onde o servidor mantém o estado das interações dos clientes ao longo do tempo. Essa escolha foi motivada pela necessidade de rastrear informações importantes sobre cada cliente, como dados de cadastro, passagens compradas e outras transações realizadas durante a sessão. Dessa forma, a cada nova conexão, o cliente pode continuar de onde parou, sem perder o contexto anterior.

### Protocolo de Comunicação

O protocolo de comunicação desenvolvido entre os componentes do sistema utiliza o **protocolo TCP (Transmission Control Protocol)**. A troca de mensagens segue uma ordem específica, sendo:

1. O cliente envia uma requisição de ação (como listar voos, comprar ou cancelar passagem).
2. O servidor processa a requisição e envia uma resposta contendo os dados solicitados ou confirmações das ações.
3. Se o cliente precisar de mais informações ou realizar novas ações, o processo se repete, mantendo a conexão ativa até a finalização.

A comunicação é estruturada para garantir a confiabilidade da troca de dados, sendo a ordem das mensagens essencial para o fluxo correto das operações.

### Formatação e Tratamento de Dados

Os dados trocados entre cliente e servidor são formatados em **JSON**, um formato amplamente utilizado e aceito para a troca de dados em sistemas distribuídos. O uso de JSON garante uma estrutura legível, fácil de implementar e compatível com diferentes plataformas e linguagens de programação. Isso permite a comunicação eficaz entre componentes distintos do sistema, mantendo a integridade e consistência das informações. Além disso, a escolha por JSON facilita a serialização e desserialização dos dados, simplificando o tratamento das mensagens enviadas e recebidas pelo sistema.

### Tratamento de Conexões Simultâneas

O sistema foi desenvolvido para permitir a **realização de compras de passagens de forma paralela ou simultânea**. A otimização do paralelismo é alcançada utilizando **goroutines**, que são threads leves e nativas da linguagem Go, possibilitando que várias operações sejam executadas simultaneamente sem bloquear a execução do servidor.

Além disso, o uso de **channels** e outros mecanismos de sincronização do Go melhora o gerenciamento das conexões concorrentes, garantindo que múltiplos clientes possam comprar passagens sem conflitos.

### Tratamento de Concorrência

Sim, o uso de conexões simultâneas pode gerar problemas de concorrência, especialmente quando várias operações tentam acessar ou modificar os mesmos dados ao mesmo tempo. Para evitar condições de corrida (race conditions), onde duas ou mais goroutines tentam modificar os mesmos recursos simultaneamente, o sistema implementa mecanismos de controle de concorrência.

## Desempenho e Avaliação

Para avaliar o desempenho do sistema, foram realizadas várias simulações e testes de carga. As métricas de desempenho consideraram o tempo de resposta do servidor, a capacidade de lidar com conexões simultâneas e a integridade dos dados após operações concorrentes.

Os testes mostraram que o servidor pode lidar eficientemente com múltiplas requisições de clientes sem degradação significativa no tempo de resposta. A implementação de goroutines e mostrou eficaz na mitigação de problemas de concorrência, garantindo que os dados permanecessem consistentes mesmo em cenários de alta carga.

## Testes

Os testes foram realizados em duas frentes: testes funcionais e testes de carga.

1. **Testes Funcionais**: Incluíram a verificação de todas as funcionalidades do sistema, como cadastro de clientes, consulta de voos, compra e cancelamento de passagens. Todos os testes funcionais passaram, confirmando que as funcionalidades estão operando conforme o esperado.

2. **Testes de Carga**: Avaliaram o comportamento do sistema sob múltiplas conexões simultâneas. O servidor demonstrou resiliência, mantendo tempos de resposta aceitáveis e garantindo a integridade dos dados durante as operações concorrentes.

## Conclusão

O sistema de venda de passagens aéreas desenvolvido em Go demonstrou ser eficiente na gestão de conexões e na realização de operações de compra e cancelamento. A utilização de TCP-IP para a comunicação entre cliente e servidor, juntamente com os conceitos de concorrência do Go, permitiu a construção de uma solução robusta e escalável.

Futuros desenvolvimentos podem incluir a adição de novas funcionalidades, como um painel de administração para gerenciar voos e usuários, bem como a implementação de autenticação e autorização de usuários para aumentar a segurança do sistema.

## Referências

1. [Go Samples: Local IP Address](https://gosamples.dev/local-ip-address/)
2. [The Complete Guide to TCP/IP Connections in Golang](https://okanexe.medium.com/the-complete-guide-to-tcp-ip-connections-in-golang-1216dae27b5a)
3. [W3Schools: Go Language](https://www.w3schools.com/go/index.php)
4. [iMasters: Trabalhando com o Protocolo TCP/IP Usando Go](https://imasters.com.br/back-end/trabalhando-com-o-protocolo-tcpip-usando-go)

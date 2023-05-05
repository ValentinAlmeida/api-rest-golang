# api-rest-golang
O código é um programa Go que inicia um servidor HTTP para lidar com requisições de usuários.
O package main é uma declaração obrigatória em qualquer arquivo Go que define um programa executável.
O arquivo é chamado de "pacote principal" porque contém a função main, que é o ponto de entrada do programa.

O arquivo importa diversos pacotes que fornecem funcionalidades como o uso de banco de dados MySQL,
criptografia de senha, manipulação de rotas HTTP e codificação/decodificação de JSON.
Alguns desses pacotes são importados apenas para registrar seus drivers (por exemplo, github.com/go-sql-driver/mysql)
e não são usados diretamente no código.

A estrutura User define um modelo de dados para um usuário que contém um ID, nome, email e senha (criptografada).

A estrutura server representa o servidor e tem um único campo, db, que é um objeto de conexão com o banco de dados.

A função main é a primeira a ser executada quando o programa é iniciado e é responsável por configurar o servidor HTTP e iniciar a escuta por requisições.

Além disso, o código define outras funções que são usadas para lidar com diferentes tipos de requisições HTTP,
como listar, criar, atualizar e excluir usuários. Cada uma dessas funções manipula a requisição recebida,
processa as informações necessárias a partir do banco de dados e envia uma resposta de volta ao cliente em formato JSON.

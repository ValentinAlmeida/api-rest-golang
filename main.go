package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// User é a estrutura que representa um usuário no banco de dados.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

// server é uma implementação que usa o MySQL como banco de dados.
type server struct {
	db *sql.DB
}

func main() {
	// Conecta ao banco de dados MySQL.
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cria uma instância da nossa aplicação server.
	s := server{db}

	// Cria um roteador usando a biblioteca Gorilla Mux.
	r := mux.NewRouter()
	r.HandleFunc("/users", s.handleGetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", s.handleGetUser).Methods("GET")
	r.HandleFunc("/users", s.handleCreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", s.handleUpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", s.handleDeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// handleGetUsers é uma função que lida com a requisição GET para obter uma lista de usuários.
func (s *server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// Executa uma consulta SQL para obter uma lista de usuários.
	rows, err := s.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Cria uma lista de usuários vazia.
	var users []User

	// Itera sobre cada linha do resultado da consulta e adiciona um usuário à lista de usuários.
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}

	// Verifica se ocorreu algum erro durante a iteração sobre as linhas do resultado da consulta.
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Converte a lista de usuários para um formato JSON e envia como resposta para o cliente.
	json.NewEncoder(w).Encode(users)
}

// handleGetUser é uma função que lida com a requisição GET para obter um usuário específico.
func (s *server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do usuário da URL da requisição.
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Executa uma consulta SQL para obter o usuário com o ID especificado.
	var u User
	err = s.db.QueryRow("SELECT id, name, email FROM users WHERE id=?", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		log.Fatal(err)
	}

	// Converte a lista de usuários para um formato JSON e envia como resposta para o cliente.
	json.NewEncoder(w).Encode(u)
}

// handleCreateUser é responsável por criar um novo usuário no banco de dados
func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Decodifica a requisição JSON em um objeto User
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Fatal(err)
	}

	// Cria um hash da senha do usuário
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Insere o novo usuário no banco de dados
	result, err := s.db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, passwordHash)
	if err != nil {
		log.Fatal(err)
	}

	// Obtém o ID do usuário criado
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Define o ID do usuário criado no objeto User
	u.ID = int(id)

	// Codifica o objeto User em JSON e retorna na resposta HTTP
	json.NewEncoder(w).Encode(u)
}

// handleUpdateUser é responsável por atualizar um usuário existente no banco de dados
func (s *server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Obtém o ID do usuário a ser atualizado a partir dos parâmetros da URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Decodifica a requisição JSON em um objeto User
	var u User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Fatal(err)
	}

	// Atualiza os dados do usuário no banco de dados
	_, err = s.db.Exec("UPDATE users SET name=?, email=?, password=? WHERE id=?", u.Name, u.Email, u.Password, id)
	if err != nil {
		log.Fatal(err)
	}

	// Define o ID do usuário no objeto User e retorna na resposta HTTP
	u.ID = id
	json.NewEncoder(w).Encode(u)
}

// handleDeleteUser é responsável por excluir um usuário do banco de dados
func (s *server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Obtém o ID do usuário a ser excluído a partir dos parâmetros da URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Exclui o usuário do banco de dados
	_, err = s.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	// Define o status da resposta HTTP como "Sem conteúdo"
	w.WriteHeader(http.StatusNoContent)
}

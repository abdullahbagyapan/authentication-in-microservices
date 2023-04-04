package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"microservices/auth-service/internal/ports"
	"time"
)

type Adapter struct {
	db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "authservice"
)

func NewAdapter() *Adapter {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}

	// CONNECTION TEST
	err = db.Ping()

	if err != nil {
		log.Fatalf("Error testing to db %v", err)
	}

	// CREATE TABLE
	query := "CREATE TABLE IF NOT EXISTS users(id VARCHAR(255) PRIMARY KEY,name VARCHAR(255) NOT NULL,username VARCHAR(255) NOT NULL , email  VARCHAR(50) NOT NULL,password VARCHAR(255) NOT NULL,date timestamp NOT NULL)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	_, err = db.ExecContext(ctx, query)

	if err != nil {
		log.Fatalf("Error %s when creating table", err)
	}

	log.Printf("Succesfully connected to the database")
	return &Adapter{db: db}
}

func (da Adapter) CloseDbConnection() {
	err := da.db.Close()

	if err != nil {
		log.Fatalf("Error closing to db connection %v", err)
	}
}

func (da Adapter) SaveUser(user *ports.User) error {

	query := "INSERT INTO users(id, name,username,email,password, date) VALUES ($1,$2,$3,$4,$5,$6)"

	_, err := da.db.Exec(query, user.Id, user.Name, user.Username, user.Email, user.Password, time.Now())

	return err

}

func (da Adapter) FindByUsername(username string) (*ports.User, error) {

	query := "SELECT Id,password FROM users WHERE username = $1"

	rows, err := da.db.Query(query, username)

	if err != nil {
		return nil, err
	}

	var user ports.User

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Password)

		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

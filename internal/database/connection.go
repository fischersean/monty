package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const DBNAME = "monty"

type Connection struct {
	db *sql.DB
}

type NewConnectionInput struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewConnection(input NewConnectionInput) (conn Connection, err error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		input.Host, input.Port, input.User, input.Password, DBNAME)

	conn.db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return conn, err
	}
	// check connection
	err = conn.db.Ping()
	return conn, err
}

func (conn *Connection) Close() error {
	return conn.db.Close()
}

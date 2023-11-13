package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=nick dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists accounts (
    		id serial primary key,
    		first_name varchar(50),
    		last_name varchar(50),
    		number serial,
    		balance serial,
    		created_at timestamp default current_timestamp
                                    		);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) updateAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) deleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) getAccountById(id int) (*Account, error) {
	return nil, nil
}

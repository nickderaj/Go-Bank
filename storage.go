package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
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
    		id serial primary key not null,
    		first_name varchar(50) not null,
    		last_name varchar(50) not null,
    		number serial not null, 
    		balance serial not null,
    		created_at timestamp default current_timestamp
                                    		);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into accounts
    		(first_name, last_name, number, balance)
    		values ($1, $2, $3, $4);`

	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("delete from accounts where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0)
	for rows.Next() {
		a, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("select * from accounts where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	a := new(Account)
	if err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Number, &a.Balance, &a.CreatedAt); err != nil {
		return nil, err
	}
	return a, nil
}

package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repository interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	ListAccounts() ([]*Account, error)
	GetAccount(string) (*Account, error)
	ValidateAccount(agency, number int) (*Account, error)
	DeleteAccount(string) error
	DebitAccount(account *Account, amount int) error
	CreditAccount(account *Account, amount int) error
}

type PostGresRepository struct {
	db *sql.DB
}

func NewPostGresRepository() (*PostGresRepository, error) {
	connStr := "user=postgres dbname=gobank password=pg123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostGresRepository{
		db: db,
	}, nil
}

func (r *PostGresRepository) CreateAccount(account *Account) error {

	query := `INSERT INTO accounts 
	(id, first_name, last_name, document, agency, number, balance, created_at, modified_at) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Query(
		query,
		account.Id,
		account.FirstName,
		account.LastName,
		account.Document,
		account.Agency,
		account.Number,
		account.Balance,
		account.CreatedAt,
		account.ModifiedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostGresRepository) UpdateAccount(*Account) error {
	return nil
}

func (r *PostGresRepository) ListAccounts() ([]*Account, error) {
	rows, err := r.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *PostGresRepository) GetAccount(id string) (*Account, error) {
	rows, err := r.db.Query(`SELECT * FROM accounts WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanAccount(rows)
	}

	return nil, fmt.Errorf("Account %s not found", id)
}

func (r *PostGresRepository) ValidateAccount(agency, number int) (*Account, error) {
	rows, err := r.db.Query(`SELECT * FROM accounts WHERE agency = $1 AND number = $2`, agency, number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found", number)
}

func (r *PostGresRepository) DeleteAccount(id string) error {
	_, err := r.db.Query(`DELETE FROM accounts WHERE id = $1`, id)
	return err
}

func (r *PostGresRepository) DebitAccount(account *Account, amount int) error {
	_, err := r.db.Query(`UPDATE accounts SET balance = (balance - $2) WHERE id = $1`, account.Id, amount)
	return err
}

func (r *PostGresRepository) CreditAccount(account *Account, amount int) error {
	_, err := r.db.Query(`UPDATE accounts SET balance = (balance + $2) WHERE id = $1`, account.Id, amount)
	return err
}

func scanAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.Id,
		&account.FirstName,
		&account.LastName,
		&account.Document,
		&account.Agency,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.ModifiedAt,
	)
	return account, err
}

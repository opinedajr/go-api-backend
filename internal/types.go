package internal

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id         string    `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Document   string    `json:"document"`
	Agency     int       `json:"agency"`
	Number     int       `json:"number"`
	Balance    int64     `json:"balance"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type CreateAccountInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Document  string `json:"document"`
}

type TransferInput struct {
	ToAgency  int `json:"toAgency"`
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type SiginInput struct {
	Agency int `json:"agency"`
	Number int `json:"number"`
}

type SiginOutput struct {
	Token string `json:"token"`
}

type ApiErrorOutput struct {
	Error string `json:"error"`
}

func NewAccount(firstName, lastName, document string) *Account {
	return &Account{
		Id:         uuid.NewString(),
		FirstName:  firstName,
		LastName:   lastName,
		Document:   document,
		Agency:     rand.Intn(2) + 1,
		Number:     rand.Intn(100000),
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
}

func NewSiginOutput(tokenString string) *SiginOutput {
	return &SiginOutput{
		Token: tokenString,
	}
}

func NewApiError(err string) *ApiErrorOutput {
	return &ApiErrorOutput{
		Error: err,
	}
}

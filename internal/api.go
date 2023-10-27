package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

const Secret = "jwt-super-secret"

type ApiServer struct {
	listenAddress string
	repo          Repository
}

func InitializeApiServer(listenAddress string, repo Repository) *ApiServer {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
	return &ApiServer{
		listenAddress: listenAddress,
		repo:          repo,
	}
}

func EncodeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (s *ApiServer) Run() {

	router := chi.NewRouter()

	// Protected routes
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/accounts", s.AccountHandler)
		r.Get("/accounts/{id}", s.GetAccountHandler)
		r.Delete("/accounts/{id}", s.DeleteAccountHandler)
		r.Post("/transfers", s.TransferHandler)
	})

	// Public routes
	router.Group(func(r chi.Router) {
		r.Post("/signin", s.SigninHandler)
		r.Post("/accounts", s.CreateAccountHandler)
	})

	fmt.Println("JSON API Server running on port:", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiServer) AccountHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.repo.ListAccounts()
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, err.Error())
		return
	}
	EncodeJson(w, http.StatusOK, accounts)
}

// Implements the API Handlers
func (s *ApiServer) TransferHandler(w http.ResponseWriter, r *http.Request) {
	transferInput := new(TransferInput)
	if err := json.NewDecoder(r.Body).Decode(transferInput); err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	toAccount, err := s.repo.ValidateAccount(transferInput.ToAgency, transferInput.ToAccount)
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, err.Error())
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	fromAccountId := claims["id"].(string)
	fromAccount, err := s.repo.GetAccount(fromAccountId)
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	if fromAccountId == toAccount.Id {
		EncodeJson(w, http.StatusBadRequest, "Same account transfer is not allowed")
		return
	}
	if debitErr := s.repo.DebitAccount(fromAccount, transferInput.Amount); debitErr != nil {
		EncodeJson(w, http.StatusBadRequest, debitErr)
		return
	}
	s.repo.CreditAccount(toAccount, transferInput.Amount)
	EncodeJson(w, http.StatusOK, "OK")
}

func (s *ApiServer) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	account, err := s.repo.GetAccount(idParam)
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	if account.Id != claims["id"] {
		EncodeJson(w, http.StatusForbidden, NewApiError("Invalid Account"))
		return
	}
	EncodeJson(w, http.StatusOK, account)
}

func (s *ApiServer) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountInput := new(CreateAccountInput)
	if err := json.NewDecoder(r.Body).Decode(accountInput); err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}

	account := NewAccount(accountInput.FirstName, accountInput.LastName, accountInput.Document)
	if err := s.repo.CreateAccount(account); err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	EncodeJson(w, http.StatusOK, account)
}

func (s *ApiServer) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	account, err := s.repo.GetAccount(idParam)
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	err = s.repo.DeleteAccount(account.Id)
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	EncodeJson(w, http.StatusOK, account)
}

func (s *ApiServer) SigninHandler(w http.ResponseWriter, r *http.Request) {
	signintInput := new(SiginInput)
	if err := json.NewDecoder(r.Body).Decode(signintInput); err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	account, err := s.repo.ValidateAccount(signintInput.Agency, signintInput.Number)
	if err != nil {
		EncodeJson(w, http.StatusBadRequest, NewApiError(err.Error()))
		return
	}
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"accoutNumber": account.Number,
		"id":           account.Id,
	})

	EncodeJson(w, http.StatusOK, NewSiginOutput(tokenString))
}

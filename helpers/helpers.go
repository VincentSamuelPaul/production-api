package helpers

import (
	"encoding/json"
	"net/http"
	"time"

	structTypes "github.com/VincentSamuelPaul/production-api/types"
	"golang.org/x/crypto/bcrypt"
)

type UserAccount struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password_hash"`
	Created_at    time.Time `json:"created_at"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (a *UserAccount) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password_hash), []byte(pw)) == nil
}

func NewAccount(username, email, password string) (*structTypes.UserAccount, error) {
	encPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &structTypes.UserAccount{
		Username:      username,
		Email:         email,
		Password_hash: string(encPw),
	}, nil
}

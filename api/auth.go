package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VincentSamuelPaul/production-api/helpers"
	structTypes "github.com/VincentSamuelPaul/production-api/types"
)

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return helpers.WriteJSON(w, http.StatusForbidden, structTypes.ErrorMSG{Error: fmt.Sprintf("%s, method not allowed", r.Method)})
	}
	var user structTypes.UserAccount
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}
	account, err := helpers.NewAccount(user.Username, user.Email, user.Password_hash)
	if err != nil {
		return err
	}
	account.Created_at = time.Now()
	if err := s.store.CreateUser(account); err != nil {
		return helpers.WriteJSON(w, http.StatusAccepted, structTypes.ErrorMSG{Error: err.Error()})
	}
	return helpers.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "success"})
}

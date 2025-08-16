package api

import (
	"net/http"
	"strconv"

	"github.com/VincentSamuelPaul/production-api/helpers"
	structTypes "github.com/VincentSamuelPaul/production-api/types"
	"github.com/gorilla/mux"
)

func (s *APIServer) handleTest(w http.ResponseWriter, r *http.Request) error {
	return helpers.WriteJSON(w, http.StatusOK, "hello")
}

// PRODUCT FUNCTIONS

func (s *APIServer) handleGetAllProducts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Forbidden"})
	}
	data, err := s.store.GetAllProducts()
	if err != nil {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
	}
	return helpers.WriteJSON(w, http.StatusOK, data)
}

func (s *APIServer) handleGetProductByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Forbidden"})
	}
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Invalid id type"})
	}
	data, err := s.store.GetProductByID(id)
	if err != nil {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
	}
	return helpers.WriteJSON(w, http.StatusOK, data)
}

// CART FUNCTIONS

func (s *APIServer) handleCart(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["userid"]
	userid, err := strconv.Atoi(idStr)
	if err != nil {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Invalid id type"})
	}
	if r.Method == "GET" {
		data, err := s.store.GetCartByID(userid)
		if err != nil {
			return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
		return helpers.WriteJSON(w, http.StatusOK, data)
	}
	return nil
}

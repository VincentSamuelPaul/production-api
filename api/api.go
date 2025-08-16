package api

import (
	"log"
	"net/http"

	"github.com/VincentSamuelPaul/production-api/helpers"
	structTypes "github.com/VincentSamuelPaul/production-api/types"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      structTypes.Storage
}

func NewAPIServer(listenAddr string, store structTypes.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (server *APIServer) Run() {
	router := mux.NewRouter()
	// TEST
	router.HandleFunc("/test", makeHTTPHandleFunc(server.handleTest))
	// AUTH ROUTES
	router.HandleFunc("/createuser", makeHTTPHandleFunc(server.handleCreateUser))
	// PRODUCT ROUTES
	router.HandleFunc("/products", makeHTTPHandleFunc(server.handleGetAllProducts))
	router.HandleFunc("/products/{id}", makeHTTPHandleFunc(server.handleGetProductByID))
	// CART ROUTES
	router.HandleFunc("/cart/{userid}", makeHTTPHandleFunc(server.handleCart))

	log.Printf("API running on: %s\n", server.listenAddr)

	http.ListenAndServe(server.listenAddr, router)
}

func makeHTTPHandleFunc(f structTypes.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
	}
}

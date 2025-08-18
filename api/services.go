package api

import (
	"encoding/json"
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
	if r.Method == "POST" {
		var req struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}
		err := s.store.AddToCart(userid, req.ProductID, req.Quantity)
		if err != nil {
			return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
		return helpers.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "added to cart"})
	}
	if r.Method == "DELETE" {
		idStr := mux.Vars(r)["productid"]
		if idStr != "" {
			productid, err := strconv.Atoi(idStr)
			if err != nil {
				return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Invalid id type"})
			}
			err = s.store.DeleteFromCart(userid, productid)
			if err != nil {
				return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
			}
			return helpers.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "item removed from cart"})
		} else {
			err := s.store.EmptyCart(userid)
			if err != nil {
				return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
			}
			return helpers.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "cart empty"})
		}
	}
	return nil
}

// ORDER FUNCTIONS

func (s *APIServer) handleOrders(w http.ResponseWriter, r *http.Request) error {

	idStr := mux.Vars(r)["userid"]
	userid, err := strconv.Atoi(idStr)
	if err != nil {
		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Invalid userid type"})
	}

	status := mux.Vars(r)["status"]

	if r.Method == "GET" {
		data, err := s.store.GetAllOrdersByUserID(userid)
		if err != nil {
			return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
		return helpers.WriteJSON(w, http.StatusOK, data)
	}

	// if r.Method == "GET" && orderStr != "" && idStr == "" {
	// 	orderid, err := strconv.Atoi(orderStr)
	// 	if err != nil {
	// 		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: "Invalid orderid type"})
	// 	}
	// 	data, err := s.store.GetOrderByID(orderid)
	// 	if err != nil {
	// 		return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
	// 	}
	// 	return helpers.WriteJSON(w, http.StatusOK, data)
	// }

	if r.Method == "POST" {
		var orders []structTypes.OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&orders); err != nil {
			return err
		}
		err := s.store.CreateOrder(userid, orders)
		if err != nil {
			return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
		return helpers.WriteJSON(w, http.StatusOK, map[string]string{"status": "orders placed"})
	}

	if r.Method == "PUT" {
		err := s.store.UpdateOrderStatus(userid, status)
		if err != nil {
			return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
		return helpers.WriteJSON(w, http.StatusOK, map[string]string{"status": "orders status updated"})
	}

	if r.Method == "DELETE" {
		err = s.store.DeleteOrder(userid)
		if err != nil {
			return helpers.WriteJSON(w, http.StatusBadRequest, structTypes.ErrorMSG{Error: err.Error()})
		}
		return helpers.WriteJSON(w, http.StatusOK, map[string]string{"status": "order deleted"})
	}

	return nil
}

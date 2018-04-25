package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getProducts(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)

		products, err := s.GetProducts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]interface{}{
				"msg":   "could not get products",
				"error": err.Error(),
			})
			return
		}

		enc.Encode(map[string]interface{}{
			"products": products,
		})
	}
}

func getProductByID(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		product, err := s.GetProductByID(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(map[string]interface{}{
				"msg":   "product not found",
				"error": err.Error(),
			})
			return
		}

		enc.Encode(map[string]interface{}{
			"product": product,
		})
	}
}

func decrementProductQuantityByID(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		err := s.DecrementProductQuantityByID(id)
		if err != nil {
			if err.Error() == "product not found" {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"msg":   "cannot decrement quantity of product that does not exist",
					"error": err.Error(),
				})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"msg":   "could not decrement",
				"error": err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

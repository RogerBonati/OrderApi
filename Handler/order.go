package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RogerBonati/OrderApi/Model"
	"github.com/RogerBonati/OrderApi/Repository/Order"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"time"
)

type Order struct {
	Repo *order.RedisRepo
}

// create the http interface

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")
	var body struct {
		customerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	order := model.Order{
		OrderId:    rand.Uint64(),
		CustomerId: body.customerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err := o.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("Failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}
func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all orders")
}
func (o *Order) GetById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an order by ID")
}
func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")
}
func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")
}

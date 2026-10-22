
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type OrderRepository interface {
	Save(ord order) error
	FindByID(id int) (order, error)
	FindAll() ([]order, error)
	UpdateStatus(id int, next statustype) error
}

type FileOrderRepository struct {
	filename string
}

func NewFileOrderRepository(filename string) *FileOrderRepository {
	return &FileOrderRepository{filename: filename}
}

// Exported fields allow JSON to save the data.
type storedCustomer struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type storedItem struct {
	Name        string      `json:"name"`
	Quantity    int         `json:"quantity"`
	ServiceType serviceType `json:"service_type"`
	UnitPrice   float64     `json:"unit_price"`
	TotalPrice  float64     `json:"total_price"`
}

type storedOrder struct {
	ID         int            `json:"id"`
	Customer   storedCustomer `json:"customer"`
	Items      []storedItem   `json:"items"`
	TotalItems int            `json:"total_items"`
	TotalPrice float64        `json:"total_price"`
	Status     statustype     `json:"status"`
}

func toStoredOrder(o order) storedOrder {
	result := storedOrder{
		ID:         o.id,
		Customer:   storedCustomer{Name: o.customer.name, Phone: o.customer.phone},
		TotalItems: o.totalItems,
		TotalPrice: o.totalPrice,
		Status:     o.status,
	}

	for _, it := range o.items {
		result.Items = append(result.Items, storedItem{
			Name:        it.name,
			Quantity:    it.quantity,
			ServiceType: it.serviceType,
			UnitPrice:   it.unitPrice,
			TotalPrice:  it.totalPrice,
		})
	}

	return result
}

func fromStoredOrder(o storedOrder) order {
	result := order{
		id: o.ID,
		customer: customer{
			name:  o.Customer.Name,
			phone: o.Customer.Phone,

			
		},
		totalItems: o.TotalItems,
		totalPrice: o.TotalPrice,
		status:     o.Status,
	}

	for _, it := range o.Items {
		result.items = append(result.items, item{
			name:        it.Name,
			quantity:    it.Quantity,
			serviceType: it.ServiceType,
			unitPrice:   it.UnitPrice,
			totalPrice:  it.TotalPrice,
		})
	}

	return result
}

func (r *FileOrderRepository) FindAll() ([]order, error) {
	data, err := os.ReadFile(r.filename)
	if os.IsNotExist(err) {
		return []order{}, nil
	}
	if err != nil {
		return nil, err
	}

	var saved []storedOrder
	if err := json.Unmarshal(data, &saved); err != nil {
		return nil, err
	}

	orders := make([]order, 0, len(saved))
	for _, o := range saved {
		orders = append(orders, fromStoredOrder(o))
	}

	return orders, nil
}

func (r *FileOrderRepository) Save(ord order) error {
	orders, err := r.FindAll()
	if err != nil {
		return err
	}

	found := false

	for i := range orders {
		if orders[i].id == ord.id {
			orders[i] = ord
			found = true
			break
		}
	}

	if !found {
		orders = append(orders, ord)
	}

	saved := make([]storedOrder, 0, len(orders))
	for _, o := range orders {
		saved = append(saved, toStoredOrder(o))
	}

	data, err := json.MarshalIndent(saved, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filename, data, 0644)
}

func (r *FileOrderRepository) FindByID(id int) (order, error) {
	orders, err := r.FindAll()
	if err != nil {
		return order{}, err
	}

	for _, o := range orders {
		if o.id == id {
			return o, nil
		}
	}

	return order{}, fmt.Errorf("order not found: %d", id)
}

func (r *FileOrderRepository) UpdateStatus(id int, next statustype) error {
	ord, err := r.FindByID(id)
	if err != nil {
		return err
	}

	if err := ord.UpdateStatus(next); err != nil {
		return err
	}

	return r.Save(ord)
}

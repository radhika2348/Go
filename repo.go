
package main

import (
	"encoding/json"
	//"fmt"
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

func (r *FileOrderRepository) FindAll() ([]order, error) {
	data, err := os.ReadFile(r.filename)
	if os.IsNotExist(err) {
		return []order{}, nil
	}
	if err != nil {
		return nil, err
	}

	var saved []order
	 err = json.Unmarshal(data, &saved); if err != nil {
		return nil, err
	}

	orders := make([]order, 0, len(saved))
	for _, o := range saved {
		orders = append(orders, order(o))
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
		if orders[i].ID == ord.ID {
			orders[i] = ord
			found = true
			break
		}
	}

	if !found {
		return ErrorOrderNotFound
	}

	saved := make([]order, 0, len(orders))
	for _, o := range orders {
		saved = append(saved, order(o))
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
		if o.ID== id {
			return o, nil
		}
	}

	return order{}, ErrorOrderNotFound
}

func (r *FileOrderRepository) UpdateStatus(id int, next statustype) error {
	ord, err := r.FindByID(id)
	if err != nil {
		return ErrorOrderNotFound
	}

	if err := ord.UpdateStatus(next); err != nil {
		return ErrInvalidOrderStatus
	}

	return r.Save(ord)
}
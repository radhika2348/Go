package main

import (
	"encoding/json"
	//"fmt"
	"os"
	"path/filepath"
	//"errors"
   
	//"strings"
	"testing"
)

func newTestService(t *testing.T) *orderService {
	t.Helper()

	repo := NewFileOrderRepository(
		filepath.Join(t.TempDir(), "orders.json"),
	)

	return NewLaundryService(repo)
}

func loadExpectedOrder(t *testing.T) order{
	t.Helper()

	data, err := os.ReadFile("testdata/expected_order.json")
	if err != nil {
		t.Error("Error reading expected JSON:", err)
		return order{}
	}

	var expected order

	err = json.Unmarshal(data, &expected)
	if err != nil {
		t.Error("Error decoding expected JSON")
		return order{}
	}

	return expected
}

func TestCustomerValidate(t *testing.T) {
	c := customer{
		Name:  "Ram",
		Phone: "555-0199",
	}

	err := c.validate()

	if err != nil {
		t.Error("Expected valid customer")
	}
}

func TestCustomerEmptyName(t *testing.T) {
	c := customer{
		Name:  "   ",
		Phone: "555-0199",
	}

	err := c.validate()

	if err == nil {
		t.Error("Expected error for empty customer name")
	}
}

func TestCustomerEmptyPhone(t *testing.T) {
	c := customer{
		Name:  "Ram",
		Phone: "   ",
	}

	err := c.validate()

	if err == nil {
		t.Error("Expected error for empty customer phone")
	}
}

func TestItemCalculatePrice(t *testing.T) {
	it := item{
		Name:        "Kurta",
		Quantity:    3,
		ServiceType: washAndFold,
	}

	actual, err := it.calculatePrice()

	if err != nil {
		t.Error("Unexpected error")
		return
	}

	if actual.UnitPrice != 2.50 {
		t.Errorf("Expected unit price 2.50, got %.2f", actual.UnitPrice)
	}

	if actual.TotalPrice != 7.50 {
		t.Errorf("Expected total price 7.50, got %.2f", actual.TotalPrice)
	}
}

func TestItemInvalidQuantity(t *testing.T) {
	it := item{
		Name:        "Shirts",
		Quantity:    0,
		ServiceType: washAndFold,
	}

	_, err := it.calculatePrice()

	if err == nil {
		t.Error("Expected error for quantity less than 1")
	}
}

func TestItemUnsupportedService(t *testing.T) {
	it := item{
		Name:        "Shirts",
		Quantity:    2,
		ServiceType: serviceType("Express Cleaning"),
	}

	_, err := it.calculatePrice()
    
	
	if err == nil {
		t.Error("Expected error for unsupported service")
	}
}

func TestItemEmptyName(t *testing.T) {
	it := item{
		Name:        "   ",
		Quantity:    2,
		ServiceType: washAndFold,
	}

	_, err := it.calculatePrice()

	if err == nil {
		t.Error("Expected error for empty item name")
	}
}

func TestCreateEmptyOrder(t *testing.T) {
	svc := newTestService(t)

	c := customer{
		Name:  "Ram",
		Phone: "555-0199",
	}

	_, err := svc.createOrder(c, []item{})

	if err == nil {
		t.Error("Expected error when order has no items")
	}

	orders, err := svc.GetAllOrders()

	if err != nil {
		t.Error("Error retrieving orders:", err)
		return
	}

	if len(orders) != 0 {
		t.Errorf("Expected 0 orders, got %d", len(orders))
	}
}

func TestOrderStatusTransitions(t *testing.T) {
	svc := newTestService(t)

	expectedJSON := loadExpectedOrder(t)
	expected := order(expectedJSON)

	created, err := svc.createOrder(
		expected.Customer,
		expected.Items,
	)

	

	if err != nil {
		t.Error("Error creating order:", err)
		return
	}

	statuses := []statustype{
		WASHING,
		READY,
		COLLECTED,
	}

	for _, status := range statuses {
		updated, err := svc.updateOrderStatus(created.ID, status)

		if err != nil {
			t.Error("Error updating status:", err)
			return
		}

		if updated.Status != status {
			t.Errorf(
				"Expected status %q, got %q",
				status,
				updated.Status,
			)
		}
	}
}

func TestRejectInvalidTransition(t *testing.T) {
	svc := newTestService(t)

	expectedJSON := loadExpectedOrder(t)
	expected := order(expectedJSON)

	created, err := svc.createOrder(
		expected.Customer,
		expected.Items,
	)

	if err != nil {
		//fmt.Print(err)
		t.Error("Error creating order:", err)
		
		return
	}

	_, err = svc.updateOrderStatus(created.ID, READY)

	if err == nil {
		t.Error("Expected error for invalid status transition")
	}

	saved, err := svc.FindByID(created.ID)

	if err != nil {
		t.Error("Error finding order:", err)
		return
	}

	if saved.Status != RECEIVED {
		t.Errorf(
			"Expected status %q, got %q",
			RECEIVED,
			saved.Status,
		)
	}
}

func TestUpdateNonexistentOrder(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.updateOrderStatus(9999, WASHING)

	if err == nil {
		t.Error("Expected error for nonexistent order")
	}
}

func TestInvalidStatus(t *testing.T) {
	err := validateStatus(statustype("Delivered"))

	 

	if err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestEmptyOrderErrorMessage(t *testing.T) {
	svc := newTestService(t)

	c := customer{
		Name:  "Ram",
		Phone: "555-0199",
	}

	_, err := svc.createOrder(c, nil)

	if err == nil {
		t.Error("Expected error when order has no items")
		return
	}

}

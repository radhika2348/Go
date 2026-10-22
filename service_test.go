package main

import (
	 "encoding/json"
	 "os"
	"path/filepath"
	
	"strings"
	"testing"
)

 func newTestService(t *testing.T) *orderService {
 	t.Helper()

 	repo := NewFileOrderRepository(
 		filepath.Join(t.TempDir(), "orders.json"),
 	)

 	return NewLaundryService(repo)
 }

 func loadExpectedOrder(t *testing.T) storedOrder {
     t.Helper()

     data, err := os.ReadFile("testdata/expected_order.json")
     if err != nil {
         t.Error("Error reading expected JSON:", err)
         return storedOrder{}
     }

     var expected storedOrder

     err = json.Unmarshal(data, &expected)
     if err != nil {
         t.Error("Error decoding expected JSON:", err)
         return storedOrder{}
     }

     return expected
 }

  func TestCustomerValidate(t *testing.T) {
	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}

	err := c.validate()

	if err != nil {
		t.Error("Expected valid customer, got:", err)
	}
}


func TestCustomerEmptyName(t *testing.T) {
	c := customer{
		name:  "   ",
		phone: "555-0199",
	}

	err := c.validate()

	if err == nil {
		t.Error("Expected error for empty customer name")
	}
}


  func TestCustomerEmptyPhone(t *testing.T) {
	c := customer{
		name:  "Ram",
		phone: "   ",
	}

	err := c.validate()

	if err == nil {
		t.Error("Expected error for empty customer phone")
	}
}

func TestItemCalculatePrice(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    3,
		serviceType: washAndFold,
	}

	   actual, err := it.calculatePrice()

	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if actual.unitPrice != 2.50 {
		t.Errorf("Expected unit price 2.50, got %.2f", actual.unitPrice)
	}

	if actual.totalPrice != 7.50 {
		t.Errorf("Expected total price 7.50, got %.2f", actual.totalPrice)
	}
}


func TestItemInvalidQuantity(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    0,
		serviceType: washAndFold,
	}

	_, err := it.calculatePrice()

	if err == nil {
		t.Error("Expected error for quantity less than 1")
	}
}



func TestItemUnsupportedService(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    2,
		serviceType: serviceType("Express Cleaning"),
	}

	_, err := it.calculatePrice()

	if err == nil {
		t.Error("Expected error for unsupported service")
	}
}

func TestItemEmptyName(t *testing.T) {
	it := item{
		name:        "   ",
		quantity:    2,
		serviceType: washAndFold,
	}

	_, err := it.calculatePrice()

	if err == nil {
		t.Error("Expected error for empty item name")
	}
}



func TestCreateEmptyOrder(t *testing.T) {
	svc := newTestService(t)

	c := customer{
		name:  "Ram",
		phone: "555-0199",
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
	expected := fromStoredOrder(expectedJSON)

	created, err := svc.createOrder(
		expected.customer,
		expected.items,
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
		updated, err := svc.updateOrderStatus(created.id, status)

		if err != nil {
			t.Error("Error updating status:", err)
			return
		}

		if updated.status != status {
			t.Errorf(
				"Expected status %q, got %q",
				status,
				updated.status,
			)
		}
	}
}

func TestRejectInvalidTransition(t *testing.T) {
	svc := newTestService(t)

	expectedJSON := loadExpectedOrder(t)
	expected := fromStoredOrder(expectedJSON)

	created, err := svc.createOrder(
		expected.customer,
		expected.items,
	)

	if err != nil {
		t.Error("Error creating order:", err)
		return
	}

	_, err = svc.updateOrderStatus(created.id, READY)

	if err == nil {
		t.Error("Expected error for invalid status transition")
	}

	saved, err := svc.FindByID(created.id)

	if err != nil {
		t.Error("Error finding order:", err)
		return
	}

	if saved.status != RECEIVED {
		t.Errorf(
			"Expected status %q, got %q",
			RECEIVED,
			saved.status,
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
		name:  "Ram",
		phone: "555-0199",
	}

	_, err := svc.createOrder(c, nil)

	if err == nil {
		t.Error("Expected error when order has no items")
		return
	}

	 
	  if !strings.Contains(
		strings.ToLower(err.Error()),
		"at least one item",
	) {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}
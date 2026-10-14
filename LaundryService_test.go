package main

import "testing"

func TestValidateCustomer(t *testing.T) {
	c := customer{
		name:  "Alice",
		phone: "555-0199",
	}

	err := validateCustomer(c)

	if err != nil {
		t.Fatalf("expected valid customer, got error: %v", err)
	}
}

func TestValidateCustomerMissingName(t *testing.T) {
	c := customer{
		name:  "",
		phone: "555-0199",
	}

	err := validateCustomer(c)

	if err == nil {
		t.Fatal("expected an error for missing customer name")
	}
}

func TestValidateItems(t *testing.T) {
	items := []item{
		{
			name:        "Shirts",
			quantity:    3,
			serviceType: washAndFold,
		},
	}

	err := validateItems(items)

	if err != nil {
		t.Fatalf("expected valid items, got error: %v", err)
	}
}

func TestValidateItemsEmpty(t *testing.T) {
	err := validateItems(nil)

	if err == nil {
		t.Fatal("expected an error for an empty item list")
	}
}

func TestValidateItemsInvalidQuantity(t *testing.T) {
	items := []item{
		{
			name:        "Shirts",
			quantity:    0,
			serviceType: washAndFold,
		},
	}

	err := validateItems(items)

	if err == nil {
		t.Fatal("expected an error for invalid quantity")
	}
}

func TestCreateOrder(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Alice",
		phone: "555-0199",
	}

	items := []item{
		{name: "Shirts", quantity: 3, serviceType: washAndFold},
		{name: "Coat", quantity: 1, serviceType: dryCleaning},
	}

	ord, err := svc.createOrder(c, items)

	if err != nil {
		t.Fatalf("expected order creation to succeed, got error: %v", err)
	}

	if ord.id != "ORD-1001" {
		t.Errorf("expected order ID ORD-1001, got %s", ord.id)
	}

	// 3 shirts * $2.50 + 1 coat * $8.00 = $15.50
	if ord.totalPrice != 15.50 {
		t.Errorf("expected total price $15.50, got $%.2f", ord.totalPrice)
	}

	if ord.totalItems != 4 {
		t.Errorf("expected 4 total pieces, got %d", ord.totalItems)
	}

	if len(ord.items) != 2 {
		t.Errorf("expected 2 item entries, got %d", len(ord.items))
	}
}

func TestCreateOrderInvalidCustomer(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "",
		phone: "555-0199",
	}

	items := []item{
		{name: "Shirts", quantity: 1, serviceType: washAndFold},
	}

	_, err := svc.createOrder(c, items)

	if err == nil {
		t.Fatal("expected an error for invalid customer")
	}

	if len(svc.getAllOrders()) != 0 {
		t.Error("expected no order to be saved")
	}
}

func TestGetAllOrders(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Alice",
		phone: "555-0199",
	}

	items := []item{
		{name: "Shirts", quantity: 1, serviceType: washAndFold},
	}

	_, err := svc.createOrder(c, items)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	orders := svc.getAllOrders()

	if len(orders) != 1 {
		t.Errorf("expected 1 saved order, got %d", len(orders))
	}
}
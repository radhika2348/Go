
package main

import (
	"strings"
	"testing"
)

// 1. Test valid customer
func TestCustomerValidate(t *testing.T) {
	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}

	err := c.validate()

	if err != nil {
		t.Fatalf("expected valid customer, got error: %v", err)
	}
}

// 2. Test empty customer name
func TestCustomerEmptyName(t *testing.T) {
	c := customer{
		name:  "   ",
		phone: "555-0199",
	}

	err := c.validate()

	if err == nil {
		t.Fatal("expected error for empty customer name")
	}
}

// 3. Test empty phone number
func TestCustomerEmptyPhone(t *testing.T) {
	c := customer{
		name:  "Ram",
		phone: " ",
	}

	err := c.validate()

	if err == nil {
		t.Fatal("expected error for empty phone number")
	}
}

// 4. Test valid item
func TestItemValidate(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    3,
		serviceType: washAndFold,
		//status:      received,
	}

	err := it.validate()

	if err != nil {
		t.Fatalf("expected valid item, got error: %v", err)
	}
}

// 5. Test invalid item quantity
func TestItemInvalidQuantity(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    0,
		serviceType: washAndFold,
		//status:      received,
	}

	err := it.validate()

	if err == nil {
		t.Fatal("expected error for zero quantity")
	}
}

// 6. Test unsupported service
func TestItemUnsupportedService(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    2,
		serviceType: serviceType("Express Wash"),
		//status:      received,
	}

	err := it.validate()

	if err == nil {
		t.Fatal("expected error for unsupported service")
	}
}

// 7. Test item price calculation
func TestItemCalculatePrice(t *testing.T) {
	it := item{
		name:        "Shirts",
		quantity:    3,
		serviceType: washAndFold,
		//status:      received,
	}

	pricedItem, err := it.calculatePrice()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pricedItem.unitPrice != 2.50 {
		t.Errorf("expected unit price 2.50, got %.2f",
			pricedItem.unitPrice)
	}

	if pricedItem.totalPrice != 7.50 {
		t.Errorf("expected total price 7.50, got %.2f",
			pricedItem.totalPrice)
	}
}

// 8. Test order creation and total calculation
func TestCreateOrder(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}

	items := []item{
		{
			name:        "Shirts",
			quantity:    3,
			serviceType: washAndFold,
		},
		{
			name:        "Coat",
			quantity:    1,
			serviceType: dryCleaning,
		},
	}
    nextID:=1001
	ord, err := svc.createOrder(c, items,nextID)

	if err != nil {
		t.Fatalf("unexpected error creating order: %v", err)
	}

	if ord.id != 1001 {
		t.Errorf("expected ID ORD-1001, got %d", ord.id)
	}

	if ord.totalItems != 4 {
		t.Errorf("expected 4 items, got %d", ord.totalItems)
	}

	if ord.totalPrice != 15.50 {
		t.Errorf("expected total 15.50, got %.2f", ord.totalPrice)
	}

	if ord.status != received {
		t.Errorf("expected Received status, got %s", ord.status)
	}

	if len(svc.getAllOrders()) != 1 {
		t.Errorf("expected 1 saved order, got %d",
			len(svc.getAllOrders()))
	}
}

// 9. Test rejection of an empty order
func TestCreateEmptyOrder(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}

	_, err := svc.createOrder(c, []item{},1001)

	if err == nil {
		t.Fatal("expected error for empty order")
	}

	if len(svc.getAllOrders()) != 0 {
		t.Fatal("invalid order should not be saved")
	}
}

// 10. Test valid status transitions
func TestOrderStatusTransitions(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}

	items := []item{
		{
			name:        "Shirts",
			quantity:    2,
			serviceType: washAndFold,
		},
	}
	nextID:=svc.nextID

	ord, err := svc.createOrder(c, items,nextID)
	if err != nil {
		t.Fatalf("createOrder failed: %v", err)
	}

	ord, err = svc.updateOrderStatus(ord.id, washing)
	if err != nil {
		t.Fatalf("Received -> Washing should succeed: %v", err)
	}

	if ord.status != washing {
		t.Errorf("expected Washing, got %s", ord.status)
	}

	ord, err = svc.updateOrderStatus(ord.id, ready)
	if err != nil {
		t.Fatalf("Washing -> Ready should succeed: %v", err)
	}

	if ord.status != ready {
		t.Errorf("expected Ready, got %s", ord.status)
	}
}

// 11. Test rejection of an invalid status transition
func TestRejectInvalidTransition(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}

	items := []item{
		{
			name:        "Shirts",
			quantity:    2,
			serviceType: washAndFold,
		},
	}
    nextID:=svc.nextID
	ord, err := svc.createOrder(c, items,nextID)
	if err != nil {
		t.Fatalf("createOrder failed: %v", err)
	}

	_, err = svc.updateOrderStatus(ord.id, ready)

	// if err != nil {
	// 	return fmt.Printf("Error")
	// }
	// t.Fatalf("not allowed")

	// The failed update must not change the saved order.
	savedOrders := svc.getAllOrders()

	if savedOrders[0].status != received {
		t.Errorf("expected status to remain Received, got %s",
			savedOrders[0].status)
	}
}

// 12. Test nonexistent order ID
func TestUpdateNonexistentOrder(t *testing.T) {
	svc := newOrderService()

	_, err := svc.updateOrderStatus(1002, washing)

	if err == nil {
		t.Fatal("expected error for nonexistent order")
	}
}

// 13. Test unsupported status
func TestInvalidStatus(t *testing.T) {
	err := validateStatus(statustype("Delivered"))

	if err == nil {
		t.Fatal("expected error for unsupported status")
	}
}

// 14. Test error message for empty order
func TestEmptyOrderErrorMessage(t *testing.T) {
	svc := newOrderService()

	c := customer{
		name:  "Ram",
		phone: "555-0199",
	}
    nextID:=svc.nextID
	_, err := svc.createOrder(c, []item{},nextID)

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "at least one item") {
		t.Errorf("unexpected error message: %v", err)
	}
}

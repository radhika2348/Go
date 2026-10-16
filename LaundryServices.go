
package main

import (
	"fmt"
	"strings"
	"os"
	"log"
)


type statustype string

const (
	received  statustype = "Received"
	washing   statustype = "Washing"
	ready     statustype = "Ready"
	collected statustype = "Collected"
)
func validateStatus(status statustype) error {
 	switch status {
 	case received, washing, ready, collected:
 		return nil
 	default:
 		return fmt.Errorf("invalid status: %s", status)
 	}
 }


type serviceType string

const (
	washAndFold serviceType = "Wash & Fold"
	washAndIron serviceType = "Wash & Iron"
	dryCleaning serviceType = "Dry Cleaning"
)

var serviceRates = map[serviceType]float64{
	washAndFold: 2.50,
	washAndIron: 4.00,
	dryCleaning: 8.00,
}

// 2. CUSTOMER
// Responsibility: Store and validate customer details

type customer struct {
	name  string
	phone string
}

func (c customer) validate() error {
	if strings.TrimSpace(c.name) == "" {
		return fmt.Errorf("customer name is required")
	}

	if strings.TrimSpace(c.phone) == "" {
		return fmt.Errorf("customer phone number is required")
	}

	return nil
}


// 3. ITEM
// Responsibility: Validate an item and calculate its price

type item struct {
	name        string
	quantity    int
	serviceType serviceType
	unitPrice   float64
	totalPrice  float64
	
}

func (it item) validate() error {
	if strings.TrimSpace(it.name) == "" {
		return fmt.Errorf("item name cannot be empty")
	}

	if it.quantity < 1 {
		return fmt.Errorf("item quantity must be at least 1")
	}

	if _, ok := serviceRates[it.serviceType]; !ok {
		return fmt.Errorf("unsupported service type: %s", it.serviceType)
	}

	

	return nil
}
func (s *orderService)validID(nextID int)error{
	for i := range s.orders {
		if s.orders[i].id == nextID{
           return fmt.Errorf("ID exist")
		}
	}
	return nil
}

func (it item) calculatePrice() (item, error) {
	if err := it.validate(); err != nil {
		return item{}, err
	}

	rate := serviceRates[it.serviceType]

	it.unitPrice = rate
	it.totalPrice = float64(it.quantity) * rate

	return it, nil
}

// 4. ORDER
// Responsibility: Validate an order and manage its status

type order struct {
	id         int
	customer   customer
	items      []item
	totalItems int
	totalPrice float64
	status     statustype
}

func (o order) validate() error {
	if err := o.customer.validate(); err != nil {
		return err
	}

	if len(o.items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	for i, it := range o.items {
		if err := it.validate(); err != nil {
			return fmt.Errorf("item #%d: %w", i+1, err)
		}
	}
	if err:=newOrderService().validID();err!=nil{
		return err
	}

	return nil
}


func (o order) canTransitionTo(next statustype) bool {
	switch o.status {
	case received:
		return next == washing
	case washing:
		return next == ready
	case ready:
		return next == collected
	default:
		return false
	}
}

func (o *order) updateStatus(next statustype) error {
	if err := validateStatus(next); err != nil {
		return err
	}

	if !o.canTransitionTo(next) {
		return fmt.Errorf(
			"invalid status transition: %s -> %s",
			o.status, next,
		)
	}

	o.status = next
	return nil
}

func (s *orderService) updateOrderStatus(
	orderID int,
	next statustype,
) (order, error) {

	for i := range s.orders {
		if s.orders[i].id == orderID {
			if err := s.orders[i].updateStatus(next); err != nil {
				return order{}, err
			}

			return s.orders[i], nil
		}
	}

	return order{}, fmt.Errorf("order not found: %s", orderID)
}


// 6. ORDER SERVICE
// Responsibility: Create, store, retrieve and update orders

type orderService struct {
	orders []order
	nextID int
}

func newOrderService() *orderService {
	return &orderService{
		orders: []order{},
		nextID: 1001,
	}
}

func (s *orderService) createOrder(
	c customer,
	items []item,
	nextID int,
) (order, error) {

	// New orders always begin in the Received state.
	

	newOrder := order{
		customer: c,
		items:    items,
		status:   received,
	}

	// The order validates its own data.
	if err := newOrder.validate(); err != nil {
		return order{}, err
	}

	// Each item calculates its own price.
	var totalBill float64
	var totalCount int
	var pricedItems []item

	for _, it := range newOrder.items {
		pricedItem, err := it.calculatePrice()
		if err != nil {
			return order{}, err
		}

		pricedItems = append(pricedItems, pricedItem)
		totalBill += pricedItem.totalPrice
		totalCount += pricedItem.quantity
	}

	newOrder.items = pricedItems
	newOrder.totalPrice = totalBill
	newOrder.totalItems = totalCount
	newOrder.id = s.nextID

	// Incrementing the ID for the next order.
	s.nextID++

	// Save the order in memory.
	s.orders = append(s.orders, newOrder)

	return newOrder, nil
}

func (s *orderService) getAllOrders() []order {
	return s.orders
}




func (s *orderService) printOrdersByStatus(status statustype) {
	fmt.Printf("\n======== %s ORDERS ========\n", status)

	found := false

	for _, ord := range s.orders {
		if ord.status == status {
			fmt.Printf(
				"Order ID: %s | Customer: %s | Status: %s\n",
				ord.id,
				ord.customer.name,
				ord.status,
			)
			found = true
		}
	}

	if !found {
		fmt.Printf("No orders with status %s\n", status)
	}

	fmt.Println("==============================")
}

func (s *orderService) printDailyOrderStatus(){
	fmt.Println("\n=============DailyOrderStatus=============")
	count1,count2,count3,count4:=0,0,0,0
    for _,ord:=range s.orders{
		if ord.status==received{
			count1++
		}
		
		if ord.status==washing{
			count2++
		}
	
		if ord.status==ready{
			count3++
		}

		if ord.status==collected{
			count4++
		}
	
		fmt.Printf("Received %d\nWashing %d\nReady %d\nCollected %d",count1,count2,count3,count4)

	}

	fmt.Println("\n=======================================")


}

func main() {
	file,err:=os.OpenFile("laundry.txt",os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err!=nil{
		log.Fatal(err)
	}
	defer file.Close()
	svc := newOrderService()
	 
   
	c := customer{
		name:  "Ram",
		phone: "555-019",
	
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
	nextID:=svc.nextID
	// Create an order.
	ord, err := svc.createOrder(c, items,nextID)
	if err != nil {
		fmt.Println("Error creating order:", err)
		return
	}
    allOrders := svc.getAllOrders()
	sv := allOrders
	_,err=fmt.Fprint(file, sv)
	if err!=nil{
		log.Fatal(err)
	}
	// Print the order summary.
	fmt.Println("\n========== ORDER SUMMARY ==========")
	fmt.Println("Order ID:", ord.id)
	fmt.Printf(
		"Customer: %s (%s)\n",
		ord.customer.name,
		ord.customer.phone,
	)
	fmt.Println("-----------------------------------")

	for _, it := range ord.items {
		fmt.Printf(
			"- %s: %d pieces (%s) = $%.2f\n",
			it.name,
			it.quantity,
			it.serviceType,
			it.totalPrice,
		)
	}

	fmt.Println("-----------------------------------")
	fmt.Println("Order Status:", ord.status)
	fmt.Println("Total Items:", ord.totalItems)
	fmt.Printf("TOTAL BILL: $%.2f\n", ord.totalPrice)
	fmt.Println("===================================")

	
	ord, err = svc.updateOrderStatus(ord.id,ready)
	if err != nil {
		fmt.Println("Error updating status:", err)
		return
	}

	fmt.Println("\nUpdated Order Status:", ord.status)

	// Display only orders that are Ready.
	svc.printOrdersByStatus(washing)

	fmt.Println(
		"\nTotal orders saved in system:",
		len(svc.getAllOrders()),
	)

	svc.printDailyOrderStatus()
	

	
}

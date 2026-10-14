package main

import (
	"fmt"
	"strings"
)

type statustype string
const (
	received statustype="Recieved"
	washing statustype="Washing"
	ready statustype="Ready"
	collected statustype="Collected"

)

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

type customer struct {
	name  string
	phone string
}

type item struct {
	name        string
	quantity    int
	serviceType serviceType
	unitPrice   float64
	totalPrice  float64
	status statustype
}

type order struct {
	id         string
	customer   customer
	items      []item
	totalItems int
	totalPrice float64
	statustype statustype
}

func validateCustomer(c customer) error {
	if strings.TrimSpace(c.name) == "" {
		return fmt.Errorf("customer name is required")
	}

	if strings.TrimSpace(c.phone) == "" {
		return fmt.Errorf("customer phone number  is required")
	}

	return nil
}

func validatestatus(s statustype) error {
	if s==received||s==washing||s==ready||s==collected{
		return nil
	
	}
		return fmt.Errorf("invalid status ")
	
}

func validateServiceType(s serviceType) error {
	if _, ok := serviceRates[s]; !ok {
		return fmt.Errorf("unsupported service type: %s", s)
	}

	return nil
}

func validateItems(items []item) error {
	if len(items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	for i, it := range items {
		if strings.TrimSpace(it.name) == "" {
			return fmt.Errorf("item name cannot be empty at item #%d", i+1)
		}

		if it.quantity < 1 {
			return fmt.Errorf("item quantity must be at least 1")
		}

		if err := validateServiceType(it.serviceType); err != nil {
			return err
		}

		if err:=validatestatus(it.status);err!=nil{
			return err
		}	

	}
	return nil
}



func validateOrderRequest(c customer, items []item) error {
	if err := validateCustomer(c); err != nil {
		return err
	}

	if err := validateItems(items); err != nil {
		return err
	}

	return nil
}
func (s *orderService) PrintOrderStatus(status statustype) {
    fmt.Println("\n============ ORDER STATUS ============")
    found := false

    for _, ord := range s.orders {
        if ord.statustype == status {
            fmt.Printf("Order ID: %s is %s!\n", ord.id, ord.statustype)
            found = true
        }
    }

    if !found {
        fmt.Println("No orders found with status:", status)
    }

    fmt.Println("=======================================")
}


func isValidTransaction(current statustype, next statustype) bool {
	switch current {
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

func (s *orderService) UpdateStatus(orderID string, nextStatus statustype) error {
	if err := validatestatus(nextStatus); err != nil {
		return err
	}

	for i := range s.orders {
		if s.orders[i].id == orderID {
			if !isValidTransaction(s.orders[i].statustype, nextStatus) {
				return fmt.Errorf("invalid transition: cannot change status from '%s' to '%s'", s.orders[i].statustype, nextStatus)
			}

			s.orders[i].statustype = nextStatus
			return nil
		}
	}

	return fmt.Errorf("order with ID %s not found", orderID)
}




type orderService struct {
	orders []order
	nextId int
}

func newOrderService() *orderService {
	return &orderService{
		nextId: 1001,
	}
}

func (s *orderService) createOrder(c customer, items []item) (order, error) {
	if err := validateOrderRequest(c, items); err != nil {
		return order{}, err
	}

	var totalBill float64
	var totalCount int
	var finalItems []item


	for _, it := range items {
		rate := serviceRates[it.serviceType]
		itemCost := float64(it.quantity) * rate

		it.unitPrice = rate
		it.totalPrice = itemCost

		finalItems = append(finalItems, it)

		totalBill += itemCost
		totalCount += it.quantity
	}

	orderId := fmt.Sprintf("ORD-%d", s.nextId)
	s.nextId++

	newOrder := order{
		id:         orderId,
		customer:   c,
		items:      finalItems,
		totalItems: totalCount,
		totalPrice: totalBill,
        statustype:washing,
	}

	s.orders = append(s.orders, newOrder)

	return newOrder, nil
}

 func (s *orderService) getAllOrders() []order {
 	return s.orders
 }

 

func main() {
	svc := newOrderService()

	
	c := customer{
		name:  " ram",
		phone: "555-0199",
	}

	items := []item{
		{name: "Shirts", quantity: 3, serviceType: washAndFold,status:washing},
		{name: "Coat", quantity: 1, serviceType: dryCleaning,status:washing},
	}

	

	ord, err := svc.createOrder(c, items)
	if err != nil {
		fmt.Println("Error creating order:", err)
		return
	}

	
	fmt.Println("             ORDER SUMMARY              ")
	fmt.Println("Order ID:   ", ord.id)
	fmt.Println("Customer:   ", ord.customer.name, "(", ord.customer.phone, ")")
	fmt.Println("----------------------------------------")
	fmt.Println("Items:")

	for _, it := range ord.items {
		fmt.Printf(
			"- %s: %d pieces (%s) = $%.2f\n",
			it.name,
			it.quantity,
			it.serviceType,
			it.totalPrice,
			
		)
	}

	fmt.Println("Orderstatus",ord.statustype)

	fmt.Println("----------------------------------------")
	fmt.Println("Total Items:", ord.totalItems)
	fmt.Printf("TOTAL BILL:  $%.2f\n", ord.totalPrice)
	fmt.Println("========================================")

	allOrders := svc.getAllOrders()
	fmt.Println("\nTotal orders saved in system:", len(allOrders))
	
    svc.UpdateStatus("ORD-1001",washing)

     svc.PrintOrderStatus(washing)
    
    

}
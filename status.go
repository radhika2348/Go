package main

import (
	"fmt"
	"strings"
)

type status_type string
const (
	received status_type="Recieved"
	washing status_type="Washing"
	ready status_type="Ready"
	collected status_type="Collected"

)

type service_type string

const (
	washAndFold service_type = "Wash & Fold"
	washAndIron service_type = "Wash & Iron"
	dryCleaning service_type = "Dry Cleaning"
)

var validstatus=map[status_type]bool{
	received:true,
	washing:true,
	ready:true,
	collected:true,
}

var serviceRates = map[service_type]float64{
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
	service service_type
	unitPrice   float64
	totalPrice  float64
	//status status_type
}

type order struct {
	id         string
	customer   customer
	items      []item
	totalItems int
	totalPrice float64
	status status_type
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

func validatestatus(status status_type) error {
	if _,true:=validstatus[status];!true{
         return fmt.Errorf("not supported")
	}
	return nil
}

func validateServiceType(s service_type) error {
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

		if err := validateServiceType(it.service); err != nil {
			return err
		}

		// if err:=validatestatus(it.status);err!=nil{
		// 	return err
		// }	

	}
	return nil
}

func isValidTransaction(current status_type, next status_type) bool {
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
var j status_type
func (s *orderService) UpdateStatus(orderID string, nextStatus status_type) {
	
	for i := range s.orders {
		if s.orders[i].id == orderID {
			j=s.orders[i].status
			if !isValidTransaction(s.orders[i].status, nextStatus) {
				
				fmt.Println("invalid transition: cannot change status from ", s.orders[i].status, nextStatus)
		       
		   }else{
			
			s.orders[i].status = nextStatus
			fmt.Println(j,"---->", nextStatus)
			fmt.Printf("Status Updated successfully")
		}
		
			
			
		}
	}
     
	
}


func validateOrderRequest(c customer, items []item,status status_type) error {
	if err := validateCustomer(c); err != nil {
		return err
	}

	if err := validateItems(items); err != nil {
		return err
	}
    if err :=validatestatus(status);err !=nil{
		return err
	}
	return nil
}
func (s *orderService) PrintOrderStatus(status status_type) {
    fmt.Println("\n============ ORDER STATUS ============")
    found := false

    for _, ord := range s.orders {
        if ord.status == status {
            fmt.Printf("Order ID: %s is %s!\n", ord.id, ord.status)
            found = true
        }
    }

    if !found {
        fmt.Println("No orders found with status:", status)
    }

    fmt.Println("=======================================")
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

func (s *orderService) createOrder(c customer, items []item,status status_type) (order, error) {
	if err := validateOrderRequest(c, items,status); err != nil {
		return order{}, err
	}

	var totalBill float64
	var totalCount int
	var finalItems []item


	for _, it := range items {
		rate := serviceRates[it.service]
		itemCost := float64(it.quantity) * rate

		it.unitPrice = rate
		it.totalPrice = itemCost

		finalItems = append(finalItems, it)

		totalBill += itemCost
		totalCount += it.quantity
	}

	orderId := fmt.Sprintf("ORD-%d", s.nextId)
	s.nextId++

	//status := 0
	newOrder := order{
		id:         orderId,
		customer:   c,
		items:      finalItems,
		totalItems: totalCount,
		totalPrice: totalBill,
        status:     status,

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
		{name: "Shirts", quantity: 3, service: washAndFold},
		{name: "Coat", quantity: 1, service: dryCleaning},
	}

	status:=collected

	ord, err := svc.createOrder(c, items,status)
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
			it.service,
			it.totalPrice,
			
		)
	}

	fmt.Println("Orderstatus",ord.status)

	fmt.Println("----------------------------------------")
	fmt.Println("Total Items:", ord.totalItems)
	fmt.Printf("TOTAL BILL:  $%.2f\n", ord.totalPrice)
	fmt.Println("========================================")

	allOrders := svc.getAllOrders()
	fmt.Println("\nTotal orders saved in system:", len(allOrders))
	
   
      svc.UpdateStatus("ORD-1001",ready)
     svc.PrintOrderStatus(ready)
    
    

}
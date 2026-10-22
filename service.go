
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Order status
type statustype string

const (
	RECEIVED  statustype = "Received"
	WASHING   statustype = "Washing"
	READY     statustype = "Ready"
	COLLECTED statustype = "Collected"
)

func validateStatus(status statustype) error {
	switch status {
	case RECEIVED, WASHING, READY, COLLECTED:
		return nil
	default:
		return fmt.Errorf("invalid status: %s", status)
	}
}

// Laundry service
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

// Customer
type customer struct {
	name  string
	phone string
}

func (c customer) validate() error {
	if strings.TrimSpace(c.name) == "" {
		return fmt.Errorf("customer name is required")
	}
	if strings.TrimSpace(c.phone) == "" {
		return fmt.Errorf("customer phone is required")
	}
	return nil
}

// Laundry item
type item struct {
	name        string
	quantity    int
	serviceType serviceType
	unitPrice   float64
	totalPrice  float64
}

func (it item) calculatePrice() (item, error) {
	if strings.TrimSpace(it.name) == "" {
		return item{}, fmt.Errorf("item name is required")
	}
	if it.quantity < 1 {
		return item{}, fmt.Errorf("quantity must be at least 1")
	}

	rate, ok := serviceRates[it.serviceType]
	if !ok {
		return item{}, fmt.Errorf(
			"unsupported service: %s", it.serviceType,
		)
	}

	it.unitPrice = rate
	it.totalPrice = float64(it.quantity) * rate
	return it, nil
}

// Order
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
		if _, err := it.calculatePrice(); err != nil {
			return fmt.Errorf("item #%d: %w", i+1, err)
		}
	}

	return nil
}

func (o order) canTransitionTo(next statustype) bool {
	switch o.status {
	case RECEIVED:
		return next == WASHING
	case WASHING:
		return next == READY
	case READY:
		return next == COLLECTED
	default:
		return false
	}
}

func (o *order) UpdateStatus(next statustype) error {
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

// Order service
type orderService struct {
	repo OrderRepository
}

func NewLaundryService(repo OrderRepository) *orderService {
	return &orderService{repo: repo}
}

// Create order: business logic
func (s *orderService) createOrder(
	c customer,
	items []item,
) (order, error) {
	if err := c.validate(); err != nil {
		return order{}, err
	}
	if len(items) == 0 {
		return order{}, fmt.Errorf(
			"order must contain at least one item",
		)
	}

	if err := (order{
		customer: c,
		items:    items,
	}).validate(); err != nil {
		return order{}, err
	}

	orders, err := s.repo.FindAll()
	if err != nil {
		return order{}, err
	}

	ord := order{
		id:       1001,
		customer: c,
		status:   RECEIVED,
	}

	for _, saved := range orders {
		if saved.id >= ord.id {
			ord.id = saved.id + 1
		}
	}

	for _, it := range items {
		priced, err := it.calculatePrice()
		if err != nil {
			return order{}, err
		}

		ord.items = append(ord.items, priced)
		ord.totalItems += priced.quantity
		ord.totalPrice += priced.totalPrice
	}

	if err := s.repo.Save(ord); err != nil {
		return order{}, err
	}

	return ord, nil
}

// Find order by ID
func (s *orderService) FindByID(id int) (order, error) {
	return s.repo.FindByID(id)
}

// Get all orders
func (s *orderService) GetAllOrders() ([]order, error) {
	return s.repo.FindAll()
}

// Update existing order status
func (s *orderService) updateOrderStatus(
	id int,
	next statustype,
) (order, error) {
	if err := s.repo.UpdateStatus(id, next); err != nil {
		return order{}, err
	}

	return s.repo.FindByID(id)
}

// Read a complete line, including spaces
func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	return strings.TrimSpace(value)
}

// Read an integer safely
func readInt(reader *bufio.Reader, prompt string) (int, error) {
	value := readLine(reader, prompt)
	return strconv.Atoi(value)
}

// Display order details
func printOrder(ord order) {
	fmt.Println("\n--- Order Summary ---")
	fmt.Println("Order ID:", ord.id)
	fmt.Printf("Customer: %s (%s)\n",
		ord.customer.name, ord.customer.phone)

	for _, it := range ord.items {
		fmt.Printf(
			"%s: %d pieces (%s) = $%.2f\n",
			it.name, it.quantity,
			it.serviceType, it.totalPrice,
		)
	}

	fmt.Println("Status:", ord.status)
	fmt.Println("Total Items:", ord.totalItems)
	fmt.Printf("Total Bill: $%.2f\n", ord.totalPrice)
}

// Create order: terminal input
func createOrderInput(
	reader *bufio.Reader,
	svc *orderService,
) {
	name := readLine(reader, "Enter customer name: ")
	phone := readLine(reader, "Enter customer phone: ")

	count, err := readInt(
		reader, "Enter number of item types: ",
	)
	if err != nil || count < 1 {
		fmt.Println("Invalid number of item types")
		return
	}

	items := make([]item, 0, count)

	for i := 0; i < count; i++ {
		fmt.Printf("\nItem %d\n", i+1)

		it := item{
			name: readLine(reader, "Enter item name: "),
		}

		it.quantity, err = readInt(reader, "Enter quantity: ")
		if err != nil || it.quantity < 1 {
			fmt.Println("Quantity must be at least 1")
			return
		}

		fmt.Println("1. Wash & Fold ($2.50 per piece)")
		fmt.Println("2. Wash & Iron ($4.00 per piece)")
		fmt.Println("3. Dry Cleaning ($8.00 per piece)")

		choice, err := readInt(reader, "Choose service: ")
		if err != nil {
			fmt.Println("Invalid service choice")
			return
		}

		switch choice {
		case 1:
			it.serviceType = washAndFold
		case 2:
			it.serviceType = washAndIron
		case 3:
			it.serviceType = dryCleaning
		default:
			fmt.Println("Invalid service choice")
			return
		}

		items = append(items, it)
	}

	c := customer{
		name:  name,
		phone: phone,
	}

	ord, err := svc.createOrder(c, items)
	if err != nil {
		fmt.Println("Error creating order:", err)
		return
	}

	fmt.Println("\nOrder created successfully!")
	printOrder(ord)
}

// Update order status: terminal input
func updateOrderInput(
	reader *bufio.Reader,
	svc *orderService,
) {
	id, err := readInt(reader, "Enter Order ID: ")
	if err != nil {
		fmt.Println("Invalid order ID")
		return
	}

	ord, err := svc.FindByID(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Current status:", ord.status)
	fmt.Println("1. Washing")
	fmt.Println("2. Ready")
	fmt.Println("3. Collected")

	choice, err := readInt(reader, "Choose new status: ")
	if err != nil {
		fmt.Println("Invalid status choice")
		return
	}

	var next statustype

	switch choice {
	case 1:
		next = WASHING
	case 2:
		next = READY
	case 3:
		next = COLLECTED
	default:
		fmt.Println("Invalid status choice")
		return
	}




	updated, err := svc.updateOrderStatus(id, next)
	if err != nil {
		fmt.Println("Error updating status:", err)
		return
	}

	fmt.Println("\nOrder status updated successfully!")
	printOrder(updated)
}

// Display orders by status
func printOrdersByStatus(svc *orderService, status statustype) {
	orders, err := svc.GetAllOrders()
	if err != nil {
		fmt.Println("Error:", err)
	  return
	}

	fmt.Printf("\n--- %s Orders ---\n", status)
	found := false

	for _, ord := range orders {
		if ord.status == status {
			fmt.Printf(
				"Order ID: %d | Customer: %s\n",
				ord.id, ord.customer.name,
			)
			found = true
		}
	}

	if !found {
		fmt.Println("No orders found")
	}
}

// Display daily status counts
func printDailyOrderStatus(svc *orderService) {
	orders, err := svc.GetAllOrders()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	counts := map[statustype]int{}

	for _, ord := range orders {
		counts[ord.status]++
	}

	fmt.Println("\n--- Order Status Counts ---")
	fmt.Println("Received:", counts[RECEIVED])
	fmt.Println("Washing:", counts[WASHING])
	fmt.Println("Ready:", counts[READY])
	fmt.Println("Collected:", counts[COLLECTED])
}

// Main menu
func main() {
	repo := NewFileOrderRepository("service.json")
	svc := NewLaundryService(repo)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n===== Laundry Management =====")
		fmt.Println("1. Create Order")
		fmt.Println("2. Update Order Status")
		fmt.Println("3. View All Orders")
		fmt.Println("4. Find Order by ID")
		fmt.Println("5. View Orders by Status")
		fmt.Println("6. View Order Status Counts")
		fmt.Println("7. Exit")

		choice, err := readInt(reader, "Enter your choice: ")
		if err != nil {
			fmt.Println("Invalid choice. Enter a number.")
			continue
		}

		switch choice {
		case 1:
			createOrderInput(reader, svc)

		case 2:
			updateOrderInput(reader, svc)

		case 3:
			orders, err := svc.GetAllOrders()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if len(orders) == 0 {
				fmt.Println("No orders found")
				continue
			}

			for _, ord := range orders {
				printOrder(ord)
			}

		case 4:
			id, err := readInt(reader, "Enter Order ID: ")
			if err != nil {
				fmt.Println("Invalid order ID")
				continue
			}

			ord, err := svc.FindByID(id)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			printOrder(ord)

		case 5:
			fmt.Println("1. Received")
			fmt.Println("2. Washing")
			fmt.Println("3. Ready")
			fmt.Println("4. Collected")

			statusChoice, err := readInt(reader, "Choose status: ")
			if err != nil {
				fmt.Println("Invalid status choice")
				continue
			}

			statuses := map[int]statustype{
				1: RECEIVED,
				2: WASHING,
				3: READY,
				4: COLLECTED,
			}

			status, ok := statuses[statusChoice]
			if !ok {
				fmt.Println("Invalid status choice")
				continue
			}

			printOrdersByStatus(svc, status)

		case 6:
			printDailyOrderStatus(svc)

		case 7:
			fmt.Println("Exiting Laundry Management.")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

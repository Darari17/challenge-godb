package menu

import (
	"bufio"
	"challenge-godb/utils"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	OrderID        	int
	CustomerID     	int
	OrderDate      	time.Time
	CompletionDate 	sql.NullTime
	ReceivedBy     	string
	CreatedAt      	time.Time
	UpdatedAt      	time.Time
}

type OrderDetail struct {
	OrderDetailID 	int
	OrderID       	int
	ServiceID     	int
	Qty           	int
}

//* 1. CREATE ORDER
func AddOrder(db *sql.DB){
	const (
		sqlCheckOrder = "SELECT EXISTS (SELECT 1 FROM \"order\" WHERE order_id = $1);"
		sqlCheckCustomer = "SELECT EXISTS (SELECT 1 FROM customer WHERE customer_id = $1);"
		sqlCheckService = "SELECT EXISTS (SELECT 1 FROM service WHERE service_id = $1);"
		sqlAddOrder = "INSERT INTO \"order\" (order_id, customer_id, order_date, received_by) VALUES ($1, $2, $3, $4);"
		sqlAddOrderDetail = "INSERT INTO order_detail (order_id, service_id, qty) VALUES ($1, $2, $3);"
	)

	o := Order{}
	od := OrderDetail{}
	var odList []OrderDetail
	var isIdExists bool
	scanner := bufio.NewScanner(os.Stdin)

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("%s Error starting transaction: %v\n", utils.ERROR, err)
		return
	}

	for {
		fmt.Print("Enter order ID: ")
		scanner.Scan()
		o.OrderID, err = strconv.Atoi(scanner.Text())
		if err != nil || o.OrderID <= 0 {
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}
		err = db.QueryRow(sqlCheckOrder, o.OrderID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking order ID existence: %v\n", utils.ERROR, err)
			return
		}
		if isIdExists {
			fmt.Printf("%s Order ID already exists. Please enter a different ID.\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter customer ID: ")
		scanner.Scan()
		o.CustomerID, err = strconv.Atoi(scanner.Text())
		if err != nil || o.CustomerID <= 0 {
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}
		err = db.QueryRow(sqlCheckCustomer, o.CustomerID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking customer ID existence: %v\n", utils.ERROR, err)
			return
		}
		if !isIdExists {
			fmt.Printf("%s Customer ID not found. Please enter a different ID.\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter receiver: ")
		scanner.Scan()
		o.ReceivedBy = scanner.Text()
		if o.ReceivedBy == "" {
			fmt.Printf("%s Receiver cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		od.OrderID = o.OrderID

		fmt.Print("Enter service ID: ")
		scanner.Scan()
		od.ServiceID, err = strconv.Atoi(scanner.Text())
		if err != nil || od.ServiceID <= 0 {
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}
		err = db.QueryRow(sqlCheckService, od.ServiceID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking service ID existence: %v\n", utils.ERROR, err)
			return
		}
		if !isIdExists {
			fmt.Printf("%s Service ID not found. Please enter a different ID.\n", utils.WARNING)
			continue
		}
		fmt.Print("Enter quantity: ")
		scanner.Scan()
		od.Qty, err = strconv.Atoi(scanner.Text())
		if err != nil || od.Qty <= 0{
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}
		odList = append(odList, od)
		fmt.Print("Do you want to add another service? (y/n): ")
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			break
		}
	}

	fmt.Println("")
	fmt.Println(utils.RESULT)
	fmt.Printf("Table 'order':\nOrder ID: %d | Customer ID: %d | Receiver By: %s\n", o.OrderID, o.CustomerID, o.ReceivedBy)
	fmt.Printf("Table 'order_detail':\n")
	for _, od := range odList {
    fmt.Printf("Order ID: %d | Service ID: %d | Quantity: %d\n", od.OrderID, od.ServiceID, od.Qty)
	}


	for {
		fmt.Print("Are you sure you want add Customer (y/n)? ")
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = tx.Exec(sqlAddOrder, o.OrderID, o.CustomerID, time.Now(), o.ReceivedBy)
			if err != nil {
				fmt.Printf("%s Transaction Rollback! %v\n", utils.ERROR, err)
				tx.Rollback()
				return
			}
			fmt.Printf("%s Add order successfully!\n", utils.SUCCESS)

			for _, od := range odList {
				_, err = tx.Exec(sqlAddOrderDetail, od.OrderID, od.ServiceID, od.Qty)
				if err != nil {
					fmt.Printf("%s Transaction Rollback! %v\n", utils.ERROR, err)
					tx.Rollback()
					return
				}
			}

			fmt.Printf("%s Add order detail successfully!\n", utils.SUCCESS)
			break
		}  else if choice == "n"{
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'.\n", utils.WARNING)
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("%s Error committing transaction: %v\n", utils.ERROR, err)
		return
	}
	fmt.Printf("%s Transaction committed!\n", utils.SUCCESS)
	fmt.Println("")
}

//* 2. COMPLETE ORDER 
func CompleteOrder(db *sql.DB){
	const (
		sqlCheckOrder = "SELECT EXISTS (SELECT 1 FROM \"order\" WHERE order_id = $1);"
		sqlCompleteOrder = "UPDATE \"order\" SET completion_date = $1, updated_at = CURRENT_TIMESTAMP WHERE order_id = $2;"
	)

	o := Order{}
	var err error
	var isIdExists bool
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter order ID: ")
		scanner.Scan()
		o.OrderID, err = strconv.Atoi(scanner.Text())
		if err != nil || o.OrderID <= 0 {
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckOrder, o.OrderID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking order ID existence: %v\n", utils.ERROR, err)
			return
		}
		if !isIdExists {
			fmt.Printf("%s Order with ID: %d not found!.\n", utils.WARNING, o.OrderID)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter complete order date (yyyy-mm-dd): ")
		scanner.Scan()
		CompletionDate, err := time.Parse("2006-01-02", scanner.Text())
		if err != nil {
			fmt.Printf("%s Invalid date format: %v. Please use yyyy-mm-dd.\n", utils.WARNING, err)
			continue
		}
		o.CompletionDate = sql.NullTime{Time: CompletionDate, Valid: true}
		break
	}
	
	fmt.Println("")
	for {
		fmt.Print("Are you sure you want complete this order (y/n)? ")
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = db.Exec(sqlCompleteOrder, o.CompletionDate, o.OrderID)
			if err != nil {
				fmt.Printf("%s Error updating order: %v\n", utils.ERROR, err)
				return
			}
			fmt.Printf("%s Complete order successfully!", utils.SUCCESS)
			break
		} else if choice == "n"{
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'.\n", utils.WARNING)
		}
	}
	fmt.Println("")
}

//* 3. VIEW OF LIST ORDER
func ViewOfListOrder(db *sql.DB){
	const sqlGetAllOrder = "SELECT order_id, customer_id, order_date, completion_date, received_by, created_at, updated_at FROM \"order\";"
	rows, err := db.Query(sqlGetAllOrder)
	if err != nil {
		fmt.Printf("%s Failed to view of list order: %v\n", utils.ERROR, err)
		return
	}
	defer rows.Close()

	fmt.Println("")
	fmt.Println(utils.RESULT)
	for rows.Next(){
		o := Order{}
		if err = rows.Scan(&o.OrderID, &o.CustomerID, &o.OrderDate, &o.CompletionDate, &o.ReceivedBy, &o.CreatedAt, &o.UpdatedAt);
		err != nil {
			fmt.Printf("%s Failed to scan data order: %v\n", utils.ERROR, err)
			return
		}
		fmt.Printf(
			"Order ID: %d | Customer ID: %d | Order Date: %s | Completion Date: %s | Received By: %s | Created At: %s | Updated AT: %s\n",
			o.OrderID, 
			o.CustomerID, 
			o.OrderDate.Format("2006-01-02 15:04:05"),
			func() string {
				if o.CompletionDate.Valid {
					return o.CompletionDate.Time.Format("2006-01-02 15:04:05")
				}
				return "NULL"
			}(),
			o.ReceivedBy,
			o.CreatedAt.Format("2006-01-02 15:04:05"), 
			o.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	fmt.Printf("%s Get all data order succesfully!\n", utils.SUCCESS)
	fmt.Println("")
}

//* 4. VIEW ORDER DETAILS BY ID
func ViewOrderDetailsById(db *sql.DB){
	const (
		sqlCheckOrder = "SELECT EXISTS (SELECT 1 FROM \"order\" WHERE order_id = $1);"
		sqlGetOrderDetailById = "SELECT order_detail_id, order_id, service_id, qty FROM order_detail WHERE order_id = $1;"
	)

	var isIdExists bool
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter order ID: ")
		scanner.Scan()
		orderID, err := strconv.Atoi(scanner.Text())
		if err != nil || orderID <= 0 {
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckOrder, orderID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking order ID existence: %v\n", utils.ERROR, err)
			return
		}
		if !isIdExists {
			fmt.Printf("%s Order with ID: %d not found!.\n", utils.WARNING, orderID)
			continue
		}
		
		rows, err := db.Query(sqlGetOrderDetailById, orderID)
		if err != nil {
			fmt.Printf("%s Error querying order detail data: %s\n", utils.ERROR, err)
			return
		}
		defer rows.Close()

		fmt.Println("")
		fmt.Println(utils.RESULT)
		fmt.Printf("Order Details for Order ID: %d\n", orderID)

		for rows.Next(){
			od := OrderDetail{}
			err := rows.Scan(&od.OrderDetailID, &od.OrderID, &od.ServiceID, &od.Qty)
			if err != nil {
				fmt.Printf("%s Error reading order detail row: %v\n", utils.ERROR, err)
				return
			}

			fmt.Printf("Order Detail ID: %d | Order ID: %d | Service ID: %d | Quantity: %d\n", od.OrderDetailID, od.OrderID, od.ServiceID, od.Qty)
		}

		if err = rows.Err(); err != nil {
			fmt.Printf("%s Error processing rows: %v\n", utils.ERROR, err)
			return
		}
		fmt.Printf("%s Retrieved all order details successfully!\n", utils.SUCCESS)
		fmt.Println("")
		break
	}
}

//* 5. MENU ORDER
func MenuOrder(db *sql.DB){
	totalLength := 50
	title := "MENU ORDER"

	padding := (totalLength - len(title) - 4) / 2
	leftPadding := strings.Repeat("=", padding)
	rightPadding := strings.Repeat("=", totalLength-len(title)-len(leftPadding)-4)

	fmt.Println(strings.Repeat("=", 48))
	fmt.Printf("%s %s %s\n", leftPadding, title, rightPadding)

	for {
		fmt.Println("1. Create Order")
		fmt.Println("2. Complete Order")
		fmt.Println("3. View of List Order")
		fmt.Println("4. View Order Details By ID")
		fmt.Println("5. Back to Main Menu")
		fmt.Print("Select Menu: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice, _ := strconv.Atoi(scanner.Text())

		switch choice {
			case 1: AddOrder(db)
			case 2: CompleteOrder(db)
			case 3: ViewOfListOrder(db)
			case 4: ViewOrderDetailsById(db)
			case 5: return
			default: fmt.Println("Input invalid, Try again!")
		}
		fmt.Println(strings.Repeat("=", 48))
		fmt.Printf("%s %s %s\n", leftPadding, title, rightPadding)
	}
}
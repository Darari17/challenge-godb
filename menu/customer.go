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

//* 1. STRUCT CUSTOMER
type Customer struct {
	CustomerID int
	Name        string
	Phone       string
	Address     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

//* 2. ADD CUSTOMER
func AddCustomer(db *sql.DB){
	const (
		sqlCheckCustomer = "SELECT EXISTS (SELECT 1 FROM customer WHERE customer_id = $1);"
		sqlAddCustomer = "INSERT INTO customer (customer_id, name, phone, address) VALUES ($1, $2, $3, $4);"
	)

	c := Customer{}
	var isIdExists bool
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter customer ID: ")
		scanner.Scan()
		c.CustomerID, err = strconv.Atoi(scanner.Text())
		if err != nil || c.CustomerID <= 0{
			fmt.Printf("%s Please enter a valid positive number!\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckCustomer, c.CustomerID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking customer ID existence: %v\n", utils.ERROR, err)
			return
		}

		if isIdExists {
			fmt.Printf("%s Customer ID already exists. Please enter a different ID!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter customer name: ")
		scanner.Scan()
		c.Name = scanner.Text()
		if c.Name == "" {
			fmt.Printf("%s Name cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter phone: ")
		scanner.Scan()
		c.Phone = scanner.Text()
		if c.Phone == "" {
			fmt.Printf("%s Phone cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	fmt.Print("Enter Address: ")
	scanner.Scan()
	c.Address = scanner.Text()

	fmt.Println(utils.RESULT)
	fmt.Printf("Customer ID: %d | Name: %s | Phone: %s | Address: %s\n", c.CustomerID, c.Name, c.Phone, c.Address)
	for {
		fmt.Print("Are you sure you want add Customer (y/n)? ")
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = db.Exec(sqlAddCustomer, c.CustomerID, c.Name, c.Phone, c.Address)
			if err != nil {
				fmt.Printf("%s Failed to add customer: %v\n", utils.ERROR, err)
				return
			}
			fmt.Printf("%s Add customer successfully!\n", utils.SUCCESS)
			break
		} else if choice == "n"{
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'!\n", utils.WARNING)
		}
	}
}

//* 3. VIEW OF LIST CUSTOMER
func ViewOfListCustomer(db *sql.DB){
	const sqlGetAllCustomer = "SELECT customer_id, name, phone, address, created_at, updated_at FROM customer;"
	rows, err := db.Query(sqlGetAllCustomer)
	if err != nil {
		fmt.Printf("%s Failed to view of list customer: %v\n", utils.ERROR, err)
		return
	}
	defer rows.Close()

	fmt.Println(utils.RESULT)
	for rows.Next(){
		c := Customer{}
		if err = rows.Scan(&c.CustomerID, &c.Name, &c.Phone, &c.Address, &c.CreatedAt, &c.UpdatedAt);
		err != nil {
			fmt.Printf("%s Failed to scan data Customer: %v\n", utils.ERROR, err)
			return
		}
		fmt.Printf("Customer ID: %d | Name: %s | Phone: %s | Address: %s | Created At: %s | Updated At: %s\n", c.CustomerID, c.Name, c.Phone, c.Address, c.CreatedAt.Format("2006-01-02"), c.UpdatedAt.Format("2006-01-02"))
	}
	fmt.Printf("%s Get all data Customer succesfully!\n", utils.SUCCESS)
}

//* 4. VIEW DETAILS CUSTOMER BY ID
func ViewDetailsCustomerById(db *sql.DB){
	const sqlGetCustomerById = "SELECT customer_id, name, phone, address, created_at, updated_at FROM customer WHERE customer_id = $1;"

	c := Customer{}
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter Customer ID: ")
		scanner.Scan()
		c.CustomerID, err = strconv.Atoi(scanner.Text())
		if err != nil || c.CustomerID <= 0 {
			fmt.Printf("%s Please enter a valid positive number!\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlGetCustomerById, c.CustomerID).Scan(&c.CustomerID, &c.Name, &c.Phone, &c.Address, &c.CreatedAt, &c.UpdatedAt)
		if err == sql.ErrNoRows {
			fmt.Printf("%s Customer not found!\n", utils.ERROR)
			continue
		}
		if err != nil {
			fmt.Printf("%s Error querying Customer data: %s\n", utils.ERROR, err)
			return
		}
		break
	}

	fmt.Println(utils.RESULT)
	fmt.Printf("Customer ID: %d | Name: %s | Phone: %s | Address: %s | Created At: %s | Updated At: %s\n", c.CustomerID, c.Name, c.Phone, c.Address, c.CreatedAt.Format("2006-01-02"), c.UpdatedAt.Format("2006-01-02"))

	fmt.Printf("%s Get data Customer by ID succesfully!\n", utils.SUCCESS)
}

//* 5. UPDATE CUSTOMER
func UpdateCustomer(db *sql.DB){
	const (
		sqlCheckCustomer = "SELECT EXISTS (SELECT 1 FROM customer WHERE customer_id = $1);"
		sqlUpdateCustomer = "UPDATE customer SET name = $2, phone = $3, address = $4, updated_at = CURRENT_TIMESTAMP WHERE customer_id =$1;"
	)

	c := Customer{}
	var err error
	var isIdExists bool
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter customer ID: ")
		scanner.Scan()
		c.CustomerID, err = strconv.Atoi(scanner.Text())
		if err != nil || c.CustomerID <= 0{
			fmt.Printf("%s Please enter a valid positive number!\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckCustomer, c.CustomerID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking customer ID existence: %v\n", utils.ERROR, err)
			return
		}
		
		if !isIdExists {
			fmt.Printf("%s Customer with ID: %d not found!\n", utils.WARNING, c.CustomerID)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Customer Name: ")
		scanner.Scan()
		c.Name = scanner.Text()
		if c.Name == "" {
			fmt.Printf("%s Name cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Phone: ")
		scanner.Scan()
		c.Phone = scanner.Text()
		if c.Phone == "" {
			fmt.Printf("%s Phone cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	fmt.Print("Enter Address: ")
	scanner.Scan()
	c.Address = scanner.Text()

	fmt.Println(utils.RESULT)
	fmt.Printf("Customer ID: %d | Name: %s | Phone: %s | Address: %s\n", c.CustomerID, c.Name, c.Phone, c.Address)
	for {
		fmt.Print("Are you sure you want update Customer (y/n)? ")
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = db.Exec(sqlUpdateCustomer, c.CustomerID, c.Name, c.Phone, c.Address)
			if err != nil {
				fmt.Printf("%s Failed to Update Customer: %v\n", utils.ERROR, err)
				return
			}
			fmt.Printf("%s Update Customer successfully!\n", utils.SUCCESS)
			break
		} else if choice == "n"{
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'.\n", utils.WARNING)
		}
	}
}

//* 6. DELETE CUSTOMER
func DeleteCustomer(db *sql.DB){
	const (
		sqlCheckCustomer = "SELECT EXISTS (SELECT 1 FROM customer WHERE customer_id = $1);"
		sqlDeleteCustomer = "DELETE FROM customer WHERE customer_id = $1;"
	)

	c := Customer{}
	var err error
	var isIdExists bool
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter customer ID: ")
		scanner.Scan()
		c.CustomerID, err = strconv.Atoi(scanner.Text())
		if err != nil || c.CustomerID <= 0{
			fmt.Printf("%s Please enter a valid positive number!\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckCustomer, c.CustomerID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking customer ID existence: %v\n", utils.ERROR, err)
			return
		}
		
		if !isIdExists {
			fmt.Printf("%s Customer with ID: %d not found!\n", utils.WARNING, c.CustomerID)
			continue
		}
		break
	}

	for {
		fmt.Printf("Are you sure want to delete Customer ID: %d (y/n)? ", c.CustomerID)
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = db.Exec(sqlDeleteCustomer, c.CustomerID)
			if err != nil {
				fmt.Printf("%s Failed to Delete Customer: %v\n", utils.ERROR, err)
			}
			fmt.Printf("%s Delete Customer Successfully!\n", utils.SUCCESS)
			break
		} else if choice == "n" {
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'!\n", utils.WARNING)
		}
	}	
}

//* 7. MENU CUSTOMER
func MenuCustomer(db *sql.DB) {
	totalLength := 50
	title := "MENU CUSTOMER"

	padding := (totalLength - len(title) - 4) / 2
	leftPadding := strings.Repeat("=", padding)
	rightPadding := strings.Repeat("=", totalLength-len(title)-len(leftPadding)-4)

	fmt.Println(strings.Repeat("=", 48))
	fmt.Printf("%s %s %s\n", leftPadding, title, rightPadding)

	for {
		fmt.Println("1. Create Customer")
		fmt.Println("2. View of List Customer")
		fmt.Println("3. View Details Customer By ID")
		fmt.Println("4. Update Customer")
		fmt.Println("5. Delete Customer")
		fmt.Println("6. Back to Main Menu")
		fmt.Print("Select Menu: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice, _ := strconv.Atoi(scanner.Text())

		switch choice {
		case 1: AddCustomer(db)
		case 2:	ViewOfListCustomer(db)
		case 3: ViewDetailsCustomerById(db)
		case 4: UpdateCustomer(db)
		case 5: DeleteCustomer(db)
		case 6: return
		default: fmt.Println("Input invalid, Try again!")
	}
	fmt.Println(strings.Repeat("=", 48))
	}
}

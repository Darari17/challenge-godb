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

//* 1. STRUCT SERVICE
type Service struct {
	ServiceID   int
	ServiceName string
	Unit        string
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//* 2. ADD SERVICE
func AddService(db *sql.DB){
	const (
		sqlCheckService = "SELECT EXISTS (SELECT 1 FROM service WHERE service_id = $1);"
		sqlAddService = "INSERT INTO service (service_id, service_name, unit, price) VALUES ($1, $2, $3, $4);"
	)

	s := Service{}
	var isIdExists bool
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter service ID: ")
		scanner.Scan()
		s.ServiceID, err = strconv.Atoi(scanner.Text())
		if err != nil || s.ServiceID <= 0{
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckService, s.ServiceID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking service ID existence: %v\n", utils.ERROR, err)
			return
		}

		if isIdExists {
			fmt.Printf("%s Service ID already exists. Please enter a different ID.\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Service Name: ")
		scanner.Scan()
		s.ServiceName = scanner.Text()
		if s.ServiceName == "" {
			fmt.Printf("%s Name cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Unit: ")
		scanner.Scan()
		s.Unit = scanner.Text()
		if s.Unit == "" {
			fmt.Printf("%s Unit cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Price: ")
		scanner.Scan()
		s.Price, err = strconv.Atoi(scanner.Text())
		if err != nil || s.Price <= 0 {
			fmt.Printf("%s Price must be a valid positive number!\n", utils.WARNING)
			continue
		}
		break
	}

	fmt.Println(utils.RESULT)
	fmt.Printf("Service ID: %d | Service Name: %s | Unit: %s | Price: %d\n", s.ServiceID, s.ServiceName, s.Unit, s.Price)

	for {
		fmt.Print("Are you sure you want to add Service (y/n)? ")
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err := db.Exec(sqlAddService, s.ServiceID, s.ServiceName, s.Unit, s.Price)
			if err != nil {
				fmt.Printf("%s Failed to add Service: %v\n", utils.ERROR, err)
				return
			}
			fmt.Printf("%s Add Service successfully!\n", utils.SUCCESS)
			break
		} else if choice == "n"{
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'.\n", utils.WARNING)
		}
	}
}

//* 3. VIEW OF LIST SERVICE
func ViewOfListService(db *sql.DB){
	const sqlGetAllService = "SELECT service_id, service_name, unit, price, created_at, updated_at FROM service;"
	rows, err := db.Query(sqlGetAllService)
	if err != nil {
		fmt.Printf("%s Failed to get all data Service: %v\n", utils.ERROR, err)
		return
	}
	defer rows.Close()

	fmt.Println(utils.RESULT)
	for rows.Next(){
		s := Service{}
		if err = rows.Scan(&s.ServiceID, &s.ServiceName, &s.Unit, &s.Price, &s.CreatedAt, &s.UpdatedAt);
		err != nil {
			fmt.Printf("%s Failed to scan data Service: %v\n", utils.ERROR, err)
			return
		}
		fmt.Printf("Service ID: %d | Service Name: %s | Unit: %s | Price: %d | Created At: %s | Updated At: %s\n", s.ServiceID, s.ServiceName, s.Unit, s.Price, s.CreatedAt.Format("2006-01-02"), s.UpdatedAt.Format("2006-01-02"))
	}
	fmt.Printf("%s Get all data service succesfully!\n", utils.SUCCESS)
}

//* 4. VIEW DETAILS SERVICE BY ID
func ViewDetailsServiceById(db *sql.DB){
	const sqlGetServiceById = "SELECT service_id, service_name, unit, price, created_at, updated_at FROM service WHERE service_id = $1;"

	s := Service{}
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter service ID: ")
		scanner.Scan()
		s.ServiceID, err = strconv.Atoi(scanner.Text())
		if err != nil || s.ServiceID <= 0 {
			fmt.Printf("%s Please enter a valid positive number!\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlGetServiceById, s.ServiceID).Scan(&s.ServiceID, &s.ServiceName, &s.Unit, &s.Price, &s.CreatedAt, &s.UpdatedAt)
		if err == sql.ErrNoRows {
			fmt.Printf("%s Service not found!\n", utils.ERROR)
			continue
		}
		if err != nil {
			fmt.Printf("%s Error querying Service data: %s\n", utils.ERROR, err)
			return
		}
		break
	}

	fmt.Println(utils.RESULT)
	fmt.Printf("Service ID: %d | Service Name: %s | Unit: %s | Price: %d | Created At: %s | Updated At: %s\n", s.ServiceID, s.ServiceName, s.Unit, s.Price, s.CreatedAt.Format("2006-01-02"), s.UpdatedAt.Format("2006-01-02"))
	
	fmt.Printf("%s Get data Service by ID successfully!\n", utils.SUCCESS)
}

//* 5. UPDATE SERVICE
func UpdateService(db *sql.DB){
	const (
		sqlCheckService = "SELECT EXISTS (SELECT 1 FROM service WHERE service_id = $1);"
		sqlUpdateService = "UPDATE service SET service_name = $2, unit = $3, price = $4, updated_at = CURRENT_TIMESTAMP WHERE service_id = $1;"
	)

	s := Service{}
	var isIdExists bool
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter service ID: ")
		scanner.Scan()
		s.ServiceID, err = strconv.Atoi(scanner.Text())
		if err != nil || s.ServiceID <= 0{
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckService, s.ServiceID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking service ID existence: %v\n", utils.ERROR, err)
			return
		}

		if !isIdExists {
			fmt.Printf("%s Service with ID: %d not found!.\n", utils.WARNING, s.ServiceID)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Service Name: ")
		scanner.Scan()
		s.ServiceName = scanner.Text()
		if s.ServiceName == "" {
			fmt.Printf("%s Name cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Unit: ")
		scanner.Scan()
		s.Unit = scanner.Text()
		if s.Unit == "" {
			fmt.Printf("%s Unit cannot be empty!\n", utils.WARNING)
			continue
		}
		break
	}

	for {
		fmt.Print("Enter Price: ")
		scanner.Scan()
		s.Price, err = strconv.Atoi(scanner.Text())
		if err != nil || s.Price <= 0 {
			fmt.Printf("%s Price must be a valid positive number!\n", utils.WARNING)
			continue
		}
		break
	}

	fmt.Println(utils.RESULT)
	fmt.Printf("Service ID: %d | Service Name: %s | Unit: %s | Price: %d\n", s.ServiceID, s.ServiceName, s.Unit, s.Price)

	for {
		fmt.Print("Are you sure you want Update Service (y/n)? ")
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = db.Exec(sqlUpdateService, s.ServiceID, s.ServiceName, s.Unit, s.Price)
			if err != nil {
				fmt.Printf("%s Failed to Update Service: %v\n", utils.ERROR, err)
				return
			}
			fmt.Printf("%s Update Service successfully!\n", utils.SUCCESS)
			break
		} else if choice == "n" {
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'.\n", utils.WARNING)
		}
	}
}

//* 6. DELETE SERVICE
func DeleteService(db *sql.DB){
	const (
		sqlCheckService = "SELECT EXISTS (SELECT 1 FROM service WHERE service_id = $1);"
		sqlDeleteService = "DELETE FROM service WHERE service_id = $1;"
	)

	s := Service{}
	var isIdExists bool
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter service ID: ")
		scanner.Scan()
		s.ServiceID, err = strconv.Atoi(scanner.Text())
		if err != nil || s.ServiceID <= 0{
			fmt.Printf("%s Please enter a valid positive number.\n", utils.WARNING)
			continue
		}

		err = db.QueryRow(sqlCheckService, s.ServiceID).Scan(&isIdExists)
		if err != nil {
			fmt.Printf("%s Error checking service ID existence: %v\n", utils.ERROR, err)
			return
		}

		if !isIdExists {
			fmt.Printf("%s Service with ID: %d not found!.\n", utils.WARNING, s.ServiceID)
			continue
		}
		break
	}

	for {
		fmt.Printf("Are you sure want to delete Service ID: %d (y/n)? ", s.ServiceID)
		scanner.Scan()
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if choice == "y" {
			_, err = db.Exec(sqlDeleteService, s.ServiceID)
			if err != nil {
				fmt.Printf("%s Failed to Delete Service: %v", utils.ERROR, err)
				return
			}
			fmt.Printf("%s Delete Service Successfully!", utils.SUCCESS)
			break
		} else if choice == "n" {
			fmt.Println("Canceled!")
			return
		} else {
			fmt.Printf("%s Invalid choice! Please enter 'y' or 'n'.\n", utils.WARNING)
		}
	}
}

//* 7. MENU SERVICE
func MenuService(db *sql.DB){
	totalLength := 50
	title := "MENU SERVICE"

	padding := (totalLength - len(title) - 4) / 2
	leftPadding := strings.Repeat("=", padding)
	rightPadding := strings.Repeat("=", totalLength-len(title)-len(leftPadding)-4)

	fmt.Println(strings.Repeat("=", 48))
	fmt.Printf("%s %s %s\n", leftPadding, title, rightPadding)

	for {
		fmt.Println("1. Create Service")
		fmt.Println("2. View of List Service")
		fmt.Println("3. View Details Service By ID")
		fmt.Println("4. Update Service")
		fmt.Println("5. Delete Service")
		fmt.Println("6. Back to Main Menu")
		fmt.Print("Select Menu: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice, _ := strconv.Atoi(scanner.Text())

		switch choice {
			case 1: AddService(db)
			case 2:	ViewOfListService(db)
			case 3: ViewDetailsServiceById(db)
			case 4: UpdateService(db)
			case 5: DeleteService(db)
			case 6: return
			default: fmt.Println("Input invalid, Try again!")
		}
		fmt.Println(strings.Repeat("=", 48))
	}
}

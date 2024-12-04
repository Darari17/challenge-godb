package main

import (
	"bufio"
	"challenge-godb/menu"
	"challenge-godb/utils"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "12345"
	dbname = "enigma_laundry"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func connect() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%s Error connecting to the database: %v\n", utils.ERROR, err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("%s Error pinging the database: %v\n", utils.ERROR, err)
		return nil
	}
	fmt.Printf("%s Succesfully Connected!\n", utils.SUCCESS)
	return db
}

func main() {
	db := connect()
	if db == nil {
		fmt.Println("Failed to establish database connection. Exiting program.")
		return
	}
	defer db.Close()

	totalLength := 50
	title := "MAIN MENU"

	padding := (totalLength - len(title) - 4) / 2
	leftPadding := strings.Repeat("=", padding)
	rightPadding := strings.Repeat("=", totalLength-len(title)-len(leftPadding)-4)

	fmt.Println(strings.Repeat("=", 48))
	fmt.Printf("%s %s %s\n", leftPadding, title, rightPadding)

	for {
		fmt.Println("1. CUSTOMER")
		fmt.Println("2. SERVICE")
		fmt.Println("3. ORDER")
		fmt.Println("4. EXIT")
		fmt.Print("SELECT MENU: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice, _ := strconv.Atoi(scanner.Text())

		switch choice {
		case 1 : menu.MenuCustomer(db)
		case 2 : menu.MenuService(db)
		case 3 : menu.MenuOrder(db)
		case 4 : return
		default: fmt.Println("Input invalid, Try again!")
		}
		fmt.Println(strings.Repeat("=", 48))
	}
}
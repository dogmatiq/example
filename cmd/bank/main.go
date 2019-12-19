package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/database"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/testkit/engine"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func businessDayFromTime(t time.Time) string {
	return t.Format(messages.BusinessDateFormat)
}

func readString(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.Replace(text, "\n", "", -1)
}

func generateID() string {
	return uuid.New().String()
}

type appState struct {
	db     *sql.DB
	engine *engine.Engine
	reader *bufio.Reader
	time   time.Time
}

type menuItem struct {
	// option is the short menu item code, usually a single character.
	option string

	// description is the description of the menu item.
	description string

	// command is the command to return when this menu item is selected.
	command string
}

func main() {
	db := database.MustNew()
	defer db.Close()

	app, err := example.NewApp(db)
	if err != nil {
		panic(err)
	}

	en, err := engine.New(app)
	if err != nil {
		panic(err)
	}

	as := appState{
		db:     db,
		engine: en,
		reader: bufio.NewReader(os.Stdin),
		time:   time.Date(2001, 10, 20, 11, 22, 33, 44, time.UTC),
	}

	fmt.Println()
	fmt.Println("Dogma Example Bank")
	fmt.Println("==================")

	as.runMainMenu()
}

// showMenu shows a title and a list of options.
func (as appState) showMenu(title, description string, items []menuItem) string {
	fmt.Println()
	fmt.Println("================================================================================")
	fmt.Printf("Time: %s", as.time.Format("2006 Jan 2 3:04pm MST"))
	fmt.Println()
	fmt.Println()
	fmt.Println(title)
	fmt.Println("--------------------")
	fmt.Println(description)
	fmt.Println()
	fmt.Println("Options:")
	for _, item := range items {
		fmt.Printf("  [%s] - %s\n", item.option, item.description)
	}

	var input string
	for {
		fmt.Println()
		fmt.Print("Choice: ")
		input = readString(as.reader)

		for _, item := range items {
			if item.option == input {
				return item.command
			}
		}

		fmt.Printf("Invalid selection: %s\n", input)
	}
}

func (as appState) runMainMenu() {
	for {
		var menuItems []menuItem

		rows, err := as.db.Query("SELECT id, name FROM customer")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var (
			id   string
			name string
		)
		number := 0

		for rows.Next() {
			if err := rows.Scan(
				&id,
				&name,
			); err != nil {
				panic(err)
			}

			number++
			menuItems = append(
				menuItems,
				menuItem{
					option:      strconv.Itoa(number),
					description: fmt.Sprintf("%s (%s)", name, id),
					command:     id,
				},
			)
		}

		menuItems = append(
			menuItems,
			menuItem{
				option:      "n",
				description: "Open account for new customer",
				command:     "new",
			},
			menuItem{
				option:      "x",
				description: "Exit banking sytem",
				command:     "exit",
			},
		)

		command := as.showMenu("Main Menu", "Select a customer to login as, or create a new customer.", menuItems)
		switch command {
		case "exit":
			fmt.Println()
			fmt.Println("Thankyou for banking with Dogma Example Bank.")
			return
		case "new":
			customerID := as.runNewCustomerMenu()
			if customerID != "" {
				as.runCustomerMenu(customerID)
			}
		default:
			as.runCustomerMenu(command)
		}
	}
}

func (as appState) runNewCustomerMenu() string {
	fmt.Println()
	fmt.Println("================================================================================")
	fmt.Println()
	fmt.Println("New Customer Signup")
	fmt.Println("--------------------")
	fmt.Println("Enter new customer details.")

	fmt.Println()
	fmt.Print("Customer Name (empty to cancel): ")
	customerName := readString(as.reader)
	if customerName == "" {
		return ""
	}

	fmt.Println()
	fmt.Print("Account Name (default 'Savings'): ")
	accountName := readString(as.reader)
	if accountName == "" {
		accountName = "Savings"
	}

	customerID := generateID()

	as.engine.Dispatch(
		context.Background(),
		commands.OpenAccountForNewCustomer{
			CustomerID:   customerID,
			CustomerName: customerName,
			AccountID:    generateID(),
			AccountName:  accountName,
		},
	)

	fmt.Println()
	fmt.Println("Customer signed up.")

	return customerID
}

func (as appState) runCustomerMenu(customerID string) {
	rows, err := as.db.Query(
		`SELECT
			name
		FROM customer
		WHERE id = ?`,
		customerID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var customerName string

	rows.Next()
	if err := rows.Scan(
		&customerName,
	); err != nil {
		panic(err)
	}

	menuItems := []menuItem{
		{
			option:      "l",
			description: "List accounts and blances",
			command:     "list",
		},
		{
			option:      "o",
			description: "Open another account",
			command:     "open",
		},
		{
			option:      "d",
			description: "Deposit funds",
			command:     "deposit",
		},
		{
			option:      "w",
			description: "Withdraw funds",
			command:     "withdraw",
		},
		{
			option:      "t",
			description: "Transfer funds",
			command:     "transfer",
		},
		{
			option:      "x",
			description: "Logout customer",
			command:     "logout",
		},
	}

	for {
		switch as.showMenu("Customer Menu", "Welcome "+customerName+", please select an action.", menuItems) {
		case "list":
			as.runListAccounts(customerID)
		case "open":
			as.runNewAccountMenu(customerID)
		case "deposit":
			as.runDepositMenu(customerID)
		case "withdraw":
			fmt.Println()
			fmt.Println("TODO: Withdraw...")
		case "transfer":
			fmt.Println()
			fmt.Println("TODO: Transfer...")
		case "logout":
			fmt.Println()
			fmt.Println("Goodbye")
			return
		}
	}
}

func (as appState) runListAccounts(customerID string) {
	rows, err := as.db.Query(
		`SELECT
			id,
			name,
			balance
		FROM account
		WHERE customer_id = ?`,
		customerID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var (
		id      string
		name    string
		balance int64
	)

	fmt.Println()
	fmt.Println("Here are you accounts:")
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&name,
			&balance,
		); err != nil {
			panic(err)
		}

		dollars := float64(balance) / 100.0
		fmt.Printf("  $%.2f %s (%s)\n", dollars, name, id)
	}

	fmt.Println()
	fmt.Print("Press enter to continue ...")
	readString(as.reader)
}

func (as appState) runSelectAccountMenu(forCustomerID string) string {
	var (
		menuItems []menuItem
		rows      *sql.Rows
		err       error
	)

	// Customer specific or all...
	if forCustomerID == "" {
		rows, err = as.db.Query(
			`SELECT
				a.id,
				a.name,
				a.customer_id,
				c.name AS customer_name,
				a.balance
			FROM account AS a
			INNER JOIN customer AS c
			ON c.id = a.customer_id
			ORDER BY a.customer_id, a.name`,
		)
	} else {
		rows, err = as.db.Query(
			`SELECT
				a.id,
				a.name,
				a.customer_id,
				c.name AS customer_name,
				a.balance
			FROM account AS a
			INNER JOIN customer AS c
			ON c.id = a.customer_id
			WHERE customer_id = ?
			ORDER BY a.customer_id, a.name`,
			forCustomerID,
		)
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var (
		accountID    string
		accountName  string
		customerID   string
		customerName string
		balance      int64
	)
	number := 0

	for rows.Next() {
		if err := rows.Scan(
			&accountID,
			&accountName,
			&customerID,
			&customerName,
			&balance,
		); err != nil {
			panic(err)
		}

		dollars := float64(balance) / 100.0
		customerDisplay := ""
		if forCustomerID == "" {
			customerDisplay = fmt.Sprintf(" [%s]", customerName)
		}

		number++
		menuItems = append(
			menuItems,
			menuItem{
				option:      strconv.Itoa(number),
				description: fmt.Sprintf("%.2f %s (%s)%s", dollars, accountName, accountID, customerDisplay),
				command:     accountID,
			},
		)
	}

	fmt.Println()
	fmt.Println("Options:")
	for _, item := range menuItems {
		fmt.Printf("  [%s] - %s\n", item.option, item.description)
	}

	var input string
	for {
		fmt.Println()
		fmt.Print("Choice: ")
		input = readString(as.reader)

		for _, item := range menuItems {
			if item.option == input {
				return item.command
			}
		}

		fmt.Printf("Invalid selection: %s\n", input)
	}
}

func (as appState) runNewAccountMenu(customerID string) {
	fmt.Println()
	fmt.Println("================================================================================")
	fmt.Println()
	fmt.Println("Open New Account")
	fmt.Println("--------------------")
	fmt.Println("Enter account details.")

	fmt.Println()
	fmt.Print("Account Name (empty to cancel): ")
	accountName := readString(as.reader)
	if accountName == "" {
		return
	}

	as.engine.Dispatch(
		context.Background(),
		commands.OpenAccount{
			CustomerID:  customerID,
			AccountID:   generateID(),
			AccountName: accountName,
		},
	)

	fmt.Println()
	fmt.Println("Account opened.")
}

func (as appState) runDepositMenu(customerID string) {
	fmt.Println()
	fmt.Println("================================================================================")
	fmt.Println()
	fmt.Println("Deposit Funds")
	fmt.Println("--------------------")
	fmt.Println("Enter deposit details.")

	fmt.Println()
	fmt.Print("Select Account: ")
	accountID := as.runSelectAccountMenu(customerID)

	fmt.Println()
	fmt.Print("Deposit amount in dollars (empty to cancel): ")
	amount := readString(as.reader)
	if amount == "" {
		return
	}

	dollarsAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println("Invalid amount.")
		return
	}

	centsAmount := int64(dollarsAmount * 100)
	if centsAmount == 0 {
		return
	}

	as.engine.Dispatch(
		context.Background(),
		commands.Deposit{
			TransactionID: generateID(),
			AccountID:     accountID,
			Amount:        centsAmount,
		},
	)

	fmt.Println()
	fmt.Println("Deposit sent.")
}

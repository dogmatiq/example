package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/database"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/testkit/engine"
	"github.com/gdamore/tcell"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

const refreshInterval = 1000 * time.Millisecond

func currentTimeString(t time.Time) string {
	return fmt.Sprintf(t.Format("Date: 2006-01-02 15:04:05"))
}

func updateTime(as *appState) {
	for {
		time.Sleep(refreshInterval)
		as.guiApp.QueueUpdateDraw(func() {
			as.time = as.time.Add(time.Minute) // accelerated example bank time
			as.clockView.SetText(currentTimeString(as.time))

			as.engine.Tick(
				context.Background(),
				engine.WithCurrentTime(as.time),
			)
		})
	}
}

func businessDayFromTime(t time.Time) string {
	return t.Format(messages.BusinessDateFormat)
}

func generateID() string {
	return uuid.New().String()
}

func generateCustomerNumber() string {
	return fmt.Sprintf("%d%03d", rand.Intn(8)+1, rand.Intn(900))
}

func generateAccountNumber() string {
	return fmt.Sprintf("%d%04d", rand.Intn(8)+1, rand.Intn(9000))
}

type appState struct {
	db     *sql.DB
	engine *engine.Engine
	time   time.Time

	// ui ...
	guiApp          *tview.Application
	currentMainView tview.Primitive
	layoutView      *tview.Grid
	clockView       *tview.TextView
	infoView        *tview.TextView
	logView         *tview.TextView
	mainMenu        *tview.List
	customerMenu    *tview.List
	advanceTimeMenu *tview.List
	customerList    *tview.List
	accountList     *tview.List
}

type customerData struct {
	id      string
	name    string
	display string
}

type accountData struct {
	id      string
	name    string
	balance int64
	display string
}

type customerAccountData struct {
	customerID   string
	customerName string
	accountID    string
	accountName  string
	balance      int64
	display      string
}

// -----------------------

func main() {
	rand.Seed(time.Now().Unix())

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
		time:   time.Date(2001, 10, 20, 11, 22, 33, 44, time.UTC),
	}

	as.createViews()

	go updateTime(&as)

	if err := as.guiApp.SetRoot(as.layoutView, true).Run(); err != nil {
		panic(err)
	}
}

func (as *appState) switchMainView(view tview.Primitive) {
	as.layoutView.RemoveItem(as.currentMainView)

	as.layoutView.AddItem(view, 2, 0, 1, 2, 0, 0, true)
	as.guiApp.SetFocus(view) // the bool above doesn't seem to actually set focus
	as.currentMainView = view
}

func (as *appState) createViews() {
	as.guiApp = tview.NewApplication()

	bankTitle := tview.NewTextView()
	bankTitle.SetText("Dogma Example Bank")
	bankTitle.SetTextColor(tcell.ColorYellow)
	bankTitle.SetBackgroundColor(tcell.ColorBlue)

	as.clockView = tview.NewTextView()
	as.clockView.SetText(currentTimeString(as.time))
	as.clockView.SetTextAlign(tview.AlignRight)
	as.clockView.SetTextColor(tcell.ColorWhite)
	as.clockView.SetBackgroundColor(tcell.ColorBlue)

	as.infoView = tview.NewTextView()
	as.infoView.SetBorder(true)
	as.infoView.SetTitle("[ Info ]")
	as.infoView.SetText("Welcome to Dogma Example bank!")

	as.logView = tview.NewTextView()
	as.logView.SetBorder(true)
	as.logView.SetTitle("[ Dogma Logs ]")
	as.logView.SetText("TODO: logs...")

	// Create static views. Others will be recreated as needed for fresh data
	as.createMainMenu()
	as.createAdvanceTimeMenu()

	as.layoutView = tview.NewGrid()
	as.layoutView.SetRows(1, 3, 0, 10)
	as.layoutView.SetColumns(0, 30)
	as.layoutView.AddItem(bankTitle, 0, 0, 1, 1, 3, 0, false)
	as.layoutView.AddItem(as.clockView, 0, 1, 1, 1, 3, 0, false)
	as.layoutView.AddItem(as.infoView, 1, 0, 1, 2, 10, 0, false)
	as.layoutView.AddItem(as.mainMenu, 2, 0, 1, 2, 0, 0, true)
	as.layoutView.AddItem(as.logView, 3, 0, 1, 2, 10, 0, false)
}

func (as *appState) createMainMenu() {
	as.mainMenu = tview.NewList()
	as.mainMenu.SetBorder(true)
	as.mainMenu.SetTitle("[ Main Menu ]")
	as.mainMenu.SetTitleAlign(tview.AlignLeft)
	as.mainMenu.AddItem("New Customer", "Open account for new customer", 'n', func() {
		as.switchMainView(as.createAddCustomerForm())
	})
	as.mainMenu.AddItem("List Customers", "List all customers", 'l', func() {
		as.createCustomerList()
		as.switchMainView(as.customerList)
	})
	as.mainMenu.AddItem("Advance Time", "Advance the bank clock time", 'a', func() {
		as.infoView.SetText("Use the options below to advance time to help test things like future scheduled transfers.")
		as.switchMainView(as.advanceTimeMenu)
	})
	as.mainMenu.AddItem("Quit", "Quit example bank application", 'q', func() {
		as.guiApp.Stop()
	})
}

func (as *appState) createAdvanceTimeMenu() {
	as.advanceTimeMenu = tview.NewList()
	as.advanceTimeMenu.SetBorder(true)
	as.advanceTimeMenu.SetTitle("[ Advance Time ]")
	as.advanceTimeMenu.SetTitleAlign(tview.AlignLeft)
	as.advanceTimeMenu.AddItem("One Hour", "Advance time by one hour", 'h', func() {
		as.time = as.time.Add(time.Hour)
	})
	as.advanceTimeMenu.AddItem("One Day", "Advance time by one day", 'd', func() {
		as.time = as.time.Add(time.Hour * 24)
	})
	as.advanceTimeMenu.AddItem("Quit to Main Menu", "Quit and return to the main menu", 'q', func() {
		as.infoView.SetText("")
		as.switchMainView(as.mainMenu)
	})
}

func (as *appState) createCustomerList() {
	as.infoView.SetText("Please select a customer.")

	as.customerList = tview.NewList()
	as.customerList.SetBorder(true)
	as.customerList.SetTitle("[ Customer List ]")
	as.customerList.SetTitleAlign(tview.AlignLeft)
	customers := as.fetchCustomers()
	for i, c := range customers {
		// Use numbered shortcuts for the first 10 customers and empty rune `0`
		// for the rest
		var shortcut rune
		if i < 10 {
			shortcut = '0' + rune(i)
		}
		cust := c
		as.customerList.AddItem(c.display, "View details for this customer", shortcut, func() {
			as.createCustomerMenu(cust)
			as.switchMainView(as.customerMenu)
		})
	}
	as.customerList.AddItem("Quit to Main Menu", "Quit and return to the main menu", 'q', func() {
		as.infoView.SetText("")
		as.switchMainView(as.mainMenu)
	})
}

func (as *appState) createAccountList(customer customerData) {
	as.infoView.SetText("Here are the customer's accounts.")

	as.accountList = tview.NewList()
	as.accountList.SetBorder(true)
	as.accountList.SetTitle(fmt.Sprintf("[ Account List: %s ]", customer.name))
	as.accountList.SetTitleAlign(tview.AlignLeft)
	accounts := as.fetchAccountsForCustomer(customer.id)
	for _, a := range accounts {
		as.accountList.AddItem(a.display, "", 0, nil)
	}
	as.accountList.AddItem("Quit to Customer Menu", "Quit and return to the customer menu", 'q', func() {
		as.infoView.SetText("")
		as.switchMainView(as.customerMenu)
	})
}

func (as *appState) createCustomerMenu(customer customerData) {
	as.infoView.SetText("Please select an action.")

	as.customerMenu = tview.NewList()
	as.customerMenu.SetBorder(true)
	as.customerMenu.SetTitle(fmt.Sprintf("[ Customer Menu: %s ]", customer.name))
	as.customerMenu.SetTitleAlign(tview.AlignLeft)
	as.customerMenu.AddItem("List Accounts", "List accounts and balances", 'l', func() {
		as.createAccountList(customer)
		as.switchMainView(as.accountList)
	})
	as.customerMenu.AddItem("Add Account", "Add a new account for this customer", 'a', func() {
		as.switchMainView(as.createAddAccountForm(customer))
	})
	as.customerMenu.AddItem("Deposit", "Deposit funds into an account", 'd', func() {
		as.switchMainView(as.createDepositForm(customer))
	})
	as.customerMenu.AddItem("Withdraw", "Withdraw funds from an account", 'w', func() {
		as.switchMainView(as.createWithdrawForm(customer))
	})
	as.customerMenu.AddItem("Transfer", "Transfer funds into another account", 't', func() {
		as.switchMainView(as.createTransferForm(customer))
	})
	as.customerMenu.AddItem("Quit to Customer List", "Quit and return to the customer list", 'q', func() {
		as.switchMainView(as.customerList)
	})
}

func (as *appState) createAddCustomerForm() *tview.Form {
	as.infoView.SetText("Enter new customer details.")

	addCustomerForm := tview.NewForm()
	addCustomerForm.SetBorder(true)
	addCustomerForm.SetTitle("[ New Customer Signup ]")
	addCustomerForm.SetTitleAlign(tview.AlignLeft)
	addCustomerForm.AddInputField("Customer name", "", 20, nil, nil)
	addCustomerForm.AddInputField("Account name", "Savings", 20, nil, nil)
	addCustomerForm.AddButton("Create", func() {
		customerField := addCustomerForm.GetFormItemByLabel("Customer name").(*tview.InputField)
		accountField := addCustomerForm.GetFormItemByLabel("Account name").(*tview.InputField)

		customerName := customerField.GetText()
		accountName := accountField.GetText()

		if customerName == "" {
			// Form not complete.
			as.infoView.SetText("Customer name must not be blank.")
			return
		}

		if accountName == "" {
			// Form not complete.
			as.infoView.SetText("Account name must not be blank.")
		}

		as.submitNewCustomerSignup(customerName, accountName)

		as.infoView.SetText("New customer details submitted.")
		as.switchMainView(as.mainMenu)
	})
	addCustomerForm.AddButton("Cancel", func() {
		as.infoView.SetText("New customer signup cancelled.")
		as.switchMainView(as.mainMenu)
	})
	addCustomerForm.SetCancelFunc(func() {
		as.infoView.SetText("New customer signup cancelled.")
		as.switchMainView(as.mainMenu)
	})

	return addCustomerForm
}

func (as *appState) createAddAccountForm(customer customerData) *tview.Form {
	as.infoView.SetText("Enter new account details.")

	addAccountForm := tview.NewForm()
	addAccountForm.SetBorder(true)
	addAccountForm.SetTitle(fmt.Sprintf("[ Open New Account: %s]", customer.name))
	addAccountForm.SetTitleAlign(tview.AlignLeft)
	addAccountForm.AddInputField("Account name", "Savings", 20, nil, nil)
	addAccountForm.AddButton("Create", func() {
		accountField := addAccountForm.GetFormItemByLabel("Account name").(*tview.InputField)

		accountName := accountField.GetText()

		if accountName == "" {
			// Form not complete.
			as.infoView.SetText("Account name must not be blank.")
			return
		}

		as.submitOpenNewAccount(customer.id, accountName)

		as.infoView.SetText("Open new account submitted.")
		as.switchMainView(as.customerMenu)
	})
	addAccountForm.AddButton("Cancel", func() {
		as.infoView.SetText("Open new account cancelled.")
		as.switchMainView(as.customerMenu)
	})
	addAccountForm.SetCancelFunc(func() {
		as.infoView.SetText("Open new account cancelled.")
		as.switchMainView(as.customerMenu)
	})

	return addAccountForm
}

func (as *appState) createDepositForm(customer customerData) *tview.Form {
	as.infoView.SetText("Enter deposit details.")

	depositForm := tview.NewForm()
	depositForm.SetBorder(true)
	depositForm.SetTitle(fmt.Sprintf("[ Deposit: %s ]", customer.name))
	depositForm.SetTitleAlign(tview.AlignLeft)
	accounts := as.fetchAccountsForCustomer(customer.id)
	options := make([]string, 0, len(accounts))
	for _, a := range accounts {
		options = append(options, a.display)
	}
	depositForm.AddDropDown("Account", options, 0, nil)
	depositForm.AddInputField("Amount", "0.00", 10, tview.InputFieldFloat, nil)
	depositForm.AddButton("Deposit", func() {
		accountField := depositForm.GetFormItemByLabel("Account").(*tview.DropDown)
		amountField := depositForm.GetFormItemByLabel("Amount").(*tview.InputField)

		accountIndex, _ := accountField.GetCurrentOption()
		amountValue := amountField.GetText()

		accountID := accounts[accountIndex].id

		dollarsAmount, err := strconv.ParseFloat(amountValue, 64)
		if err != nil || dollarsAmount == 0 {
			// Invalid amount
			as.infoView.SetText("Amount must not be zero.")
			return
		}

		centsAmount := int64(dollarsAmount * 100)
		as.submitDeposit(accountID, centsAmount)

		as.infoView.SetText(fmt.Sprintf("Deposit submitted for $%.2f.", dollarsAmount))
		as.switchMainView(as.customerMenu)
	})
	depositForm.AddButton("Cancel", func() {
		as.infoView.SetText("Deposit cancelled.")
		as.switchMainView(as.customerMenu)
	})
	depositForm.SetCancelFunc(func() {
		as.infoView.SetText("Deposit cancelled.")
		as.switchMainView(as.customerMenu)
	})

	return depositForm
}

func (as *appState) createWithdrawForm(customer customerData) *tview.Form {
	as.infoView.SetText("Enter withdrawal details.")

	withdrawForm := tview.NewForm()
	withdrawForm.SetBorder(true)
	withdrawForm.SetTitle(fmt.Sprintf("[ Withdraw: %s ]", customer.name))
	withdrawForm.SetTitleAlign(tview.AlignLeft)
	accounts := as.fetchAccountsForCustomer(customer.id)
	options := make([]string, 0, len(accounts))
	for _, a := range accounts {
		options = append(options, a.display)
	}
	withdrawForm.AddDropDown("Account", options, 0, nil)
	withdrawForm.AddInputField("Amount", "0.00", 10, tview.InputFieldFloat, nil)
	withdrawForm.AddButton("Withdraw", func() {
		accountField := withdrawForm.GetFormItemByLabel("Account").(*tview.DropDown)
		amountField := withdrawForm.GetFormItemByLabel("Amount").(*tview.InputField)

		accountIndex, _ := accountField.GetCurrentOption()
		amountValue := amountField.GetText()

		accountID := accounts[accountIndex].id

		dollarsAmount, err := strconv.ParseFloat(amountValue, 64)
		if err != nil || dollarsAmount == 0 {
			// Invalid amount
			as.infoView.SetText("Amount must not be zero.")
			return
		}

		centsAmount := int64(dollarsAmount * 100)
		if centsAmount > accounts[accountIndex].balance {
			// Insufficient funds
			as.infoView.SetText("Insufficient funds available for withdrawal.")
			return
		}

		as.submitWithdraw(accountID, centsAmount, businessDayFromTime(as.time))

		as.infoView.SetText(fmt.Sprintf("Withdraw submitted for $%.2f.", dollarsAmount))
		as.switchMainView(as.customerMenu)
	})
	withdrawForm.AddButton("Cancel", func() {
		as.infoView.SetText("Withdraw cancelled.")
		as.switchMainView(as.customerMenu)
	})
	withdrawForm.SetCancelFunc(func() {
		as.infoView.SetText("Withdraw cancelled.")
		as.switchMainView(as.customerMenu)
	})

	return withdrawForm
}

func (as *appState) createTransferForm(customer customerData) *tview.Form {
	as.infoView.SetText("Enter transfer details. You may schedule it for 'Today', 'Tomorrow' or a date formatted as 'YYYY-MM-DD'.")

	transferForm := tview.NewForm()
	transferForm.SetBorder(true)
	transferForm.SetTitle(fmt.Sprintf("[ Transfer: %s ]", customer.name))
	transferForm.SetTitleAlign(tview.AlignLeft)
	myAccounts := as.fetchAccountsForCustomer(customer.id)
	fromOptions := make([]string, 0, len(myAccounts))
	for _, a := range myAccounts {
		fromOptions = append(fromOptions, a.display)
	}
	transferForm.AddDropDown("From Account", fromOptions, 0, nil)
	allAccounts := as.fetchAccountsForAllCustomers()
	toOptions := make([]string, 0, len(allAccounts))
	for _, a := range allAccounts {
		toOptions = append(toOptions, a.display)
	}
	transferForm.AddDropDown("To Account", toOptions, 0, nil)
	transferForm.AddInputField("Amount", "0.00", 10, tview.InputFieldFloat, nil)
	transferForm.AddInputField("Scheduled Date", "Today", 20, nil, nil)
	transferForm.AddButton("Transfer", func() {
		fromAccountField := transferForm.GetFormItemByLabel("From Account").(*tview.DropDown)
		toAccountField := transferForm.GetFormItemByLabel("To Account").(*tview.DropDown)
		amountField := transferForm.GetFormItemByLabel("Amount").(*tview.InputField)
		dateField := transferForm.GetFormItemByLabel("Scheduled Date").(*tview.InputField)

		fromAccountIndex, _ := fromAccountField.GetCurrentOption()
		toAccountIndex, _ := toAccountField.GetCurrentOption()
		amountValue := amountField.GetText()
		scheduledDate := strings.ToLower(dateField.GetText())

		fromAccountID := myAccounts[fromAccountIndex].id
		toAccountID := allAccounts[toAccountIndex].accountID

		if fromAccountID == toAccountID {
			// Invalid account
			as.infoView.SetText("Transfer To account must be different to the From account.")
			return
		}

		dollarsAmount, err := strconv.ParseFloat(amountValue, 64)
		if err != nil || dollarsAmount == 0 {
			// Invalid amount
			as.infoView.SetText("Amount must not be zero.")
			return
		}

		centsAmount := int64(dollarsAmount * 100)
		if centsAmount > myAccounts[fromAccountIndex].balance {
			// Insufficient funds
			as.infoView.SetText("Insufficient funds available for transfer.")
			return
		}

		if scheduledDate == "" {
			// Invalid date
			as.infoView.SetText("Scheduled date must not be empty. You may use 'today', 'tomorrow' or a date formatted as 'YYYY-MM-DD'.")
			return
		} else if scheduledDate == "today" {
			scheduledDate = businessDayFromTime(as.time)
		} else if scheduledDate == "tomorrow" {
			scheduledDate = businessDayFromTime(as.time.Add(time.Hour * 24))
		} else {
			t, err := time.Parse(messages.BusinessDateFormat, scheduledDate)
			if err != nil {
				// Invalid date
				as.infoView.SetText("Scheduled date is invalid. You may use 'today', 'tomorrow' or a date formatted as 'YYYY-MM-DD'.")
				return
			}
			scheduledDate = businessDayFromTime(t)
		}

		as.submitTransfer(fromAccountID, toAccountID, centsAmount, scheduledDate)

		as.infoView.SetText(fmt.Sprintf("Transfer submitted for $%.2f to occur at %s.", dollarsAmount, scheduledDate))
		as.switchMainView(as.customerMenu)
	})
	transferForm.AddButton("Cancel", func() {
		as.infoView.SetText("Transfer cancelled.")
		as.switchMainView(as.customerMenu)
	})
	transferForm.SetCancelFunc(func() {
		as.infoView.SetText("Transfer cancelled.")
		as.switchMainView(as.customerMenu)
	})

	return transferForm
}

func (as *appState) submitNewCustomerSignup(customerName, accountName string) {
	as.engine.Dispatch(
		context.Background(),
		commands.OpenAccountForNewCustomer{
			CustomerID:   generateCustomerNumber(),
			CustomerName: customerName,
			AccountID:    generateAccountNumber(),
			AccountName:  accountName,
		},
		engine.WithCurrentTime(as.time),
	)
}

func (as *appState) submitOpenNewAccount(customerID string, accountName string) {
	as.engine.Dispatch(
		context.Background(),
		commands.OpenAccount{
			CustomerID:  customerID,
			AccountID:   generateAccountNumber(),
			AccountName: accountName,
		},
		engine.WithCurrentTime(as.time),
	)
}

func (as *appState) submitDeposit(accountID string, centsAmount int64) {
	as.engine.Dispatch(
		context.Background(),
		commands.Deposit{
			TransactionID: generateID(),
			AccountID:     accountID,
			Amount:        centsAmount,
		},
		engine.WithCurrentTime(as.time),
	)
}

func (as *appState) submitWithdraw(accountID string, centsAmount int64, scheduledDate string) {
	as.engine.Dispatch(
		context.Background(),
		commands.Withdraw{
			TransactionID: generateID(),
			AccountID:     accountID,
			Amount:        centsAmount,
			ScheduledDate: scheduledDate,
		},
		engine.WithCurrentTime(as.time),
	)
}

func (as *appState) submitTransfer(fromAccountID, toAccountID string, centsAmount int64, scheduledDate string) {
	as.engine.Dispatch(
		context.Background(),
		commands.Transfer{
			TransactionID: generateID(),
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        centsAmount,
			ScheduledDate: scheduledDate,
		},
		engine.WithCurrentTime(as.time),
	)
}

func (as *appState) fetchCustomers() []customerData {
	rows, err := as.db.Query(
		`SELECT
			id,
			name
		FROM customer
		ORDER BY name`,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	results := make([]customerData, 0, 8)

	for rows.Next() {
		data := customerData{}

		if err := rows.Scan(
			&data.id,
			&data.name,
		); err != nil {
			panic(err)
		}

		data.display = fmt.Sprintf("%s (#%s)", data.name, data.id)

		results = append(results, data)
	}

	return results
}

func (as *appState) fetchAccountsForCustomer(customerID string) []accountData {
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

	results := make([]accountData, 0, 8)

	for rows.Next() {
		data := accountData{}

		if err := rows.Scan(
			&data.id,
			&data.name,
			&data.balance,
		); err != nil {
			panic(err)
		}

		dollars := float64(data.balance) / 100.0
		data.display = fmt.Sprintf("#%s (%s) $%.2f", data.id, data.name, dollars)

		results = append(results, data)
	}

	return results
}

func (as *appState) fetchAccountsForAllCustomers() []customerAccountData {
	rows, err := as.db.Query(
		`SELECT
				a.id,
				a.name,
				a.customer_id,
				c.name AS customer_name,
				a.balance
			FROM account AS a
			INNER JOIN customer AS c
			ON c.id = a.customer_id
			ORDER BY c.name, a.name`,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	results := make([]customerAccountData, 0, 8)

	for rows.Next() {
		data := customerAccountData{}

		if err := rows.Scan(
			&data.accountID,
			&data.accountName,
			&data.customerID,
			&data.customerName,
			&data.balance,
		); err != nil {
			panic(err)
		}

		dollars := float64(data.balance) / 100.0
		data.display = fmt.Sprintf("%s - #%s (%s) $%.2f", data.customerName, data.accountID, data.accountName, dollars)

		results = append(results, data)
	}

	return results
}

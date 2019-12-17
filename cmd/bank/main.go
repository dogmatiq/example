package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/database"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/testkit/engine"
	_ "github.com/mattn/go-sqlite3"
)

func businessDayFromTime(t time.Time) string {
	return t.Format(messages.BusinessDateFormat)
}

func readString(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.Replace(text, "\n", "", -1)
}

type appState struct {
	input  *bufio.Reader
	engine *engine.Engine
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
		input:  bufio.NewReader(os.Stdin),
		engine: en,
		time:   time.Date(2001, 10, 20, 11, 22, 33, 44, time.UTC),
	}

	fmt.Println()
	fmt.Println("Dogma Example Bank")
	fmt.Println("==================")

	as.runMainMenu()
}

// showMenu shows a title and a list of options.
func (as appState) showMenu(title string, items []menuItem) string {
	fmt.Println()
	fmt.Printf("Time: %s", as.time)
	fmt.Println()
	fmt.Println()
	fmt.Println(title)
	fmt.Println("--------------------")
	fmt.Println()
	fmt.Println("Options:")
	for _, item := range items {
		fmt.Printf("  [%s] - %s\n", item.option, item.description)
	}

	var input string
	for {
		fmt.Println()
		fmt.Print("Choice: ")
		input = readString(as.input)

		for _, item := range items {
			if item.option == input {
				return item.command
			}
		}

		fmt.Printf("Invalid selection: %s\n", input)
	}
}

func (as appState) runMainMenu() {
	mainMenuItems := []menuItem{
		{
			option:      "h",
			description: "Say Hello World!",
			command:     "hello",
		},
		{
			option:      "x",
			description: "Exit banking sytem",
			command:     "exit",
		},
	}

	for {
		switch as.showMenu("Main Menu", mainMenuItems) {
		case "exit":
			fmt.Println()
			fmt.Println("Goodbye")
			return
		case "hello":
			fmt.Println()
			fmt.Println("Hello World!")
		}
	}
}

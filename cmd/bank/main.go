package main

import (
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/cmd/bank/ui"
	"github.com/dogmatiq/example/database"
	"github.com/dogmatiq/testkit/engine"
	_ "github.com/mattn/go-sqlite3"
)

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

	ui.Run(db, en)
}

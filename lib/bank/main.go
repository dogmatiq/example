package main

import (
	"context"
	"errors"
	"io"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/database"
	_ "github.com/mattn/go-sqlite3"
)

// DogmaV1 is a Dogma plugin.
var DogmaV1 v1

type v1 struct{}

func (p *v1) ListApplications() []string {
	return []string{"bank"}
}

func (p *v1) OpenApplication(
	ctx context.Context,
	name string,
) (dogma.Application, io.Closer, error) {
	if name != "bank" {
		return nil, nil, errors.New("unknown application")
	}

	db := database.MustNew()

	app, err := example.NewApp(db)
	if err != nil {
		return nil, nil, err
	}

	return app, db, nil
}

func main() {}

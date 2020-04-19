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

func NewDogmaPluginV1(ctx context.Context) (interface{}, error) {
	return plugin{}, nil
}

type plugin struct{}

func (plugin) ApplicationKeys() []string {
	return []string{example.AppKey}
}

func (plugin) NewApplication(
	ctx context.Context,
	k string,
) (dogma.Application, io.Closer, error) {
	if k != example.AppKey {
		return nil, nil, errors.New("unrecognized application")
	}

	db, err := database.New()
	if err != nil {
		return nil, nil, err
	}

	app, err := example.NewApp(db)
	if err != nil {
		return nil, nil, err
	}

	return app, db, nil
}

func main() {}

package main

import (
	"fmt"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/graphkit"
)

func main() {
	app := &example.App{}

	cfg := configkit.FromApplication(app)

	g, err := graphkit.Generate(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(g.String())
}

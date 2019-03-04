package main

import (
	"fmt"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/graphkit"
)

func main() {
	g, err := graphkit.Generate(&example.App{})
	if err != nil {
		panic(err)
	}

	fmt.Println(g.String())
}

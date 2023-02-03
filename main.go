package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	raw, err := os.ReadFile("./inputs/kitchen-cabinets.json")
	if err != nil {
		panic(err)
	}

	res, err := ProcessCabinets(raw)
	if err != nil {
		panic(err)
	}

	spew.Dump(res)
}

package main

import (
	"io"
	"lflpo/pkg/lflpo"
	"lflpo/pkg/resolver"
	"os"
)

func main() {
	planning, err := os.Open("planning-reel.json")
	if err != nil {
		panic(err)
	}
	defer planning.Close()
	byteValue, err := io.ReadAll(planning)
	if err != nil {
		panic(err)
	}

	split, err := lflpo.NewSplit(byteValue)
	if err != nil {
		panic(err)
	}

	resolver := resolver.NewResolver(split)
	results := resolver.ComputeScenarios()
	results.Print()
}

package main

import (
	"fmt"
	"io"
	"lflpo/pkg/lflpo"
	"math/rand"
	"os"
	"time"
)

func main() {
	planning, err := os.Open("planning-reel-vierge.json")
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
	rand.Seed(time.Now().Unix())
	for !split.IsFinished() {
		meta, match, err := split.GetNextMatch()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Day %d, Match %d: %s vs %s\n", meta.Day+1, meta.Match+1, match.BlueTeam, match.RedTeam)
		teams := []string{match.BlueTeam, match.RedTeam}
		err = split.ResolveNextMatch(teams[rand.Intn(len(teams))])
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()
	fmt.Println("Final ranking:")
	for i, team := range split.Ranking() {
		fmt.Printf("%d: %s (%d - %d)\n", i, team.Name, team.Win.Total, team.Lose)
	}
}

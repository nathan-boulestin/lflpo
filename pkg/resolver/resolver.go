package resolver

import (
	"fmt"
	"lflpo/pkg/lflpo"
	"math"
	"sort"
)

type winner int

const (
	winnerBlue winner = iota
	winnerRed
)

type Resolver struct {
	split           lflpo.Split
	matchsRemaining int
	scenarios       int
}

func NewResolver(split lflpo.Split) Resolver {
	matchsPlayed := 0
num:
	for _, day := range split.Days {
		for _, match := range day.Matchs {
			if match.Winner != "" {
				matchsPlayed++
			} else {
				break num
			}
		}
	}
	totalMatchs := len(split.Days) * len(split.Days[0].Matchs)
	matchsRemaining := totalMatchs - matchsPlayed
	totalScenarios := math.Pow(2, float64(matchsRemaining))
	if totalScenarios > 9223372036854775807 {
		panic(fmt.Sprintf("too many scenarios (%f), cannot compute.", totalScenarios))
	}

	return Resolver{
		split:           split,
		matchsRemaining: matchsRemaining,
		scenarios:       int(totalScenarios),
	}
}

type Result struct {
	totalScenarios     int
	teamQualifications map[string]int
}

func (r *Resolver) ComputeScenarios() Result {
	result := Result{
		totalScenarios:     r.scenarios,
		teamQualifications: make(map[string]int, len(r.split.Teams)),
	}
	for _, team := range r.split.Teams {
		result.teamQualifications[team.Name] = 0
	}
	for i := 0; i < r.scenarios; i++ {
		fmt.Printf("\rScenarios compute: %d%% (%d)", i*100/r.scenarios, i)
		scenario := make([]winner, 0, r.matchsRemaining)
		for j := 0; j < r.matchsRemaining; j++ {
			if hasBit(i, uint(j)) {
				scenario = append(scenario, winnerBlue)
			} else {
				scenario = append(scenario, winnerRed)
			}
		}
		split := r.split.Copy()
		for _, gameResult := range scenario {
			_, nextMatch, err := split.GetNextMatch()
			if err != nil {
				panic(err)
			}
			if gameResult == winnerBlue {
				split.ResolveNextMatch(nextMatch.BlueTeam)
			} else {
				split.ResolveNextMatch(nextMatch.RedTeam)
			}
		}
		if !split.IsFinished() {
			panic("split should be finished")
		}
		for _, team := range split.Qualified() {
			if n, ok := result.teamQualifications[team.Name]; ok {
				n += 1
				result.teamQualifications[team.Name] = n
			}
		}
	}

	return result
}

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

type TeamQualification struct {
	Key   string
	Value int
}

type TeamQualificationList []TeamQualification

func (t TeamQualificationList) Len() int           { return len(t) }
func (t TeamQualificationList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TeamQualificationList) Less(i, j int) bool { return t[i].Value > t[j].Value }

func (r Result) Print() {
	fmt.Printf("Number of possible scenarios: %d\n", r.totalScenarios)
	t := make(TeamQualificationList, len(r.teamQualifications))
	i := 0
	for team, qualifications := range r.teamQualifications {
		t[i] = TeamQualification{Key: team, Value: qualifications * 100 / r.totalScenarios}
		i++
	}
	sort.Sort(t)
	for _, tq := range t {
		fmt.Printf("Team %s: %d%%\n", tq.Key, tq.Value)

	}

	fmt.Printf("Unsupported tie break scenario: %d\n", lflpo.UnsupportedTieBreakScenario)
}

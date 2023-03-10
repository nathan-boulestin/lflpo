package lflpo

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

var UnsupportedTieBreakScenario int = 0

type Split struct {
	Days            []Day           `json:"days"`
	Teams           map[string]Team `json:"teams"`
	next            *NextMatch
	PlacesInPlayoff int `json:"placesInPlayoff"`
}

type Day struct {
	Matchs []Match `json:"matchs"`
}

type Match struct {
	BlueTeam string `json:"blue"`
	RedTeam  string `json:"red"`
	Winner   string `json:"winner"`
}

type NextMatch struct {
	Day          int
	Match        int
	IsSecondHalf bool
}

func NewSplit(file []byte) (Split, error) {
	split := Split{}
	if err := json.Unmarshal(file, &split); err != nil {
		return Split{}, err
	}
	split.next = &NextMatch{
		Day:          0,
		Match:        0,
		IsSecondHalf: false,
	}
days:
	for _, day := range split.Days {
		for _, match := range day.Matchs {
			if match.Winner == "" {
				break days
			}
			err := split.ResolveNextMatch(match.Winner)
			if err != nil {
				return Split{}, err
			}
		}
	}
	return split, nil
}

func (s Split) GetNextMatch() (NextMatch, Match, error) {
	if s.next == nil {
		return NextMatch{}, Match{}, errors.New("no more matchs to play")
	}
	match := s.Days[s.next.Day].Matchs[s.next.Match]
	return *s.next, Match{
		BlueTeam: match.BlueTeam,
		RedTeam:  match.RedTeam,
	}, nil
}

func (s *Split) ResolveNextMatch(winner string) error {
	if s.next == nil {
		return errors.New("no more matchs to play")
	}
	match := s.Days[s.next.Day].Matchs[s.next.Match]
	if winner != match.BlueTeam && winner != match.RedTeam {
		return fmt.Errorf("team %s is not part of %s vs %s", winner, match.BlueTeam, match.RedTeam)
	}
	match.Winner = winner
	s.Days[s.next.Day].Matchs[s.next.Match] = match
	var loser string
	if match.Winner == match.BlueTeam {
		loser = match.RedTeam
	} else {
		loser = match.BlueTeam
	}
	if team, ok := s.Teams[winner]; ok {
		team.WinVersus(loser, s.next.IsSecondHalf)
		s.Teams[winner] = team
	}
	if team, ok := s.Teams[loser]; ok {
		team.Lose += 1
		s.Teams[loser] = team
	}
	if s.next.Match+1 >= len(s.Days[s.next.Day].Matchs) {
		if s.next.Day+1 >= len(s.Days) {
			s.next = nil
		} else {
			s.next = &NextMatch{
				Day:          s.next.Day + 1,
				Match:        0,
				IsSecondHalf: s.next.Day+1 >= len(s.Days)/2,
			}
		}
	} else {
		s.next = &NextMatch{
			Day:          s.next.Day,
			Match:        s.next.Match + 1,
			IsSecondHalf: s.next.Day+1 >= len(s.Days)/2,
		}
	}
	return nil
}

func (s Split) IsFinished() bool {
	return s.next == nil
}

func (s Split) Ranking() []Team {
	ranking := make([]Team, 0, len(s.Teams))
	for _, team := range s.Teams {
		ranking = append(ranking, team)
	}
	sort.Slice(ranking, func(i, j int) bool {
		if ranking[i].Win.Total == ranking[j].Win.Total {
			return resolveTieBreak(ranking[i], ranking[j])
		}
		return ranking[i].Win.Total > ranking[j].Win.Total
	})
	return ranking
}

func (s Split) Qualified() []Team {
	return s.Ranking()[:s.PlacesInPlayoff]
}

func (s Split) Copy() Split {
	days := make([]Day, 0, len(s.Days))
	days = append(days, s.Days...)
	teams := make(map[string]Team)
	for key, team := range s.Teams {
		teams[key] = Team{
			Name: team.Name,
			Win:  team.Win,
			Lose: team.Lose,
		}
	}
	return Split{
		Days:  days,
		Teams: teams,
		next: &NextMatch{
			Day:   s.next.Day,
			Match: s.next.Match,
		},
		PlacesInPlayoff: s.PlacesInPlayoff,
	}
}

func resolveTieBreak(team1, team2 Team) bool {
	var result bool
	team1Versus := team1.Win.Versus[team2.Name]
	team2Versus := team2.Win.Versus[team1.Name]
	if team1Versus == team2Versus {
		team1SecondHalf := team1.Win.SecondHalfWins
		team2SecondHalf := team2.Win.SecondHalfWins
		if team1SecondHalf == team2SecondHalf {
			UnsupportedTieBreakScenario += 1
		}
		result = team1SecondHalf > team2SecondHalf
	} else {
		result = team1.Win.Versus[team2.Name] > team2.Win.Versus[team1.Name]
	}
	return result
}

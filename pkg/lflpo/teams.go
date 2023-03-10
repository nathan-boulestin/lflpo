package lflpo

type Team struct {
	Name string `json:"name"`
	Win  WinDetails
	Lose int
}

type WinDetails struct {
	Total          int
	Versus         map[string]int
	SecondHalfWins int
}

func (t *Team) WinVersus(team string, isSecondHalf bool) {
	if t.Win.Versus == nil {
		t.Win.Versus = make(map[string]int)
	}
	t.Win.Total += 1
	if current, ok := t.Win.Versus[team]; ok {
		new := current + 1
		t.Win.Versus[team] = new
	} else {
		t.Win.Versus[team] = 1
	}
	if isSecondHalf {
		t.Win.SecondHalfWins += 1
	}
}

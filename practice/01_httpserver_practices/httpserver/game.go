package poker

import "time"

type TexasHoldem struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewGame(store PlayerStore, alerter BlindAlerter) *TexasHoldem {
	return &TexasHoldem{
		store:   store,
		alerter: alerter,
	}
}

func (g *TexasHoldem) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, amount := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, amount)
		blindTime = blindTime + blindIncrement
	}
}

func (g *TexasHoldem) Finish(winner string) {
	g.store.RecordWin(winner)
}

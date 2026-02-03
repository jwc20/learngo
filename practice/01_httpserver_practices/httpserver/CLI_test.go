package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/jwc20/learngowithtests/practice/01_httpserver_practices/httpserver"
)

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips in %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := strings.NewReader("3\nChris wins\n")
		cli := poker.NewCLI(game, in, stdout)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := strings.NewReader("8\nCleo wins\n")
		cli := poker.NewCLI(game, in, dummyStdOut)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

}
func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}
	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayers int) {
	t.Helper()
	if game.StartedWith != numberOfPlayers {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayers, game.StartedWith)
	}
}

func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	if game.FinishedWith != winner {
		t.Errorf("wanted Finish called with %q but got %q", winner, game.FinishedWith)
	}
}

func assertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

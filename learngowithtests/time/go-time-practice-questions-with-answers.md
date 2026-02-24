# Go CLI Practice Questions with Answers - Time & Scheduling

## Blind Alert Scheduling Questions

### Q1: Create SpyBlindAlerter Test
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "undefined: SpyBlindAlerter"

**Question:**
Create a new test that verifies blind alerts are scheduled:
1. Create a test `"it schedules printing of blind values"`
2. Set up user input with `strings.NewReader("Chris wins\n")`
3. Create a `StubPlayerStore` and a `SpyBlindAlerter`
4. Pass the `blindAlerter` as a third argument to `NewCLI`
5. Call `cli.PlayPoker()`
6. Check that `len(blindAlerter.alerts)` equals 1

**Answer:**
```go
// CLI_test.go
t.Run("it schedules printing of blind values", func(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &poker.StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(playerStore, in, blindAlerter)
	cli.PlayPoker()

	if len(blindAlerter.alerts) != 1 {
		t.Fatal("expected a blind alert to be scheduled")
	}
})
```

---

### Q2: Define SpyBlindAlerter
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "too many arguments in call to poker.NewCLI"

**Question:**
Create the `SpyBlindAlerter` type with:
- A field `alerts` that is a slice of structs containing:
  - `scheduledAt time.Duration`
  - `amount int`
- A method `ScheduleAlertAt(duration time.Duration, amount int)` that appends to alerts

**Answer:**
```go
// CLI_test.go
type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{duration, amount})
}
```

---

### Q3: Define BlindAlerter Interface
**File to edit:** `CLI.go`  
**Expected test result:** Test should fail with "expected a blind alert to be scheduled"

**Question:**
Create a `BlindAlerter` interface with one method:
```go
ScheduleAlertAt(duration time.Duration, amount int)
```

**Answer:**
```go
// CLI.go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}
```

---

### Q4: Add BlindAlerter to CLI Constructor
**File to edit:** `CLI.go`  
**Expected test result:** Other tests should fail (need updating)

**Question:**
Update the `CLI` struct and constructor:
1. Add `alerter BlindAlerter` field to the CLI struct
2. Update `NewCLI` to accept `alerter BlindAlerter` as the third parameter
3. Initialize the alerter field in the constructor

**Answer:**
```go
// CLI.go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	alerter     BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		alerter:     alerter,
	}
}
```

---

### Q5: Add Dummy Alerter for Other Tests
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should compile again

**Question:**
Create a dummy alerter variable:
```go
var dummySpyAlerter = &SpyBlindAlerter{}
```

Update all other tests to pass `dummySpyAlerter` to `NewCLI`.

**Answer:**
```go
// CLI_test.go
var dummySpyAlerter = &SpyBlindAlerter{}

// Then in each test that doesn't care about alerts:
cli := poker.NewCLI(playerStore, in, dummySpyAlerter)
```

---

### Q6: Make First Alert Test Pass
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

**Question:**
In the `PlayPoker()` method, schedule one alert before reading user input:
```go
cli.alerter.ScheduleAlertAt(5*time.Second, 100)
```

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

---

## Multiple Alerts Questions

### Q7: Add Table-Based Test for Multiple Alerts
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - only one alert scheduled

**Question:**
Expand the test to check multiple scheduled alerts:
1. Create a slice of test cases with expected schedule times and amounts
2. Loop through the cases
3. For each case, verify the alert at that index has correct `scheduledAt` and `amount`

Expected alerts:
- 0 seconds, 100 chips
- 10 minutes, 200 chips
- 20 minutes, 300 chips
- 30 minutes, 400 chips
- 40 minutes, 500 chips
- 50 minutes, 600 chips
- 60 minutes, 800 chips
- 70 minutes, 1000 chips
- 80 minutes, 2000 chips
- 90 minutes, 4000 chips
- 100 minutes, 8000 chips

**Answer:**
```go
// CLI_test.go
t.Run("it schedules printing of blind values", func(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &poker.StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(playerStore, in, blindAlerter)
	cli.PlayPoker()

	cases := []struct {
		expectedScheduleTime time.Duration
		expectedAmount       int
	}{
		{0 * time.Second, 100},
		{10 * time.Minute, 200},
		{20 * time.Minute, 300},
		{30 * time.Minute, 400},
		{40 * time.Minute, 500},
		{50 * time.Minute, 600},
		{60 * time.Minute, 800},
		{70 * time.Minute, 1000},
		{80 * time.Minute, 2000},
		{90 * time.Minute, 4000},
		{100 * time.Minute, 8000},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			alert := blindAlerter.alerts[i]

			amountGot := alert.amount
			if amountGot != c.expectedAmount {
				t.Errorf("got amount %d, want %d", amountGot, c.expectedAmount)
			}

			gotScheduledTime := alert.scheduledAt
			if gotScheduledTime != c.expectedScheduleTime {
				t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, c.expectedScheduleTime)
			}
		})
	}
})
```

---

### Q8: Implement Multiple Alert Scheduling
**File to edit:** `CLI.go`  
**Expected test result:** All tests should pass

**Question:**
Update `PlayPoker()` to schedule all blind alerts:
1. Create a slice `blinds` with values: 100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000
2. Initialize `blindTime` to 0
3. Loop through blinds, scheduling each alert
4. Increment `blindTime` by 10 minutes after each alert

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}

	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

---

### Q9: Extract scheduleBlindAlerts Method
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Refactor by creating a `scheduleBlindAlerts()` method:
1. Move the scheduling logic from `PlayPoker()` into this new method
2. Call the method from `PlayPoker()`

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}
```

---

## Test Refactoring Questions

### Q10: Create ScheduledAlert Type
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Define a new type to replace the anonymous struct:
1. Create `type scheduledAlert struct` with `at time.Duration` and `amount int` fields
2. Add a `String()` method that returns a formatted string like "100 chips at 0s"
3. Update `SpyBlindAlerter` to use `[]scheduledAlert` instead of the anonymous struct slice

**Answer:**
```go
// CLI_test.go
type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}
```

---

### Q11: Refactor Test to Use ScheduledAlert
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Update the test cases:
1. Change the test cases slice from anonymous structs to `[]scheduledAlert`
2. Create an `assertScheduledAlert` helper function
3. Use the helper in the test loop

**Answer:**
```go
// CLI_test.go
t.Run("it schedules printing of blind values", func(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &poker.StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(playerStore, in, blindAlerter)
	cli.PlayPoker()

	cases := []scheduledAlert{
		{0 * time.Second, 100},
		{10 * time.Minute, 200},
		{20 * time.Minute, 300},
		{30 * time.Minute, 400},
		{40 * time.Minute, 500},
		{50 * time.Minute, 600},
		{60 * time.Minute, 800},
		{70 * time.Minute, 1000},
		{80 * time.Minute, 2000},
		{90 * time.Minute, 4000},
		{100 * time.Minute, 8000},
	}

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
})

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}
	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}
```

---

## Integration Questions

### Q12: Create BlindAlerter Implementation
**File to edit:** `blind_alerter.go` (new file)  
**Expected test result:** N/A (setup for integration)

**Question:**
Create a new file with:
1. Move the `BlindAlerter` interface definition here
2. Create `BlindAlerterFunc` type: `type BlindAlerterFunc func(duration time.Duration, amount int)`
3. Implement `ScheduleAlertAt` method on `BlindAlerterFunc` to satisfy the interface
4. Create `StdOutAlerter` function that uses `time.AfterFunc` to print to `os.Stdout`

**Answer:**
```go
// blind_alerter.go
package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
```

Remove the `BlindAlerter` interface from `CLI.go`.

---

### Q13: Wire Up BlindAlerter in Main
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** N/A (integration verification)

**Question:**
Update main to use the real alerter:
```go
poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter)).PlayPoker()
```

Test by running the application manually.

**Answer:**
```go
// cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/your-username/your-repo" // adjust to your module path
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	
	poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter)).PlayPoker()
}
```

Test by running:
```bash
go run cmd/cli/main.go
```

---

## Player Input Questions

### Q14: Add Dummy Variables
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (setup)

**Question:**
Add dummy variables for cleaner test setup:
```go
var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
```

**Answer:**
```go
// CLI_test.go
var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
```

These can now be used in tests where those dependencies aren't the focus.

---

### Q15: Test Player Prompt
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "too many arguments in call to poker.NewCLI"

**Question:**
Create a test that verifies the player prompt:
1. Create a new test `"it prompts the user to enter the number of players"`
2. Create a `bytes.Buffer` for stdout
3. Create CLI with dummy inputs and the stdout buffer
4. Call `PlayPoker()`
5. Assert that stdout contains "Please enter the number of players: "

**Answer:**
```go
// CLI_test.go
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
	stdout := &bytes.Buffer{}
	cli := poker.NewCLI(dummyPlayerStore, dummyStdIn, stdout, dummyBlindAlerter)
	cli.PlayPoker()

	got := stdout.String()
	want := "Please enter the number of players: "

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
})
```

---

### Q16: Add io.Writer to CLI
**File to edit:** `CLI.go`  
**Expected test result:** Test should fail - nothing printed to stdout

**Question:**
Update CLI to accept and store an `io.Writer`:
1. Add `out io.Writer` field to CLI struct
2. Update `NewCLI` signature to accept `out io.Writer` as third parameter (before alerter)
3. Initialize the field in constructor

Update other tests to pass the new parameter.

**Answer:**
```go
// CLI.go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	alerter     BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		out:         out,
		alerter:     alerter,
	}
}
```

Update all tests to pass `dummyStdOut` as the third argument.

---

### Q17: Implement Player Prompt
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

**Question:**
At the start of `PlayPoker()`, print the prompt:
```go
fmt.Fprint(cli.out, "Please enter the number of players: ")
```

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, "Please enter the number of players: ")
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

---

### Q18: Extract Prompt as Constant
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Create a constant:
```go
const PlayerPrompt = "Please enter the number of players: "
```

Use it in both the code and test.

**Answer:**
```go
// CLI.go
const PlayerPrompt = "Please enter the number of players: "

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

Update the test:
```go
// CLI_test.go
want := poker.PlayerPrompt
```

---

### Q19: Test Dynamic Scheduling Based on Players
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - scheduling doesn't use player count

**Question:**
Update the test to:
1. Provide input with "7\n" (7 players)
2. Check that alerts are scheduled with 12-minute intervals instead of 10
3. Verify at least the first 4 alerts: 0s/100, 12m/200, 24m/300, 36m/400

**Answer:**
```go
// CLI_test.go
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
	stdout := &bytes.Buffer{}
	in := strings.NewReader("7\n")
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(dummyPlayerStore, in, stdout, blindAlerter)
	cli.PlayPoker()

	got := stdout.String()
	want := poker.PlayerPrompt

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	cases := []scheduledAlert{
		{0 * time.Second, 100},
		{12 * time.Minute, 200},
		{24 * time.Minute, 300},
		{36 * time.Minute, 400},
	}

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
})
```

---

### Q20: Read Player Count and Use for Scheduling
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (some may need input updates)

**Question:**
Update `PlayPoker()` to:
1. After printing the prompt, read the player count with `cli.readLine()`
2. Convert it to an integer with `strconv.Atoi()`
3. Update `scheduleBlindAlerts` to accept `numberOfPlayers int` parameter
4. Calculate `blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute`
5. Use `blindIncrement` instead of the hardcoded 10 minutes

Fix other tests by providing appropriate player input (e.g., "5\n").

**Answer:**
```go
// CLI.go
import (
	"strconv"
	"strings"
	// ... other imports
)

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, _ := strconv.Atoi(cli.readLine())

	cli.scheduleBlindAlerts(numberOfPlayers)

	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}
```

Update other tests to provide player input:
```go
// For tests that were checking specific scheduling
in := strings.NewReader("5\nChris wins\n")
```

---

## Separation of Concerns Questions

### Q21: Create Game Type
**File to edit:** `game.go` (new file)  
**Expected test result:** Tests should still pass

**Question:**
Create a new `Game` struct that:
1. Has fields: `alerter BlindAlerter` and `store PlayerStore`
2. Has a `Start(numberOfPlayers int)` method - move scheduling logic here
3. Has a `Finish(winner string)` method - move RecordWin logic here

**Answer:**
```go
// game.go
package poker

import "time"

type Game struct {
	alerter BlindAlerter
	store   PlayerStore
}

func (g *Game) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (g *Game) Finish(winner string) {
	g.store.RecordWin(winner)
}
```

---

### Q22: Update CLI to Use Game
**File to edit:** `CLI.go`  
**Expected test result:** Tests should still pass

**Question:**
Refactor CLI to:
1. Replace `playerStore` and `alerter` fields with a single `game *Game` field
2. Update `NewCLI` to accept `alerter BlindAlerter` and `store PlayerStore`, create Game internally
3. Update `PlayPoker()` to call `game.Start()` and `game.Finish()` methods

**Answer:**
```go
// CLI.go
type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game *Game
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		in:  bufio.NewScanner(in),
		out: out,
		game: &Game{
			alerter: alerter,
			store:   store,
		},
	}
}

const PlayerPrompt = "Please enter the number of players: "

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersInput)

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
```

Remove the `scheduleBlindAlerts` method as it's now in `Game.Start()`.

---

### Q23: Create Game Constructor
**File to edit:** `game.go`  
**Expected test result:** Tests should compile more easily

**Question:**
Add a constructor:
```go
func NewGame(alerter BlindAlerter, store PlayerStore) *Game
```

**Answer:**
```go
// game.go
func NewGame(alerter BlindAlerter, store PlayerStore) *Game {
	return &Game{
		alerter: alerter,
		store:   store,
	}
}
```

---

### Q24: Update Tests to Use NewGame
**File to edit:** `CLI_test.go` and `cmd/cli/main.go`  
**Expected test result:** Tests should pass (cleanup)

**Question:**
Update test setup to use `NewGame`:
```go
game := poker.NewGame(blindAlerter, dummyPlayerStore)
cli := poker.NewCLI(in, stdout, game)
```

Also update `main.go` to use the new constructor.

**Answer:**
```go
// CLI_test.go - example test update
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
	stdout := &bytes.Buffer{}
	in := strings.NewReader("7\n")
	blindAlerter := &SpyBlindAlerter{}
	game := poker.NewGame(blindAlerter, dummyPlayerStore)

	cli := poker.NewCLI(in, stdout, game)
	cli.PlayPoker()

	// ... assertions
})
```

Update `NewCLI` signature:
```go
// CLI.go
func NewCLI(in io.Reader, out io.Writer, game *Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}
```

Update main:
```go
// cmd/cli/main.go
game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
cli := poker.NewCLI(os.Stdin, os.Stdout, game)
cli.PlayPoker()
```

---

### Q25: Create Game Tests File
**File to edit:** `game_test.go` (new file)  
**Expected test result:** New tests should pass

**Question:**
Move game-specific tests from CLI tests:
1. Create `TestGame_Start` with tests for 5 and 7 players
2. Create `TestGame_Finish` to test winner recording
3. These tests should focus on the `Game` type directly

**Answer:**
```go
// game_test.go
package poker

import (
	"fmt"
	"testing"
	"time"
)

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	store := &StubPlayerStore{}
	game := NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(cases []scheduledAlert, t *testing.T, blindAlerter *SpyBlindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}
```

---

### Q26: Introduce Game Interface
**File to edit:** `CLI.go` and `game.go`  
**Expected test result:** Tests should still pass

**Question:**
Create a `Game` interface with:
```go
Start(numberOfPlayers int)
Finish(winner string)
```

Change CLI's `game` field type from `*Game` to `Game` (the interface).

Rename the concrete implementation from `Game` to `TexasHoldem`.

**Answer:**
```go
// CLI.go
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}
```

```go
// game.go
type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}

func NewTexasHoldem(alerter BlindAlerter, store PlayerStore) *TexasHoldem {
	return &TexasHoldem{
		alerter: alerter,
		store:   store,
	}
}

func (g *TexasHoldem) Start(numberOfPlayers int) {
	// ... same implementation
}

func (g *TexasHoldem) Finish(winner string) {
	// ... same implementation
}
```

Update tests and main to use `NewTexasHoldem`.

---

### Q27: Create GameSpy
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests will need updating

**Question:**
Create a spy that implements the `Game` interface:
```go
type GameSpy struct {
    StartedWith  int
    FinishedWith string
}
```

Implement `Start` and `Finish` methods that record the arguments.

**Answer:**
```go
// CLI_test.go
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
```

---

### Q28: Simplify CLI Tests with GameSpy
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass with cleaner assertions

**Question:**
Refactor CLI tests to:
1. Use `GameSpy` instead of real game
2. Remove blind alert checking from CLI tests
3. Focus on: reading player count, calling `Start` with correct number, calling `Finish` with winner
4. Create helper functions like `assertGameStartedWith` and `assertFinishCalledWith`

**Answer:**
```go
// CLI_test.go
func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := strings.NewReader("3\nChris wins\n")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := strings.NewReader("8\nCleo wins\n")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})
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
```

---

## Error Handling Questions

### Q29: Test Invalid Player Input
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - game starts even with invalid input

**Question:**
Create a test for invalid player input:
1. Test name: "it prints an error when a non numeric value is entered and does not start the game"
2. Provide input "Pies\n"
3. Check that `game.StartCalled` is false (you'll need to add this field to GameSpy)

**Answer:**
```go
// CLI_test.go
t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
	game := &GameSpy{}
	stdout := &bytes.Buffer{}
	in := strings.NewReader("Pies\n")

	cli := poker.NewCLI(in, stdout, game)
	cli.PlayPoker()

	assertGameNotStarted(t, game)
})
```

---

### Q30: Handle Atoi Error
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

**Question:**
In `PlayPoker()`, check the error from `strconv.Atoi`:
```go
numberOfPlayers, err := strconv.Atoi(cli.readLine())
if err != nil {
    return
}
```

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		return
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}
```

---

### Q31: Test Error Message to User
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - no error message printed

**Question:**
Update the invalid input test to also check stdout:
1. Assert that stdout contains the prompt AND an error message
2. Use a temporary error message like "you're so silly" for now

**Answer:**
```go
// CLI_test.go
t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
	game := &GameSpy{}
	stdout := &bytes.Buffer{}
	in := strings.NewReader("Pies\n")

	cli := poker.NewCLI(in, stdout, game)
	cli.PlayPoker()

	assertGameNotStarted(t, game)

	gotPrompt := stdout.String()
	wantPrompt := poker.PlayerPrompt + "you're so silly"

	if gotPrompt != wantPrompt {
		t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
	}
})
```

---

### Q32: Print Error Message
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

**Question:**
Update the error handling to print to stdout:
```go
if err != nil {
    fmt.Fprint(cli.out, "you're so silly")
    return
}
```

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, "you're so silly")
		return
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}
```

---

### Q33: Create Error Constant
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Create a constant for the error message:
```go
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
```

Use it in both the code and test.

**Answer:**
```go
// CLI.go
const PlayerPrompt = "Please enter the number of players: "
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}
```

Update test:
```go
// CLI_test.go
wantPrompt := poker.PlayerPrompt + poker.BadPlayerInputErrMsg
```

---

### Q34: Create assertMessagesSentToUser Helper
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Create a helper function:
```go
func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string)
```

This should join all messages and compare against stdout contents. Use it to clean up tests.

**Answer:**
```go
// CLI_test.go
func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
```

Use it in tests:
```go
assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
```

---

### Q35: Refactor All CLI Tests
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass with improved readability

**Question:**
Refactor all CLI tests to:
1. Use helper functions for common assertions
2. Create a `userSends()` helper that creates input readers
3. Focus each test on a single concern
4. Remove redundant checks now that Game is separated

Example helpers:
- `assertMessagesSentToUser()`
- `assertGameStartedWith()`
- `assertGameNotStarted()`
- `assertFinishCalledWith()`
- `userSends(messages ...string) io.Reader`

**Answer:**
```go
// CLI_test.go
package poker_test

import (
	"bytes"
	"strings"
	"testing"
	
	poker "github.com/your-username/your-repo"
)

var dummyStdOut = &bytes.Buffer{}

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

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})
}

func userSends(messages ...string) *strings.Reader {
	return strings.NewReader(strings.Join(messages, "\n") + "\n")
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
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
```

---

## Final Integration Question

### Q36: Verify Complete Application
**Files to test:** All  
**Expected test result:** All tests pass, application runs correctly

**Question:**
1. Run `go test ./...` - all tests should pass
2. Run the CLI application
3. Test valid input (enter a number, then "{Name} wins")
4. Test invalid input (enter non-numeric value)
5. Verify blind alerts print at correct intervals (you may want to reduce timing for testing)
6. Confirm the web server still works alongside the CLI

**Answer:**
```bash
# 1. Run all tests
go test ./...
# Expected: PASS for all packages

# 2. Test the CLI (you may want to change blind intervals to seconds in code)
go run cmd/cli/main.go
# Enter: 5
# Wait and observe blind alerts printing every 10 minutes (or seconds if modified)
# Enter: Alice wins

# 3. Test invalid input
go run cmd/cli/main.go
# Enter: hello
# Should see error message and program exits

# 4. Test the web server still works
go run cmd/webserver/main.go
# In another terminal:
curl http://localhost:5000/league
# Should see league data including Alice
```

For easier testing of timing, temporarily change in `game.go`:
```go
blindIncrement := time.Duration(5+numberOfPlayers) * time.Second  // Changed from Minute
```

---

## Additional Challenge Question

### Q37: Handle Invalid Winner Input
**File to edit:** `CLI_test.go` and `CLI.go`  
**Expected test result:** Should gracefully handle bad input

**Question:**
The tutorial mentions: "What happens if instead of putting `Ruth wins` the user puts in `Lloyd is a killer`?"

Write a test for this scenario and implement error handling:
1. Test that invalid winner format is detected
2. Print an appropriate error message
3. Don't record a win for invalid input

**Answer:**
```go
// CLI_test.go
t.Run("it prints an error when winner input is invalid", func(t *testing.T) {
	game := &GameSpy{}
	stdout := &bytes.Buffer{}

	in := userSends("5", "Lloyd is a killer")
	cli := poker.NewCLI(in, stdout, game)

	cli.PlayPoker()

	assertGameStartedWith(t, game, 5)
	
	// Check that Finish was not called or was called with empty string
	if game.FinishedWith != "" {
		t.Errorf("should not have recorded a winner with invalid input")
	}
	
	assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputErrMsg)
})
```

```go
// CLI.go
const PlayerPrompt = "Please enter the number of players: "
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
const BadWinnerInputErrMsg = "Invalid winner format, please use '{Name} wins'"

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	if winner == "" || winner == winnerInput {
		// extractWinner returns empty or unchanged string if format is wrong
		fmt.Fprint(cli.out, BadWinnerInputErrMsg)
		return
	}

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	if !strings.HasSuffix(userInput, " wins") {
		return ""
	}
	return strings.Replace(userInput, " wins", "", 1)
}
```

Test it:
```bash
go test ./...
go run cmd/cli/main.go
# Enter: 5
# Enter: Bob is awesome
# Should see error message
```

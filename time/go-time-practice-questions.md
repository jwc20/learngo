# Go CLI Practice Questions - Time & Scheduling

## Blind Alert Scheduling Questions

### Q1: Create SpyBlindAlerter Test
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "undefined: SpyBlindAlerter"

Create a new test that verifies blind alerts are scheduled:
1. Create a test `"it schedules printing of blind values"`
2. Set up user input with `strings.NewReader("Chris wins\n")`
3. Create a `StubPlayerStore` and a `SpyBlindAlerter`
4. Pass the `blindAlerter` as a third argument to `NewCLI`
5. Call `cli.PlayPoker()`
6. Check that `len(blindAlerter.alerts)` equals 1

---

### Q2: Define SpyBlindAlerter
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "too many arguments in call to poker.NewCLI"

Create the `SpyBlindAlerter` type with:
- A field `alerts` that is a slice of structs containing:
  - `scheduledAt time.Duration`
  - `amount int`
- A method `ScheduleAlertAt(duration time.Duration, amount int)` that appends to alerts

---

### Q3: Define BlindAlerter Interface
**File to edit:** `CLI.go`  
**Expected test result:** Test should fail with "expected a blind alert to be scheduled"

Create a `BlindAlerter` interface with one method:
```go
ScheduleAlertAt(duration time.Duration, amount int)
```

---

### Q4: Add BlindAlerter to CLI Constructor
**File to edit:** `CLI.go`  
**Expected test result:** Other tests should fail (need updating)

Update the `CLI` struct and constructor:
1. Add `alerter BlindAlerter` field to the CLI struct
2. Update `NewCLI` to accept `alerter BlindAlerter` as the third parameter
3. Initialize the alerter field in the constructor

---

### Q5: Add Dummy Alerter for Other Tests
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should compile again

Create a dummy alerter variable:
```go
var dummySpyAlerter = &SpyBlindAlerter{}
```

Update all other tests to pass `dummySpyAlerter` to `NewCLI`.

---

### Q6: Make First Alert Test Pass
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

In the `PlayPoker()` method, schedule one alert before reading user input:
```go
cli.alerter.ScheduleAlertAt(5*time.Second, 100)
```

---

## Multiple Alerts Questions

### Q7: Add Table-Based Test for Multiple Alerts
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - only one alert scheduled

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

---

### Q8: Implement Multiple Alert Scheduling
**File to edit:** `CLI.go`  
**Expected test result:** All tests should pass

Update `PlayPoker()` to schedule all blind alerts:
1. Create a slice `blinds` with values: 100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000
2. Initialize `blindTime` to 0
3. Loop through blinds, scheduling each alert
4. Increment `blindTime` by 10 minutes after each alert

---

### Q9: Extract scheduleBlindAlerts Method
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring)

Refactor by creating a `scheduleBlindAlerts()` method:
1. Move the scheduling logic from `PlayPoker()` into this new method
2. Call the method from `PlayPoker()`

---

## Test Refactoring Questions

### Q10: Create ScheduledAlert Type
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (refactoring)

Define a new type to replace the anonymous struct:
1. Create `type scheduledAlert struct` with `at time.Duration` and `amount int` fields
2. Add a `String()` method that returns a formatted string like "100 chips at 0s"
3. Update `SpyBlindAlerter` to use `[]scheduledAlert` instead of the anonymous struct slice

---

### Q11: Refactor Test to Use ScheduledAlert
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (refactoring)

Update the test cases:
1. Change the test cases slice from anonymous structs to `[]scheduledAlert`
2. Create an `assertScheduledAlert` helper function
3. Use the helper in the test loop

---

## Integration Questions

### Q12: Create BlindAlerter Implementation
**File to edit:** `blind_alerter.go` (new file)  
**Expected test result:** N/A (setup for integration)

Create a new file with:
1. Move the `BlindAlerter` interface definition here
2. Create `BlindAlerterFunc` type: `type BlindAlerterFunc func(duration time.Duration, amount int)`
3. Implement `ScheduleAlertAt` method on `BlindAlerterFunc` to satisfy the interface
4. Create `StdOutAlerter` function that uses `time.AfterFunc` to print to `os.Stdout`

---

### Q13: Wire Up BlindAlerter in Main
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** N/A (integration verification)

Update main to use the real alerter:
```go
poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter)).PlayPoker()
```

Test by running the application manually.

---

## Player Input Questions

### Q14: Add Dummy Variables
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (setup)

Add dummy variables for cleaner test setup:
```go
var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
```

---

### Q15: Test Player Prompt
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "too many arguments in call to poker.NewCLI"

Create a test that verifies the player prompt:
1. Create a new test `"it prompts the user to enter the number of players"`
2. Create a `bytes.Buffer` for stdout
3. Create CLI with dummy inputs and the stdout buffer
4. Call `PlayPoker()`
5. Assert that stdout contains "Please enter the number of players: "

---

### Q16: Add io.Writer to CLI
**File to edit:** `CLI.go`  
**Expected test result:** Test should fail - nothing printed to stdout

Update CLI to accept and store an `io.Writer`:
1. Add `out io.Writer` field to CLI struct
2. Update `NewCLI` signature to accept `out io.Writer` as third parameter (before alerter)
3. Initialize the field in constructor

Update other tests to pass the new parameter.

---

### Q17: Implement Player Prompt
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

At the start of `PlayPoker()`, print the prompt:
```go
fmt.Fprint(cli.out, "Please enter the number of players: ")
```

---

### Q18: Extract Prompt as Constant
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring)

Create a constant:
```go
const PlayerPrompt = "Please enter the number of players: "
```

Use it in both the code and test.

---

### Q19: Test Dynamic Scheduling Based on Players
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - scheduling doesn't use player count

Update the test to:
1. Provide input with "7\n" (7 players)
2. Check that alerts are scheduled with 12-minute intervals instead of 10
3. Verify at least the first 4 alerts: 0s/100, 12m/200, 24m/300, 36m/400

---

### Q20: Read Player Count and Use for Scheduling
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (some may need input updates)

Update `PlayPoker()` to:
1. After printing the prompt, read the player count with `cli.readLine()`
2. Convert it to an integer with `strconv.Atoi()`
3. Update `scheduleBlindAlerts` to accept `numberOfPlayers int` parameter
4. Calculate `blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute`
5. Use `blindIncrement` instead of the hardcoded 10 minutes

Fix other tests by providing appropriate player input (e.g., "5\n").

---

## Separation of Concerns Questions

### Q21: Create Game Type
**File to edit:** `game.go` (new file)  
**Expected test result:** Tests should still pass

Create a new `Game` struct that:
1. Has fields: `alerter BlindAlerter` and `store PlayerStore`
2. Has a `Start(numberOfPlayers int)` method - move scheduling logic here
3. Has a `Finish(winner string)` method - move RecordWin logic here

---

### Q22: Update CLI to Use Game
**File to edit:** `CLI.go`  
**Expected test result:** Tests should still pass

Refactor CLI to:
1. Replace `playerStore` and `alerter` fields with a single `game *Game` field
2. Update `NewCLI` to accept `alerter BlindAlerter` and `store PlayerStore`, create Game internally
3. Update `PlayPoker()` to call `game.Start()` and `game.Finish()` methods

---

### Q23: Create Game Constructor
**File to edit:** `game.go`  
**Expected test result:** Tests should compile more easily

Add a constructor:
```go
func NewGame(alerter BlindAlerter, store PlayerStore) *Game
```

---

### Q24: Update Tests to Use NewGame
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (cleanup)

Update test setup to use `NewGame`:
```go
game := poker.NewGame(blindAlerter, dummyPlayerStore)
cli := poker.NewCLI(in, stdout, game)
```

Also update `main.go` to use the new constructor.

---

### Q25: Create Game Tests File
**File to edit:** `game_test.go` (new file)  
**Expected test result:** New tests should pass

Move game-specific tests from CLI tests:
1. Create `TestGame_Start` with tests for 5 and 7 players
2. Create `TestGame_Finish` to test winner recording
3. These tests should focus on the `Game` type directly

---

### Q26: Introduce Game Interface
**File to edit:** `CLI.go`  
**Expected test result:** Tests should still pass

Create a `Game` interface with:
```go
Start(numberOfPlayers int)
Finish(winner string)
```

Change CLI's `game` field type from `*Game` to `Game` (the interface).

Rename the concrete implementation from `Game` to `TexasHoldem`.

---

### Q27: Create GameSpy
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests will need updating

Create a spy that implements the `Game` interface:
```go
type GameSpy struct {
    StartedWith  int
    FinishedWith string
}
```

Implement `Start` and `Finish` methods that record the arguments.

---

### Q28: Simplify CLI Tests with GameSpy
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass with cleaner assertions

Refactor CLI tests to:
1. Use `GameSpy` instead of real game
2. Remove blind alert checking from CLI tests
3. Focus on: reading player count, calling `Start` with correct number, calling `Finish` with winner
4. Create helper functions like `assertGameStartedWith` and `assertFinishCalledWith`

---

## Error Handling Questions

### Q29: Test Invalid Player Input
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - game starts even with invalid input

Create a test for invalid player input:
1. Test name: "it prints an error when a non numeric value is entered and does not start the game"
2. Provide input "Pies\n"
3. Check that `game.StartCalled` is false (you'll need to add this field to GameSpy)

---

### Q30: Handle Atoi Error
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

In `PlayPoker()`, check the error from `strconv.Atoi`:
```go
numberOfPlayers, err := strconv.Atoi(cli.readLine())
if err != nil {
    return
}
```

---

### Q31: Test Error Message to User
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail - no error message printed

Update the invalid input test to also check stdout:
1. Assert that stdout contains the prompt AND an error message
2. Use a temporary error message like "you're so silly" for now

---

### Q32: Print Error Message
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

Update the error handling to print to stdout:
```go
if err != nil {
    fmt.Fprint(cli.out, "you're so silly")
    return
}
```

---

### Q33: Create Error Constant
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring)

Create a constant for the error message:
```go
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
```

Use it in both the code and test.

---

### Q34: Create assertMessagesSentToUser Helper
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass (refactoring)

Create a helper function:
```go
func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string)
```

This should join all messages and compare against stdout contents. Use it to clean up tests.

---

### Q35: Refactor All CLI Tests
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should pass with improved readability

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

---

## Final Integration Question

### Q36: Verify Complete Application
**Files to test:** All  
**Expected test result:** All tests pass, application runs correctly

1. Run `go test ./...` - all tests should pass
2. Run the CLI application
3. Test valid input (enter a number, then "{Name} wins")
4. Test invalid input (enter non-numeric value)
5. Verify blind alerts print at correct intervals (you may want to reduce timing for testing)
6. Confirm the web server still works alongside the CLI

---

## Additional Challenge Question

### Q37: Handle Invalid Winner Input
**File to edit:** `CLI_test.go` and `CLI.go`  
**Expected test result:** Should gracefully handle bad input

The tutorial mentions: "What happens if instead of putting `Ruth wins` the user puts in `Lloyd is a killer`?"

Write a test for this scenario and implement error handling:
1. Test that invalid winner format is detected
2. Print an appropriate error message
3. Don't record a win for invalid input

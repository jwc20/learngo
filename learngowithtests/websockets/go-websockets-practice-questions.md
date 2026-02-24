# Go WebSockets Practice Questions - WebSockets Integration

## Serving HTML Questions

### Q1: Test GET /game Returns 200
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail with 404

Create a test that verifies `/game` endpoint exists:
1. Create a test function `TestGame`
2. Add a subtest `"GET /game returns 200"`
3. Create a new player server
4. Make a GET request to `/game`
5. Assert the response status is 200

---

### Q2: Add /game Route
**File to edit:** `server.go`  
**Expected test result:** Test should pass

Add the `/game` route to the server:
1. In the router setup, add a handler for `/game`
2. Create a `game` method on PlayerServer that returns 200 status

---

### Q3: Create Test Helpers
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass (refactoring)

Create helper functions to clean up the test:
1. Create `newGameRequest()` that returns a GET request to `/game`
2. Update `assertStatus` to accept `*httptest.ResponseRecorder` instead of just the status code
3. Use these helpers in the test

---

### Q4: Create HTML Template File
**File to edit:** `game.html` (new file)  
**Expected test result:** N/A (no test written per tutorial)

Create `game.html` with:
- A form section with an input for winner name
- A submit button to declare winner
- JavaScript to open a WebSocket connection to `/ws`
- JavaScript to send the winner name when button is clicked

The tutorial provides the full HTML - use it as-is. Note: We're intentionally not testing the HTML itself.

---

### Q5: Serve HTML Template
**File to edit:** `server.go`  
**Expected test result:** Manual testing required

Update the `game` method to serve the HTML template:
1. Use `template.ParseFiles("game.html")` to load the template
2. Handle any errors from parsing
3. Execute the template to the response writer

Test manually by running the server and visiting `/game`.

---

## WebSocket Connection Questions

### Q6: Install WebSocket Library
**File to edit:** N/A (command line)  
**Expected test result:** N/A

Run the command to install the Gorilla WebSocket library:
```bash
go get github.com/gorilla/websocket
```

---

### Q7: Test WebSocket Winner Message
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail with "bad handshake"

Create a test for WebSocket communication:
1. Create subtest `"when we get a message over a websocket it is a winner of a game"`
2. Create a test server using `httptest.NewServer`
3. Convert the HTTP URL to a WebSocket URL (replace "http" with "ws")
4. Use `websocket.DefaultDialer.Dial` to connect to `/ws`
5. Send a winner name using `ws.WriteMessage`
6. Assert that the player store recorded the win

---

### Q8: Add /ws Route
**File to edit:** `server.go`  
**Expected test result:** Test should fail - no winner recorded

Add WebSocket endpoint:
1. Add a handler for `/ws` in the router
2. Create a `webSocket` method on PlayerServer
3. Use `websocket.Upgrader` to upgrade the HTTP connection
4. For now, just upgrade the connection without reading messages

---

### Q9: Read WebSocket Message and Record Win
**File to edit:** `server.go`  
**Expected test result:** Test should hang/timeout

Implement message reading:
1. After upgrading, use `conn.ReadMessage()` to read the message
2. Extract the winner from the message
3. Call `p.store.RecordWin()` with the winner

The test will hang because of timing issues.

---

### Q10: Add Sleep to Fix Timing Issue
**File to edit:** `server_test.go`  
**Expected test result:** Test should pass (but using bad practice)

Add a `time.Sleep(10 * time.Millisecond)` before the assertion to wait for async processing.

Note: This is acknowledged as bad practice but gets us to working code first.

---

## Refactoring Questions

### Q11: Extract WebSocket Upgrader
**File to edit:** `server.go`  
**Expected test result:** Tests should pass (refactoring)

Extract the upgrader to a package-level variable:
```go
var wsUpgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
```

Use this in the `webSocket` method.

---

### Q12: Parse Template Once in Constructor
**File to edit:** `server.go`  
**Expected test result:** Compilation errors in tests

Optimize template parsing:
1. Add a `template *template.Template` field to PlayerServer
2. Parse the template in `NewPlayerServer` constructor
3. Store the parsed template in the struct
4. Update `NewPlayerServer` to return `(*PlayerServer, error)`
5. Update the `game` method to use the stored template

---

### Q13: Fix Compilation Errors in Tests
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass

Create helper function to handle the new constructor signature:
```go
func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer
```

This helper should call `NewPlayerServer` and fail the test if there's an error.

Update all tests to use this helper.

---

### Q14: Fix Main.go
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (integration)

Update main to handle the new constructor signature:
- Call `NewPlayerServer` and handle the error
- Use symlink or copy `game.html` to the cmd/webserver directory if needed

---

### Q15: Create WebSocket Test Helpers
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass (refactoring)

Create helpers to clean up WebSocket test code:
1. `mustDialWS(t *testing.T, url string) *websocket.Conn` - opens WS connection
2. `writeWSMessage(t testing.TB, conn *websocket.Conn, message string)` - sends message

Use these helpers in your WebSocket test.

---

## Integrating Game Questions

### Q16: Update HTML for Number of Players
**File to edit:** `game.html`  
**Expected test result:** N/A (no tests for HTML)

Update the HTML to:
1. Add an input for number of players
2. Add a "Start" button
3. Add a div to display blind values
4. Update JavaScript to:
   - Send number of players when game starts
   - Listen for blind alert messages from server
   - Update the blind display when messages arrive

The tutorial provides the complete updated HTML.

---

### Q17: Update BlindAlerter Interface
**File to edit:** `blind_alerter.go`  
**Expected test result:** Compilation errors

Change the `BlindAlerter` interface:
1. Add `io.Writer` parameter to `ScheduleAlertAt`
2. Update `BlindAlerterFunc` type signature
3. Update `Alerter` function to accept the `io.Writer` parameter
4. Rename `StdOutAlerter` to `Alerter`

---

### Q18: Update TexasHoldem to Use io.Writer
**File to edit:** `game.go`  
**Expected test result:** Compilation errors in tests

Temporarily fix compilation:
1. Update calls to `ScheduleAlertAt` to pass `os.Stdout` as the destination

---

### Q19: Fix SpyBlindAlerter
**File to edit:** Test files using SpyBlindAlerter  
**Expected test result:** Tests should pass

Update `SpyBlindAlerter.ScheduleAlertAt` signature to match the new interface.

---

### Q20: Update Game Interface
**File to edit:** `game.go`  
**Expected test result:** Compilation errors

Update the `Game` interface:
1. Change `Start` to accept `(numberOfPlayers int, alertsDestination io.Writer)`

---

### Q21: Fix TexasHoldem Implementation
**File to edit:** `game.go`  
**Expected test result:** Compilation errors in other files

Update `TexasHoldem.Start`:
1. Update the method signature to match the interface
2. Use the `alertsDestination` parameter instead of hardcoded `os.Stdout`

---

### Q22: Fix CLI Usage
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass

Update CLI to pass the output destination:
```go
cli.game.Start(numberOfPlayers, cli.out)
```

---

### Q23: Fix Game Tests
**File to edit:** `game_test.go`  
**Expected test result:** Tests should pass

Update game tests to pass `io.Discard` as the alerts destination:
```go
game.Start(5, io.Discard)
```

---

## WebSocket Game Integration Questions

### Q24: Test Game Integration with WebSocket
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail - compilation error

Create a test using GameSpy:
1. Create subtest `"start a game with 3 players and declare Ruth the winner"`
2. Create a `GameSpy`
3. Pass it to `mustMakePlayerServer` (will need to update signature)
4. Open WebSocket connection
5. Send "3" (number of players)
6. Send "Ruth" (winner)
7. Assert `Start` was called with 3
8. Assert `Finish` was called with "Ruth"

---

### Q25: Add Dummy Game for Other Tests
**File to edit:** `server_test.go`  
**Expected test result:** Tests should compile

Create a dummy game variable:
```go
var dummyGame = &GameSpy{}
```

Update all other tests to pass this dummy.

---

### Q26: Update NewPlayerServer Signature
**File to edit:** `server.go`  
**Expected test result:** Test should fail - game not used

Add `Game` parameter to `NewPlayerServer`:
1. Update the function signature
2. Add a `game Game` field to PlayerServer
3. Store the game in the constructor

---

### Q27: Implement Game Integration in WebSocket Handler
**File to edit:** `server.go`  
**Expected test result:** Tests should pass

Update the `webSocket` method:
1. Read the number of players message
2. Convert it to an integer
3. Call `p.game.Start(numberOfPlayers, io.Discard)`
4. Read the winner message
5. Call `p.game.Finish(winner)`

Note: We're discarding blind alerts for now.

---

### Q28: Update Main.go for Game
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (integration)

Update main to create and pass a Game to the server:
```go
game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)
server, err := poker.NewPlayerServer(store, game)
```

---

## WebSocket Wrapper Questions

### Q29: Create playerServerWS Wrapper
**File to edit:** `server.go`  
**Expected test result:** Tests should pass (refactoring)

Create a wrapper type for WebSocket connection:
1. Define `type playerServerWS struct` embedding `*websocket.Conn`
2. Create `newPlayerServerWS(w, r)` constructor that upgrades the connection
3. Add `WaitForMsg()` method that reads and returns a message as string
4. Refactor `webSocket` method to use this wrapper

---

### Q30: Make playerServerWS Implement io.Writer
**File to edit:** `server.go`  
**Expected test result:** Manual testing shows it works

Add a `Write` method to `playerServerWS`:
```go
func (w *playerServerWS) Write(p []byte) (n int, err error)
```

Implement it to:
1. Use `w.WriteMessage(websocket.TextMessage, p)` to send data
2. Return appropriate length and error

---

### Q31: Send Alerts Through WebSocket
**File to edit:** `server.go`  
**Expected test result:** Manual testing required

Change the game start call:
```go
p.game.Start(numberOfPlayers, ws)  // Instead of io.Discard
```

Test manually - blind alerts should now appear in the browser!

---

### Q32: Reduce Blind Increment for Testing
**File to edit:** `game.go`  
**Expected test result:** N/A (temporary change)

Temporarily change the blind increment from minutes to seconds for easier testing:
```go
blindIncrement := time.Duration(5+numberOfPlayers) * time.Second
```

---

## Testing Blind Alerts Questions

### Q33: Add BlindAlert to GameSpy
**File to edit:** Test files with GameSpy  
**Expected test result:** Tests will need updating

Update `GameSpy`:
1. Add `BlindAlert []byte` field
2. Update `Start` to write `BlindAlert` to the `out` writer

---

### Q34: Test WebSocket Receives Blind Alert
**File to edit:** `server_test.go`  
**Expected test result:** Test should hang

Update the WebSocket game test:
1. Set `wantedBlindAlert := "Blind is 100"`
2. Configure GameSpy with this alert: `game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}`
3. After writing messages, read from WebSocket: `ws.ReadMessage()`
4. Assert the message matches `wantedBlindAlert`

Test will hang because `ReadMessage` blocks.

---

### Q35: Create within Helper for Timeouts
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail but not hang

Create a helper to handle timeouts:
```go
func within(t testing.TB, d time.Duration, assert func())
```

This should:
1. Run the assertion in a goroutine
2. Use a channel to signal completion
3. Use `select` with `time.After` to timeout if it takes too long

---

### Q36: Create assertWebsocketGotMsg Helper
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass

Create a helper for WebSocket message assertions:
```go
func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string)
```

---

### Q37: Use within Helper in Test
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail with timeout initially

Update the test to use the `within` helper:
```go
within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
```

---

### Q38: Fix Implementation to Send Alerts
**File to edit:** `server.go`  
**Expected test result:** Tests should pass

Ensure the webSocket method passes `ws` to `game.Start` instead of `io.Discard`.

This should make the test pass!

---

## Removing Sleep from Tests Questions

### Q39: Create retryUntil Helper
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass without sleep

Create a retry helper:
```go
func retryUntil(d time.Duration, f func() bool) bool
```

This should:
1. Set a deadline
2. Keep calling `f()` until it returns true or deadline passes
3. Return whether it succeeded

---

### Q40: Update assertFinishCalledWith
**File to edit:** `server_test.go`  
**Expected test result:** Can remove time.Sleep

Refactor the helper to use `retryUntil`:
```go
func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
    passed := retryUntil(500*time.Millisecond, func() bool {
        return game.FinishCalledWith == winner
    })
    
    if !passed {
        t.Errorf(...)
    }
}
```

---

### Q41: Update assertGameStartedWith
**File to edit:** `server_test.go`  
**Expected test result:** Can remove time.Sleep

Apply the same refactoring to `assertGameStartedWith` using `retryUntil`.

---

### Q42: Remove time.Sleep from Tests
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass without sleeps

Remove all `time.Sleep` calls from the WebSocket tests. The retry logic in the helpers should handle timing now.

---

## Final Integration Question

### Q43: Complete Application Test
**Files to test:** All  
**Expected test result:** Full application works

Test the complete application:
1. Run `go test ./...` - all tests should pass
2. Start the web server
3. Visit `/game` in a browser
4. Enter number of players (e.g., 5)
5. Click "Start" - should see blind alerts updating automatically
6. Enter a winner name and click "Declare winner"
7. Visit `/league` - winner should be recorded
8. Verify the CLI app still works alongside the web app

Remember to revert the blind increment back to minutes for production if you changed it to seconds for testing!

---

## Challenge Questions

### Q44: Better Error Handling
**File to edit:** `server.go`  
**Expected test result:** Should gracefully handle errors

The current implementation ignores many errors. Add proper error handling:
1. Handle errors from `ReadMessage()`
2. Handle errors from `Atoi()` conversion
3. Log appropriate messages
4. Consider how to inform the user of errors through WebSocket

---

### Q45: Test the HTML/JavaScript
**File to edit:** New test file (JavaScript)  
**Expected test result:** JavaScript tests pass

The tutorial skips JavaScript testing. As a challenge:
1. Set up a JavaScript testing framework
2. Test the WebSocket connection logic
3. Test the UI updates when messages are received
4. Test form submission behavior

This is beyond the scope of the tutorial but good practice!

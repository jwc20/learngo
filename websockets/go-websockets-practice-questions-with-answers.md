# Go WebSockets Practice Questions with Answers - WebSockets Integration

## Serving HTML Questions

### Q1: Test GET /game Returns 200
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail with 404

**Question:**
Create a test that verifies `/game` endpoint exists:
1. Create a test function `TestGame`
2. Add a subtest `"GET /game returns 200"`
3. Create a new player server
4. Make a GET request to `/game`
5. Assert the response status is 200

**Answer:**
```go
// server_test.go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

Run the test:
```bash
go test
```

Expected failure: "did not get correct status, got 404, want 200"

---

### Q2: Add /game Route
**File to edit:** `server.go`  
**Expected test result:** Test should pass

**Question:**
Add the `/game` route to the server:
1. In the router setup, add a handler for `/game`
2. Create a `game` method on PlayerServer that returns 200 status

**Answer:**
```go
// server.go
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.game))  // Add this

	p.Handler = router

	return p
}

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```

---

### Q3: Create Test Helpers
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Create helper functions to clean up the test:
1. Create `newGameRequest()` that returns a GET request to `/game`
2. Update `assertStatus` to accept `*httptest.ResponseRecorder` instead of just the status code
3. Use these helpers in the test

**Answer:**
```go
// server_test.go
func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Errorf("did not get correct status, got %d, want %d", response.Code, want)
	}
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})
}
```

---

### Q4: Create HTML Template File
**File to edit:** `game.html` (new file)  
**Expected test result:** N/A (no test written per tutorial)

**Question:**
Create `game.html` with:
- A form section with an input for winner name
- A submit button to declare winner
- JavaScript to open a WebSocket connection to `/ws`
- JavaScript to send the winner name when button is clicked

**Answer:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Let's play poker</title>
</head>
<body>
<section id="game">
    <div id="declare-winner">
        <label for="winner">Winner</label>
        <input type="text" id="winner"/>
        <button id="winner-button">Declare winner</button>
    </div>
</section>
</body>
<script type="application/javascript">

    const submitWinnerButton = document.getElementById('winner-button')
    const winnerInput = document.getElementById('winner')

    if (window['WebSocket']) {
        const conn = new WebSocket('ws://' + document.location.host + '/ws')

        submitWinnerButton.onclick = event => {
            conn.send(winnerInput.value)
        }
    }
</script>
</html>
```

---

### Q5: Serve HTML Template
**File to edit:** `server.go`  
**Expected test result:** Manual testing required

**Question:**
Update the `game` method to serve the HTML template:
1. Use `template.ParseFiles("game.html")` to load the template
2. Handle any errors from parsing
3. Execute the template to the response writer

**Answer:**
```go
// server.go
import (
	"html/template"
	// ... other imports
)

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("game.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
```

Test manually:
```bash
go run cmd/webserver/main.go
# Visit http://localhost:5000/game
```

You may need to create a symlink or copy the file to the correct location.

---

## WebSocket Connection Questions

### Q6: Install WebSocket Library
**File to edit:** N/A (command line)  
**Expected test result:** N/A

**Question:**
Run the command to install the Gorilla WebSocket library

**Answer:**
```bash
go get github.com/gorilla/websocket
```

---

### Q7: Test WebSocket Winner Message
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail with "bad handshake"

**Question:**
Create a test for WebSocket communication:
1. Create subtest `"when we get a message over a websocket it is a winner of a game"`
2. Create a test server using `httptest.NewServer`
3. Convert the HTTP URL to a WebSocket URL (replace "http" with "ws")
4. Use `websocket.DefaultDialer.Dial` to connect to `/ws`
5. Send a winner name using `ws.WriteMessage`
6. Assert that the player store recorded the win

**Answer:**
```go
// server_test.go
import (
	"github.com/gorilla/websocket"
	// ... other imports
)

t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
	store := &StubPlayerStore{}
	winner := "Ruth"
	server := httptest.NewServer(NewPlayerServer(store))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}

	AssertPlayerWin(t, store, winner)
})
```

Expected error: "could not open a ws connection on ws://127.0.0.1:xxxxx/ws websocket: bad handshake"

---

### Q8: Add /ws Route
**File to edit:** `server.go`  
**Expected test result:** Test should fail - no winner recorded

**Question:**
Add WebSocket endpoint:
1. Add a handler for `/ws` in the router
2. Create a `webSocket` method on PlayerServer
3. Use `websocket.Upgrader` to upgrade the HTTP connection

**Answer:**
```go
// server.go
import (
	"github.com/gorilla/websocket"
	// ... other imports
)

func NewPlayerServer(store PlayerStore) *PlayerServer {
	// ... existing code ...
	router.Handle("/ws", http.HandlerFunc(p.webSocket))
	// ...
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.Upgrade(w, r, nil)
}
```

---

### Q9: Read WebSocket Message and Record Win
**File to edit:** `server.go`  
**Expected test result:** Test should hang/timeout

**Question:**
Implement message reading:
1. After upgrading, use `conn.ReadMessage()` to read the message
2. Extract the winner from the message
3. Call `p.store.RecordWin()` with the winner

**Answer:**
```go
// server.go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _ := upgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.store.RecordWin(string(winnerMsg))
}
```

---

### Q10: Add Sleep to Fix Timing Issue
**File to edit:** `server_test.go`  
**Expected test result:** Test should pass (but using bad practice)

**Question:**
Add a `time.Sleep(10 * time.Millisecond)` before the assertion to wait for async processing.

**Answer:**
```go
// server_test.go
t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
	store := &StubPlayerStore{}
	winner := "Ruth"
	server := httptest.NewServer(NewPlayerServer(store))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	AssertPlayerWin(t, store, winner)
})
```

---

## Refactoring Questions

### Q11: Extract WebSocket Upgrader
**File to edit:** `server.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Extract the upgrader to a package-level variable

**Answer:**
```go
// server.go
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.store.RecordWin(string(winnerMsg))
}
```

---

### Q12: Parse Template Once in Constructor
**File to edit:** `server.go`  
**Expected test result:** Compilation errors in tests

**Question:**
Optimize template parsing:
1. Add a `template *template.Template` field to PlayerServer
2. Parse the template in `NewPlayerServer` constructor
3. Update `NewPlayerServer` to return `(*PlayerServer, error)`
4. Update the `game` method to use the stored template

**Answer:**
```go
// server.go
type PlayerServer struct {
	store    PlayerStore
	http.Handler
	template *template.Template
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.game))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}
```

---

### Q13: Fix Compilation Errors in Tests
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass

**Question:**
Create helper function to handle the new constructor signature

**Answer:**
```go
// server_test.go
func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}
```

Update all tests to use this helper:
```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{})
		
		// ... rest of test
	})
	
	t.Run("when we get a message over a websocket...", func(t *testing.T) {
		store := &StubPlayerStore{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, store))
		
		// ... rest of test
	})
}
```

Do the same for other test files.

---

### Q14: Fix Main.go
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (integration)

**Question:**
Update main to handle the new constructor signature

**Answer:**
```go
// cmd/webserver/main.go
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server, err := poker.NewPlayerServer(store)

	if err != nil {
		log.Fatalf("problem creating player server %v", err)
	}

	log.Fatal(http.ListenAndServe(":5000", server))
}
```

If needed, create a symlink to game.html:
```bash
cd cmd/webserver
ln -s ../../game.html game.html
```

---

### Q15: Create WebSocket Test Helpers
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Create helpers to clean up WebSocket test code

**Answer:**
```go
// server_test.go
func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
```

Update the test to use these:
```go
t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
	store := &StubPlayerStore{}
	winner := "Ruth"
	server := httptest.NewServer(mustMakePlayerServer(t, store))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws := mustDialWS(t, wsURL)
	defer ws.Close()

	writeWSMessage(t, ws, winner)

	time.Sleep(10 * time.Millisecond)
	AssertPlayerWin(t, store, winner)
})
```

---

## Integrating Game Questions

### Q16: Update HTML for Number of Players
**File to edit:** `game.html`  
**Expected test result:** N/A (no tests for HTML)

**Question:**
Update the HTML to handle number of players and blind alerts

**Answer:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Lets play poker</title>
</head>
<body>
<section id="game">
    <div id="game-start">
        <label for="player-count">Number of players</label>
        <input type="number" id="player-count"/>
        <button id="start-game">Start</button>
    </div>

    <div id="declare-winner">
        <label for="winner">Winner</label>
        <input type="text" id="winner"/>
        <button id="winner-button">Declare winner</button>
    </div>

    <div id="blind-value"/>
</section>

<section id="game-end">
    <h1>Another great game of poker everyone!</h1>
    <p><a href="/league">Go check the league table</a></p>
</section>

</body>
<script type="application/javascript">
    const startGame = document.getElementById('game-start')

    const declareWinner = document.getElementById('declare-winner')
    const submitWinnerButton = document.getElementById('winner-button')
    const winnerInput = document.getElementById('winner')

    const blindContainer = document.getElementById('blind-value')

    const gameContainer = document.getElementById('game')
    const gameEndContainer = document.getElementById('game-end')

    declareWinner.hidden = true
    gameEndContainer.hidden = true

    document.getElementById('start-game').addEventListener('click', event => {
        startGame.hidden = true
        declareWinner.hidden = false

        const numberOfPlayers = document.getElementById('player-count').value

        if (window['WebSocket']) {
            const conn = new WebSocket('ws://' + document.location.host + '/ws')

            submitWinnerButton.onclick = event => {
                conn.send(winnerInput.value)
                gameEndContainer.hidden = false
                gameContainer.hidden = true
            }

            conn.onclose = evt => {
                blindContainer.innerText = 'Connection closed'
            }

            conn.onmessage = evt => {
                blindContainer.innerText = evt.data
            }

            conn.onopen = function () {
                conn.send(numberOfPlayers)
            }
        }
    })
</script>
</html>
```

---

### Q17: Update BlindAlerter Interface
**File to edit:** `blind_alerter.go`  
**Expected test result:** Compilation errors

**Question:**
Add `io.Writer` parameter to `ScheduleAlertAt`

**Answer:**
```go
// blind_alerter.go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}
```

---

### Q18: Update TexasHoldem to Use io.Writer
**File to edit:** `game.go`  
**Expected test result:** Compilation errors in tests

**Question:**
Update `Alerter` function and temporarily fix compilation by passing `os.Stdout`

**Answer:**
```go
// blind_alerter.go
func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
```

```go
// game.go
func (g *TexasHoldem) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind, os.Stdout)  // Temporarily hardcoded
		blindTime = blindTime + blindIncrement
	}
}
```

---

### Q19: Fix SpyBlindAlerter
**File to edit:** Test files  
**Expected test result:** Tests should pass

**Question:**
Update `SpyBlindAlerter.ScheduleAlertAt` signature

**Answer:**
```go
// In test files (CLI_test.go, game_test.go, etc.)
type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}
```

---

### Q20: Update Game Interface
**File to edit:** `game.go` or `CLI.go` (wherever Game interface is defined)  
**Expected test result:** Compilation errors

**Question:**
Update the `Game` interface `Start` method

**Answer:**
```go
// CLI.go or game.go
type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
```

---

### Q21: Fix TexasHoldem Implementation
**File to edit:** `game.go`  
**Expected test result:** Compilation errors in other files

**Question:**
Update `TexasHoldem.Start` to match the interface

**Answer:**
```go
// game.go
func (g *TexasHoldem) Start(numberOfPlayers int, alertsDestination io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind, alertsDestination)
		blindTime = blindTime + blindIncrement
	}
}
```

---

### Q22: Fix CLI Usage
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass

**Question:**
Update CLI to pass the output destination

**Answer:**
```go
// CLI.go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(numberOfPlayers, cli.out)  // Pass cli.out

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}
```

---

### Q23: Fix Game Tests
**File to edit:** `game_test.go`  
**Expected test result:** Tests should pass

**Question:**
Update game tests to pass `io.Discard`

**Answer:**
```go
// game_test.go
func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5, io.Discard)

		// ... assertions
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7, io.Discard)

		// ... assertions
	})
}
```

---

## WebSocket Game Integration Questions

### Q24: Test Game Integration with WebSocket
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail - compilation error

**Question:**
Create a test using GameSpy

**Answer:**
```go
// server_test.go
t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
	game := &GameSpy{}
	winner := "Ruth"
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	time.Sleep(10 * time.Millisecond)
	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)
})
```

Note: You'll need GameSpy from CLI tests if not already in this file.

---

### Q25: Add Dummy Game for Other Tests
**File to edit:** `server_test.go`  
**Expected test result:** Tests should compile

**Question:**
Create a dummy game variable

**Answer:**
```go
// server_test.go
var (
	dummyGame = &GameSpy{}
)
```

Update existing tests:
```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, dummyGame)
		// ... rest
	})
}
```

---

### Q26: Update NewPlayerServer Signature
**File to edit:** `server.go`  
**Expected test result:** Test should fail - game not used

**Question:**
Add `Game` parameter to `NewPlayerServer`

**Answer:**
```go
// server.go
type PlayerServer struct {
	store    PlayerStore
	http.Handler
	template *template.Template
	game     Game
}

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.game = game
	p.template = tmpl
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.playGame))  // Renamed from 'game'
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

func (p *PlayerServer) playGame(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}
```

Update mustMakePlayerServer:
```go
func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}
```

---

### Q27: Implement Game Integration in WebSocket Handler
**File to edit:** `server.go`  
**Expected test result:** Tests should pass

**Question:**
Update the `webSocket` method to use Game

**Answer:**
```go
// server.go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, io.Discard) // TODO: Don't discard!

	_, winner, _ := conn.ReadMessage()
	p.game.Finish(string(winner))
}
```

---

### Q28: Update Main.go for Game
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (integration)

**Question:**
Update main to create and pass a Game

**Answer:**
```go
// cmd/webserver/main.go
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)

	server, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatalf("problem creating player server %v", err)
	}

	log.Fatal(http.ListenAndServe(":5000", server))
}
```

---

## WebSocket Wrapper Questions

### Q29: Create playerServerWS Wrapper
**File to edit:** `server.go`  
**Expected test result:** Tests should pass (refactoring)

**Question:**
Create a wrapper type for WebSocket connection

**Answer:**
```go
// server.go
type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, io.Discard)

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

---

### Q30: Make playerServerWS Implement io.Writer
**File to edit:** `server.go`  
**Expected test result:** Manual testing shows it works

**Question:**
Add a `Write` method to `playerServerWS`

**Answer:**
```go
// server.go
func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
```

---

### Q31: Send Alerts Through WebSocket
**File to edit:** `server.go`  
**Expected test result:** Manual testing required

**Question:**
Change the game start call to send to WebSocket

**Answer:**
```go
// server.go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)  // Send to ws instead of io.Discard

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

Test manually:
```bash
go run cmd/webserver/main.go
# Visit http://localhost:5000/game
# Enter number of players and start
# You should see blind alerts appearing!
```

---

### Q32: Reduce Blind Increment for Testing
**File to edit:** `game.go`  
**Expected test result:** N/A (temporary change)

**Question:**
Change the blind increment to seconds for easier testing

**Answer:**
```go
// game.go
func (g *TexasHoldem) Start(numberOfPlayers int, alertsDestination io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Second  // Changed from Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind, alertsDestination)
		blindTime = blindTime + blindIncrement
	}
}
```

Remember to change this back to `time.Minute` for production!

---

## Testing Blind Alerts Questions

### Q33: Add BlindAlert to GameSpy
**File to edit:** Test files with GameSpy  
**Expected test result:** Tests will need updating

**Question:**
Update `GameSpy` to send a blind alert

**Answer:**
```go
// In test files
type GameSpy struct {
	StartCalled      bool
	StartCalledWith  int
	BlindAlert       []byte

	FinishedCalled   bool
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
	out.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedCalled = true
	g.FinishCalledWith = winner
}
```

---

### Q34: Test WebSocket Receives Blind Alert
**File to edit:** `server_test.go`  
**Expected test result:** Test should hang

**Question:**
Update the test to check for blind alerts

**Answer:**
```go
// server_test.go
t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
	wantedBlindAlert := "Blind is 100"
	winner := "Ruth"

	game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	time.Sleep(10 * time.Millisecond)
	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)

	_, gotBlindAlert, _ := ws.ReadMessage()

	if string(gotBlindAlert) != wantedBlindAlert {
		t.Errorf("got blind alert %q, want %q", string(gotBlindAlert), wantedBlindAlert)
	}
})
```

The test hangs because `ReadMessage` blocks forever.

---

### Q35: Create within Helper for Timeouts
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail but not hang

**Question:**
Create a helper to handle timeouts

**Answer:**
```go
// server_test.go
func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
```

---

### Q36: Create assertWebsocketGotMsg Helper
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass

**Question:**
Create a helper for WebSocket message assertions

**Answer:**
```go
// server_test.go
func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}
```

---

### Q37: Use within Helper in Test
**File to edit:** `server_test.go`  
**Expected test result:** Test should fail with timeout

**Question:**
Update the test to use the `within` helper

**Answer:**
```go
// server_test.go
const tenMS = 10 * time.Millisecond

t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
	wantedBlindAlert := "Blind is 100"
	winner := "Ruth"

	game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	time.Sleep(tenMS)

	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)
	within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
})
```

---

### Q38: Fix Implementation to Send Alerts
**File to edit:** `server.go`  
**Expected test result:** Tests should pass

**Question:**
Ensure the webSocket method sends alerts through WebSocket

**Answer:**
The code should already be correct from Q31:
```go
// server.go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)  // This should already be 'ws', not io.Discard

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

If you still have `io.Discard`, change it to `ws`.

---

## Removing Sleep from Tests Questions

### Q39: Create retryUntil Helper
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass without sleep

**Question:**
Create a retry helper

**Answer:**
```go
// server_test.go
func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
```

---

### Q40: Update assertFinishCalledWith
**File to edit:** `server_test.go`  
**Expected test result:** Can remove time.Sleep

**Question:**
Refactor the helper to use `retryUntil`

**Answer:**
```go
// server_test.go
func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishCalledWith == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}
```

---

### Q41: Update assertGameStartedWith
**File to edit:** `server_test.go`  
**Expected test result:** Can remove time.Sleep

**Question:**
Apply the same refactoring

**Answer:**
```go
// server_test.go
func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayers int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartCalledWith == numberOfPlayers
	})

	if !passed {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayers, game.StartCalledWith)
	}
}
```

---

### Q42: Remove time.Sleep from Tests
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass without sleeps

**Question:**
Remove all `time.Sleep` calls from tests

**Answer:**
```go
// server_test.go
t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
	wantedBlindAlert := "Blind is 100"
	winner := "Ruth"

	game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	// Removed: time.Sleep(tenMS)

	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)
	within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
})
```

---

## Final Integration Question

### Q43: Complete Application Test
**Files to test:** All  
**Expected test result:** Full application works

**Question:**
Test the complete application

**Answer:**
```bash
# 1. Run all tests
go test ./...
# Expected: PASS

# 2. Start the web server
go run cmd/webserver/main.go

# 3. In a browser, visit http://localhost:5000/game

# 4. Test the workflow:
# - Enter number of players: 5
# - Click "Start"
# - Observe blind alerts updating automatically in the browser
# - Enter winner name: "Alice"
# - Click "Declare winner"
# - Should see "Another great game" message

# 5. Check league
# Visit http://localhost:5000/league
# Alice should be in the league with 1 win

# 6. Test CLI still works
# In another terminal:
go run cmd/cli/main.go
# Enter: 3
# Wait for blind alerts
# Enter: Bob wins
# Check league - Bob should be added

# Remember to revert blind increment back to minutes!
```

In `game.go`, change back:
```go
blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute  // Back to Minute
```

---

## Challenge Questions

### Q44: Better Error Handling
**File to edit:** `server.go`  
**Expected test result:** Should gracefully handle errors

**Question:**
Add proper error handling

**Answer:**
```go
// server.go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)
	if ws == nil {
		return
	}

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, err := strconv.Atoi(numberOfPlayersMsg)
	if err != nil {
		log.Printf("invalid number of players: %v", err)
		ws.Write([]byte("Invalid number of players"))
		return
	}

	p.game.Start(numberOfPlayers, ws)

	winner := ws.WaitForMsg()
	if winner == "" {
		log.Printf("no winner provided")
		return
	}

	p.game.Finish(winner)
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
		return nil
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
		return ""
	}
	return string(msg)
}
```

---

### Q45: Test the HTML/JavaScript
**File to edit:** New test file (JavaScript)  
**Expected test result:** JavaScript tests pass

**Question:**
Set up JavaScript testing (challenge - not covered in tutorial)

**Answer:**
This is beyond the scope of the Go testing tutorial, but here's a starting point:

1. Install a JavaScript testing framework (Jest, Mocha, etc.)
2. Set up a test environment that can run browser JavaScript
3. Mock the WebSocket API
4. Test the connection logic
5. Test UI updates

Example Jest test structure:
```javascript
// game.test.js
describe('Poker Game UI', () => {
  let mockWebSocket;
  
  beforeEach(() => {
    // Set up DOM
    document.body.innerHTML = `
      <input id="player-count" />
      <button id="start-game">Start</button>
      <div id="blind-value"></div>
    `;
    
    // Mock WebSocket
    mockWebSocket = {
      send: jest.fn(),
      onmessage: null,
      onopen: null,
    };
    global.WebSocket = jest.fn(() => mockWebSocket);
  });
  
  test('should send player count on start', () => {
    document.getElementById('player-count').value = '5';
    document.getElementById('start-game').click();
    
    mockWebSocket.onopen();
    
    expect(mockWebSocket.send).toHaveBeenCalledWith('5');
  });
  
  test('should update blind value on message', () => {
    const blindContainer = document.getElementById('blind-value');
    
    mockWebSocket.onmessage({ data: 'Blind is 100' });
    
    expect(blindContainer.innerText).toBe('Blind is 100');
  });
});
```

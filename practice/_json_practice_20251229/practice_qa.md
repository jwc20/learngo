# JSON, Routing and Embedding - TDD Practice Q/A

## Starting Point
You should have the following files from the previous chapter:
- `server.go`
- `server_test.go`
- `in_memory_player_store.go`
- `server_integration_test.go`
- `main.go`

---

## Exercise 1: Write Test for /league Endpoint
**File to edit:** `server_test.go`
**Expected result:** Test fails with `got 404, want 200`

**Question:** Add a new test function `TestLeague` that makes a GET request to `/league` and asserts it returns status 200.

**Answer:**
```go
func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := &PlayerServer{&store}

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

---

## Exercise 2: Add Router to ServeHTTP (Quick and Dirty)
**File to edit:** `server.go`
**Expected result:** All tests pass

**Question:** Update `ServeHTTP` to use `http.NewServeMux()` for routing. Handle `/league` to return `StatusOK` and `/players/` to use the existing player logic.

**Answer:**
```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := strings.TrimPrefix(r.URL.Path, "/players/")

		switch r.Method {
		case http.MethodPost:
			p.processWin(w, player)
		case http.MethodGet:
			p.showScore(w, player)
		}
	}))

	router.ServeHTTP(w, r)
}
```

---

## Exercise 3: Refactor - Extract Handler Methods
**File to edit:** `server.go`
**Expected result:** All tests pass (refactor)

**Question:** Extract the inline handler functions into separate methods `leagueHandler` and `playersHandler` on `PlayerServer`.

**Answer:**
```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	router.ServeHTTP(w, r)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}
```

---

## Exercise 4: Refactor - Create NewPlayerServer Constructor
**File to edit:** `server.go`
**Expected result:** Compilation error - tests need updating

**Question:** Add a `router` field to `PlayerServer` struct and create a `NewPlayerServer` constructor that sets up routing once. Simplify `ServeHTTP` to just delegate to the router.

**Answer:**
```go
type PlayerServer struct {
	store  PlayerStore
	router *http.ServeMux
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := &PlayerServer{
		store,
		http.NewServeMux(),
	}

	p.router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	p.router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	return p
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
```

---

## Exercise 5: Update Tests to Use NewPlayerServer
**Files to edit:** `server_test.go`, `server_integration_test.go`, `main.go`
**Expected result:** All tests pass

**Question:** Replace all instances of `&PlayerServer{&store}` or `PlayerServer{&store}` with `NewPlayerServer(&store)`.

**Answer:**

In `server_test.go` (multiple places):
```go
server := NewPlayerServer(&store)
```

In `server_integration_test.go`:
```go
server := NewPlayerServer(store)
```

In `main.go`:
```go
func main() {
	server := NewPlayerServer(NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

---

## Exercise 6: Refactor - Use Embedding Instead of Named Router Field
**File to edit:** `server.go`
**Expected result:** All tests pass (refactor)

**Question:** Replace the named `router *http.ServeMux` field with an embedded `http.Handler`. Update `NewPlayerServer` to assign the router to `p.Handler`. Delete the `ServeHTTP` method (it's now provided by embedding).

**Answer:**
```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}
```

(Delete the `ServeHTTP` method entirely)

---

## Exercise 7: Define Player Type
**File to edit:** `server.go`
**Expected result:** Code compiles (no test change yet)

**Question:** Create a `Player` struct with `Name` (string) and `Wins` (int) fields to represent the JSON data model.

**Answer:**
```go
type Player struct {
	Name string
	Wins int
}
```

---

## Exercise 8: Update Test to Parse JSON Response
**File to edit:** `server_test.go`
**Expected result:** Test fails with `Unable to parse response from server '' into slice of Player, 'unexpected end of JSON input'`

**Question:** Update the `TestLeague` test to decode the response body as JSON into a `[]Player` slice. Import `encoding/json`.

**Answer:**
```go
func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

---

## Exercise 9: Return Hard-coded JSON from leagueHandler
**File to edit:** `server.go`
**Expected result:** All tests pass

**Question:** Update `leagueHandler` to encode and return a hard-coded slice containing one `Player{"Chris", 20}` as JSON.

**Answer:**
```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	leagueTable := []Player{
		{"Chris", 20},
	}

	json.NewEncoder(w).Encode(leagueTable)

	w.WriteHeader(http.StatusOK)
}
```

---

## Exercise 10: Refactor - Extract getLeagueTable Method
**File to edit:** `server.go`
**Expected result:** All tests pass (refactor)

**Question:** Extract the hard-coded league data into a separate `getLeagueTable` method.

**Answer:**
```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.getLeagueTable())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) getLeagueTable() []Player {
	return []Player{
		{"Chris", 20},
	}
}
```

---

## Exercise 11: Add league Field to StubPlayerStore
**File to edit:** `server_test.go`
**Expected result:** Compilation error - struct initializers have too few values

**Question:** Add a `league []Player` field to `StubPlayerStore`.

**Answer:**
```go
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}
```

---

## Exercise 12: Update Test to Assert League Data
**File to edit:** `server_test.go`
**Expected result:** Test fails with `got [{Chris 20}] want [{Cleo 32} {Chris 20} {Tiest 14}]`

**Question:** Update `TestLeague` to create a `wantedLeague` with 3 players, pass it to the stub, and assert the response matches using `reflect.DeepEqual`. Update all other `StubPlayerStore` initializers to include `nil` for the league field.

**Answer:**
```go
func TestLeague(t *testing.T) {

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(got, wantedLeague) {
			t.Errorf("got %v want %v", got, wantedLeague)
		}
	})
}
```

Also update other tests:
```go
store := StubPlayerStore{
	map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	},
	nil,
	nil,
}
```

---

## Exercise 13: Add GetLeague to PlayerStore Interface
**File to edit:** `server.go`
**Expected result:** Compilation error - missing GetLeague method on implementations

**Question:** Add `GetLeague() []Player` to the `PlayerStore` interface.

**Answer:**
```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}
```

---

## Exercise 14: Implement GetLeague on StubPlayerStore
**File to edit:** `server_test.go`
**Expected result:** Still compilation error (InMemoryPlayerStore missing method)

**Question:** Add `GetLeague` method to `StubPlayerStore` that returns the `league` field.

**Answer:**
```go
func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}
```

---

## Exercise 15: Implement GetLeague on InMemoryPlayerStore (Minimal)
**File to edit:** `in_memory_player_store.go`
**Expected result:** All tests pass

**Question:** Add a minimal `GetLeague` method to `InMemoryPlayerStore` that just returns `nil` for now.

**Answer:**
```go
func (i *InMemoryPlayerStore) GetLeague() []Player {
	return nil
}
```

---

## Exercise 16: Update leagueHandler to Use Store
**File to edit:** `server.go`
**Expected result:** All tests pass

**Question:** Update `leagueHandler` to call `p.store.GetLeague()` instead of `p.getLeagueTable()`. Delete the `getLeagueTable` method.

**Answer:**
```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)
}
```

---

## Exercise 17: Refactor Test - Extract Helper Functions
**File to edit:** `server_test.go`
**Expected result:** All tests pass (refactor)

**Question:** Create helper functions: `getLeagueFromResponse`, `assertLeague`, and `newLeagueRequest`. Update the test to use them.

**Answer:**
```go
t.Run("it returns the league table as JSON", func(t *testing.T) {
	wantedLeague := []Player{
		{"Cleo", 32},
		{"Chris", 20},
		{"Tiest", 14},
	}

	store := StubPlayerStore{nil, nil, wantedLeague}
	server := NewPlayerServer(&store)

	request := newLeagueRequest()
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := getLeagueFromResponse(t, response.Body)
	assertStatus(t, response.Code, http.StatusOK)
	assertLeague(t, got, wantedLeague)
})

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
```

---

## Exercise 18: Add Test for Content-Type Header
**File to edit:** `server_test.go`
**Expected result:** Test fails with `response did not have content-type of application/json`

**Question:** Add an assertion to check that the response has `content-type` header set to `application/json`.

**Answer:**
```go
if response.Result().Header.Get("content-type") != "application/json" {
	t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
}
```

---

## Exercise 19: Set Content-Type Header in leagueHandler
**File to edit:** `server.go`
**Expected result:** All tests pass

**Question:** Update `leagueHandler` to set the `content-type` header to `application/json` before encoding.

**Answer:**
```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.GetLeague())
}
```

---

## Exercise 20: Refactor - Create Constant and Helper for Content-Type
**Files to edit:** `server.go`, `server_test.go`
**Expected result:** All tests pass (refactor)

**Question:** Create a `jsonContentType` constant in `server.go`. Create an `assertContentType` helper in tests and use it.

**Answer:**

In `server.go`:
```go
const jsonContentType = "application/json"

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}
```

In `server_test.go`:
```go
func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
```

Use in test:
```go
assertContentType(t, response, jsonContentType)
```

---

## Exercise 21: Add Integration Test for /league Endpoint
**File to edit:** `server_integration_test.go`
**Expected result:** Test fails with `got [] want [{Pepper 3}]`

**Question:** Update the integration test to use `t.Run` for subtests. Add a "get league" subtest that POSTs 3 wins for "Pepper" then GETs `/league` and asserts the response contains `[{Pepper, 3}]`.

**Answer:**
```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
}
```

---

## Exercise 22: Implement GetLeague Properly on InMemoryPlayerStore
**File to edit:** `in_memory_player_store.go`
**Expected result:** All tests pass

**Question:** Update `GetLeague` on `InMemoryPlayerStore` to iterate over the map and return all players with their win counts.

**Answer:**
```go
func (i *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league
}
```

---

# Questions Only (No Answers)

1. **[server_test.go] [Test fails: 404 want 200]** Add a new test function `TestLeague` that makes a GET request to `/league` and asserts it returns status 200.

2. **[server.go] [Tests pass]** Update `ServeHTTP` to use `http.NewServeMux()` for routing. Handle `/league` to return `StatusOK` and `/players/` to use the existing player logic.

3. **[server.go] [Tests pass - refactor]** Extract the inline handler functions into separate methods `leagueHandler` and `playersHandler` on `PlayerServer`.

4. **[server.go] [Compilation error]** Add a `router` field to `PlayerServer` struct and create a `NewPlayerServer` constructor that sets up routing once. Simplify `ServeHTTP` to just delegate to the router.

5. **[server_test.go, server_integration_test.go, main.go] [Tests pass]** Replace all instances of `&PlayerServer{&store}` or `PlayerServer{&store}` with `NewPlayerServer(&store)`.

6. **[server.go] [Tests pass - refactor]** Replace the named `router *http.ServeMux` field with an embedded `http.Handler`. Update `NewPlayerServer` to assign the router to `p.Handler`. Delete the `ServeHTTP` method.

7. **[server.go] [Compiles]** Create a `Player` struct with `Name` (string) and `Wins` (int) fields to represent the JSON data model.

8. **[server_test.go] [Test fails: Unable to parse JSON]** Update the `TestLeague` test to decode the response body as JSON into a `[]Player` slice.

9. **[server.go] [Tests pass]** Update `leagueHandler` to encode and return a hard-coded slice containing one `Player{"Chris", 20}` as JSON.

10. **[server.go] [Tests pass - refactor]** Extract the hard-coded league data into a separate `getLeagueTable` method.

11. **[server_test.go] [Compilation error: too few values]** Add a `league []Player` field to `StubPlayerStore`.

12. **[server_test.go] [Test fails: got Chris want Cleo,Chris,Tiest]** Update `TestLeague` to create a `wantedLeague` with 3 players, pass it to the stub, and assert the response matches using `reflect.DeepEqual`. Update all `StubPlayerStore` initializers.

13. **[server.go] [Compilation error: missing GetLeague]** Add `GetLeague() []Player` to the `PlayerStore` interface.

14. **[server_test.go] [Still compilation error]** Add `GetLeague` method to `StubPlayerStore` that returns the `league` field.

15. **[in_memory_player_store.go] [Tests pass]** Add a minimal `GetLeague` method to `InMemoryPlayerStore` that just returns `nil`.

16. **[server.go] [Tests pass]** Update `leagueHandler` to call `p.store.GetLeague()` instead of `p.getLeagueTable()`. Delete the `getLeagueTable` method.

17. **[server_test.go] [Tests pass - refactor]** Create helper functions: `getLeagueFromResponse`, `assertLeague`, and `newLeagueRequest`. Update the test to use them.

18. **[server_test.go] [Test fails: wrong content-type]** Add an assertion to check that the response has `content-type` header set to `application/json`.

19. **[server.go] [Tests pass]** Update `leagueHandler` to set the `content-type` header to `application/json` before encoding.

20. **[server.go, server_test.go] [Tests pass - refactor]** Create a `jsonContentType` constant in `server.go`. Create an `assertContentType` helper in tests.

21. **[server_integration_test.go] [Test fails: got [] want [{Pepper 3}]]** Update the integration test to add a "get league" subtest that verifies the `/league` endpoint returns recorded wins.

22. **[in_memory_player_store.go] [All tests pass]** Update `GetLeague` on `InMemoryPlayerStore` to iterate over the map and return all players with their win counts.

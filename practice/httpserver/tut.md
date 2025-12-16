## HTTP Server TDD Practice Exercises

### Exercise 1: Write the First Test

**Task:** Create `server_test.go` with a test for `PlayerServer` that expects to get Pepper's score as "20".

**Answer:**

```go
func TestGETPlayers(t *testing.T) {
t.Run("returns Pepper's score", func(t *testing.T) {
request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
response := httptest.NewRecorder()

PlayerServer(response, request)

got := response.Body.String()
want := "20"

if got != want {
t.Errorf("got %q, want %q", got, want)
}
})
}
```

---

### Exercise 2: Create Minimal PlayerServer (Doesn't Compile)

**Task:** Create `server.go` with an empty `PlayerServer` function to satisfy the compiler's "undefined: PlayerServer" error.

**Answer:**

```go
func PlayerServer() {}
```

---

### Exercise 3: Add Arguments to PlayerServer

**Task:** The compiler says "too many arguments". Add the correct parameters to `PlayerServer`.

**Answer:**

```go
import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

---

### Exercise 4: Make the Test Pass with Hard-coded Value

**Task:** Make `PlayerServer` return "20" to pass the test.

**Answer:**

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
fmt.Fprint(w, "20")
}
```

---

### Exercise 5: Create main.go Scaffolding

**Task:** Create `main.go` that wires up `PlayerServer` as an HTTP handler on port 5000.

**Answer:**

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(PlayerServer)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
```

---

### Exercise 6: Add Second Test Case for Floyd

**Task:** Add a subtest for Floyd's score (expected: "10") to break the hard-coded approach.

**Answer:**

```go
t.Run("returns Floyd's score", func(t *testing.T) {
request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
response := httptest.NewRecorder()

PlayerServer(response, request)

got := response.Body.String()
want := "10"

if got != want {
t.Errorf("got %q, want %q", got, want)
}
})
```

---

### Exercise 7: Extract Player Name from URL

**Task:** Update `PlayerServer` to parse the player name from the URL and return different scores for Pepper and Floyd.

**Answer:**

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
player := strings.TrimPrefix(r.URL.Path, "/players/")

if player == "Pepper" {
fmt.Fprint(w, "20")
return
}

if player == "Floyd" {
fmt.Fprint(w, "10")
return
}
}
```

---

### Exercise 8: Refactor - Extract GetPlayerScore Function

**Task:** Refactor by extracting score retrieval into a separate `GetPlayerScore` function.

**Answer:**

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
player := strings.TrimPrefix(r.URL.Path, "/players/")

fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
if name == "Pepper" {
return "20"
}

if name == "Floyd" {
return "10"
}

return ""
}
```

---

### Exercise 9: Refactor Tests - Add Helper Functions

**Task:** DRY up the tests by creating `newGetScoreRequest` and `assertResponseBody` helpers.

**Answer:**

```go
func TestGETPlayers(t *testing.T) {
t.Run("returns Pepper's score", func(t *testing.T) {
request := newGetScoreRequest("Pepper")
response := httptest.NewRecorder()

PlayerServer(response, request)

assertResponseBody(t, response.Body.String(), "20")
})

t.Run("returns Floyd's score", func(t *testing.T) {
request := newGetScoreRequest("Floyd")
response := httptest.NewRecorder()

PlayerServer(response, request)

assertResponseBody(t, response.Body.String(), "10")
})
}

func newGetScoreRequest(name string) *http.Request {
req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
return req
}

func assertResponseBody(t testing.TB, got, want string) {
t.Helper()
if got != want {
t.Errorf("response body is wrong, got %q want %q", got, want)
}
}
```

---

### Exercise 10: Define PlayerStore Interface

**Task:** Create a `PlayerStore` interface with a `GetPlayerScore` method.

**Answer:**

```go
type PlayerStore interface {
GetPlayerScore(name string) int
}
```

---

### Exercise 11: Convert PlayerServer to a Struct

**Task:** Convert `PlayerServer` from a function to a struct that holds a `PlayerStore`.

**Answer:**

```go
type PlayerServer struct {
store PlayerStore
}
```

---

### Exercise 12: Implement ServeHTTP Method

**Task:** Add a `ServeHTTP` method to `PlayerServer` that implements the `Handler` interface.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
player := strings.TrimPrefix(r.URL.Path, "/players/")
fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

---

### Exercise 13: Update Tests to Use PlayerServer Struct

**Task:** Update tests to create a `PlayerServer` instance and call `ServeHTTP`.

**Answer:**

```go
func TestGETPlayers(t *testing.T) {
server := &PlayerServer{}

t.Run("returns Pepper's score", func(t *testing.T) {
request := newGetScoreRequest("Pepper")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertResponseBody(t, response.Body.String(), "20")
})

t.Run("returns Floyd's score", func(t *testing.T) {
request := newGetScoreRequest("Floyd")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertResponseBody(t, response.Body.String(), "10")
})
}
```

---

### Exercise 14: Update main.go for Struct

**Task:** Update `main.go` to create a `PlayerServer` struct instance.

**Answer:**

```go
func main() {
server := &PlayerServer{}
log.Fatal(http.ListenAndServe(":5000", server))
}
```

---

### Exercise 15: Create StubPlayerStore for Tests

**Task:** Create a `StubPlayerStore` struct that implements `PlayerStore` using a map.

**Answer:**

```go
type StubPlayerStore struct {
scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
score := s.scores[name]
return score
}
```

---

### Exercise 16: Inject StubPlayerStore into Tests

**Task:** Create a `StubPlayerStore` with test data and inject it into `PlayerServer`.

**Answer:**

```go
func TestGETPlayers(t *testing.T) {
store := StubPlayerStore{
map[string]int{
"Pepper": 20,
"Floyd":  10,
},
}
server := &PlayerServer{&store}

t.Run("returns Pepper's score", func(t *testing.T) {
request := newGetScoreRequest("Pepper")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertResponseBody(t, response.Body.String(), "20")
})

t.Run("returns Floyd's score", func(t *testing.T) {
request := newGetScoreRequest("Floyd")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertResponseBody(t, response.Body.String(), "10")
})
}
```

---

### Exercise 17: Create InMemoryPlayerStore for main.go

**Task:** Create a minimal `InMemoryPlayerStore` in `main.go` that returns a hard-coded value.

**Answer:**

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
return 123
}

func main() {
server := &PlayerServer{&InMemoryPlayerStore{}}
log.Fatal(http.ListenAndServe(":5000", server))
}
```

---

### Exercise 18: Add Test for Missing Player (404)

**Task:** Add a test case that expects 404 status for a player not in the store.

**Answer:**

```go
t.Run("returns 404 on missing players", func(t *testing.T) {
request := newGetScoreRequest("Apollo")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

got := response.Code
want := http.StatusNotFound

if got != want {
t.Errorf("got status %d want %d", got, want)
}
})
```

---

### Exercise 19: Return 404 (Minimal - All Responses)

**Task:** Make the test pass by writing `StatusNotFound` on all responses (intentionally wrong).

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
player := strings.TrimPrefix(r.URL.Path, "/players/")

w.WriteHeader(http.StatusNotFound)

fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

---

### Exercise 20: Add Status Assertions to Existing Tests

**Task:** Update all test cases to assert status codes and create an `assertStatus` helper.

**Answer:**

```go
func TestGETPlayers(t *testing.T) {
store := StubPlayerStore{
map[string]int{
"Pepper": 20,
"Floyd":  10,
},
}
server := &PlayerServer{&store}

t.Run("returns Pepper's score", func(t *testing.T) {
request := newGetScoreRequest("Pepper")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertStatus(t, response.Code, http.StatusOK)
assertResponseBody(t, response.Body.String(), "20")
})

t.Run("returns Floyd's score", func(t *testing.T) {
request := newGetScoreRequest("Floyd")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertStatus(t, response.Code, http.StatusOK)
assertResponseBody(t, response.Body.String(), "10")
})

t.Run("returns 404 on missing players", func(t *testing.T) {
request := newGetScoreRequest("Apollo")
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertStatus(t, response.Code, http.StatusNotFound)
})
}

func assertStatus(t testing.TB, got, want int) {
t.Helper()
if got != want {
t.Errorf("did not get correct status, got %d, want %d", got, want)
}
}
```

---

### Exercise 21: Fix 404 Logic (Only for Missing Players)

**Task:** Update `ServeHTTP` to only return 404 when score is 0.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
player := strings.TrimPrefix(r.URL.Path, "/players/")

score := p.store.GetPlayerScore(player)

if score == 0 {
w.WriteHeader(http.StatusNotFound)
}

fmt.Fprint(w, score)
}
```

---

### Exercise 22: Add Test for POST (Status Accepted)

**Task:** Write a test for `POST /players/{name}` that expects `StatusAccepted`.

**Answer:**

```go
func TestStoreWins(t *testing.T) {
store := StubPlayerStore{
map[string]int{},
}
server := &PlayerServer{&store}

t.Run("it returns accepted on POST", func(t *testing.T) {
request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
response := httptest.NewRecorder()

server.ServeHTTP(response, request)

assertStatus(t, response.Code, http.StatusAccepted)
})
}
```

---

### Exercise 23: Handle POST Method

**Task:** Add an `if` statement to check for POST method and return `StatusAccepted`.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

if r.Method == http.MethodPost {
w.WriteHeader(http.StatusAccepted)
return
}

player := strings.TrimPrefix(r.URL.Path, "/players/")

score := p.store.GetPlayerScore(player)

if score == 0 {
w.WriteHeader(http.StatusNotFound)
}

fmt.Fprint(w, score)
}
```

---

### Exercise 24: Refactor with Switch and Extract Methods

**Task:** Refactor `ServeHTTP` to use a switch statement and extract `processWin` and `showScore` methods.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

switch r.Method {
case http.MethodPost:
p.processWin(w)
case http.MethodGet:
p.showScore(w, r)
}

}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
player := strings.TrimPrefix(r.URL.Path, "/players/")

score := p.store.GetPlayerScore(player)

if score == 0 {
w.WriteHeader(http.StatusNotFound)
}

fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter) {
w.WriteHeader(http.StatusAccepted)
}
```

---

### Exercise 25: Add RecordWin to StubPlayerStore

**Task:** Extend `StubPlayerStore` with a `winCalls` slice and `RecordWin` method to spy on calls.

**Answer:**

```go
type StubPlayerStore struct {
scores   map[strin]int
winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
score := s.scores[name]
return score
}

func (s *StubPlayerStore) RecordWin(name string) {
s.winCalls = append(s.winCalls, name)
}
```

---

### Exercise 26: Test That RecordWin Is Called

**Task:** Update `TestStoreWins` to verify that `RecordWin` is called once on POST. Add `newPostWinRequest` helper.

**Answer:**

```go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{&store}

	t.Run("it records wins when POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
```

---

### Exercise 27: Update StubPlayerStore Initializers

**Task:** Fix struct initializer errors by adding `nil` for the new `winCalls` field.

**Answer:**

```go
store := StubPlayerStore{
	map[string]int{},
	nil,
}
```

---

### Exercise 28: Add RecordWin to PlayerStore Interface

**Task:** Add `RecordWin(name string)` to the `PlayerStore` interface.

**Answer:**

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}
```

---

### Exercise 29: Add RecordWin to InMemoryPlayerStore

**Task:** Add an empty `RecordWin` method to `InMemoryPlayerStore` to satisfy the interface.

**Answer:**

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

---

### Exercise 30: Call RecordWin in processWin (Hard-coded)

**Task:** Call `RecordWin` in `processWin` with a hard-coded name "Bob".

**Answer:**

```go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
	p.store.RecordWin("Bob")
	w.WriteHeader(http.StatusAccepted)
}
```

---

### Exercise 31: Test Correct Player Name Is Recorded

**Task:** Update the test to verify the correct player name is passed to `RecordWin`.

**Answer:**

```go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{&store}

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}
```

---

### Exercise 32: Extract Player Name in processWin

**Task:** Update `processWin` to accept `http.Request` and extract the player name from URL.

**Answer:**

```go
func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

---

### Exercise 33: Refactor - Extract Player Name Once

**Task:** DRY up by extracting player name once in `ServeHTTP` and passing it to both methods.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

---

### Exercise 34: Write Integration Test

**Task:** Create `server_integration_test.go` that tests `PlayerServer` with `InMemoryPlayerStore` - POST 3 wins then GET score.

**Answer:**

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{}
	server := PlayerServer{&store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
```

---

### Exercise 35: Implement Real InMemoryPlayerStore

**Task:** Create `in_memory_player_store.go` with a working `InMemoryPlayerStore` that uses a map and a constructor.

**Answer:**

```go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
```


```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore{}
	server := PlayerServer{store}
	
	...
```

---

### Exercise 36: Update main.go to Use NewInMemoryPlayerStore

**Task:** Update `main.go` to use the constructor function.

**Answer:**

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

# HTTP Server TDD Practice - Questions & Answers

------

## Section 1: Getting Started with a Hard-coded Response

### Exercise 1

**File:** `server_test.go` (create new)
 **Expected:** Compilation error - `undefined: PlayerServer`

**Question:** Create `server_test.go` with a test for `PlayerServer` that expects to get Pepper's score as "20".

**Answer:**

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

------

### Exercise 2

**File:** `server.go` (create new)
 **Expected:** Compilation error - `too many arguments in call to PlayerServer`

**Question:** Create `server.go` with an empty `PlayerServer` function to satisfy the compiler's "undefined: PlayerServer" error.

**Answer:**

```go
package main

func PlayerServer() {}
```

------

### Exercise 3

**File:** `server.go`
 **Expected:** Test fails with `got '', want '20'`

**Question:** The compiler says "too many arguments". Add the correct parameters to `PlayerServer`.

**Answer:**

```go
package main

import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

------

### Exercise 4

**File:** `server.go`
 **Expected:** Tests pass ✓

**Question:** Make `PlayerServer` return "20" to pass the test.

**Answer:**

```go
package main

import (
	"fmt"
	"net/http"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
```

------

### Exercise 5

**File:** `main.go` (create new)
 **Expected:** Application builds and runs (manual test only)

**Question:** Create `main.go` that wires up `PlayerServer` as an HTTP handler on port 5000.

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

------

## Section 2: Breaking the Hard-coded Value

### Exercise 6

**File:** `server_test.go`
 **Expected:** Test fails with `got '20', want '10'`

**Question:** Add a subtest for Floyd's score (expected: "10") to break the hard-coded approach.

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

------

### Exercise 7

**File:** `server.go`
 **Expected:** Tests pass ✓

**Question:** Update `PlayerServer` to parse the player name from the URL and return different scores for Pepper ("20") and Floyd ("10").

**Answer:**

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
)

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

------

### Exercise 8

**File:** `server.go`
 **Expected:** Tests pass ✓ (refactor)

**Question:** Refactor by extracting score retrieval into a separate `GetPlayerScore` function.

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

------

### Exercise 9

**File:** `server_test.go`
 **Expected:** Tests pass ✓ (refactor)

**Question:** DRY up the tests by creating `newGetScoreRequest` and `assertResponseBody` helpers.

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

------

## Section 3: Introducing the PlayerStore Interface

### Exercise 10

**File:** `server.go`
 **Expected:** Tests pass ✓ (no behavior change yet)

**Question:** Create a `PlayerStore` interface with a `GetPlayerScore(name string) int` method.

**Answer:**

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
}
```

------

### Exercise 11

**File:** `server.go`
 **Expected:** Compilation error - tests need updating

**Question:** Convert `PlayerServer` from a function to a struct that holds a `PlayerStore`.

**Answer:**

```go
type PlayerServer struct {
	store PlayerStore
}
```

------

### Exercise 12

**File:** `server.go`
 **Expected:** Compilation error - tests still need updating

**Question:** Add a `ServeHTTP` method to `PlayerServer` that implements the `Handler` interface.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

Full server.go at this point:

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
}

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

------

### Exercise 13

**File:** `server_test.go`
 **Expected:** Compilation error - main.go needs updating

**Question:** Update tests to create a `PlayerServer` instance and call `ServeHTTP`.

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

------

### Exercise 14

**File:** `main.go`
 **Expected:** Runtime panic - `invalid memory address or nil pointer dereference`

**Question:** Update `main.go` to create a `PlayerServer` struct instance.

**Answer:**

```go
func main() {
	server := &PlayerServer{}
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

------

### Exercise 15

**File:** `server_test.go`
 **Expected:** Runtime panic still (store not injected yet)

**Question:** Create a `StubPlayerStore` struct that implements `PlayerStore` using a map.

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

------

### Exercise 16

**File:** `server_test.go`
 **Expected:** Tests pass ✓

**Question:** Create a `StubPlayerStore` with test data and inject it into `PlayerServer`.

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

------

### Exercise 17

**File:** `main.go`
 **Expected:** Application runs (returns "123" for all players)

**Question:** Create a minimal `InMemoryPlayerStore` in `main.go` that returns a hard-coded value (123).

**Answer:**

```go
package main

import (
	"log"
	"net/http"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

------

## Section 4: Handling Missing Players (404)

### Exercise 18

**File:** `server_test.go`
 **Expected:** Test fails with `got status 200 want 404`

**Question:** Add a test case that expects 404 status for a player not in the store.

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

------

### Exercise 19

**File:** `server.go`
 **Expected:** Tests pass ✓ (but incorrectly - all responses return 404)

**Question:** Make the test pass by writing `StatusNotFound` on all responses (intentionally wrong).

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

------

### Exercise 20

**File:** `server_test.go`
 **Expected:** Tests fail - Pepper and Floyd tests now fail with `got 404, want 200`

**Question:** Update all test cases to assert status codes and create an `assertStatus` helper.

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

------

### Exercise 21

**File:** `server.go`
 **Expected:** Tests pass ✓

**Question:** Update `ServeHTTP` to only return 404 when score is 0.

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

------

## Section 5: Handling POST Requests

### Exercise 22

**File:** `server_test.go`
 **Expected:** Test fails with `got status 404, want 202`

**Question:** Write a new test function `TestStoreWins` for `POST /players/{name}` that expects `StatusAccepted`.

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

------

### Exercise 23

**File:** `server.go`
 **Expected:** Tests pass ✓

**Question:** Add an `if` statement to check for POST method and return `StatusAccepted`.

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

------

### Exercise 24

**File:** `server.go`
 **Expected:** Tests pass ✓ (refactor)

**Question:** Refactor `ServeHTTP` to use a switch statement and extract `processWin` and `showScore` methods.

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

------

## Section 6: Recording Wins

### Exercise 25

**File:** `server_test.go`
 **Expected:** Compilation error - `too few values in struct initializer`

**Question:** Extend `StubPlayerStore` with a `winCalls []string` field and `RecordWin` method to spy on calls.

**Answer:**

```go
type StubPlayerStore struct {
	scores   map[string]int
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

------

### Exercise 26

**File:** `server_test.go`
 **Expected:** Test fails with `got 0 calls to RecordWin want 1`

**Question:** Update `TestStoreWins` to verify that `RecordWin` is called once on POST. Add `newPostWinRequest` helper. Fix struct initializers by adding `nil` for the new `winCalls` field.

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

Also update `TestGETPlayers`:

```go
store := StubPlayerStore{
	map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	},
	nil,
}
```

------

### Exercise 27

**File:** `server.go`
 **Expected:** Compilation error - `InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)`

**Question:** Add `RecordWin(name string)` to the `PlayerStore` interface.

**Answer:**

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}
```

------

### Exercise 28

**File:** `main.go`
 **Expected:** Test fails with `got 0 calls to RecordWin want 1`

**Question:** Add an empty `RecordWin` method to `InMemoryPlayerStore` to satisfy the interface.

**Answer:**

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

------

### Exercise 29

**File:** `server.go`
 **Expected:** Tests pass ✓

**Question:** Call `RecordWin` in `processWin` with a hard-coded name "Bob".

**Answer:**

```go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
	p.store.RecordWin("Bob")
	w.WriteHeader(http.StatusAccepted)
}
```

------

### Exercise 30

**File:** `server_test.go`
 **Expected:** Test fails with `did not store correct winner got 'Bob' want 'Pepper'`

**Question:** Update the test to verify the correct player name is passed to `RecordWin`.

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

------

### Exercise 31

**File:** `server.go`
 **Expected:** Tests pass ✓

**Question:** Update `processWin` to accept `http.Request` and extract the player name from URL. Update the switch statement to pass the request.

**Answer:**

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, r)
	case http.MethodGet:
		p.showScore(w, r)
	}

}

func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

------

### Exercise 32

**File:** `server.go`
 **Expected:** Tests pass ✓ (refactor)

**Question:** DRY up by extracting player name once in `ServeHTTP` and passing it to both methods.

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

------

## Section 7: Integration Test

### Exercise 33

**File:** `server_integration_test.go` (create new)
 **Expected:** Test fails with `got '123' want '3'`

**Question:** Create `server_integration_test.go` that tests `PlayerServer` with `InMemoryPlayerStore` - POST 3 wins then GET score.

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

------

### Exercise 34

**File:** `in_memory_player_store.go` (create new)
 **Expected:** Tests pass ✓

**Question:** Create `in_memory_player_store.go` with a working `InMemoryPlayerStore` that uses a map and a constructor.

**Answer:**

```go
package main

type InMemoryPlayerStore struct {
	store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
```

------

### Exercise 35

**File:** `server_integration_test.go`
 **Expected:** Tests pass ✓

**Question:** Update the integration test to use `NewInMemoryPlayerStore()`.

**Answer:**

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
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

------

### Exercise 36

**File:** `main.go`
 **Expected:** Application works correctly

**Question:** Update `main.go` to use the `NewInMemoryPlayerStore` constructor. Remove the old `InMemoryPlayerStore` from main.go.

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

------

## Final Code Listing

### server.go

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
}

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

### server_test.go

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
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

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
```

### in_memory_player_store.go

```go
package main

type InMemoryPlayerStore struct {
	store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
```

### server_integration_test.go

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
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

### main.go

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

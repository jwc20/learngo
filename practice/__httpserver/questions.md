# HTTP Server TDD Practice Questions

------

## Section 1: Getting Started with a Hard-coded Response

### Exercise 1

**File:** `server_test.go` (create new)
 **Expected:** Compilation error - `undefined: PlayerServer`

Create `server_test.go` with a test for `PlayerServer` that expects to get Pepper's score as "20".

------

### Exercise 2

**File:** `server.go` (create new)
 **Expected:** Compilation error - `too many arguments in call to PlayerServer`

Create `server.go` with an empty `PlayerServer` function to satisfy the compiler's "undefined: PlayerServer" error.

------

### Exercise 3

**File:** `server.go`
 **Expected:** Test fails with `got '', want '20'`

The compiler says "too many arguments". Add the correct parameters to `PlayerServer`.

------

### Exercise 4

**File:** `server.go`
 **Expected:** Tests pass ✓

Make `PlayerServer` return "20" to pass the test.

------

### Exercise 5

**File:** `main.go` (create new)
 **Expected:** Application builds and runs (manual test only)

Create `main.go` that wires up `PlayerServer` as an HTTP handler on port 5000.

To run this, do `go build` which will take all the `.go` files in the directory and build you a program. You can then execute it with `./myprogram`.

------







## Section 2: Breaking the Hard-coded Value

### Exercise 6

**File:** `server_test.go`
 **Expected:** Test fails with `got '20', want '10'`

Add a subtest for Floyd's score (expected: "10") to break the hard-coded approach.

------

### Exercise 7

**File:** `server.go`
 **Expected:** Tests pass ✓

Update `PlayerServer` to parse the player name from the URL and return different scores for Pepper ("20") and Floyd ("10").

------

### Exercise 8

**File:** `server.go`
 **Expected:** Tests pass ✓ (refactor)

Refactor by extracting score retrieval into a separate `GetPlayerScore` function.

------

### Exercise 9

**File:** `server_test.go`
 **Expected:** Tests pass ✓ (refactor)

DRY up the tests by creating `newGetScoreRequest` and `assertResponseBody` helpers.

------







## Section 3: Introducing the PlayerStore Interface

### Exercise 10

**File:** `server.go`
 **Expected:** Tests pass ✓ (no behavior change yet)

Create a `PlayerStore` interface with a `GetPlayerScore(name string) int` method.

------

### Exercise 11

**File:** `server.go`
 **Expected:** Compilation error - tests need updating

Convert `PlayerServer` from a function to a struct that holds a `PlayerStore`.

------

### Exercise 12

**File:** `server.go`
 **Expected:** Compilation error - tests still need updating

Add a `ServeHTTP` method to `PlayerServer` that implements the `Handler` interface. Move the existing handler logic into this method and call `p.store.GetPlayerScore(player)`.

------

### Exercise 13

**File:** `server_test.go`
 **Expected:** Compilation error - main.go needs updating

Update tests to create a `PlayerServer` instance and call `ServeHTTP` instead of `PlayerServer` function.

------

### Exercise 14 x

**File:** `main.go`
 **Expected:** Runtime panic - `invalid memory address or nil pointer dereference`

Update `main.go` to create a `PlayerServer` struct instance.

------

### Exercise 15

**File:** `server_test.go`
 **Expected:** Runtime panic still (store not injected yet)

Create a `StubPlayerStore` struct that implements `PlayerStore` using a map. And `GetPlayerScore` method.

------

### Exercise 16

**File:** `server_test.go`
 **Expected:** Tests pass ✓

Create a `StubPlayerStore` with test data and inject it into `PlayerServer`.

------

### Exercise 17 x

**File:** `main.go`
 **Expected:** Application runs (returns "123" for all players)

Create a minimal `InMemoryPlayerStore` in `main.go` that returns a hard-coded value (123).

------







## Section 4: Handling Missing Players (404)

### Exercise 18

**File:** `server_test.go`
 **Expected:** Test fails with `got status 200 want 404`

Add a test case that expects 404 status for a player not in the store.

------

### Exercise 19 x

**File:** `server.go`
 **Expected:** Tests pass ✓ (but incorrectly - all responses return 404)

Make the test pass by writing `StatusNotFound` on **all** responses (intentionally wrong to highlight test gaps). Use `WriteHeader`method.

------

### Exercise 20

**File:** `server_test.go`
 **Expected:** Tests fail - Pepper and Floyd tests now fail with `got 404, want 200`

Update all test cases to assert status codes (`StatusOK` for existing players) and create an `assertStatus` helper.

------

### Exercise 21

**File:** `server.go`
 **Expected:** Tests pass ✓

Update `ServeHTTP` to only return 404 when score is 0.

------







## Section 5: Handling POST Requests

### Exercise 22

**File:** `server_test.go`
 **Expected:** Test fails with `got status 404, want 202`

Write a new test function `TestStoreWins` for `POST /players/{name}` that expects `StatusAccepted`.

------

### Exercise 23

**File:** `server.go`
 **Expected:** Tests pass ✓

Add an `if` statement to check for POST method and return `StatusAccepted`.

------

### Exercise 24

**File:** `server.go`
 **Expected:** Tests pass ✓ (refactor)

Refactor `ServeHTTP` to use a switch statement and extract `processWin` and `showScore` methods.

The handler is looking a bit muddled now. Refactor `ServeHTTP` to make the routing clearer:

1.  Replace the `if` statement with a `switch` statement on `r.Method`
2.  Extract the GET logic (player lookup, 404 handling, score display) into a new method called `showScore(w http.ResponseWriter, r *http.Request)`
3.  Extract the POST logic (returning StatusAccepted) into a new method called `processWin(w http.ResponseWriter)`
4.  The `switch` should call `p.processWin(w)` for POST and `p.showScore(w, r)` for GET

------







## Section 6: Recording Wins

### Exercise 25

**File:** `server_test.go`
 **Expected:** Compilation error - `too few values in struct initializer`

Extend `StubPlayerStore` with a `winCalls []string` field and `RecordWin` method to spy on calls.

------

### Exercise 26

**File:** `server_test.go`
 **Expected:** Test fails with `got 0 calls to RecordWin want 1`

Update `TestStoreWins` to verify that `RecordWin` is called once on POST. Add `newPostWinRequest` helper. Fix struct initializers by adding `nil` for the new `winCalls` field.



```go
if len(store.winCalls) != 1 {
	t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
}
```





------

### Exercise 27

**File:** `server.go`
 **Expected:** Compilation error - `InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)`

Add `RecordWin(name string)` to the `PlayerStore` interface.

------

### Exercise 28

**File:** `main.go`
 **Expected:** Test fails with `got 0 calls to RecordWin want 1`

Add an empty `RecordWin` method to `InMemoryPlayerStore` to satisfy the interface.

------

### Exercise 29

**File:** `server.go`
 **Expected:** Tests pass ✓

Call `RecordWin` in `processWin` with a hard-coded name "Bob".

------

### Exercise 30

**File:** `server_test.go`
 **Expected:** Test fails with `did not store correct winner got 'Bob' want 'Pepper'`

Update the test to verify the correct player name is passed to `RecordWin`.

------

### Exercise 31

**File:** `server.go`
 **Expected:** Tests pass ✓

Update `processWin` to accept `http.Request` and extract the player name from URL. Update the switch statement to pass the request.

------

### Exercise 32

**File:** `server.go`
 **Expected:** Tests pass ✓ (refactor)

DRY up by extracting player name once in `ServeHTTP` and passing it to both methods (change method signatures to accept `player string`).

------









## Section 7: Integration Test

### Exercise 33 x

**File:** `server_integration_test.go` (create new)
 **Expected:** Test fails with `got '123' want '3'`

Create `server_integration_test.go` that tests `PlayerServer` with `InMemoryPlayerStore` - POST 3 wins for "Pepper" then GET score and expect "3".

------

### Exercise 34 x

**File:** `in_memory_player_store.go` (create new)
 **Expected:** Tests fails

Create `in_memory_player_store.go` with a working `InMemoryPlayerStore` that uses a `map[string]int` and a `NewInMemoryPlayerStore` constructor.

And then update `GetPlayerScore` and `RecordWin` methods.

------

### Exercise 35

**File:** `server_integration_test.go`
 **Expected:** Tests pass ✓

Update the integration test to use `NewInMemoryPlayerStore()`.

------

### Exercise 36

**File:** `main.go`
 **Expected:** Application works correctly

Update `main.go` to use the `NewInMemoryPlayerStore` constructor. Remove the old `InMemoryPlayerStore` from main.go.

------

## Final File Structure

After completing all exercises, you should have:

-   `server.go` - PlayerServer struct, PlayerStore interface, HTTP handlers
-   `server_test.go` - Unit tests with StubPlayerStore
-   `in_memory_player_store.go` - In-memory implementation of PlayerStore
-   `server_integration_test.go` - Integration test
-   `main.go` - Application entry point
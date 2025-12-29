# JSON, Routing and Embedding - Practice Questions Only

## Starting Point
You should have the completed files from the HTTP Server chapter.

---

## Exercise 1
**File:** `server_test.go`
**Expected:** Test fails with `got 404, want 200`

Add a new test function `TestLeague` that makes a GET request to `/league` and asserts it returns status 200.

---

## Exercise 2
**File:** `server.go`
**Expected:** All tests pass

Update `ServeHTTP` to use `http.NewServeMux()` for routing. Handle `/league` to return `StatusOK` and `/players/` to use the existing player logic.

---

## Exercise 3
**File:** `server.go`
**Expected:** All tests pass (refactor)

Extract the inline handler functions into separate methods `leagueHandler` and `playersHandler` on `PlayerServer`.

---

## Exercise 4
**File:** `server.go`
**Expected:** Compilation error - tests need updating

Add a `router` field to `PlayerServer` struct and create a `NewPlayerServer` constructor that sets up routing once. Simplify `ServeHTTP` to just delegate to the router.

---

## Exercise 5
**Files:** `server_test.go`, `server_integration_test.go`, `main.go`
**Expected:** All tests pass

Replace all instances of `&PlayerServer{&store}` or `PlayerServer{&store}` with `NewPlayerServer(&store)`.

---

## Exercise 6
**File:** `server.go`
**Expected:** All tests pass (refactor)

Replace the named `router *http.ServeMux` field with an embedded `http.Handler`. Update `NewPlayerServer` to assign the router to `p.Handler`. Delete the `ServeHTTP` method (it's now provided by embedding).

---

## Exercise 7
**File:** `server.go`
**Expected:** Code compiles (no test change)

Create a `Player` struct with `Name` (string) and `Wins` (int) fields to represent the JSON data model.

---

## Exercise 8
**File:** `server_test.go`
**Expected:** Test fails with `Unable to parse response from server '' into slice of Player, 'unexpected end of JSON input'`

Update the `TestLeague` test to decode the response body as JSON into a `[]Player` slice. Import `encoding/json`.

---

## Exercise 9
**File:** `server.go`
**Expected:** All tests pass

Update `leagueHandler` to encode and return a hard-coded slice containing one `Player{"Chris", 20}` as JSON.

---

## Exercise 10
**File:** `server.go`
**Expected:** All tests pass (refactor)

Extract the hard-coded league data into a separate `getLeagueTable` method.

---

## Exercise 11
**File:** `server_test.go`
**Expected:** Compilation error - struct initializers have too few values

Add a `league []Player` field to `StubPlayerStore`.

---

## Exercise 12
**File:** `server_test.go`
**Expected:** Test fails with `got [{Chris 20}] want [{Cleo 32} {Chris 20} {Tiest 14}]`

Update `TestLeague` to create a `wantedLeague` with 3 players, pass it to the stub, and assert the response matches using `reflect.DeepEqual`. Update all other `StubPlayerStore` initializers to include `nil` for the league field.

---

## Exercise 13
**File:** `server.go`
**Expected:** Compilation error - missing GetLeague method on implementations

Add `GetLeague() []Player` to the `PlayerStore` interface.

---

## Exercise 14
**File:** `server_test.go`
**Expected:** Still compilation error (InMemoryPlayerStore missing method)

Add `GetLeague` method to `StubPlayerStore` that returns the `league` field.

---

## Exercise 15
**File:** `in_memory_player_store.go`
**Expected:** All tests pass

Add a minimal `GetLeague` method to `InMemoryPlayerStore` that just returns `nil` for now.

---

## Exercise 16
**File:** `server.go`
**Expected:** All tests pass

Update `leagueHandler` to call `p.store.GetLeague()` instead of `p.getLeagueTable()`. Delete the `getLeagueTable` method.

---

## Exercise 17
**File:** `server_test.go`
**Expected:** All tests pass (refactor)

Create helper functions: `getLeagueFromResponse`, `assertLeague`, and `newLeagueRequest`. Update the test to use them.

---

## Exercise 18
**File:** `server_test.go`
**Expected:** Test fails with `response did not have content-type of application/json`

Add an assertion to check that the response has `content-type` header set to `application/json`.

---

## Exercise 19
**File:** `server.go`
**Expected:** All tests pass

Update `leagueHandler` to set the `content-type` header to `application/json` before encoding.

---

## Exercise 20
**Files:** `server.go`, `server_test.go`
**Expected:** All tests pass (refactor)

Create a `jsonContentType` constant in `server.go`. Create an `assertContentType` helper in tests and use it.

---

## Exercise 21
**File:** `server_integration_test.go`
**Expected:** Test fails with `got [] want [{Pepper 3}]`

Update the integration test to use `t.Run` for subtests. Add a "get league" subtest that POSTs 3 wins for "Pepper" then GETs `/league` and asserts the response contains `[{"Pepper", 3}]`.

---

## Exercise 22
**File:** `in_memory_player_store.go`
**Expected:** All tests pass

Update `GetLeague` on `InMemoryPlayerStore` to iterate over the map and return all players with their win counts.

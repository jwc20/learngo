HTTP Server TDD Practice Questions

---
# Section 1
## Exercise 1

Create `server_test.go` with a test for `PlayerServer` that expects to get Pepper's score as "20".

## Exercise 2

Create `server.go` with an empty `PlayerServer` function to satisfy the compiler's "undefined: PlayerServer" error.

## Exercise 3

The compiler says "too many arguments". Add the correct parameters to `PlayerServer`.

## Exercise 4

Make `PlayerServer` return "20" to pass the test.

## **Exercise 5**

Create `main.go` that wires up `PlayerServer` as an HTTP handler on port 5000.

To run this, do `go build` which will take all the `.go` files in the directory and build you a program. You can then execute it with `./myprogram`.

---
# Section 2
## Exercise 6

Add a subtest for Floyd's score (expected: "10") to break the hard-coded approach.

## **Exercise 7**

Update `PlayerServer` to parse the player name from the URL and return different scores for Pepper and Floyd.

## Exercise 8

Refactor by extracting score retrieval into a separate `GetPlayerScore` function.

## **Exercise 9**

DRY up the tests by creating `newGetScoreRequest` and `assertResponseBody` helpers.

---
## Exercise 10

Create a `PlayerStore` interface with a `GetPlayerScore` method.

(in `server.go`)
(should not compile)
## Exercise 11

Convert `PlayerServer` from a function to a struct that holds a `PlayerStore`.
(in `server.go`)
## Exercise 12

Add a `ServeHTTP` method to `PlayerServer` that implements the `Handler` interface.
(in `server.go`)
## Exercise 13

Update tests to create a `PlayerServer` instance and call `ServeHTTP`.
(in `server_test.go`)
## Exercise 14

Update `main.go` to create a `PlayerServer` struct instance.

(tests will not run)
- `panic: runtime error: invalid memory address or nil pointer dereference ...`
- When `ServeHTTP` runs, it tries to call `p.store.GetPlayerScore(player)` on a nil store, causing the panic.
## Exercise 15

Create a `StubPlayerStore` struct that implements `PlayerStore` using a map.

(in `server_test.go`)
(tests will **not** pass)
## Exercise 16

Create a `StubPlayerStore` with test data and inject it into `PlayerServer`.
(in `server_test.go`)
(tests will now pass)

## Exercise 17

Create a minimal `InMemoryPlayerStore` in `main.go` that returns a hard-coded value.

(just let it return 123)

If you run `go build` again and hit the same URL you should get `"123"`. Not great, but until we **store data that's the best we can do**.
It also didn't feel great that our main application was starting up but not actually working.
We had to manually test to see the problem.

We have a few options as to what to do next
- Handle the scenario where the player doesn't exist
- Handle the `POST /players/{name}` scenario

Whilst the `POST` scenario gets us closer to the "happy path", I feel it'll be easier to tackle the **missing player scenario first** as we're in that context already. We'll get to the rest later.

---
# Section 3
## Exercise 18

Add a test case that expects 404 status for a player not in the store.

(in `server_test.go`)
(tests will not pass)

## Exercise 19

Make the test pass by writing `StatusNotFound` on all responses (intentionally wrong to highlight test gaps).

(in `server.go`)
(tests will now pass)
## Exercise 20

Update all test cases to assert status codes and create an `assertStatus` helper.

(in `server_test.go`)
(use `StatusOK` and `StatusNotFound`)
(tests will not pass)
## Exercise 21

Update `ServeHTTP` to only return 404 when score is 0.
(tests will now pass)

---
# Section 4

## Exercise 22

Write a test for `POST /players/{name}` that expects `StatusAccepted`.

(tests will not pass)
## Exercise 23

Add an `if` statement to check for POST method and return `StatusAccepted`.

(in `server.go`)
(tests will now pass)

## **Exercise 24**

Refactor `ServeHTTP` to use a switch statement and extract `processWin` and `showScore` methods.

(in `server.go`)
(The test will pass)

This makes the routing aspect of `ServeHTTP` a bit clearer and means our next iterations on storing can just be inside `processWin`.

Next, we want to check that when we do our `POST /players/{name}` that our `PlayerStore` is told to record the win.


## Exercise 25

Extend `StubPlayerStore` with a `winCalls` slice and `RecordWin` method to spy on calls.

## Exercise 26

Update `TestStoreWins` to verify that `RecordWin` is called once on POST.
Add `newPostWinRequest` helper.

(test will fail)
## Exercise 27

Fix struct initializer errors by adding `nil` for the new `winCalls` field.

## Exercise 28

Add `RecordWin(name string)` to the `PlayerStore` interface.

## Exercise 29

Add an empty `RecordWin` method to `InMemoryPlayerStore` to satisfy the interface.

## Exercise 30

Call `RecordWin` in `processWin` with a hard-coded name "Bob".

(test will fail)
## Exercise 31

Update the test to verify the correct player name is passed to `RecordWin`.

## Exercise 32

Update `processWin` to accept `http.Request` and extract the player name from URL.

## Exercise 33

DRY up by extracting player name once in `ServeHTTP` and passing it to both methods.

## Exercise 34

Create `server_integration_test.go` that tests `PlayerServer` with `InMemoryPlayerStore` - POST 3 wins then GET score.

(test will not pass)

## Exercise 35

Create `in_memory_player_store.go` with a working `InMemoryPlayerStore` that uses a map and a constructor.

## Exercise 36

Update `main.go` to use the constructor function.

---
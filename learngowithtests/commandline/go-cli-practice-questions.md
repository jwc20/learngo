# Go CLI Practice Questions - Command Line & Package Structure

## Setup Questions

### Q1: Project Structure Setup
**File to edit:** Project directory structure  
**Expected test result:** N/A (setup task)

Create the necessary directory structure for separating the web server application from domain code:
- Create a `cmd` directory in the project root
- Create a `webserver` directory inside `cmd`
- Move the existing `main.go` file into `cmd/webserver/`

Verify your structure matches:
```
.
├── cmd
│   └── webserver
│       └── main.go
├── file_system_store.go
├── file_system_store_test.go
├── league.go
├── server.go
├── server_integration_test.go
├── server_test.go
├── tape.go
└── tape_test.go
```

---

### Q2: Package Name Changes
**Files to edit:** All `.go` files except `cmd/webserver/main.go`  
**Expected test result:** Tests should fail with package/import errors

Change the package declaration in all domain code files from `package main` to `package poker`.

Files to update:
- `file_system_store.go`
- `league.go`
- `server.go`
- `tape.go`
- And all test files (except those in cmd/)

---

### Q3: Update Main Import
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (not running tests yet)

Update `main.go` to:
1. Keep `package main`
2. Import the poker package (adjust the path to match your module path)
3. Use `poker.` prefix for all functions and types from the poker package

Expected imports should include:
```go
import (
    "github.com/your-username/your-repo/path"
    "log"
    "net/http"
    "os"
)
```

Update function calls to use `poker.NewFileSystemPlayerStore()` and `poker.NewPlayerServer()`.

---

### Q4: Verify Tests Still Pass
**File to edit:** N/A (verification step)  
**Expected test result:** All tests should pass

Run `go test` from the project root directory to ensure all existing tests still pass after the refactoring.

---

## CLI Application Questions

### Q5: Create CLI Main Entry Point
**File to edit:** `cmd/cli/main.go` (new file)  
**Expected test result:** N/A (not running tests yet)

Create a new command-line application entry point:
1. Create the directory `cmd/cli/`
2. Create `main.go` inside it
3. Add a simple main function that prints "Let's play poker"

---

### Q6: Write First CLI Test
**File to edit:** `CLI_test.go` (new file in project root)  
**Expected test result:** Test should fail with "undefined: CLI"

Create a test file `CLI_test.go` with `package poker` that:
1. Creates a test function `TestCLI`
2. Creates a `StubPlayerStore`
3. Creates a `CLI` with the player store
4. Calls `cli.PlayPoker()`
5. Verifies that `len(playerStore.winCalls)` equals 1

---

### Q7: Create Minimal CLI Struct
**File to edit:** `CLI.go` (new file in project root)  
**Expected test result:** Test should fail with "expected a win call but didn't get any"

Create the minimal code needed to make the test run:
1. Create a `CLI` struct with a `playerStore PlayerStore` field
2. Add a `PlayPoker()` method that does nothing (empty body)

---

### Q8: Make First Test Pass
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

Implement the `PlayPoker()` method to record a win for any hardcoded player (e.g., "Cleo").

---

### Q9: Add Input Reading to Test
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "too many values in struct initializer"

Extend the test to simulate user input:
1. Create an `io.Reader` using `strings.NewReader("Chris wins\n")`
2. Pass it as a second argument when creating the CLI: `cli := &CLI{playerStore, in}`
3. Add assertions to check that the winner recorded is "Chris" (not just any win)

Expected assertion:
```go
got := playerStore.winCalls[0]
want := "Chris"
if got != want {
    t.Errorf("didn't record correct winner, got %q, want %q", got, want)
}
```

---

### Q10: Add io.Reader Dependency
**File to edit:** `CLI.go`  
**Expected test result:** Test should fail with "didn't record the correct winner, got 'Cleo', want 'Chris'"

Update the `CLI` struct to include an `io.Reader` field:
```go
type CLI struct {
    playerStore PlayerStore
    in          io.Reader
}
```

---

### Q11: Hardcode Different Winner
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

Update `PlayPoker()` to record a win for "Chris" instead of "Cleo" (still hardcoded, we'll fix this next).

---

### Q12: Create AssertPlayerWin Helper
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass (refactoring step)

Create a helper function in `server_test.go`:
```go
func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
    t.Helper()
    
    if len(store.winCalls) != 1 {
        t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
    }
    
    if store.winCalls[0] != winner {
        t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
    }
}
```

---

### Q13: Use Helper in CLI Test
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should pass (refactoring step)

Replace the win assertions in `CLI_test.go` with `assertPlayerWin(t, playerStore, "Chris")`.

---

### Q14: Add Second Test Case
**File to edit:** `CLI_test.go`  
**Expected test result:** Second test should fail with "did not store correct winner got 'Chris' want 'Cleo'"

Refactor `TestCLI` to use subtests:
1. Wrap existing test in `t.Run("record chris win from user input", func(t *testing.T) {...})`
2. Add a second subtest `t.Run("record cleo win from user input", func(t *testing.T) {...})` with input "Cleo wins\n"

---

### Q15: Implement Actual Input Reading
**File to edit:** `CLI.go`  
**Expected test result:** Both tests should pass

Implement proper input reading in `PlayPoker()`:
1. Create a `bufio.Scanner` from `cli.in`
2. Call `Scan()` to read a line
3. Use `Text()` to get the string
4. Extract the winner's name by removing " wins" from the input
5. Call `RecordWin()` with the extracted winner

Helper function to add:
```go
func extractWinner(userInput string) string {
    return strings.Replace(userInput, " wins", "", 1)
}
```

---

## Integration Questions

### Q16: Wire Up CLI in Main
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** Should fail to compile with "implicit assignment of unexported field"

Update `main()` to:
1. Print "Let's play poker" and "Type {Name} wins to record a win"
2. Open the database file
3. Create a `FileSystemPlayerStore`
4. Create and use a `CLI` with `game := poker.CLI{store, os.Stdin}`
5. Call `game.PlayPoker()`

---

## Package Design Questions

### Q17: Change Test Package Name
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should fail with multiple "undefined" errors

Change the package declaration from `package poker` to `package poker_test` to test only the public API.

---

### Q18: Create Testing Helpers File
**File to edit:** `testing.go` (new file in project root)  
**Expected test result:** Still failing (not imported yet)

Create a new file `testing.go` with `package poker` containing:
1. The `StubPlayerStore` struct and its methods (make it exported with capital S)
2. The `AssertPlayerWin` helper function (make it exported with capital A)
3. Any other test helpers you want to make public

---

### Q19: Update Test to Use Public API
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should fail with "implicit assignment of unexported field"

Update the test to use the exported types with the `poker.` prefix:
- `poker.StubPlayerStore{}`
- `poker.CLI{...}`
- `poker.AssertPlayerWin(...)`

---

### Q20: Create CLI Constructor
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass after next step

Create a constructor function and update the struct:
1. Change the `in` field type from `io.Reader` to `*bufio.Scanner`
2. Create a `NewCLI` function that takes a `PlayerStore` and `io.Reader`, returns `*CLI`
3. In the constructor, create the `bufio.Scanner` and return the initialized CLI

---

### Q21: Refactor PlayPoker with Scanner
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring step)

Refactor `PlayPoker()` to use the scanner directly:
1. Create a helper method `readLine()` that calls `cli.in.Scan()` and returns `cli.in.Text()`
2. Update `PlayPoker()` to use `userInput := cli.readLine()`

---

### Q22: Update Tests to Use Constructor
**File to edit:** `CLI_test.go`  
**Expected test result:** All tests should pass

Replace `&poker.CLI{playerStore, in}` with `poker.NewCLI(playerStore, in)` in both subtests.

---

### Q23: Update CLI Main to Use Constructor
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** N/A (integration step)

Replace `game := poker.CLI{store, os.Stdin}` with `game := poker.NewCLI(store, os.Stdin)`.

Run the CLI application manually and test that it works by typing "Bob wins".

---

## Final Refactoring Questions

### Q24: Create File Opening Helper
**File to edit:** `file_system_store.go`  
**Expected test result:** Tests should still pass (adding new function)

Create a helper function `FileSystemPlayerStoreFromFile` that:
1. Takes a file path string as parameter
2. Returns `(*FileSystemPlayerStore, func(), error)`
3. Opens the file with `os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)`
4. Creates the store using `NewFileSystemPlayerStore(db)`
5. Returns the store, a cleanup function that closes the file, and any error

---

### Q25: Refactor CLI Main to Use Helper
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** N/A (integration step)

Simplify main by:
1. Using `store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)`
2. Adding `defer close()` after error checking
3. Removing the manual file opening code

---

### Q26: Refactor Webserver Main to Use Helper
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (integration step)

Apply the same refactoring to the webserver's main:
1. Use `FileSystemPlayerStoreFromFile` instead of manual file operations
2. Add `defer close()`
3. Remove redundant code

---

## Verification Question

### Q27: Final Integration Test
**Files to test:** All  
**Expected test result:** All tests pass, both applications run correctly

1. Run `go test` from the project root - all tests should pass
2. Run the webserver with `go run cmd/webserver/main.go` and verify http://localhost:5000/league works
3. Run the CLI with `go run cmd/cli/main.go`, type "Alice wins" and verify the database is updated
4. Check that both applications share the same database file

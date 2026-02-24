# Go CLI Practice Questions with Answers - Command Line & Package Structure

## Setup Questions

### Q1: Project Structure Setup
**File to edit:** Project directory structure  
**Expected test result:** N/A (setup task)

**Question:**
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

**Answer:**
```bash
# In your project root directory
mkdir -p cmd/webserver
mv main.go cmd/webserver/

# Verify the structure
tree  # or ls -R if tree is not installed
```

---

### Q2: Package Name Changes
**Files to edit:** All `.go` files except `cmd/webserver/main.go`  
**Expected test result:** Tests should fail with package/import errors

**Question:**
Change the package declaration in all domain code files from `package main` to `package poker`.

Files to update:
- `file_system_store.go`
- `league.go`
- `server.go`
- `tape.go`
- And all test files (except those in cmd/)

**Answer:**
In each of these files, change the first line from:
```go
package main
```
to:
```go
package poker
```

Files to update:
- `file_system_store.go`
- `file_system_store_test.go`
- `league.go`
- `server.go`
- `server_test.go`
- `server_integration_test.go`
- `tape.go`
- `tape_test.go`

---

### Q3: Update Main Import
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (not running tests yet)

**Question:**
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

**Answer:**
```go
// cmd/webserver/main.go
package main

import (
	"github.com/your-username/your-repo"  // adjust to your module path
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := poker.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", server))
}
```

---

### Q4: Verify Tests Still Pass
**File to edit:** N/A (verification step)  
**Expected test result:** All tests should pass

**Question:**
Run `go test` from the project root directory to ensure all existing tests still pass after the refactoring.

**Answer:**
```bash
# From project root
go test

# Or for verbose output
go test -v
```

All existing tests should pass. If they don't, check that you've updated all package declarations correctly.

---

## CLI Application Questions

### Q5: Create CLI Main Entry Point
**File to edit:** `cmd/cli/main.go` (new file)  
**Expected test result:** N/A (not running tests yet)

**Question:**
Create a new command-line application entry point:
1. Create the directory `cmd/cli/`
2. Create `main.go` inside it
3. Add a simple main function that prints "Let's play poker"

**Answer:**
```bash
mkdir -p cmd/cli
```

```go
// cmd/cli/main.go
package main

import "fmt"

func main() {
	fmt.Println("Let's play poker")
}
```

Test it:
```bash
go run cmd/cli/main.go
```

---

### Q6: Write First CLI Test
**File to edit:** `CLI_test.go` (new file in project root)  
**Expected test result:** Test should fail with "undefined: CLI"

**Question:**
Create a test file `CLI_test.go` with `package poker` that:
1. Creates a test function `TestCLI`
2. Creates a `StubPlayerStore`
3. Creates a `CLI` with the player store
4. Calls `cli.PlayPoker()`
5. Verifies that `len(playerStore.winCalls)` equals 1

**Answer:**
```go
// CLI_test.go
package poker

import "testing"

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &CLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}
}
```

Run the test:
```bash
go test
```

Expected error: `undefined: CLI`

---

### Q7: Create Minimal CLI Struct
**File to edit:** `CLI.go` (new file in project root)  
**Expected test result:** Test should fail with "expected a win call but didn't get any"

**Question:**
Create the minimal code needed to make the test run:
1. Create a `CLI` struct with a `playerStore PlayerStore` field
2. Add a `PlayPoker()` method that does nothing (empty body)

**Answer:**
```go
// CLI.go
package poker

type CLI struct {
	playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

Run the test:
```bash
go test
```

Expected failure message: "expected a win call but didn't get any"

---

### Q8: Make First Test Pass
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

**Question:**
Implement the `PlayPoker()` method to record a win for any hardcoded player (e.g., "Cleo").

**Answer:**
```go
// CLI.go
package poker

type CLI struct {
	playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
```

Run the test:
```bash
go test
```

Test should pass.

---

### Q9: Add Input Reading to Test
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should fail with "too many values in struct initializer"

**Question:**
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

**Answer:**
```go
// CLI_test.go
package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}

	got := playerStore.winCalls[0]
	want := "Chris"

	if got != want {
		t.Errorf("didn't record correct winner, got %q, want %q", got, want)
	}
}
```

Run the test:
```bash
go test
```

Expected error: "too many values in struct initializer"

---

### Q10: Add io.Reader Dependency
**File to edit:** `CLI.go`  
**Expected test result:** Test should fail with "didn't record the correct winner, got 'Cleo', want 'Chris'"

**Question:**
Update the `CLI` struct to include an `io.Reader` field:
```go
type CLI struct {
    playerStore PlayerStore
    in          io.Reader
}
```

**Answer:**
```go
// CLI.go
package poker

import "io"

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
```

Run the test:
```bash
go test
```

Expected failure: "didn't record the correct winner, got 'Cleo', want 'Chris'"

---

### Q11: Hardcode Different Winner
**File to edit:** `CLI.go`  
**Expected test result:** Test should pass

**Question:**
Update `PlayPoker()` to record a win for "Chris" instead of "Cleo" (still hardcoded, we'll fix this next).

**Answer:**
```go
// CLI.go
package poker

import "io"

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Chris")
}
```

Run the test:
```bash
go test
```

Test should pass.

---

### Q12: Create AssertPlayerWin Helper
**File to edit:** `server_test.go`  
**Expected test result:** Tests should pass (refactoring step)

**Question:**
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

**Answer:**
Add this function to `server_test.go`:

```go
// server_test.go
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

Run tests to ensure nothing broke:
```bash
go test
```

---

### Q13: Use Helper in CLI Test
**File to edit:** `CLI_test.go`  
**Expected test result:** Test should pass (refactoring step)

**Question:**
Replace the win assertions in `CLI_test.go` with `assertPlayerWin(t, playerStore, "Chris")`.

**Answer:**
```go
// CLI_test.go
package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	assertPlayerWin(t, playerStore, "Chris")
}
```

Run the test:
```bash
go test
```

Test should still pass.

---

### Q14: Add Second Test Case
**File to edit:** `CLI_test.go`  
**Expected test result:** Second test should fail with "did not store correct winner got 'Chris' want 'Cleo'"

**Question:**
Refactor `TestCLI` to use subtests:
1. Wrap existing test in `t.Run("record chris win from user input", func(t *testing.T) {...})`
2. Add a second subtest `t.Run("record cleo win from user input", func(t *testing.T) {...})` with input "Cleo wins\n"

**Answer:**
```go
// CLI_test.go
package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Cleo")
	})

}
```

Run the test:
```bash
go test
```

Expected failure: "did not store correct winner got 'Chris' want 'Cleo'" on the second subtest.

---

### Q15: Implement Actual Input Reading
**File to edit:** `CLI.go`  
**Expected test result:** Both tests should pass

**Question:**
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

**Answer:**
```go
// CLI.go
package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
```

Run the test:
```bash
go test
```

Both subtests should pass.

---

## Integration Questions

### Q16: Wire Up CLI in Main
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** Should fail to compile with "implicit assignment of unexported field"

**Question:**
Update `main()` to:
1. Print "Let's play poker" and "Type {Name} wins to record a win"
2. Open the database file
3. Create a `FileSystemPlayerStore`
4. Create and use a `CLI` with `game := poker.CLI{store, os.Stdin}`
5. Call `game.PlayPoker()`

**Answer:**
```go
// cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/your-username/your-repo"  // adjust to your module path
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.CLI{store, os.Stdin}
	game.PlayPoker()
}
```

Try to compile:
```bash
go run cmd/cli/main.go
```

Expected error: "implicit assignment of unexported field 'playerStore' in poker.CLI literal"

---

## Package Design Questions

### Q17: Change Test Package Name
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should fail with multiple "undefined" errors

**Question:**
Change the package declaration from `package poker` to `package poker_test` to test only the public API.

**Answer:**
```go
// CLI_test.go
package poker_test  // Changed from: package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	// ... rest of the test remains the same
}
```

Run the test:
```bash
go test
```

Expected errors: 
- `undefined: StubPlayerStore`
- `undefined: CLI`
- `undefined: assertPlayerWin`

---

### Q18: Create Testing Helpers File
**File to edit:** `testing.go` (new file in project root)  
**Expected test result:** Still failing (not imported yet)

**Question:**
Create a new file `testing.go` with `package poker` containing:
1. The `StubPlayerStore` struct and its methods (make it exported with capital S)
2. The `AssertPlayerWin` helper function (make it exported with capital A)
3. Any other test helpers you want to make public

**Answer:**
```go
// testing.go
package poker

import "testing"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}
```

Note: You can also move other test helpers from `server_test.go` to this file and export them.

---

### Q19: Update Test to Use Public API
**File to edit:** `CLI_test.go`  
**Expected test result:** Tests should fail with "implicit assignment of unexported field"

**Question:**
Update the test to use the exported types with the `poker.` prefix:
- `poker.StubPlayerStore{}`
- `poker.CLI{...}`
- `poker.AssertPlayerWin(...)`

**Answer:**
```go
// CLI_test.go
package poker_test

import (
	"strings"
	"testing"
	
	"github.com/your-username/your-repo"  // adjust to your module path
)

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.CLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.CLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

}
```

Run the test:
```bash
go test
```

Expected error: "implicit assignment of unexported field 'playerStore' in poker.CLI literal"

---

### Q20: Create CLI Constructor
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass after next step

**Question:**
Create a constructor function and update the struct:
1. Change the `in` field type from `io.Reader` to `*bufio.Scanner`
2. Create a `NewCLI` function that takes a `PlayerStore` and `io.Reader`, returns `*CLI`
3. In the constructor, create the `bufio.Scanner` and return the initialized CLI

**Answer:**
```go
// CLI.go
package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
```

---

### Q21: Refactor PlayPoker with Scanner
**File to edit:** `CLI.go`  
**Expected test result:** Tests should pass (refactoring step)

**Question:**
Refactor `PlayPoker()` to use the scanner directly:
1. Create a helper method `readLine()` that calls `cli.in.Scan()` and returns `cli.in.Text()`
2. Update `PlayPoker()` to use `userInput := cli.readLine()`

**Answer:**
```go
// CLI.go
package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}

func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
```

---

### Q22: Update Tests to Use Constructor
**File to edit:** `CLI_test.go`  
**Expected test result:** All tests should pass

**Question:**
Replace `&poker.CLI{playerStore, in}` with `poker.NewCLI(playerStore, in)` in both subtests.

**Answer:**
```go
// CLI_test.go
package poker_test

import (
	"strings"
	"testing"
	
	"github.com/your-username/your-repo"  // adjust to your module path
)

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

}
```

Run the test:
```bash
go test
```

All tests should pass.

---

### Q23: Update CLI Main to Use Constructor
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** N/A (integration step)

**Question:**
Replace `game := poker.CLI{store, os.Stdin}` with `game := poker.NewCLI(store, os.Stdin)`.

Run the CLI application manually and test that it works by typing "Bob wins".

**Answer:**
```go
// cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/your-username/your-repo"  // adjust to your module path
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}
```

Test it:
```bash
go run cmd/cli/main.go
# Type: Bob wins
# Press Enter
```

Verify "Bob" has a win recorded in game.db.json.

---

## Final Refactoring Questions

### Q24: Create File Opening Helper
**File to edit:** `file_system_store.go`  
**Expected test result:** Tests should still pass (adding new function)

**Question:**
Create a helper function `FileSystemPlayerStoreFromFile` that:
1. Takes a file path string as parameter
2. Returns `(*FileSystemPlayerStore, func(), error)`
3. Opens the file with `os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)`
4. Creates the store using `NewFileSystemPlayerStore(db)`
5. Returns the store, a cleanup function that closes the file, and any error

**Answer:**
Add this function to `file_system_store.go`:

```go
// file_system_store.go
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}
```

Run tests to verify nothing broke:
```bash
go test
```

---

### Q25: Refactor CLI Main to Use Helper
**File to edit:** `cmd/cli/main.go`  
**Expected test result:** N/A (integration step)

**Question:**
Simplify main by:
1. Using `store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)`
2. Adding `defer close()` after error checking
3. Removing the manual file opening code

**Answer:**
```go
// cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/your-username/your-repo"  // adjust to your module path
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
```

Test it:
```bash
go run cmd/cli/main.go
```

---

### Q26: Refactor Webserver Main to Use Helper
**File to edit:** `cmd/webserver/main.go`  
**Expected test result:** N/A (integration step)

**Question:**
Apply the same refactoring to the webserver's main:
1. Use `FileSystemPlayerStoreFromFile` instead of manual file operations
2. Add `defer close()`
3. Remove redundant code

**Answer:**
```go
// cmd/webserver/main.go
package main

import (
	"github.com/your-username/your-repo"  // adjust to your module path
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

Test it:
```bash
go run cmd/webserver/main.go
# In another terminal:
curl http://localhost:5000/league
```

---

## Verification Question

### Q27: Final Integration Test
**Files to test:** All  
**Expected test result:** All tests pass, both applications run correctly

**Question:**
1. Run `go test` from the project root - all tests should pass
2. Run the webserver with `go run cmd/webserver/main.go` and verify http://localhost:5000/league works
3. Run the CLI with `go run cmd/cli/main.go`, type "Alice wins" and verify the database is updated
4. Check that both applications share the same database file

**Answer:**
```bash
# 1. Run all tests
go test
# Should output: PASS

# 2. Start the webserver (in one terminal)
go run cmd/webserver/main.go
# In another terminal:
curl http://localhost:5000/league
# Should see league data

# 3. Record a win via CLI (after stopping webserver)
go run cmd/cli/main.go
# Type: Alice wins
# Press Enter

# 4. Verify database is shared
# Start webserver again:
go run cmd/webserver/main.go
# Check league:
curl http://localhost:5000/league
# Should see Alice in the league

# You can also check the game.db.json file directly:
cat game.db.json
```

Expected: Both applications successfully read from and write to the same `game.db.json` file, demonstrating proper separation of concerns and code reuse.

# IO and Sorting - TDD Practice Q/A

## Starting Point
You should have the completed files from the JSON, Routing and Embedding chapter:
- `server.go`
- `server_test.go`
- `in_memory_player_store.go`
- `server_integration_test.go`
- `main.go`
- `league.go` (if created in previous chapter, otherwise will be created here)

---

## Exercise 1: Write Test for FileSystemStore GetLeague
**File to create:** `file_system_store_test.go`
**Expected result:** Compilation error - `undefined: FileSystemPlayerStore`

**Question:** Create a new test file `file_system_store_test.go` with a test for `FileSystemPlayerStore` that reads league data from a `strings.Reader` containing JSON and calls `GetLeague()`.

**Answer:**
```go
func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)
	})
}
```

---

## Exercise 2: Create Empty FileSystemPlayerStore
**File to create:** `file_system_store.go`
**Expected result:** Compilation error - `too many values in struct initializer` and `GetLeague undefined`

**Question:** Create a new file `file_system_store.go` with an empty `FileSystemPlayerStore` struct.

**Answer:**
```go
type FileSystemPlayerStore struct{}
```

---

## Exercise 3: Add database Field and GetLeague Method
**File to edit:** `file_system_store.go`
**Expected result:** Test fails with `got [] want [{Cleo 10} {Chris 33}]`

**Question:** Add a `database` field of type `io.Reader` to `FileSystemPlayerStore` and add a `GetLeague` method that returns `nil`.

**Answer:**
```go
type FileSystemPlayerStore struct {
	database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	return nil
}
```

---

## Exercise 4: Implement GetLeague with JSON Decoding
**File to edit:** `file_system_store.go`
**Expected result:** Test passes

**Question:** Implement `GetLeague` to decode JSON from the database reader and return the league.

**Answer:**
```go
func (f *FileSystemPlayerStore) GetLeague() []Player {
	var league []Player
	json.NewDecoder(f.database).Decode(&league)
	return league
}
```

---

## Exercise 5: Refactor - Create NewLeague Helper
**File to create/edit:** `league.go`
**Expected result:** Tests pass (refactor)

**Question:** Create a `NewLeague` function in `league.go` that takes an `io.Reader` and returns `([]Player, error)`. Update `FileSystemPlayerStore.GetLeague` to use it.

**Answer:**

In `league.go`:
```go
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}
```

In `file_system_store.go`:
```go
func (f *FileSystemPlayerStore) GetLeague() []Player {
	league, _ := NewLeague(f.database)
	return league
}
```

---

## Exercise 6: Add Test for Reading League Twice
**File to edit:** `file_system_store_test.go`
**Expected result:** Test fails (second read returns empty)

**Question:** Add a second call to `store.GetLeague()` at the end of the "league from a reader" test and assert it returns the same data.

**Answer:**
```go
// read again
got = store.GetLeague()
assertLeague(t, got, want)
```

---

## Exercise 7: Change to ReadSeeker and Add Seek
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Change the `database` field type from `io.Reader` to `io.ReadSeeker`. Add a `Seek` call to the beginning of `GetLeague` to reset the reader position.

**Answer:**
```go
type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, io.SeekStart)
	league, _ := NewLeague(f.database)
	return league
}
```

---

## Exercise 8: Write Test for GetPlayerScore
**File to edit:** `file_system_store_test.go`
**Expected result:** Compilation error - `GetPlayerScore undefined`

**Question:** Add a test for `GetPlayerScore` that expects Chris's score to be 33.

**Answer:**
```go
t.Run("get player score", func(t *testing.T) {
	database := strings.NewReader(`[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)

	store := FileSystemPlayerStore{database}

	got := store.GetPlayerScore("Chris")

	want := 33

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
})
```

---

## Exercise 9: Add Empty GetPlayerScore Method
**File to edit:** `file_system_store.go`
**Expected result:** Test fails with `got 0 want 33`

**Question:** Add an empty `GetPlayerScore` method that returns 0.

**Answer:**
```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	return 0
}
```

---

## Exercise 10: Implement GetPlayerScore
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Implement `GetPlayerScore` to iterate over the league and return the player's wins.

**Answer:**
```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	var wins int

	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			break
		}
	}

	return wins
}
```

---

## Exercise 11: Refactor Test - Add assertScoreEquals Helper
**File to edit:** `file_system_store_test.go`
**Expected result:** Tests pass (refactor)

**Question:** Create an `assertScoreEquals` helper function and use it in the test.

**Answer:**
```go
func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

Update test:
```go
got := store.GetPlayerScore("Chris")
want := 33
assertScoreEquals(t, got, want)
```

---

## Exercise 12: Change to ReadWriteSeeker
**File to edit:** `file_system_store.go`
**Expected result:** Compilation error - `strings.Reader does not implement io.ReadWriteSeeker`

**Question:** Change the `database` field type to `io.ReadWriteSeeker` to prepare for writing.

**Answer:**
```go
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}
```

---

## Exercise 13: Create Temp File Helper
**File to edit:** `file_system_store_test.go`
**Expected result:** Tests pass after updating to use temp files

**Question:** Create a `createTempFile` helper that creates a temporary file with initial data and returns a cleanup function. Update all tests to use it.

**Answer:**
```go
func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
```

Update tests:
```go
t.Run("league from a reader", func(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store := FileSystemPlayerStore{database}

	got := store.GetLeague()

	want := []Player{
		{"Cleo", 10},
		{"Chris", 33},
	}

	assertLeague(t, got, want)

	// read again
	got = store.GetLeague()
	assertLeague(t, got, want)
})

t.Run("get player score", func(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store := FileSystemPlayerStore{database}

	got := store.GetPlayerScore("Chris")
	want := 33
	assertScoreEquals(t, got, want)
})
```

---

## Exercise 14: Write Test for RecordWin (Existing Player)
**File to edit:** `file_system_store_test.go`
**Expected result:** Compilation error - `RecordWin undefined`

**Question:** Add a test for `RecordWin` that records a win for Chris and expects his score to increase from 33 to 34.

**Answer:**
```go
t.Run("store wins for existing players", func(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store := FileSystemPlayerStore{database}

	store.RecordWin("Chris")

	got := store.GetPlayerScore("Chris")
	want := 34
	assertScoreEquals(t, got, want)
})
```

---

## Exercise 15: Add Empty RecordWin Method
**File to edit:** `file_system_store.go`
**Expected result:** Test fails with `got 33 want 34`

**Question:** Add an empty `RecordWin` method.

**Answer:**
```go
func (f *FileSystemPlayerStore) RecordWin(name string) {

}
```

---

## Exercise 16: Implement RecordWin
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Implement `RecordWin` to find the player, increment their wins, and write the updated league back to the file.

**Answer:**
```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
```

---

## Exercise 17: Refactor - Create League Type with Find Method
**File to edit:** `league.go`
**Expected result:** Tests pass (refactor)

**Question:** Create a `League` type alias for `[]Player` and add a `Find` method that returns a pointer to a player by name.

**Answer:**
```go
type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
```

---

## Exercise 18: Update PlayerStore Interface to Return League
**File to edit:** `server.go`
**Expected result:** Compilation errors in tests (easy to fix)

**Question:** Change the `GetLeague()` return type in `PlayerStore` interface from `[]Player` to `League`.

**Answer:**
```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}
```

---

## Exercise 19: Refactor FileSystemPlayerStore to Use League.Find
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass (refactor)

**Question:** Update `GetPlayerScore` and `RecordWin` to use `League.Find` method.

**Answer:**
```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.GetLeague().Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
```

---

## Exercise 20: Write Test for RecordWin (New Player)
**File to edit:** `file_system_store_test.go`
**Expected result:** Test fails with `got 0 want 1`

**Question:** Add a test for recording a win for a new player "Pepper" who doesn't exist in the database.

**Answer:**
```go
t.Run("store wins for new players", func(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store := FileSystemPlayerStore{database}

	store.RecordWin("Pepper")

	got := store.GetPlayerScore("Pepper")
	want := 1
	assertScoreEquals(t, got, want)
})
```

---

## Exercise 21: Handle New Player in RecordWin
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Update `RecordWin` to append a new player with 1 win if the player doesn't exist.

**Answer:**
```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
```

---

## Exercise 22: Update Integration Test to Use FileSystemPlayerStore
**File to edit:** `server_integration_test.go`
**Expected result:** Tests pass

**Question:** Replace `InMemoryPlayerStore` with `FileSystemPlayerStore` in the integration test.

**Answer:**
```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()
	store := &FileSystemPlayerStore{database}

	// rest of test...
}
```

---

## Exercise 23: Delete InMemoryPlayerStore and Update main.go
**Files to edit:** Delete `in_memory_player_store.go`, edit `main.go`
**Expected result:** Application compiles and runs with file-based storage

**Question:** Delete `in_memory_player_store.go` and update `main.go` to use `FileSystemPlayerStore` with a real file.

**Answer:**
```go
package main

import (
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

	store := &FileSystemPlayerStore{db}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

---

## Exercise 24: Create Constructor with Cached League
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass (performance refactor)

**Question:** Add a `league` field to `FileSystemPlayerStore` and create a `NewFileSystemPlayerStore` constructor that reads the league once at startup. Update `GetLeague`, `GetPlayerScore`, and `RecordWin` to use the cached league.

**Answer:**
```go
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: database,
		league:   league,
	}
}

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(f.league)
}
```

---

## Exercise 25: Update Tests to Use Constructor
**Files to edit:** `file_system_store_test.go`, `server_integration_test.go`, `main.go`
**Expected result:** Tests pass

**Question:** Update all tests and main.go to use `NewFileSystemPlayerStore` instead of struct literal.

**Answer:**
```go
store := NewFileSystemPlayerStore(database)
```

---

## Exercise 26: Create tape Type for Safe Writes
**File to create:** `tape.go`
**Expected result:** Code compiles

**Question:** Create a new file `tape.go` with a `tape` struct that wraps an `io.ReadWriteSeeker` and implements `Write` by seeking to start before writing.

**Answer:**
```go
package main

import "io"

type tape struct {
	file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
```

---

## Exercise 27: Update FileSystemPlayerStore to Use tape
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Change `FileSystemPlayerStore.database` to `io.Writer` and update the constructor to wrap the file in a `tape`. Remove the `Seek` call from `RecordWin`.

**Answer:**
```go
type FileSystemPlayerStore struct {
	database io.Writer
	league   League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := NewLeague(database)

	return &FileSystemPlayerStore{
		database: &tape{database},
		league:   league,
	}
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	json.NewEncoder(f.database).Encode(f.league)
}
```

---

## Exercise 28: Write Test for tape Truncation Issue
**File to create:** `tape_test.go`
**Expected result:** Test fails with `got 'abc45' want 'abc'`

**Question:** Create `tape_test.go` with a test that writes "abc" to a file containing "12345" and asserts the file only contains "abc".

**Answer:**
```go
func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

---

## Exercise 29: Fix tape to Use os.File with Truncate
**File to edit:** `tape.go`
**Expected result:** Tests pass

**Question:** Change `tape.file` to `*os.File` and call `Truncate(0)` before seeking and writing.

**Answer:**
```go
type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
```

---

## Exercise 30: Update Constructor to Accept *os.File
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Update `NewFileSystemPlayerStore` to accept `*os.File` instead of `io.ReadWriteSeeker`. Also store a `*json.Encoder` instead of creating one each time.

**Answer:**
```go
type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
	file.Seek(0, io.SeekStart)
	league, _ := NewLeague(file)

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
}
```

---

## Exercise 31: Update createTempFile to Return *os.File
**File to edit:** `file_system_store_test.go`
**Expected result:** Tests pass

**Question:** Update `createTempFile` to return `*os.File` instead of `io.ReadWriteSeeker`.

**Answer:**
```go
func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
```

---

## Exercise 32: Add Error Handling to Constructor
**File to edit:** `file_system_store.go`
**Expected result:** Compilation errors - multiple-value in single-value context

**Question:** Update `NewFileSystemPlayerStore` to return `(*FileSystemPlayerStore, error)` and handle the error from `NewLeague`.

**Answer:**
```go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0, io.SeekStart)
	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
```

---

## Exercise 33: Add assertNoError Helper and Fix Tests
**Files to edit:** `file_system_store_test.go`, `server_integration_test.go`, `main.go`
**Expected result:** Test fails with `problem parsing league, EOF`

**Question:** Create an `assertNoError` helper. Update all calls to `NewFileSystemPlayerStore` to handle the error. In `main.go`, log fatal on error.

**Answer:**

Helper:
```go
func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
```

In tests:
```go
store, err := NewFileSystemPlayerStore(database)
assertNoError(t, err)
```

In main.go:
```go
store, err := NewFileSystemPlayerStore(db)

if err != nil {
	log.Fatalf("problem creating file system player store, %v ", err)
}
```

---

## Exercise 34: Fix Integration Test with Valid JSON
**File to edit:** `server_integration_test.go`
**Expected result:** Tests pass

**Question:** Update the integration test to pass `[]` (empty JSON array) instead of empty string to `createTempFile`.

**Answer:**
```go
database, cleanDatabase := createTempFile(t, `[]`)
```

---

## Exercise 35: Write Test for Empty File Handling
**File to edit:** `file_system_store_test.go`
**Expected result:** Test fails with `problem parsing league, EOF`

**Question:** Add a test that creates a store from an empty file and expects no error.

**Answer:**
```go
t.Run("works with an empty file", func(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()

	_, err := NewFileSystemPlayerStore(database)

	assertNoError(t, err)
})
```

---

## Exercise 36: Handle Empty File in Constructor
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass

**Question:** Update `NewFileSystemPlayerStore` to check if the file is empty using `file.Stat()` and write `[]` if it is.

**Answer:**
```go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
```

---

## Exercise 37: Refactor - Extract initialisePlayerDBFile
**File to edit:** `file_system_store.go`
**Expected result:** Tests pass (refactor)

**Question:** Extract the file initialization logic into a separate `initialisePlayerDBFile` function.

**Answer:**
```go
func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
```

---

## Exercise 38: Write Test for Sorted League
**File to edit:** `file_system_store_test.go`
**Expected result:** Test fails with `got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]`

**Question:** Add a test that expects `GetLeague` to return players sorted by wins (highest first).

**Answer:**
```go
t.Run("league sorted", func(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)

	assertNoError(t, err)

	got := store.GetLeague()

	want := League{
		{"Chris", 33},
		{"Cleo", 10},
	}

	assertLeague(t, got, want)

	// read again
	got = store.GetLeague()
	assertLeague(t, got, want)
})
```

---

## Exercise 39: Implement Sorting in GetLeague
**File to edit:** `file_system_store.go`
**Expected result:** All tests pass

**Question:** Update `GetLeague` to sort the league by wins in descending order using `sort.Slice`.

**Answer:**
```go
func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}
```

---

# Questions Only (No Answers)

1. **[file_system_store_test.go] [Compilation error: undefined FileSystemPlayerStore]** Create a new test file with a test for `FileSystemPlayerStore` that reads league data from a `strings.Reader`.

2. **[file_system_store.go] [Compilation error: too many values, GetLeague undefined]** Create an empty `FileSystemPlayerStore` struct.

3. **[file_system_store.go] [Test fails: got [] want [{Cleo 10} {Chris 33}]]** Add a `database io.Reader` field and a `GetLeague` method that returns `nil`.

4. **[file_system_store.go] [Tests pass]** Implement `GetLeague` to decode JSON from the database reader.

5. **[league.go, file_system_store.go] [Tests pass - refactor]** Create a `NewLeague` helper function and use it in `GetLeague`.

6. **[file_system_store_test.go] [Test fails on second read]** Add a second `GetLeague()` call to test reading twice.

7. **[file_system_store.go] [Tests pass]** Change to `io.ReadSeeker` and add `Seek` to reset reader position.

8. **[file_system_store_test.go] [Compilation error: GetPlayerScore undefined]** Add a test for `GetPlayerScore`.

9. **[file_system_store.go] [Test fails: got 0 want 33]** Add an empty `GetPlayerScore` method returning 0.

10. **[file_system_store.go] [Tests pass]** Implement `GetPlayerScore` by iterating over the league.

11. **[file_system_store_test.go] [Tests pass - refactor]** Create `assertScoreEquals` helper.

12. **[file_system_store.go] [Compilation error: strings.Reader doesn't implement ReadWriteSeeker]** Change database type to `io.ReadWriteSeeker`.

13. **[file_system_store_test.go] [Tests pass]** Create `createTempFile` helper and update tests to use temp files.

14. **[file_system_store_test.go] [Compilation error: RecordWin undefined]** Add test for `RecordWin` on existing player.

15. **[file_system_store.go] [Test fails: got 33 want 34]** Add empty `RecordWin` method.

16. **[file_system_store.go] [Tests pass]** Implement `RecordWin` to update and persist the league.

17. **[league.go] [Tests pass - refactor]** Create `League` type with `Find` method.

18. **[server.go] [Compilation errors]** Change `GetLeague()` return type to `League` in interface.

19. **[file_system_store.go] [Tests pass - refactor]** Update methods to use `League.Find`.

20. **[file_system_store_test.go] [Test fails: got 0 want 1]** Add test for `RecordWin` on new player.

21. **[file_system_store.go] [Tests pass]** Handle new player in `RecordWin` by appending.

22. **[server_integration_test.go] [Tests pass]** Replace `InMemoryPlayerStore` with `FileSystemPlayerStore`.

23. **[Delete in_memory_player_store.go, edit main.go] [App runs]** Delete old store and update main to use file-based storage.

24. **[file_system_store.go] [Tests pass - performance refactor]** Add constructor with cached league field.

25. **[file_system_store_test.go, server_integration_test.go, main.go] [Tests pass]** Update all code to use constructor.

26. **[tape.go] [Compiles]** Create `tape` type that seeks before writing.

27. **[file_system_store.go] [Tests pass]** Update store to use `tape` and `io.Writer`.

28. **[tape_test.go] [Test fails: got 'abc45' want 'abc']** Write test exposing truncation issue.

29. **[tape.go] [Tests pass]** Fix tape to use `*os.File` with `Truncate`.

30. **[file_system_store.go] [Tests pass]** Update constructor to accept `*os.File` and store `*json.Encoder`.

31. **[file_system_store_test.go] [Tests pass]** Update `createTempFile` to return `*os.File`.

32. **[file_system_store.go] [Compilation errors]** Add error handling to constructor.

33. **[file_system_store_test.go, server_integration_test.go, main.go] [Test fails: EOF error]** Create `assertNoError` helper and handle errors everywhere.

34. **[server_integration_test.go] [Tests pass]** Fix integration test with valid empty JSON array.

35. **[file_system_store_test.go] [Test fails: EOF error]** Add test for empty file handling.

36. **[file_system_store.go] [Tests pass]** Handle empty file by writing `[]`.

37. **[file_system_store.go] [Tests pass - refactor]** Extract `initialisePlayerDBFile` function.

38. **[file_system_store_test.go] [Test fails: wrong order]** Add test for sorted league output.

39. **[file_system_store.go] [All tests pass]** Implement sorting with `sort.Slice`.

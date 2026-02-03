# IO and Sorting - TDD Practice Questions

## Starting Point

You should have the completed files from the JSON, Routing and Embedding chapter:

- `server.go` - HTTP server with PlayerServer
- `server_test.go` - Server tests
- `in_memory_player_store.go` - In-memory implementation of PlayerStore
- `server_integration_test.go` - Integration tests
- `main.go` - Application entry point
- `league.go` - Player struct (if created in previous chapter)

## Goal

By the end of this chapter, you will replace the in-memory store with a file-based store that persists data as JSON. You'll learn about `io.Reader`, `io.Writer`, `io.Seeker`, working with files, and sorting.

---

## Part 1: Reading from Files

This section introduces the FileSystemPlayerStore and implements reading league data from an io.Reader.

### Exercise 1: Write the First Test for FileSystemPlayerStore

**File:** `file_system_store_test.go` (create new) **TDD Phase:** Write the test first **Expected Result:** Compilation error - `undefined: FileSystemPlayerStore`

#### Why This Step?

We're following TDD - before writing any production code, we write a failing test. We want a new `FileSystemPlayerStore` that can read league data from a data source. We'll start simple by using `strings.NewReader` which implements `io.Reader`.

#### Steps

1. Create a new file `file_system_store_test.go`
2. Write a test function `TestFileSystemStore` with a subtest "league from a reader"
3. Create a `strings.Reader` containing JSON player data
4. Instantiate a `FileSystemPlayerStore` passing the reader as a struct field
5. Call `GetLeague()` and assert it returns the expected players
6. Run the test and confirm you get the compilation error

#### Concepts

- `strings.NewReader` creates an `io.Reader` from a string - perfect for testing without real files
- We're designing the API we want before implementing it

---

### Exercise 2: Create the Empty Struct

**File:** `file_system_store.go` (create new) **TDD Phase:** Write minimal code to compile **Expected Result:** Compilation error - `too many values in struct initializer` and `GetLeague undefined`

#### Why This Step?

The compiler tells us `FileSystemPlayerStore` is undefined. We write just enough code to move past this error, which will reveal the next one.

#### Steps

1. Create a new file `file_system_store.go`
2. Define an empty `FileSystemPlayerStore` struct with no fields
3. Run the test and observe the new compilation errors

#### Concepts

- In TDD, we take tiny steps - just enough to change the error message
- The compiler guides us to what we need to implement next

---

### Exercise 3: Add the Database Field and Empty Method

**File:** `file_system_store.go` **TDD Phase:** Write minimal code to compile and run **Expected Result:** Test fails with `got [] want [{Cleo 10} {Chris 33}]`

#### Why This Step?

The errors tell us we're passing a value to the struct (the reader) but it has no fields, and `GetLeague` doesn't exist. We need to add both, but keep the implementation minimal.

#### Steps

1. Add a `database` field of type `io.Reader` to the struct
2. Add a `GetLeague` method that takes no arguments and returns `[]Player`
3. Have the method return `nil` for now
4. Run the test - it should now compile and run, but fail

#### Concepts

- `io.Reader` is a minimal interface that just requires `Read(p []byte) (n int, err error)`
- Returning `nil` is the simplest implementation - we want to see the test fail for the right reason

---

### Exercise 4: Implement GetLeague with JSON Decoding

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** Test passes ✓

#### Why This Step?

Now we implement the actual logic. We need to read JSON from the `database` reader and decode it into a slice of players.

#### Steps

1. In `GetLeague`, create a `var league []Player`
2. Use `json.NewDecoder(f.database).Decode(&league)` to parse the JSON
3. Return the league
4. Run the test - it should pass

#### Concepts

- `json.NewDecoder` takes any `io.Reader` and creates a decoder
- We're ignoring the error for now - we'll handle it later (pragmatic TDD)

---

### Exercise 5: Refactor - Extract NewLeague Helper

**Files:** `league.go`, `file_system_store.go` **TDD Phase:** Refactor (tests must stay green) **Expected Result:** Tests pass ✓

#### Why This Step?

We've done this JSON decoding before in our server tests. DRY (Don't Repeat Yourself) - extract this into a reusable function. This is a refactor, so tests must pass before and after.

#### Steps

1. In `league.go`, create a function `NewLeague(rdr io.Reader) ([]Player, error)`
2. Move the JSON decoding logic there, but also handle and wrap the error
3. Update `GetLeague` in `file_system_store.go` to call `NewLeague`
4. Ignore the error for now with `_`
5. Run tests to ensure they still pass
6. Update `getLeagueFromResponse` in `server_test.go` to also use `NewLeague` (if applicable)

#### Concepts

- Refactoring means changing code structure without changing behavior
- Wrapping errors with context (e.g., `fmt.Errorf("problem parsing league, %v", err)`) makes debugging easier

---

## Part 2: Handling Multiple Reads with io.Seeker

This section addresses the problem of reading from a reader multiple times.

### Exercise 6: Expose the Read-Once Bug

**File:** `file_system_store_test.go` **TDD Phase:** Write a test that exposes a bug **Expected Result:** Test fails (second read returns empty)

#### Why This Step?

What happens when we call `GetLeague()` twice? An `io.Reader` is like a cursor - once you've read to the end, there's nothing left to read. We need to expose this problem with a test.

#### Steps

1. At the end of the "league from a reader" test, add another call to `store.GetLeague()`
2. Assert that `got` still equals `want`
3. Run the test - observe it fails because the second read returns empty

#### Concepts

- `io.Reader` is a streaming interface - once data is consumed, it's gone
- We need a way to "rewind" the reader to the beginning

---

### Exercise 7: Fix with ReadSeeker

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** Tests pass ✓

#### Why This Step?

We need to reset the reader position before each read. The `io.Seeker` interface lets us move the read position. `io.ReadSeeker` combines `Reader` and `Seeker`.

#### Steps

1. Change the `database` field type from `io.Reader` to `io.ReadSeeker`
2. At the beginning of `GetLeague`, add `f.database.Seek(0, io.SeekStart)`
3. Run tests - they should pass

#### Concepts

- `io.ReadSeeker` embeds both `io.Reader` and `io.Seeker` interfaces
- `Seek(0, io.SeekStart)` moves to the beginning (offset 0 from start)
- `strings.NewReader` implements `ReadSeeker`, so our test still works without changes

---

## Part 3: Implementing GetPlayerScore

This section implements the GetPlayerScore method to complete basic read functionality.

### Exercise 8: Write Test for GetPlayerScore

**File:** `file_system_store_test.go` **TDD Phase:** Write the test first **Expected Result:** Compilation error - `GetPlayerScore undefined`

#### Why This Step?

Our `PlayerStore` interface requires `GetPlayerScore`. We need to implement it for `FileSystemPlayerStore`. Start with a test.

#### Steps

1. Add a new subtest "get player score"
2. Create a reader with JSON containing players Cleo (10 wins) and Chris (33 wins)
3. Create the store and call `store.GetPlayerScore("Chris")`
4. Assert the result is 33
5. Run tests - confirm compilation error

#### Concepts

- We're building up the `PlayerStore` interface implementation method by method
- Each method gets its own test

---

### Exercise 9: Add Empty GetPlayerScore Method

**File:** `file_system_store.go` **TDD Phase:** Make it compile **Expected Result:** Test fails with `got 0 want 33`

#### Why This Step?

Add the method signature to make it compile, but with a trivial implementation.

#### Steps

1. Add method `GetPlayerScore(name string) int` to `FileSystemPlayerStore`
2. Return `0`
3. Run tests - observe the meaningful failure message

#### Concepts

- Returning a zero value is the simplest implementation
- The failing test tells us exactly what we need to implement

---

### Exercise 10: Implement GetPlayerScore

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** Tests pass ✓

#### Why This Step?

Implement the actual logic - find the player in the league and return their wins.

#### Steps

1. Call `f.GetLeague()` to get all players
2. Iterate over the league with a for loop
3. If the player name matches, store their wins and break
4. Return the wins (0 if not found)
5. Run tests - they should pass

#### Concepts

- We reuse `GetLeague()` instead of duplicating the JSON parsing logic
- Returning 0 for unknown players matches our HTTP 404 behavior

---

### Exercise 11: Refactor - Add assertScoreEquals Helper

**File:** `file_system_store_test.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Why This Step?

The score comparison code can be extracted into a helper for cleaner tests.

#### Steps

1. Create `assertScoreEquals(t testing.TB, got, want int)` helper function
2. Include `t.Helper()` so test failures report the correct line
3. Update the test to use this helper
4. Run tests to confirm they still pass

---

## Part 4: Writing to Files

This section introduces writing functionality by implementing RecordWin and transitioning to real files.

### Exercise 12: Prepare for Writing - Change to ReadWriteSeeker

**File:** `file_system_store.go` **TDD Phase:** Change interface to prepare for next feature **Expected Result:** Compilation error - `strings.Reader does not implement io.ReadWriteSeeker`

#### Why This Step?

We need to implement `RecordWin` which writes data. We'll need an interface that supports reading, writing, and seeking. `io.ReadWriteSeeker` combines all three.

#### Steps

1. Change the `database` field type to `io.ReadWriteSeeker`
2. Run tests - observe compilation error because `strings.Reader` can't write

#### Concepts

- `io.ReadWriteSeeker` embeds `io.Reader`, `io.Writer`, and `io.Seeker`
- We need to update our tests to use a type that supports writing

---

### Exercise 13: Create Temp File Helper

**File:** `file_system_store_test.go` **TDD Phase:** Update test infrastructure **Expected Result:** Tests pass ✓

#### Why This Step?

`strings.Reader` can't write. We need real files. `*os.File` implements `ReadWriteSeeker`. We'll create a helper that makes temp files for testing.

#### Steps

1. Create `createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func())`
2. Use `os.CreateTemp("", "db")` to create a temp file
3. Write the initial data to the file
4. Return the file and a cleanup function that closes and removes it
5. Update all tests to use `database, cleanDatabase := createTempFile(t, jsonData)`
6. Add `defer cleanDatabase()` to each test
7. Run tests - they should pass

#### Concepts

- `os.CreateTemp` creates a file in the OS temp directory with a random name
- Returning a cleanup function as a closure keeps test code clean
- `defer` ensures cleanup runs even if the test fails

---

### Exercise 14: Write Test for RecordWin

**File:** `file_system_store_test.go` **TDD Phase:** Write the test first **Expected Result:** Compilation error - `RecordWin undefined`

#### Why This Step?

Now we can test writing. `RecordWin` should increment a player's win count and persist it to the file.

#### Steps

1. Add a new subtest "store wins for existing players"
2. Create a temp file with Chris having 33 wins
3. Create the store, call `store.RecordWin("Chris")`
4. Call `GetPlayerScore("Chris")` and assert it equals 34
5. Run tests - observe compilation error

#### Concepts

- We test the behavior (score increased) not the implementation (file contents)
- Testing an existing player first is simpler than testing a new player

---

### Exercise 15: Add Empty RecordWin Method

**File:** `file_system_store.go` **TDD Phase:** Make it compile **Expected Result:** Test fails with `got 33 want 34`

#### Why This Step?

Add the method signature with an empty body.

#### Steps

1. Add method `RecordWin(name string)` with an empty body
2. Run tests - observe the test fails because score is still 33

---

### Exercise 16: Implement RecordWin

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** Tests pass ✓

#### Why This Step?

We need to find the player, increment their wins, and write the updated league back to the file.

#### Steps

1. Call `f.GetLeague()` to get the current league
2. Iterate to find the player by name
3. **Important:** Use `league[i].Wins++` not `player.Wins++` (range gives copies!)
4. Seek to the beginning of the file
5. Encode the updated league as JSON back to the file
6. Run tests - they should pass

#### Concepts

- When ranging over a slice, you get copies of elements - modifying `player` won't affect the original
- We must write the entire league back - we can't update just one "row"
- This is inefficient but works for our prototype

---

## Part 5: Refactoring with Custom Types

This section introduces the League type and refactors the code to use it for cleaner, more expressive code.

### Exercise 17: Refactor - Create League Type with Find Method

**File:** `league.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Why This Step?

We're iterating over players to find by name in multiple places. Create a `League` type with a `Find` method to encapsulate this.

#### Steps

1. Create `type League []Player`
2. Add method `(l League) Find(name string) *Player`
3. Iterate over the league, return pointer to player if found (use `&l[i]`)
4. Return `nil` if not found
5. Run tests to ensure they still pass

#### Concepts

- Returning `*Player` allows callers to modify the player directly
- Custom types on slices can have methods - a powerful Go pattern
- `&l[i]` gets a pointer to the actual slice element, not a copy

---

### Exercise 18: Update Interface to Return League

**File:** `server.go` **TDD Phase:** Update interface **Expected Result:** Compilation errors (easy fixes)

#### Why This Step?

Change `GetLeague()` to return our new `League` type instead of `[]Player`.

#### Steps

1. In `server.go`, change `GetLeague() []Player` to `GetLeague() League` in the `PlayerStore` interface
2. Fix any compilation errors (should be straightforward type changes)
3. Run tests

---

### Exercise 19: Refactor Methods to Use League.Find

**File:** `file_system_store.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Why This Step?

Now use the new `Find` method to simplify `GetPlayerScore` and `RecordWin`.

#### Steps

1. Update `GetPlayerScore` to use `f.GetLeague().Find(name)`
2. Check if the result is not nil before returning wins
3. Update `RecordWin` to use `league.Find(name)`
4. Since `Find` returns a pointer, `player.Wins++` now works correctly
5. Run tests

#### Concepts

- The pointer from `Find` points to the actual element in the slice
- Modifying through the pointer modifies the underlying data

---

### Exercise 20: Write Test for Recording Win for New Player

**File:** `file_system_store_test.go` **TDD Phase:** Write the test first **Expected Result:** Test fails with `got 0 want 1`

#### Why This Step?

What happens when we record a win for a player who doesn't exist? They should be added with 1 win.

#### Steps

1. Add a new subtest "store wins for new players"
2. Create temp file with existing players
3. Call `store.RecordWin("Pepper")` (a new player)
4. Assert `GetPlayerScore("Pepper")` equals 1
5. Run tests - observe failure

---

### Exercise 21: Handle New Player in RecordWin

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** Tests pass ✓

#### Why This Step?

If `Find` returns nil, append a new player to the league.

#### Steps

1. After calling `Find`, check if player is nil
2. If nil, append a new `Player{name, 1}` to the league
3. Otherwise increment existing player's wins
4. Run tests

---

## Part 6: Integration and Deployment

This section integrates the FileSystemPlayerStore into the application and replaces the in-memory store.

### Exercise 22: Update Integration Test to Use FileSystemPlayerStore

**File:** `server_integration_test.go` **TDD Phase:** Integration **Expected Result:** Tests pass ✓

#### Why This Step?

Our new store should satisfy the same interface as `InMemoryPlayerStore`. Let's prove it by using it in our integration test.

#### Steps

1. Replace `NewInMemoryPlayerStore()` with `FileSystemPlayerStore{database}` using `createTempFile`
2. Add `defer cleanDatabase()`
3. Run integration test - it should pass

#### Concepts

- The integration test doesn't care which implementation we use
- This proves our new store is compatible with the existing system

---

### Exercise 23: Update main.go and Delete Old Store

**Files:** Delete `in_memory_player_store.go`, edit `main.go` **TDD Phase:** Ship it! **Expected Result:** Application compiles and runs

#### Why This Step?

Time to use our file-based store in the real application.

#### Steps

1. Delete `in_memory_player_store.go`
2. In `main.go`, use `os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)` to open/create a file
3. Handle the error
4. Create `FileSystemPlayerStore{db}`
5. Run the application and test it manually

#### Concepts

- `os.O_RDWR` = read+write, `os.O_CREATE` = create if doesn't exist
- `0666` permission means all users can read/write
- Data now persists between restarts!

---

## Part 7: Performance Optimization

This section optimizes the store by caching data in memory and introduces proper constructors.

### Exercise 24: Performance Refactor - Cache the League

**File:** `file_system_store.go` **TDD Phase:** Refactor for performance **Expected Result:** Tests pass ✓

#### Why This Step?

Every call to `GetLeague()` or `GetPlayerScore()` reads the entire file. Since our store owns the data, we can cache it in memory and only read once at startup.

#### Steps

1. Add a `league League` field to the struct
2. Create a constructor `NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore`
3. In the constructor, seek to start, call `NewLeague`, and store in the field
4. Update `GetLeague` to return `f.league`
5. Update `GetPlayerScore` to use `f.league`
6. Update `RecordWin` to use `f.league`
7. Run tests

---

### Exercise 25: Update All Code to Use Constructor

**Files:** `file_system_store_test.go`, `server_integration_test.go`, `main.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Steps

1. Update all tests to call `NewFileSystemPlayerStore(database)`
2. Update `main.go` to use the constructor
3. Run all tests

---

## Part 8: The tape Abstraction

This section introduces the tape type to handle automatic file rewinding and truncation when writing.

### Exercise 26: Create tape Type for Write-From-Start Behavior

**File:** `tape.go` (create new) **TDD Phase:** Prepare for refactor **Expected Result:** Code compiles

#### Why This Step?

When we write, we always seek to start first. Let's encapsulate this in a type called `tape` (like a cassette tape that rewinds before writing).

#### Steps

1. Create `tape.go` with a `tape` struct containing `file io.ReadWriteSeeker`
2. Add a `Write(p []byte) (n int, err error)` method
3. In `Write`, seek to start, then write
4. This makes `tape` implement `io.Writer`

---

### Exercise 27: Update Store to Use tape

**File:** `file_system_store.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Steps

1. Change `database` field type to `io.Writer`
2. In constructor, wrap the file in `&tape{database}`
3. Remove the `Seek` call from `RecordWin` (tape handles it)
4. Run tests

---

### Exercise 28: Expose the Truncation Bug x

**File:** `tape_test.go` (create new) **TDD Phase:** Write test exposing a bug **Expected Result:** Test fails with `got 'abc45' want 'abc'`

#### Why This Step?


What if we write data that's shorter than what was there before? The old data will remain at the end! We need a test to expose this.

#### Steps

1. Create a test that writes "12345" to a file
2. Create a tape and write "abc"
3. Read the entire file contents
4. Assert contents equal "abc" (not "abc45")
5. Run test - observe it fails

#### Concepts

- Writing doesn't automatically shrink a file
- We need to truncate before writing

---

### Exercise 29: Fix tape with Truncate x

**File:** `tape.go` **TDD Phase:** Make the test pass **Expected Result:** Tests pass ✓

#### Why This Step?

`*os.File` has a `Truncate` method that can shrink the file. We need to use a concrete type now.

#### Steps

1. Change `file` field from `io.ReadWriteSeeker` to `*os.File`
2. In `Write`, call `t.file.Truncate(0)` before seeking
3. Run tests

#### Concepts

- `Truncate(0)` sets file size to 0, effectively clearing it
- We sacrifice interface flexibility for functionality we need

---

### Exercise 30: Update Constructor for *os.File and json.Encoder* x

**File:** `file_system_store.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Why This Step?

Update the constructor to accept `*os.File` and store a `*json.Encoder` instead of creating one each time.

#### Steps

1. Change constructor parameter to `file *os.File`
2. Change `database` field to `*json.Encoder`
3. In constructor, create `json.NewEncoder(&tape{file})` and store it
4. Update `RecordWin` to use `f.database.Encode(f.league)`
5. Run tests

---

### Exercise 31: Update createTempFile to Return *os.File* x

**File:** `file_system_store_test.go` **TDD Phase:** Update test helper **Expected Result:** Tests pass ✓

#### Steps

1. Change `createTempFile` return type from `io.ReadWriteSeeker` to `*os.File`
2. Return `tmpfile` directly
3. Run tests

---

## Part 9: Robust Error Handling

This section adds comprehensive error handling throughout the store.

### Exercise 32: Add Error Handling to Constructor

**File:** `file_system_store.go` **TDD Phase:** Add error handling **Expected Result:** Compilation errors - multiple-value in single-value context

#### Why This Step?

We've been ignoring the error from `NewLeague`. Time to handle it properly.

#### Steps

1. Change constructor signature to return `(*FileSystemPlayerStore, error)`
2. Check error from `NewLeague`
3. Return descriptive error including filename: `fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)`
4. Run tests - observe compilation errors

#### Concepts

- Good error messages include context (what you were doing, which file)
- This is idiomatic Go, not just `return err`

---

### Exercise 33: Handle Errors in All Call Sites

**Files:** `file_system_store_test.go`, `server_integration_test.go`, `main.go` **TDD Phase:** Fix compilation, add error handling **Expected Result:** Test fails with `problem parsing league, EOF`

#### Steps

1. Create helper `assertNoError(t testing.TB, err error)`
2. Update all test calls to `store, err := NewFileSystemPlayerStore(database)` and add `assertNoError(t, err)`
3. In `main.go`, use `log.Fatalf` on error
4. Run tests - observe failure when initializing with empty string

---

### Exercise 34: Fix Integration Test with Valid Empty JSON

**File:** `server_integration_test.go` **TDD Phase:** Fix test data **Expected Result:** Tests pass ✓

#### Why This Step?

An empty string isn't valid JSON. An empty array `[]` is.

#### Steps

1. Change `createTempFile(t, "")` to `createTempFile(t, "[]")`
2. Run tests

---

### Exercise 35: Write Test for Empty File Handling

**File:** `file_system_store_test.go` **TDD Phase:** Write the test first **Expected Result:** Test fails with `problem parsing league, EOF`

#### Why This Step?

When someone runs our app for the first time, the file will be empty. We should handle this gracefully.

#### Steps

1. Add test "works with an empty file" that creates store from empty file
2. Assert no error is returned
3. Run tests - observe failure

---

### Exercise 36: Handle Empty File in Constructor

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** Tests pass ✓

#### Steps

1. After seeking to start, call `file.Stat()` to get file info
2. Check `info.Size()` - if 0, write `[]` and seek back to start
3. Handle error from `Stat`
4. Run tests

---

### Exercise 37: Refactor - Extract initialisePlayerDBFile

**File:** `file_system_store.go` **TDD Phase:** Refactor **Expected Result:** Tests pass ✓

#### Why This Step?

The constructor is getting complex. Extract the file initialization logic.

#### Steps

1. Create `initialisePlayerDBFile(file *os.File) error`
2. Move the seek, stat, and empty-file-handling logic there
3. Call it from the constructor
4. Run tests

---

## Part 10: Sorting

This final section implements sorting functionality to display the league in order of wins.

### Exercise 38: Write Test for Sorted League

**File:** `file_system_store_test.go` **TDD Phase:** Write the test first **Expected Result:** Test fails with `got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]`

#### Why This Step?

Our product owner wants `/league` to show players sorted by wins (highest first).

#### Steps

1. Add test "league sorted" with Cleo (10 wins) and Chris (33 wins) in JSON
2. Assert `GetLeague()` returns them in order: Chris first, then Cleo
3. Also test reading twice returns same order
4. Run tests - observe failure

---

### Exercise 39: Implement Sorting

**File:** `file_system_store.go` **TDD Phase:** Make the test pass **Expected Result:** All tests pass ✓ 🎉

#### Steps

1. In `GetLeague`, use `sort.Slice(f.league, func(i, j int) bool { ... })`
2. Sort by wins descending: return `f.league[i].Wins > f.league[j].Wins`
3. Return the sorted league
4. Run all tests

#### Concepts

- `sort.Slice` takes a slice and a comparison function
- The function should return true if element i should come before element j

---

## Summary

### What You've Learned

- **io.Reader, io.Writer, io.Seeker** - composable interfaces for I/O
- **Working with files** - creating temp files for tests, proper cleanup
- **sort.Slice** - sorting slices with custom comparison functions
- **Type aliases with methods** - `type League []Player` with `Find` method
- **Error handling** - returning and wrapping errors with context
- **Refactoring in TDD** - safe changes with test coverage

### Files Created

- `file_system_store.go` - The persistent store implementation
- `file_system_store_test.go` - Tests for the store
- `tape.go` - Helper type for write-with-truncate behavior
- `tape_test.go` - Tests for tape

### Key Design Decisions

- Cache league in memory, only read on startup
- Write entire league on each update (simple but not scalable)
- Handle empty files gracefully for first-time users
- Sort by wins for better UX

### Learning Path

1. **Part 1-3:** Reading fundamentals (io.Reader, io.Seeker, basic operations)
2. **Part 4-5:** Writing and custom types (io.Writer, League type abstraction)
3. **Part 6:** Integration with existing system
4. **Part 7-8:** Optimization and abstraction (caching, tape pattern)
5. **Part 9:** Production-ready error handling
6. **Part 10:** Feature completion (sorting)
# IO and Sorting - Practice Questions Only

## Starting Point
You should have the completed files from the JSON, Routing and Embedding chapter.

---

## Exercise 1
**File:** `file_system_store_test.go` (create new)
**Expected:** Compilation error - `undefined: FileSystemPlayerStore`

Create a new test file with a test for `FileSystemPlayerStore` that reads league data from a `strings.Reader` containing JSON and calls `GetLeague()`.

---

## Exercise 2
**File:** `file_system_store.go` (create new)
**Expected:** Compilation error - `too many values in struct initializer` and `GetLeague undefined`

Create an empty `FileSystemPlayerStore` struct.

---

## Exercise 3
**File:** `file_system_store.go`
**Expected:** Test fails with `got [] want [{Cleo 10} {Chris 33}]`

Add a `database` field of type `io.Reader` and a `GetLeague` method that returns `nil`.

---

## Exercise 4
**File:** `file_system_store.go`
**Expected:** Tests pass

Implement `GetLeague` to decode JSON from the database reader.

---

## Exercise 5
**Files:** `league.go`, `file_system_store.go`
**Expected:** Tests pass (refactor)

Create a `NewLeague` function that takes an `io.Reader` and returns `([]Player, error)`. Update `GetLeague` to use it.

---

## Exercise 6
**File:** `file_system_store_test.go`
**Expected:** Test fails (second read returns empty)

Add a second `GetLeague()` call and assert it returns the same data.

---

## Exercise 7
**File:** `file_system_store.go`
**Expected:** Tests pass

Change `database` to `io.ReadSeeker` and add `Seek(0, io.SeekStart)` at the beginning of `GetLeague`.

---

## Exercise 8
**File:** `file_system_store_test.go`
**Expected:** Compilation error - `GetPlayerScore undefined`

Add a test for `GetPlayerScore` that expects Chris's score to be 33.

---

## Exercise 9
**File:** `file_system_store.go`
**Expected:** Test fails with `got 0 want 33`

Add an empty `GetPlayerScore` method returning 0.

---

## Exercise 10
**File:** `file_system_store.go`
**Expected:** Tests pass

Implement `GetPlayerScore` by iterating over the league.

---

## Exercise 11
**File:** `file_system_store_test.go`
**Expected:** Tests pass (refactor)

Create `assertScoreEquals` helper function.

---

## Exercise 12
**File:** `file_system_store.go`
**Expected:** Compilation error - `strings.Reader does not implement io.ReadWriteSeeker`

Change `database` type to `io.ReadWriteSeeker`.

---

## Exercise 13
**File:** `file_system_store_test.go`
**Expected:** Tests pass

Create `createTempFile` helper that returns `(io.ReadWriteSeeker, func())`. Update all tests to use temp files with `defer cleanDatabase()`.

---

## Exercise 14
**File:** `file_system_store_test.go`
**Expected:** Compilation error - `RecordWin undefined`

Add a test for `RecordWin` - record win for Chris and expect score to go from 33 to 34.

---

## Exercise 15
**File:** `file_system_store.go`
**Expected:** Test fails with `got 33 want 34`

Add an empty `RecordWin` method.

---

## Exercise 16
**File:** `file_system_store.go`
**Expected:** Tests pass

Implement `RecordWin` to find player, increment wins, seek to start, and encode league back to file.

---

## Exercise 17
**File:** `league.go`
**Expected:** Tests pass (refactor)

Create `League` type (`type League []Player`) with a `Find(name string) *Player` method.

---

## Exercise 18
**File:** `server.go`
**Expected:** Compilation errors (easy fixes)

Change `GetLeague()` return type from `[]Player` to `League` in `PlayerStore` interface.

---

## Exercise 19
**File:** `file_system_store.go`
**Expected:** Tests pass (refactor)

Update `GetPlayerScore` and `RecordWin` to use `league.Find(name)`.

---

## Exercise 20
**File:** `file_system_store_test.go`
**Expected:** Test fails with `got 0 want 1`

Add test for recording a win for new player "Pepper".

---

## Exercise 21
**File:** `file_system_store.go`
**Expected:** Tests pass

Update `RecordWin` to append new player if `Find` returns nil.

---

## Exercise 22
**File:** `server_integration_test.go`
**Expected:** Tests pass

Replace `InMemoryPlayerStore` with `FileSystemPlayerStore` using `createTempFile`.

---

## Exercise 23
**Files:** Delete `in_memory_player_store.go`, edit `main.go`
**Expected:** Application compiles and runs

Delete old store. Update `main.go` to use `os.OpenFile` and `FileSystemPlayerStore`.

---

## Exercise 24
**File:** `file_system_store.go`
**Expected:** Tests pass (performance refactor)

Add `league League` field and create `NewFileSystemPlayerStore` constructor that reads league once. Update all methods to use cached `f.league`.

---

## Exercise 25
**Files:** `file_system_store_test.go`, `server_integration_test.go`, `main.go`
**Expected:** Tests pass

Update all code to use `NewFileSystemPlayerStore(database)`.

---

## Exercise 26
**File:** `tape.go` (create new)
**Expected:** Code compiles

Create `tape` struct wrapping `io.ReadWriteSeeker` with `Write` method that seeks to start before writing.

---

## Exercise 27
**File:** `file_system_store.go`
**Expected:** Tests pass

Change `database` to `io.Writer`. Update constructor to wrap file in `tape`. Remove `Seek` from `RecordWin`.

---

## Exercise 28
**File:** `tape_test.go` (create new)
**Expected:** Test fails with `got 'abc45' want 'abc'`

Write test: create file with "12345", write "abc", assert file contains only "abc".

---

## Exercise 29
**File:** `tape.go`
**Expected:** Tests pass

Change `file` to `*os.File`. Add `Truncate(0)` before `Seek` in `Write`.

---

## Exercise 30
**File:** `file_system_store.go`
**Expected:** Tests pass

Update constructor to accept `*os.File`. Store `*json.Encoder` instead of `io.Writer`.

---

## Exercise 31
**File:** `file_system_store_test.go`
**Expected:** Tests pass

Update `createTempFile` to return `*os.File` instead of `io.ReadWriteSeeker`.

---

## Exercise 32
**File:** `file_system_store.go`
**Expected:** Compilation errors - multiple-value in single-value context

Update constructor to return `(*FileSystemPlayerStore, error)`. Handle error from `NewLeague`.

---

## Exercise 33
**Files:** `file_system_store_test.go`, `server_integration_test.go`, `main.go`
**Expected:** Test fails with `problem parsing league, EOF`

Create `assertNoError` helper. Handle errors from constructor in all tests. In main, `log.Fatalf` on error.

---

## Exercise 34
**File:** `server_integration_test.go`
**Expected:** Tests pass

Change `createTempFile(t, "")` to `createTempFile(t, "[]")`.

---

## Exercise 35
**File:** `file_system_store_test.go`
**Expected:** Test fails with `problem parsing league, EOF`

Add test "works with an empty file" that expects no error.

---

## Exercise 36
**File:** `file_system_store.go`
**Expected:** Tests pass

In constructor, use `file.Stat()` to check size. If 0, write `[]` and seek back.

---

## Exercise 37
**File:** `file_system_store.go`
**Expected:** Tests pass (refactor)

Extract file initialization into `initialisePlayerDBFile(file *os.File) error`.

---

## Exercise 38
**File:** `file_system_store_test.go`
**Expected:** Test fails with `got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]`

Add "league sorted" test expecting players sorted by wins descending.

---

## Exercise 39
**File:** `file_system_store.go`
**Expected:** All tests pass

Add `sort.Slice` to `GetLeague` to sort by wins descending.

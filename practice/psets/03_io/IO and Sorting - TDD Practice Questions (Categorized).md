
## 1. Initial FileSystemPlayerStore Setup (Exercises 1-7)

**Focus: Creating the basic file-based store with read functionality**

- Exercise 1: Write the First Test for FileSystemPlayerStore
- Exercise 2: Create the Empty Struct
- Exercise 3: Add the Database Field and Empty Method
- Exercise 4: Implement GetLeague with JSON Decoding
- Exercise 5: Refactor - Extract NewLeague Helper
- Exercise 6: Expose the Read-Once Bug
- Exercise 7: Fix with ReadSeeker

## 2. Implementing GetPlayerScore (Exercises 8-11)

**Focus: Reading individual player data**

- Exercise 8: Write Test for GetPlayerScore
- Exercise 9: Add Empty GetPlayerScore Method
- Exercise 10: Implement GetPlayerScore
- Exercise 11: Refactor - Add assertScoreEquals Helper

## 3. Preparing for Write Operations (Exercises 12-13)

**Focus: Transitioning from read-only to read-write capability**

- Exercise 12: Prepare for Writing - Change to ReadWriteSeeker
- Exercise 13: Create Temp File Helper

## 4. Implementing RecordWin (Exercises 14-21)

**Focus: Writing player data and handling updates**

- Exercise 14: Write Test for RecordWin
- Exercise 15: Add Empty RecordWin Method
- Exercise 16: Implement RecordWin
- Exercise 17: Refactor - Create League Type with Find Method
- Exercise 18: Update Interface to Return League
- Exercise 19: Refactor Methods to Use League.Find
- Exercise 20: Write Test for Recording Win for New Player
- Exercise 21: Handle New Player in RecordWin

## 5. Integration and Deployment (Exercises 22-25)

**Focus: Connecting the store to the application and optimizing**

- Exercise 22: Update Integration Test to Use FileSystemPlayerStore
- Exercise 23: Update main.go and Delete Old Store
- Exercise 24: Performance Refactor - Cache the League
- Exercise 25: Update All Code to Use Constructor

## 6. File Write Refinement with Tape (Exercises 26-31)

**Focus: Handling file truncation and write behavior**

- Exercise 26: Create tape Type for Write-From-Start Behavior
- Exercise 27: Update Store to Use tape
- Exercise 28: Expose the Truncation Bug
- Exercise 29: Fix tape with Truncate
- Exercise 30: Update Constructor for *os.File and json.Encoder
- Exercise 31: Update createTempFile to Return *os.File

## 7. Error Handling and Edge Cases (Exercises 32-37)

**Focus: Robust error handling and empty file scenarios**

- Exercise 32: Add Error Handling to Constructor
- Exercise 33: Handle Errors in All Call Sites
- Exercise 34: Fix Integration Test with Valid Empty JSON
- Exercise 35: Write Test for Empty File Handling
- Exercise 36: Handle Empty File in Constructor
- Exercise 37: Refactor - Extract initialisePlayerDBFile

## 8. Sorting Feature (Exercises 38-39)

**Focus: Implementing league sorting by wins**

- Exercise 38: Write Test for Sorted League
- Exercise 39: Implement Sorting
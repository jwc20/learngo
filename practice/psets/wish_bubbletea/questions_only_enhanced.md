# Go Tutorial: Basic TUI App Served Over SSH - Practice Questions

Based on the video tutorial progression with Test-Driven Development integrated throughout.

---

## Part 1: Project Setup (0:00 - 4:00)

### Exercise 1
**File:** Terminal  
**Expected:** New directory and Go module created

Create a new directory for your project and initialize a Go module.

---

### Exercise 2
**File:** Terminal  
**Expected:** Packages downloaded successfully

Install the Bubble Tea package:
```
go get github.com/charmbracelet/bubbletea
```

---

### Exercise 3
**File:** `main.go` (create new)  
**Expected:** File compiles with empty main

Create `main.go` with the package declaration, import for bubbletea (aliased as `tea`), and an empty main function.

---

## Part 2: Creating the Model (4:00 - 8:00)

### Exercise 4
**File:** `main.go`  
**Expected:** Compiles (struct unused warning)

Create a `model` struct. In Bubble Tea, the model holds all your application state. For now, just create an empty struct.

---

### Exercise 5
**File:** `main.go`  
**Expected:** Compiles

Bubble Tea models must implement three methods. Add the `Init` method stub to your model:
```go
func (m model) Init() tea.Cmd
```
Return `nil` since we don't need any initial commands.

---

### Exercise 6
**File:** `main.go`  
**Expected:** Compiles

Add the `Update` method stub to your model:
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```
For now, just return the model unchanged and `nil` for the command.

---

### Exercise 7
**File:** `main.go`  
**Expected:** Compiles

Add the `View` method stub to your model:
```go
func (m model) View() string
```
Return a simple string like `"Hello Bubble Tea!"` for now.

---

### Exercise 8
**File:** `main.go`  
**Expected:** App runs and displays "Hello Bubble Tea!"

In your `main` function, create a new Bubble Tea program and run it:
1. Create an instance of your model
2. Pass it to `tea.NewProgram()`
3. Call `Run()` on the program
4. Handle any error from `Run()`

Run with `go run main.go` to test. (Note: you won't be able to quit yet!)

---

## Part 3: Setting Up the Test Framework

### Exercise 9
**File:** Terminal  
**Expected:** Packages downloaded successfully

Install the teatest package and termenv for testing:
```
go get github.com/charmbracelet/x/exp/teatest@latest
go get github.com/muesli/termenv
```

---

### Exercise 10
**File:** `main_test.go` (create new)  
**Expected:** File compiles

Create `main_test.go` with the package declaration, required imports, and an `init()` function that sets the color profile to ASCII for consistent test output.

Required imports:
- `testing`
- `tea "github.com/charmbracelet/bubbletea"`
- `"github.com/charmbracelet/lipgloss"`
- `"github.com/muesli/termenv"`

Use `lipgloss.SetColorProfile(termenv.Ascii)` in `init()`.

---

### Exercise 11
**File:** `.gitattributes` (create new)  
**Expected:** File created

Create a `.gitattributes` file to prevent Git from modifying golden file line endings. Mark `*.golden` files as binary (`-text`).

---

## Part 4: Handling Window Size with TDD (8:00 - 12:00)

### Exercise 12
**File:** `main_test.go`  
**Expected:** Test fails (width/height fields don't exist yet)

**Write the test first!** Create a test function `TestUpdateWindowSize` that:
1. Creates a model with `width: 0, height: 0`
2. Creates a `tea.WindowSizeMsg{Width: 120, Height: 40}`
3. Calls `m.Update(msg)`
4. Asserts the returned model has `width == 120` and `height == 40`
5. Asserts the returned command is `nil`

Run `go test` - it should fail because the fields don't exist yet.

---

### Exercise 13
**File:** `main.go`  
**Expected:** Compiles (test still fails)

Add `width` and `height` fields (both `int`) to your model struct to store the terminal dimensions.

Run `go test` - it should still fail because `Update` doesn't handle the message yet.

---

### Exercise 14
**File:** `main.go`  
**Expected:** Test passes!

Update your `Update` method to handle `tea.WindowSizeMsg`. Use a type switch:
```go
switch msg := msg.(type) {
case tea.WindowSizeMsg:
    // update m.width and m.height from msg.Width and msg.Height
}
```

Run `go test` - your test should pass now!

---

### Exercise 15
**File:** `main_test.go`  
**Expected:** Test fails (View doesn't show dimensions yet)

**Write the test first!** Create a test function `TestViewContainsDimensions` that:
1. Creates a model with `width: 80, height: 24`
2. Calls `m.View()`
3. Asserts the output contains "80" and "24" using `strings.Contains`

Add `"strings"` to your test imports.

---

### Exercise 16
**File:** `main.go`  
**Expected:** Test passes!

Update your `View` method to display the current window size:
```
"Window size: [width] x [height]"
```

Use `fmt.Sprintf` to format the string.

---

## Part 5: Handling Key Input with TDD (12:00 - 16:00)

### Exercise 17
**File:** `main_test.go`  
**Expected:** Test fails (quit not implemented yet)

**Write the test first!** Create a test function `TestUpdateQuitOnQ` that:
1. Creates an empty model
2. Creates a `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}`
3. Calls `m.Update(msg)`
4. Asserts the returned command is NOT nil (we expect `tea.Quit`)

---

### Exercise 18
**File:** `main.go`  
**Expected:** Test passes!

Add handling for `tea.KeyMsg` in your `Update` method. When the user presses "q" or "ctrl+c", return `tea.Quit` as the command to exit the program.

Hint: Check `msg.String() == "q"` or `msg.String() == "ctrl+c"`.

---

### Exercise 19
**File:** `main_test.go`  
**Expected:** Test fails (help text not added yet)

**Write the test first!** Create a test function `TestViewContainsHelpText` that:
1. Creates a model with `width: 80, height: 24`
2. Calls `m.View()`
3. Asserts the output contains "Press q to quit"

---

### Exercise 20
**File:** `main.go`  
**Expected:** Test passes!

Update your `View` to include a help line at the bottom:
```
"Press q to quit"
```

---

### Exercise 21
**File:** `main_test.go`  
**Expected:** All tests pass

Create a table-driven test `TestUpdateKeyHandling` that verifies multiple key scenarios:
- "q" should quit (command not nil)
- ctrl+c should quit (command not nil)
- "a" should NOT quit (command is nil)
- enter should NOT quit (command is nil)

Use the standard Go table-driven test pattern with `t.Run()` for subtests.

---

### Exercise 22
**File:** `main.go`  
**Expected:** App runs in fullscreen mode

Modify your `tea.NewProgram()` call to use alternate screen (fullscreen) mode:
```go
tea.NewProgram(m, tea.WithAltScreen())
```

---

## Part 6: Integration Tests with teatest

### Exercise 23
**File:** `main_test.go`  
**Expected:** Test passes, golden file created

Create an integration test `TestFullOutput` using teatest that:
1. Creates a model with width/height set
2. Wraps it with `teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))`
3. Sends a "q" key press with `tm.Send()`
4. Reads the final output with `tm.FinalOutput()`
5. Compares against golden file with `teatest.RequireEqualOutput(t, out)`

Add new imports: `"io"`, `"time"`, `"github.com/charmbracelet/x/exp/teatest"`

Run `go test -update` first to create the golden file, then `go test` for subsequent runs.

---

### Exercise 24
**File:** `main_test.go`  
**Expected:** Test passes

Create a test `TestFinalModelState` that:
1. Creates a model with specific width/height
2. Wraps with `teatest.NewTestModel`
3. Sends "q" to quit
4. Gets the final model with `tm.FinalModel()`
5. Type asserts to `model` and verifies width/height match initial values

---

### Exercise 25
**File:** `main_test.go`  
**Expected:** Test passes

Create a test `TestWindowResizeIntegration` that:
1. Creates a test model with initial size 80x24
2. Sends a `tea.WindowSizeMsg{Width: 200, Height: 50}`
3. Sends "q" to quit
4. Verifies the final model has the new dimensions (200x50)

---

### Exercise 26
**File:** `main_test.go`  
**Expected:** Test passes

Create a test `TestOutputContainsExpectedText` using `teatest.WaitFor` that:
1. Creates a test model
2. Uses `teatest.WaitFor()` to wait until output contains "80", "24", and "Press q to quit"
3. Sends "q" to quit
4. Calls `tm.WaitFinished()` to ensure clean exit

Add `"bytes"` to imports and use `bytes.Contains()` in the WaitFor callback.

---

### Exercise 27
**File:** `main_test.go`  
**Expected:** All tests pass

Create a parameterized test `TestVariousTerminalSizes` that tests multiple terminal dimensions:
- small: 40x10
- medium: 80x24
- large: 200x60
- wide: 300x20
- tall: 40x100

For each size, verify the output contains the correct dimensions.

---

## Part 7: Adding Style with Lip Gloss (16:00 - 24:00)

### Exercise 28
**File:** Terminal  
**Expected:** Package downloaded

Install the Lip Gloss styling library:
```
go get github.com/charmbracelet/lipgloss
```

---

### Exercise 29
**File:** `main.go`  
**Expected:** Compiles

Add the lipgloss import to your file and create a package-level style variable:
```go
var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```
Color "205" is a nice pink color.

---

### Exercise 30
**File:** `main.go`  
**Expected:** App displays styled text

Use your style in the `View` method to render the window size text:
```go
style.Render("your text here")
```

---

### Exercise 31
**File:** `main.go`  
**Expected:** App displays centered, bordered content

Enhance your style to add more visual appeal:
- Add `Bold(true)`
- Add `Border(lipgloss.RoundedBorder())`
- Add `BorderForeground(lipgloss.Color("63"))`
- Add `Padding(1, 2)`

---

### Exercise 32
**File:** `main.go`  
**Expected:** Content is centered in the terminal

Use lipgloss's `Place` function to center your content:
```go
lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
```

---

### Exercise 33
**File:** Terminal  
**Expected:** Golden file updated

Update your golden file to reflect the new styled output:
```bash
go test -update
```

Then verify tests still pass with `go test`.

---

## Part 8: Setting Up the SSH Server with Wish (24:00 - 32:00)

### Exercise 34
**File:** Terminal  
**Expected:** Packages downloaded

Install the Wish SSH server library and its dependencies:
```
go get github.com/charmbracelet/wish
go get github.com/charmbracelet/ssh
go get github.com/charmbracelet/log
```

---

### Exercise 35
**File:** `main.go`  
**Expected:** Compiles

Add the new imports to your file:
- `github.com/charmbracelet/wish`
- `github.com/charmbracelet/wish/bubbletea`
- `github.com/charmbracelet/wish/activeterm`
- `github.com/charmbracelet/wish/logging`
- `github.com/charmbracelet/ssh`
- `github.com/charmbracelet/log`

Also add standard library imports: `net`, `os/signal`, `syscall`, `context`, `time`, `errors`

---

### Exercise 36
**File:** `main.go`  
**Expected:** Compiles

Define constants for your SSH server:
```go
const (
    host = "localhost"
    port = "23234"
)
```

---

### Exercise 37
**File:** `main.go`  
**Expected:** Compiles

1. Add a `term` field (string) to your model to store the terminal type.
2. Create a `teaHandler` function:
```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption)
```

Inside teaHandler:
- Get PTY info with `sess.Pty()`
- Create a model with `term`, `width`, and `height` from the PTY
- Return the model and `[]tea.ProgramOption{tea.WithAltScreen()}`

Also update `View` to display the terminal type.

---

## Part 9: Creating the Wish Server (32:00 - 36:00)

### Exercise 38
**File:** `main.go`  
**Expected:** Compiles

Replace your `main` function to create a Wish server:
```go
s, err := wish.NewServer(
    wish.WithAddress(net.JoinHostPort(host, port)),
    wish.WithHostKeyPath(".ssh/id_ed25519"),
    wish.WithMiddleware(
        bubbletea.Middleware(teaHandler),
        activeterm.Middleware(),
        logging.Middleware(),
    ),
)
```

Handle the error with `log.Error()`.

---

### Exercise 39
**File:** `main.go`  
**Expected:** Compiles

Set up signal handling for graceful shutdown:
```go
done := make(chan os.Signal, 1)
signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
```

---

### Exercise 40
**File:** `main.go`  
**Expected:** Server starts listening

Add the startup log and start the server in a goroutine:
```go
log.Info("Starting SSH server", "host", host, "port", port)

go func() {
    if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
        log.Error("Could not start server", "error", err)
        done <- nil
    }
}()
```

---

### Exercise 41
**File:** `main.go`  
**Expected:** Server shuts down gracefully on Ctrl+C

Block on the done channel and implement graceful shutdown:
```go
<-done
log.Info("Stopping SSH server")
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := s.Shutdown(ctx); err != nil {
    log.Error("Could not stop server", "error", err)
}
```

---

## Part 10: Testing Your SSH App (36:00+)

### Exercise 42
**File:** Terminal  
**Expected:** Server starts successfully

Run your application:
```bash
go run main.go
```

You should see: `INFO Starting SSH server host=localhost port=23234`

---

### Exercise 43
**File:** Terminal (new window)  
**Expected:** TUI appears over SSH

In a new terminal, SSH into your app:
```bash
ssh localhost -p 23234
```

Your Bubble Tea app should appear!

---

### Exercise 44
**File:** SSH session  
**Expected:** Dimensions update

Resize your terminal window while connected. The displayed dimensions should update in real-time.

---

### Exercise 45
**File:** SSH session  
**Expected:** Connection closes

Press `q` to quit. The SSH connection should close gracefully.

---

### Exercise 46
**File:** Server terminal  
**Expected:** Server stops

In the server terminal, press Ctrl+C. The server should log "Stopping SSH server" and exit cleanly.

---

## Bonus Exercises

### Exercise 47
**File:** `main.go`  
**Expected:** Styled output works correctly over SSH

The styles might look different or broken over SSH because each session needs its own renderer. Fix this by:

1. Add a `renderer *lipgloss.Renderer` field to your model
2. In `teaHandler`, create a renderer with `bubbletea.MakeRenderer(sess)`
3. In `View`, create styles using `m.renderer.NewStyle()` instead of the global style

---

### Exercise 48
**File:** `~/.ssh/config`  
**Expected:** No known_hosts warnings

During development, add this to avoid host key warnings:
```
Host localhost
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
```

---

## Summary

You've built a TUI application that:
1. Uses the Elm Architecture (Model-Update-View)
2. Has comprehensive unit tests for Update and View functions
3. Has integration tests using teatest with golden files
4. Handles window resize events
5. Handles keyboard input
6. Uses Lip Gloss for styling
7. Runs as an SSH server using Wish
8. Supports multiple simultaneous connections
9. Shuts down gracefully

Each SSH connection gets its own independent Bubble Tea program!

---

## Test Summary

| Test | Type | What it tests |
|------|------|---------------|
| `TestUpdateWindowSize` | Unit | WindowSizeMsg updates model dimensions |
| `TestViewContainsDimensions` | Unit | View displays width and height |
| `TestUpdateQuitOnQ` | Unit | "q" key returns quit command |
| `TestViewContainsHelpText` | Unit | View includes help text |
| `TestUpdateKeyHandling` | Unit (table) | Multiple key bindings |
| `TestFullOutput` | Integration | Golden file comparison |
| `TestFinalModelState` | Integration | Model state after quit |
| `TestWindowResizeIntegration` | Integration | Resize via Send() |
| `TestOutputContainsExpectedText` | Integration | WaitFor content |
| `TestVariousTerminalSizes` | Integration | Multiple terminal sizes |

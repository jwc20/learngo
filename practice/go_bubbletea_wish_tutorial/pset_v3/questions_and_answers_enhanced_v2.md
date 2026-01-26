# Go Tutorial: Basic TUI App Served Over SSH - Questions & Answers

Based on the video tutorial progression with Test-Driven Development integrated throughout.

---

## Part 1: Project Setup (0:00 - 4:00)

### Exercise 1
**File:** Terminal  
**Expected:** New directory and Go module created

**Question:** Create a new directory for your project and initialize a Go module.

**Why:** Every Go project needs a module to manage dependencies and define the project's import path. The `go mod init` command creates a `go.mod` file that tracks your project's dependencies. This is the foundation of Go's dependency management system introduced in Go 1.11.

**Answer:**
```bash
mkdir ssh-tui
cd ssh-tui
go mod init ssh-tui
```

---

### Exercise 2
**File:** Terminal  
**Expected:** Packages downloaded successfully

**Question:** Install the Bubble Tea package.

**Why:** Bubble Tea is a powerful framework for building terminal user interfaces in Go. It implements the Elm Architecture pattern (Model-Update-View), which provides a clean, predictable way to manage application state and handle events. Without this dependency, we'd have to manually handle terminal raw mode, input parsing, and screen rendering ourselves.

**Answer:**
```bash
go get github.com/charmbracelet/bubbletea
```

---

### Exercise 3
**File:** `main.go` (create new)  
**Expected:** File compiles with empty main

**Question:** Create `main.go` with the package declaration, import for bubbletea (aliased as `tea`), and an empty main function.

**Why:** We alias `bubbletea` as `tea` because it's the conventional shorthand used in all Bubble Tea documentation and examples. The shorter alias makes code more readable since we'll be referencing it frequently (e.g., `tea.Cmd`, `tea.Msg`, `tea.Model`). The underscore import (`_ "github.com/..."`) would be wrong here since we need to actually use the package.

**Answer:**
```go
package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
}
```

---

## Part 2: Creating the Model (4:00 - 8:00)

### Exercise 4
**File:** `main.go`  
**Expected:** Compiles (struct unused warning)

**Question:** Create a `model` struct.

**Why:** In the Elm Architecture, the Model is the single source of truth for your application's state. Every piece of data your app needs to track—user input, computed values, UI state—lives here. By centralizing state in a struct, we avoid global variables and make the data flow explicit and testable. Right now it's empty, but we'll add fields as our app grows.

**Answer:**
```go
package main

import tea "github.com/charmbracelet/bubbletea"

type model struct{}

func main() {
}
```

---

### Exercise 5
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the `Init` method stub to your model.

**Why:** `Init()` is the first of three required methods for the `tea.Model` interface. It runs once when the program starts and returns an optional `tea.Cmd` (a function that performs I/O and returns a message). Common uses include: fetching initial data, starting timers, or sending initial commands. Returning `nil` means "do nothing on startup"—perfect for simple apps that just wait for user input.

**Answer:**
```go
package main

import tea "github.com/charmbracelet/bubbletea"

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {
}
```

---

### Exercise 6
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the `Update` method stub to your model.

**Why:** `Update()` is the heart of your application's logic. It receives messages (keyboard input, window resize, custom events) and returns an updated model plus an optional command. This pure-function approach (input → output, no side effects) makes your code predictable and easy to test. The signature `(tea.Model, tea.Cmd)` returns an interface type, allowing Bubble Tea to work with any model that implements the required methods.

**Answer:**
```go
package main

import tea "github.com/charmbracelet/bubbletea"

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func main() {
}
```

---

### Exercise 7
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the `View` method stub to your model.

**Why:** `View()` is a pure function that transforms your model into a string for display. It's called after every `Update()` and should have no side effects—just read the model and return what to show. This separation means rendering logic never accidentally modifies state. The returned string can include ANSI escape codes for colors/styling, which we'll add later with Lip Gloss.

**Answer:**
```go
package main

import tea "github.com/charmbracelet/bubbletea"

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return "Hello Bubble Tea!"
}

func main() {
}
```

---

### Exercise 8
**File:** `main.go`  
**Expected:** App runs and displays "Hello Bubble Tea!"

**Question:** In your `main` function, create a new Bubble Tea program and run it.

**Why:** `tea.NewProgram()` creates a new Bubble Tea application with your model. The `Run()` method starts the event loop: it puts the terminal in raw mode, renders the initial view, and begins listening for input. The returned error handles cases like terminal initialization failures. We exit with code 1 on error because that's the Unix convention for "something went wrong." Note: you won't be able to quit yet since we haven't added key handling!

**Answer:**
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return "Hello Bubble Tea!"
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

## Part 3: Setting Up the Test Framework

### Exercise 9
**File:** Terminal  
**Expected:** Packages downloaded successfully

**Question:** Install the teatest package and termenv for testing.

**Why:** `teatest` is Charm's official testing library for Bubble Tea applications. It provides utilities for testing TUI apps including golden file comparisons, model state assertions, and simulating user input. `termenv` is needed to normalize color output so tests pass consistently across different terminals and CI environments.

**Answer:**
```bash
go get github.com/charmbracelet/x/exp/teatest@latest
go get github.com/muesli/termenv
```

---

### Exercise 10
**File:** `main_test.go` (create new)  
**Expected:** File compiles

**Question:** Create `main_test.go` with the package declaration, required imports, and an `init()` function that sets the color profile to ASCII.

**Why:** The `init()` function runs before any tests and sets a consistent color profile. Without this, golden files generated on your machine (with 256-color support) would differ from CI environments (often with limited color support). Setting `termenv.Ascii` disables colors entirely, ensuring consistent output everywhere. This is a common gotcha that causes CI failures.

**Answer:**
```go
package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}
```

---

### Exercise 11
**File:** `.gitattributes` (create new)  
**Expected:** File created

**Question:** Create a `.gitattributes` file to prevent Git from modifying golden file line endings.

**Why:** Golden files contain exact byte-for-byte expected output. Git's automatic line ending conversion (CRLF on Windows, LF on Unix) can corrupt these files, causing tests to fail mysteriously. Marking `*.golden` files as `-text` tells Git to treat them as binary, preserving their exact contents across platforms.

**Answer:**
```
*.golden -text
```

---

## Part 4: Handling Window Size with TDD (8:00 - 12:00)

### Exercise 12
**File:** `main_test.go`  
**Expected:** Test fails (width/height fields don't exist yet)

**Question:** Write a test that verifies the `Update` method correctly handles `tea.WindowSizeMsg` by updating the model's width and height.

**Why:** This is Test-Driven Development (TDD) in action: write the test first, watch it fail, then implement the feature. Testing `Update` is straightforward because it's a pure function—given an input message, it returns a predictable output. This test documents the expected behavior and will catch regressions if someone accidentally breaks window resize handling.

**Answer:**
```go
package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestUpdateWindowSize(t *testing.T) {
	m := model{width: 0, height: 0}
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	newModel, cmd := m.Update(msg)

	updatedModel := newModel.(model)
	if updatedModel.width != 120 {
		t.Errorf("expected width 120, got %d", updatedModel.width)
	}
	if updatedModel.height != 40 {
		t.Errorf("expected height 40, got %d", updatedModel.height)
	}
	if cmd != nil {
		t.Error("expected nil command")
	}
}
```

Run `go test` - it should fail because `width` and `height` fields don't exist yet.

---

### Exercise 13
**File:** `main.go`  
**Expected:** Compiles (test still fails)

**Question:** Add `width` and `height` fields (both `int`) to your model struct to store the terminal dimensions.

**Why:** Terminal applications need to know the window dimensions to properly lay out content—centering text, creating borders that fit, handling responsive designs. Now our test compiles, but still fails because `Update` doesn't handle the message yet.

**Answer:**
```go
type model struct {
	width  int
	height int
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return "Hello Bubble Tea!"
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

### Exercise 14
**File:** `main.go`  
**Expected:** Test passes!

**Question:** Update your `Update` method to handle `tea.WindowSizeMsg` by storing the width and height in the model.

**Why:** Now we implement just enough code to make the test pass. The type switch pattern `switch msg := msg.(type)` is idiomatic Go for handling different message types. Bubble Tea sends `WindowSizeMsg` both at startup and whenever the user resizes their terminal. Run `go test` - it should pass now!

**Answer:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	return "Hello Bubble Tea!"
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

### Exercise 15
**File:** `main_test.go`  
**Expected:** Test fails (View doesn't show dimensions yet)

**Question:** Write a test that verifies the `View` method displays the window dimensions.

**Why:** Continuing TDD, we write the test for View before implementing it. `View` is also a pure function, making it easy to test. We use `strings.Contains` for flexible matching—we don't care about exact formatting, just that the dimensions appear somewhere in the output.

**Answer:**
```go
func TestViewContainsDimensions(t *testing.T) {
	m := model{
		width:  80,
		height: 24,
	}

	output := m.View()

	if !strings.Contains(output, "80") {
		t.Error("view should contain width")
	}
	if !strings.Contains(output, "24") {
		t.Error("view should contain height")
	}
}
```

Add `"strings"` to your test file imports:
```go
import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)
```

---

### Exercise 16
**File:** `main.go`  
**Expected:** Test passes!

**Question:** Update your `View` method to display the current window size in the format "Window size: [width] x [height]".

**Why:** This exercises the data flow: WindowSizeMsg → Update stores values → View reads and displays them. Using `fmt.Sprintf` lets us format the integers into a readable string. Run `go test` to verify the test passes.

**Answer:**
```go
func (m model) View() string {
	return fmt.Sprintf("Window size: %d x %d", m.width, m.height)
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Window size: %d x %d", m.width, m.height)
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

## Part 5: Handling Key Input with TDD (12:00 - 16:00)

### Exercise 17
**File:** `main_test.go`  
**Expected:** Test fails (quit not implemented yet)

**Question:** Write a test that verifies pressing "q" returns `tea.Quit` command.

**Why:** Every TUI needs a way to exit! We test this critical functionality before implementing it. The test creates a `KeyMsg` with the rune 'q' and verifies that `Update` returns a non-nil command. Note: We can't directly compare functions in Go, so we check that a command is returned (not nil).

**Answer:**
```go
func TestUpdateQuitOnQ(t *testing.T) {
	m := model{}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}
```

Full test file context:
```go
package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestUpdateWindowSize(t *testing.T) {
	m := model{width: 0, height: 0}
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	newModel, cmd := m.Update(msg)

	updatedModel := newModel.(model)
	if updatedModel.width != 120 {
		t.Errorf("expected width 120, got %d", updatedModel.width)
	}
	if updatedModel.height != 40 {
		t.Errorf("expected height 40, got %d", updatedModel.height)
	}
	if cmd != nil {
		t.Error("expected nil command")
	}
}

func TestViewContainsDimensions(t *testing.T) {
	m := model{
		width:  80,
		height: 24,
	}

	output := m.View()

	if !strings.Contains(output, "80") {
		t.Error("view should contain width")
	}
	if !strings.Contains(output, "24") {
		t.Error("view should contain height")
	}
}

func TestUpdateQuitOnQ(t *testing.T) {
	m := model{}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}
```

---

### Exercise 18
**File:** `main.go`  
**Expected:** Test passes!

**Question:** Add handling for `tea.KeyMsg` in your `Update` method. When the user presses "q" or "ctrl+c", return `tea.Quit` as the command.

**Why:** Now we implement the quit functionality to make our test pass. `tea.KeyMsg` is sent for each keypress, and `msg.String()` returns a human-readable representation like "q", "enter", "ctrl+c". Returning `tea.Quit` tells Bubble Tea to gracefully shut down—it restores the terminal to normal mode and exits.

**Answer:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Window size: %d x %d", m.width, m.height)
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

### Exercise 19
**File:** `main_test.go`  
**Expected:** Test fails (help text not added yet)

**Question:** Write a test that verifies the `View` method includes help text telling users how to quit.

**Why:** Good UX means users know how to interact with your app. This test ensures the help text is always present and won't be accidentally removed during refactoring.

**Answer:**
```go
func TestViewContainsHelpText(t *testing.T) {
	m := model{width: 80, height: 24}

	output := m.View()

	if !strings.Contains(output, "Press q to quit") {
		t.Error("view should contain help text")
	}
}
```

---

### Exercise 20
**File:** `main.go`  
**Expected:** Test passes!

**Question:** Update your `View` to include a help line at the bottom: "Press q to quit".

**Why:** A help line at the bottom is a TUI convention—think of vim's `:help` prompt or htop's key shortcuts. The `\n\n` creates visual separation between content and help.

**Answer:**
```go
func (m model) View() string {
	s := fmt.Sprintf("Window size: %d x %d\n\n", m.width, m.height)
	s += "Press q to quit"
	return s
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Window size: %d x %d\n\n", m.width, m.height)
	s += "Press q to quit"
	return s
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

### Exercise 21
**File:** `main_test.go`  
**Expected:** All tests pass

**Question:** Write a table-driven test that verifies multiple key bindings: "q" should quit, "ctrl+c" should quit, and other keys like "a" should not quit.

**Why:** Table-driven tests are idiomatic Go. They make it easy to add new test cases and keep tests DRY. This pattern scales well as you add more key bindings. Each test case is self-documented and failures clearly indicate which case failed.

**Answer:**
```go
func TestUpdateKeyHandling(t *testing.T) {
	tests := []struct {
		name       string
		keyType    tea.KeyType
		runes      []rune
		shouldQuit bool
	}{
		{"q quits", tea.KeyRunes, []rune{'q'}, true},
		{"ctrl+c quits", tea.KeyCtrlC, nil, true},
		{"a doesn't quit", tea.KeyRunes, []rune{'a'}, false},
		{"enter doesn't quit", tea.KeyEnter, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{}
			msg := tea.KeyMsg{Type: tt.keyType, Runes: tt.runes}

			_, cmd := m.Update(msg)

			if tt.shouldQuit && cmd == nil {
				t.Error("expected quit command")
			}
			if !tt.shouldQuit && cmd != nil {
				t.Error("expected no command")
			}
		})
	}
}
```

---

### Exercise 22
**File:** `main.go`  
**Expected:** App runs in fullscreen mode

**Question:** Modify your `tea.NewProgram()` call to use alternate screen (fullscreen) mode.

**Why:** Terminals have two screen buffers: the normal buffer (where your shell history lives) and the alternate screen. Programs like vim, htop, and less use the alternate screen so they can take over the entire terminal without destroying your command history. When the program exits, the terminal automatically restores the previous buffer.

**Answer:**
```go
func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Window size: %d x %d\n\n", m.width, m.height)
	s += "Press q to quit"
	return s
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

## Part 6: Integration Tests with teatest

### Exercise 23
**File:** `main_test.go`  
**Expected:** Test passes, golden file created

**Question:** Write an integration test using teatest that captures the full output of your app and compares it against a golden file.

**Why:** Golden file testing captures the entire visual output of your app. When you run with `-update`, it creates a snapshot in `testdata/`. Future tests compare against this snapshot, catching unintended visual regressions. This is like screenshot testing for TUIs. The `FinalOutput` method waits for the program to exit before capturing output.

**Answer:**
```go
func TestFullOutput(t *testing.T) {
	m := model{
		width:  80,
		height: 24,
	}
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if err != nil {
		t.Fatal(err)
	}
	teatest.RequireEqualOutput(t, out)
}
```

Add new imports to test file:
```go
import (
	"io"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
)
```

Run first time to create golden file:
```bash
go test -update
```

Run subsequent times to compare:
```bash
go test
```

---

### Exercise 24
**File:** `main_test.go`  
**Expected:** Test passes

**Question:** Write a test using teatest that verifies the model's state after the program exits.

**Why:** `FinalModel` returns the model's internal state when the program quit. This lets you verify internal state, not just visual output. For example: are dimensions stored correctly after initialization? This catches bugs where the output looks right but internal state is wrong.

**Answer:**
```go
func TestFinalModelState(t *testing.T) {
	initialWidth, initialHeight := 120, 40
	tm := teatest.NewTestModel(
		t,
		model{width: initialWidth, height: initialHeight},
		teatest.WithInitialTermSize(initialWidth, initialHeight),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	fm := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second))
	m, ok := fm.(model)
	if !ok {
		t.Fatalf("final model has wrong type: %T", fm)
	}

	if m.width != initialWidth {
		t.Errorf("expected width %d, got %d", initialWidth, m.width)
	}
	if m.height != initialHeight {
		t.Errorf("expected height %d, got %d", initialHeight, m.height)
	}
}
```

---

### Exercise 25
**File:** `main_test.go`  
**Expected:** Test passes

**Question:** Write a test that simulates a window resize event and verifies the model updates correctly.

**Why:** `tm.Send()` can send any `tea.Msg`, not just key presses. This tests the reactive resize behavior by sending `WindowSizeMsg` and checking the model updates. It's like simulating a user resizing their terminal mid-session.

**Answer:**
```go
func TestWindowResizeIntegration(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		model{width: 80, height: 24},
		teatest.WithInitialTermSize(80, 24),
	)

	tm.Send(tea.WindowSizeMsg{Width: 200, Height: 50})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	fm := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second))
	m := fm.(model)

	if m.width != 200 {
		t.Errorf("expected width 200 after resize, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("expected height 50 after resize, got %d", m.height)
	}
}
```

---

### Exercise 26
**File:** `main_test.go`  
**Expected:** Test passes

**Question:** Write a test using `teatest.WaitFor` to verify that specific content appears in the output stream.

**Why:** `WaitFor` is powerful for testing dynamic UIs. Instead of waiting for the program to end, you assert that certain content appears within a time window. If the condition isn't met within the timeout, the test fails. This is useful for apps with animations, loading states, or real-time updates.

**Answer:**
```go
func TestOutputContainsExpectedText(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		model{width: 80, height: 24},
		teatest.WithInitialTermSize(80, 24),
	)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("80")) &&
			bytes.Contains(bts, []byte("24")) &&
			bytes.Contains(bts, []byte("Press q to quit"))
	}, teatest.WithCheckInterval(100*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
```

Add `"bytes"` to imports.

---

### Exercise 27
**File:** `main_test.go`  
**Expected:** All tests pass

**Question:** Write a parameterized test that verifies your app works correctly at various terminal sizes.

**Why:** Users have wildly different terminal sizes—from small laptop screens to ultrawide monitors. Testing multiple sizes catches layout bugs like content overflow, incorrect centering, or elements that disappear at small sizes. Each test case runs independently, so one failure doesn't block others.

**Answer:**
```go
func TestVariousTerminalSizes(t *testing.T) {
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"small", 40, 10},
		{"medium", 80, 24},
		{"large", 200, 60},
		{"wide", 300, 20},
		{"tall", 40, 100},
	}

	for _, sz := range sizes {
		t.Run(sz.name, func(t *testing.T) {
			tm := teatest.NewTestModel(
				t,
				model{width: sz.width, height: sz.height},
				teatest.WithInitialTermSize(sz.width, sz.height),
			)

			teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
				widthStr := fmt.Sprintf("%d", sz.width)
				heightStr := fmt.Sprintf("%d", sz.height)
				return bytes.Contains(bts, []byte(widthStr)) &&
					bytes.Contains(bts, []byte(heightStr))
			}, teatest.WithDuration(time.Second))

			tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
		})
	}
}
```

---

## Part 7: Adding Style with Lip Gloss (16:00 - 24:00)

### Exercise 28
**File:** Terminal  
**Expected:** Package downloaded

**Question:** Install the Lip Gloss styling library.

**Why:** Raw terminal output is plain and hard to read. Lip Gloss is Charm's styling library that generates ANSI escape codes for colors, borders, padding, and alignment. It's like CSS for the terminal. Instead of manually writing escape sequences like `\033[1;35m`, you use a fluent API.

**Answer:**
```bash
go get github.com/charmbracelet/lipgloss
```

---

### Exercise 29
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the lipgloss import to your file and create a style variable at the package level for styling text with a pink foreground color.

**Why:** Lip Gloss styles are created by chaining method calls on a `Style` object. Defining styles at the package level means they're created once at startup, not recreated on every render. Color "205" is from the 256-color palette (a nice pink).

**Answer:**
```go
import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```

---

### Exercise 30
**File:** `main.go`  
**Expected:** App displays styled text

**Question:** Use your style in the `View` method to render the window size text with styling.

**Why:** `style.Render()` takes a string and returns it wrapped in ANSI escape codes for the configured styling. We keep the help text unstyled for contrast, making the main content stand out.

**Answer:**
```go
func (m model) View() string {
	s := style.Render(fmt.Sprintf("Window size: %d x %d", m.width, m.height))
	s += "\n\nPress q to quit"
	return s
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := style.Render(fmt.Sprintf("Window size: %d x %d", m.width, m.height))
	s += "\n\nPress q to quit"
	return s
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

### Exercise 31
**File:** `main.go`  
**Expected:** App displays centered, bordered content

**Question:** Enhance your style with bold text, a rounded border with a different color, and padding.

**Why:** Lip Gloss's fluent API lets you chain multiple style properties. `Bold(true)` makes text bold. `Border()` draws a box around content—`RoundedBorder()` uses Unicode box-drawing characters for smooth corners. `Padding(1, 2)` adds 1 line vertical and 2 characters horizontal padding inside the border.

**Answer:**
```go
var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(1, 2)
```

---

### Exercise 32
**File:** `main.go`  
**Expected:** Content is centered in the terminal

**Question:** Use lipgloss's `Place` function to center your content in the terminal window.

**Why:** `lipgloss.Place()` positions content within a bounding box. By passing `m.width` and `m.height` (our terminal dimensions), we center content in the full window. This is why we stored window size earlier—layout functions need these dimensions.

**Answer:**
```go
func (m model) View() string {
	content := style.Render(fmt.Sprintf("Window size: %d x %d", m.width, m.height))
	help := "\n\nPress q to quit"

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content+help,
	)
}
```

Full file context:
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(1, 2)

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	content := style.Render(fmt.Sprintf("Window size: %d x %d", m.width, m.height))
	help := "\n\nPress q to quit"

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content+help,
	)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

---

### Exercise 33
**File:** Terminal  
**Expected:** Golden file updated

**Question:** Update your golden file to reflect the new styled output.

**Why:** Since we changed the visual output (added styling, centering), our golden file is now out of date. Running with `-update` regenerates it with the new expected output. Always review the diff in your golden file to ensure the changes are intentional!

**Answer:**
```bash
go test -update
```

Then verify:
```bash
go test
```

---

## Part 8: Setting Up the SSH Server with Wish (24:00 - 32:00)

### Exercise 34
**File:** Terminal  
**Expected:** Packages downloaded

**Question:** Install the Wish SSH server library and its dependencies.

**Why:** Wish is an SSH server library that makes it trivial to serve TUI applications over SSH. Instead of users running your binary locally, they can `ssh yourserver.com` and interact with your app remotely.

**Answer:**
```bash
go get github.com/charmbracelet/wish
go get github.com/charmbracelet/ssh
go get github.com/charmbracelet/log
```

---

### Exercise 35
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add all the new imports to your file for SSH server functionality.

**Why:** Each package serves a specific purpose: `wish` provides the SSH server framework; `wish/bubbletea` integrates Bubble Tea apps with SSH sessions; `wish/activeterm` ensures clients have a proper PTY; `wish/logging` logs connections; `ssh` provides session types; `log` is Charm's structured logger.

**Answer:**
```go
import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)
```

---

### Exercise 36
**File:** `main.go`  
**Expected:** Compiles

**Question:** Define constants for your SSH server host and port.

**Why:** Constants make configuration clear and changeable from one place. Port 23234 is a high port (>1024) so we don't need root privileges. Using `localhost` restricts connections to the local machine—in production you'd use `0.0.0.0`.

**Answer:**
```go
const (
	host = "localhost"
	port = "23234"
)
```

Place after imports, before the style variable.

---

### Exercise 37
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add a `term` field to your model struct to store the terminal type, and create a `teaHandler` function that initializes a model from SSH session info.

**Why:** The `teaHandler` is called by Wish for each SSH connection to create a fresh Bubble Tea program. Each user gets their own model instance. `sess.Pty()` returns PTY information including terminal type and initial dimensions.

**Answer:**
```go
type model struct {
	term   string
	width  int
	height int
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()
	m := model{
		term:   pty.Term,
		width:  pty.Window.Width,
		height: pty.Window.Height,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
```

Also update View to display terminal type:
```go
func (m model) View() string {
	info := fmt.Sprintf("Terminal: %s\nWindow size: %d x %d", m.term, m.width, m.height)
	content := style.Render(info)
	help := "\n\nPress q to quit"

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content+help,
	)
}
```

---

## Part 9: Creating the Wish Server (32:00 - 36:00)

### Exercise 38
**File:** `main.go`  
**Expected:** Compiles

**Question:** Replace your `main` function to create a Wish server with address, host key path, and middleware stack.

**Why:** `wish.NewServer()` creates an SSH server. `WithAddress` sets the listen address. `WithHostKeyPath` specifies where to store the server's SSH host key. The middleware stack (logging → activeterm → bubbletea) processes connections in reverse order.

**Answer:**
```go
func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
	}
}
```

---

### Exercise 39
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add signal handling for graceful shutdown with SIGINT and SIGTERM.

**Why:** Production servers need to shut down cleanly when receiving interrupt signals (Ctrl+C) or termination signals (from systemd, Docker). Using a buffered channel ensures the signal isn't lost.

**Answer:**
```go
func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}
```

---

### Exercise 40
**File:** `main.go`  
**Expected:** Server starts listening

**Question:** Add the startup log message and start the server in a goroutine.

**Why:** `ListenAndServe()` blocks forever—running it in a goroutine lets the main goroutine continue to handle shutdown signals. The error check excludes `ssh.ErrServerClosed` because that's expected when we call `Shutdown()`.

**Answer:**
```go
func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()
}
```

---

### Exercise 41
**File:** `main.go`  
**Expected:** Server shuts down gracefully on Ctrl+C

**Question:** Implement graceful shutdown by blocking on the done channel and using a context with timeout.

**Why:** `<-done` blocks until a signal arrives. On receiving a signal, we create a context with 30-second timeout (giving active connections time to finish) and call `Shutdown()`. This prevents hanging forever if a connection won't close.

**Answer:**
```go
func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}
```

---

## Part 10: Testing Your SSH App (36:00+)

### Exercise 42
**File:** Terminal  
**Expected:** Server starts successfully

**Question:** Run your application.

**Why:** This tests that everything compiles and the server can bind to the port. The first run will generate an SSH host key file.

**Answer:**
```bash
go run main.go
```

Output:
```
INFO Starting SSH server host=localhost port=23234
```

---

### Exercise 43
**File:** Terminal (new window)  
**Expected:** TUI appears over SSH

**Question:** SSH into your app.

**Why:** This is the moment of truth—connecting to your TUI over SSH. The `-p` flag specifies the non-standard port.

**Answer:**
```bash
ssh localhost -p 23234
```

---

### Exercise 44
**File:** SSH session  
**Expected:** Dimensions update

**Question:** Test window resize handling by resizing your terminal window while connected.

**Why:** Resize events are forwarded over SSH via the PTY. Seeing the dimensions update confirms the full event pipeline works.

**Answer:**
Resize your terminal window while connected. The dimensions should update automatically.

---

### Exercise 45
**File:** SSH session  
**Expected:** Connection closes

**Question:** Quit the app by pressing "q".

**Why:** Pressing "q" triggers your KeyMsg handler, which returns `tea.Quit`. The SSH connection closes but the server keeps running.

**Answer:**
Press `q`. Output: `Connection to localhost closed.`

---

### Exercise 46
**File:** Server terminal  
**Expected:** Server stops

**Question:** Stop the server by pressing Ctrl+C.

**Why:** The graceful shutdown code runs, closing the listener and waiting for active connections.

**Answer:**
Press Ctrl+C in the server terminal.
Output:
```
INFO Stopping SSH server
```

---

## Bonus Exercises

### Exercise 47
**File:** `main.go`  
**Expected:** Styled output works correctly over SSH

**Question:** Fix styles for SSH by using a per-session renderer instead of the global style.

**Why:** Lip Gloss defaults to detecting terminal capabilities from `os.Stdout`. But over SSH, each session has its own output stream! `bubbletea.MakeRenderer(sess)` creates a renderer aware of the specific session's terminal capabilities.

**Answer:**
```go
type model struct {
	term     string
	width    int
	height   int
	renderer *lipgloss.Renderer
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()
	renderer := bubbletea.MakeRenderer(sess)
	m := model{
		term:     pty.Term,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		renderer: renderer,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) View() string {
	style := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	info := fmt.Sprintf("Terminal: %s\nWindow size: %d x %d", m.term, m.width, m.height)
	content := style.Render(info)
	help := "\n\nPress q to quit"

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content+help,
	)
}
```

---

### Exercise 48
**File:** `~/.ssh/config`  
**Expected:** No known_hosts warnings

**Question:** Configure SSH to skip host key checking for localhost during development.

**Why:** During development, you'll regenerate the server's host key often. Disabling this check for localhost only saves frustration. Never do this for remote servers!

**Answer:**
```
Host localhost
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
```

---

## Complete Final Code

### main.go
```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

type model struct {
	term     string
	width    int
	height   int
	renderer *lipgloss.Renderer
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	style := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	info := fmt.Sprintf("Terminal: %s\nWindow size: %d x %d", m.term, m.width, m.height)
	content := style.Render(info)
	help := "\n\nPress q to quit"

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content+help,
	)
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()
	renderer := bubbletea.MakeRenderer(sess)
	m := model{
		term:     pty.Term,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		renderer: renderer,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}
```

### main_test.go
```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestUpdateWindowSize(t *testing.T) {
	m := model{width: 0, height: 0}
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	newModel, cmd := m.Update(msg)

	updatedModel := newModel.(model)
	if updatedModel.width != 120 {
		t.Errorf("expected width 120, got %d", updatedModel.width)
	}
	if updatedModel.height != 40 {
		t.Errorf("expected height 40, got %d", updatedModel.height)
	}
	if cmd != nil {
		t.Error("expected nil command")
	}
}

func TestViewContainsDimensions(t *testing.T) {
	m := model{
		width:    80,
		height:   24,
		renderer: lipgloss.DefaultRenderer(),
	}

	output := m.View()

	if !strings.Contains(output, "80") {
		t.Error("view should contain width")
	}
	if !strings.Contains(output, "24") {
		t.Error("view should contain height")
	}
}

func TestUpdateQuitOnQ(t *testing.T) {
	m := model{}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}

func TestViewContainsHelpText(t *testing.T) {
	m := model{
		width:    80,
		height:   24,
		renderer: lipgloss.DefaultRenderer(),
	}

	output := m.View()

	if !strings.Contains(output, "Press q to quit") {
		t.Error("view should contain help text")
	}
}

func TestUpdateKeyHandling(t *testing.T) {
	tests := []struct {
		name       string
		keyType    tea.KeyType
		runes      []rune
		shouldQuit bool
	}{
		{"q quits", tea.KeyRunes, []rune{'q'}, true},
		{"ctrl+c quits", tea.KeyCtrlC, nil, true},
		{"a doesn't quit", tea.KeyRunes, []rune{'a'}, false},
		{"enter doesn't quit", tea.KeyEnter, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{}
			msg := tea.KeyMsg{Type: tt.keyType, Runes: tt.runes}

			_, cmd := m.Update(msg)

			if tt.shouldQuit && cmd == nil {
				t.Error("expected quit command")
			}
			if !tt.shouldQuit && cmd != nil {
				t.Error("expected no command")
			}
		})
	}
}

func TestFullOutput(t *testing.T) {
	m := model{
		width:    80,
		height:   24,
		renderer: lipgloss.DefaultRenderer(),
	}
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if err != nil {
		t.Fatal(err)
	}
	teatest.RequireEqualOutput(t, out)
}

func TestFinalModelState(t *testing.T) {
	initialWidth, initialHeight := 120, 40
	tm := teatest.NewTestModel(
		t,
		model{width: initialWidth, height: initialHeight, renderer: lipgloss.DefaultRenderer()},
		teatest.WithInitialTermSize(initialWidth, initialHeight),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	fm := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second))
	m, ok := fm.(model)
	if !ok {
		t.Fatalf("final model has wrong type: %T", fm)
	}

	if m.width != initialWidth {
		t.Errorf("expected width %d, got %d", initialWidth, m.width)
	}
	if m.height != initialHeight {
		t.Errorf("expected height %d, got %d", initialHeight, m.height)
	}
}

func TestWindowResizeIntegration(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		model{width: 80, height: 24, renderer: lipgloss.DefaultRenderer()},
		teatest.WithInitialTermSize(80, 24),
	)

	tm.Send(tea.WindowSizeMsg{Width: 200, Height: 50})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

	fm := tm.FinalModel(t, teatest.WithFinalTimeout(time.Second))
	m := fm.(model)

	if m.width != 200 {
		t.Errorf("expected width 200 after resize, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("expected height 50 after resize, got %d", m.height)
	}
}

func TestOutputContainsExpectedText(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		model{width: 80, height: 24, renderer: lipgloss.DefaultRenderer()},
		teatest.WithInitialTermSize(80, 24),
	)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("80")) &&
			bytes.Contains(bts, []byte("24")) &&
			bytes.Contains(bts, []byte("Press q to quit"))
	}, teatest.WithCheckInterval(100*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestVariousTerminalSizes(t *testing.T) {
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"small", 40, 10},
		{"medium", 80, 24},
		{"large", 200, 60},
		{"wide", 300, 20},
		{"tall", 40, 100},
	}

	for _, sz := range sizes {
		t.Run(sz.name, func(t *testing.T) {
			tm := teatest.NewTestModel(
				t,
				model{width: sz.width, height: sz.height, renderer: lipgloss.DefaultRenderer()},
				teatest.WithInitialTermSize(sz.width, sz.height),
			)

			teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
				widthStr := fmt.Sprintf("%d", sz.width)
				heightStr := fmt.Sprintf("%d", sz.height)
				return bytes.Contains(bts, []byte(widthStr)) &&
					bytes.Contains(bts, []byte(heightStr))
			}, teatest.WithDuration(time.Second))

			tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
		})
	}
}
```

---

## Key Concepts Learned

1. **Test-Driven Development (TDD)**: Write tests first, watch them fail, implement to make them pass. This ensures every feature is tested and documents expected behavior.

2. **Unit Testing Pure Functions**: `Update` and `View` are pure functions (no side effects), making them ideal for unit testing. Given the same input, they always produce the same output.

3. **Golden File Testing**: Capture expected output as snapshots. Future tests compare against these snapshots to catch visual regressions.

4. **teatest Library**: Charm's official testing library for Bubble Tea apps. Provides `NewTestModel`, `FinalOutput`, `FinalModel`, `WaitFor`, and `Send` for comprehensive TUI testing.

5. **Table-Driven Tests**: Idiomatic Go pattern for testing multiple cases with the same logic. Self-documenting and easy to extend.

6. **Integration vs Unit Tests**: Unit tests verify individual functions in isolation. Integration tests (teatest) verify the full program behavior including message routing.

7. **CI Consistency**: Set `lipgloss.SetColorProfile(termenv.Ascii)` and use `.gitattributes` to ensure tests pass identically on all platforms.

8. **Elm Architecture**: Model holds state, Update handles events, View renders UI. This separation makes code predictable and testable.

9. **Per-Session Rendering**: SSH sessions need their own Lip Gloss renderer to detect terminal capabilities correctly.

10. **Graceful Shutdown**: Signal handling with context timeout ensures clean server shutdown without dropping connections.

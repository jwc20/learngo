# Go Tutorial: Basic TUI App Served Over SSH - Questions & Answers

Based on the video tutorial progression.

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

## Part 3: Handling Window Size (8:00 - 12:00)

### Exercise 9
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add `width` and `height` fields to your model struct.

**Why:** Terminal applications need to know the window dimensions to properly lay out content—centering text, creating borders that fit, handling responsive designs. Bubble Tea automatically sends `WindowSizeMsg` when the terminal resizes, but we need fields to store these values so `View()` can use them for layout calculations. Without storing these, we'd have no way to adapt our UI to different terminal sizes.

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

### Exercise 10
**File:** `main.go`  
**Expected:** Compiles

**Question:** Update your `Update` method to handle `tea.WindowSizeMsg`.

**Why:** Bubble Tea sends a `WindowSizeMsg` both at startup and whenever the user resizes their terminal. The type switch pattern `switch msg := msg.(type)` is idiomatic Go for handling different message types—it's like pattern matching in functional languages. By capturing the width/height into our model, we ensure the `View()` method always has current dimensions for layout. This is reactive programming: state changes flow through Update, then View re-renders.

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

### Exercise 11
**File:** `main.go`  
**Expected:** App displays window dimensions

**Question:** Update your `View` method to display the current window size.

**Why:** This exercises the data flow: WindowSizeMsg → Update stores values → View reads and displays them. Using `fmt.Sprintf` lets us format the integers into a readable string. This pattern of "store in Update, read in View" is fundamental—View should never modify state, only read it. Try resizing your terminal to see the values update in real-time!

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

## Part 4: Handling Key Input (12:00 - 16:00)

### Exercise 12
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add handling for `tea.KeyMsg` in your `Update` method to quit on "q".

**Why:** Every TUI needs a way to exit! `tea.KeyMsg` is sent for each keypress, and `msg.String()` returns a human-readable representation like "q", "enter", "ctrl+c". Returning `tea.Quit` as the command tells Bubble Tea to gracefully shut down—it restores the terminal to normal mode and exits. Handling both "q" and "ctrl+c" follows TUI conventions; users expect both to work.

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

### Exercise 13
**File:** `main.go`  
**Expected:** App shows help text

**Question:** Update your `View` to include a help line.

**Why:** Good UX means users know how to interact with your app. A help line at the bottom is a TUI convention—think of vim's `:help` prompt or htop's key shortcuts. The `\n\n` creates visual separation between content and help. As your app grows, you might move this to a dedicated help bar component, but for now inline text works fine.

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

### Exercise 14
**File:** `main.go`  
**Expected:** App runs in fullscreen mode

**Question:** Modify your `tea.NewProgram()` call to use alternate screen mode.

**Why:** Terminals have two screen buffers: the normal buffer (where your shell history lives) and the alternate screen. Programs like vim, htop, and less use the alternate screen so they can take over the entire terminal without destroying your command history. When the program exits, the terminal automatically restores the previous buffer. `tea.WithAltScreen()` enables this behavior—essential for any "fullscreen" TUI application.

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

## Part 5: Adding Style with Lip Gloss (16:00 - 24:00)

### Exercise 15
**File:** Terminal  
**Expected:** Package downloaded

**Question:** Install the Lip Gloss styling library.

**Why:** Raw terminal output is plain and hard to read. Lip Gloss is Charm's styling library that generates ANSI escape codes for colors, borders, padding, and alignment. It's like CSS for the terminal. Instead of manually writing escape sequences like `\033[1;35m`, you use a fluent API: `style.Bold(true).Foreground(lipgloss.Color("205"))`. This makes styling maintainable and readable.

**Answer:**
```bash
go get github.com/charmbracelet/lipgloss
```

---

### Exercise 16
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the lipgloss import to your file.

**Why:** Go requires explicit imports for every package used. Adding lipgloss to our imports makes its types and functions available. We import it without an alias (unlike bubbletea→tea) because `lipgloss` is already reasonably short and descriptive.

**Answer:**
```go
import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
```

---

### Exercise 17
**File:** `main.go`  
**Expected:** Compiles

**Question:** Create a style variable for styling text.

**Why:** Lip Gloss styles are created by chaining method calls on a `Style` object. Defining styles at the package level (outside any function) means they're created once at startup, not recreated on every render—this is more efficient. Color "205" is from the 256-color palette (a nice pink). You can also use hex colors like `lipgloss.Color("#FF79C6")` or named colors.

**Answer:**
```go
var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```

Full file context (showing placement at package level):
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

### Exercise 18
**File:** `main.go`  
**Expected:** App displays styled text

**Question:** Use your style in the `View` method.

**Why:** `style.Render()` takes a string and returns it wrapped in ANSI escape codes for the configured styling. The result is still a string—Bubble Tea doesn't know or care about the styling, it just outputs whatever string View returns. We keep the help text unstyled for contrast, making the main content stand out.

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

### Exercise 19
**File:** `main.go`  
**Expected:** App displays centered, bordered content

**Question:** Enhance your style with bold, border, and padding.

**Why:** Lip Gloss's fluent API lets you chain multiple style properties. `Bold(true)` makes text bold. `Border()` draws a box around content—`RoundedBorder()` uses Unicode box-drawing characters for smooth corners. `BorderForeground()` colors the border separately from text. `Padding(1, 2)` adds 1 line vertical and 2 characters horizontal padding inside the border. This creates a professional-looking widget.

**Answer:**
```go
var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(1, 2)
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

### Exercise 20
**File:** `main.go`  
**Expected:** Content is centered in the terminal

**Question:** Use lipgloss's `Place` function to center your content.

**Why:** `lipgloss.Place()` positions content within a bounding box. It takes the box dimensions (width, height), horizontal alignment, vertical alignment, and the content to place. By passing `m.width` and `m.height` (our terminal dimensions), we center content in the full window. This is why we stored window size in Exercise 9—layout functions need these dimensions. The result is a polished, centered UI regardless of terminal size.

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

## Part 6: Setting Up the SSH Server with Wish (24:00 - 32:00)

### Exercise 21
**File:** Terminal  
**Expected:** Packages downloaded

**Question:** Install the Wish SSH server library and its dependencies.

**Why:** Wish is an SSH server library that makes it trivial to serve TUI applications over SSH. Instead of users running your binary locally, they can `ssh yourserver.com` and interact with your app remotely. This is how services like charm.sh work. We also need the `ssh` package for session handling and `log` for structured logging. Together, these enable building production-ready SSH applications.

**Answer:**
```bash
go get github.com/charmbracelet/wish
go get github.com/charmbracelet/ssh
go get github.com/charmbracelet/log
```

---

### Exercise 22
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add all the new imports to your file.

**Why:** Each package serves a specific purpose: `wish` provides the SSH server framework; `wish/bubbletea` integrates Bubble Tea apps with SSH sessions; `wish/activeterm` ensures clients have a proper PTY (pseudo-terminal); `wish/logging` logs connections for debugging; `ssh` provides session types; `log` is Charm's structured logger. The standard library imports (`context`, `net`, `os/signal`, etc.) are for server lifecycle management, signal handling, and graceful shutdown.

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

### Exercise 23
**File:** `main.go`  
**Expected:** Compiles

**Question:** Define constants for your SSH server.

**Why:** Constants make configuration clear and changeable from one place. Port 23234 is a high port (>1024) so we don't need root privileges to bind to it. Using `localhost` restricts connections to the local machine—in production you'd use `0.0.0.0` to accept remote connections. Separating these as constants follows the principle of making magic values explicit.

**Answer:**
```go
const (
	host = "localhost"
	port = "23234"
)
```

Full file context (showing placement after imports, before type definitions):
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

### Exercise 24
**File:** `main.go`  
**Expected:** Compiles

**Question:** Create a `teaHandler` function.

**Why:** The `teaHandler` is called by Wish for each SSH connection to create a fresh Bubble Tea program. This is crucial: each user gets their own model instance, so one user's actions don't affect another's. The function signature `func(ssh.Session) (tea.Model, []tea.ProgramOption)` is required by the bubbletea middleware. Returning options like `tea.WithAltScreen()` configures how the program runs for that session.

**Answer:**
```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
```

Full file context (showing placement before main):
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

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
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

### Exercise 25
**File:** `main.go`  
**Expected:** Compiles

**Question:** Enhance `teaHandler` to get terminal info and add a `term` field to model.

**Why:** SSH sessions carry terminal information: the client's terminal type (xterm-256color, etc.) and initial window dimensions. `sess.Pty()` returns this PTY (pseudo-terminal) information. By capturing `pty.Term`, we can display what terminal the user is connecting with—useful for debugging or adapting behavior. Initializing `width`/`height` from `pty.Window` ensures the first render has correct dimensions before any WindowSizeMsg arrives.

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

Also update View to display the terminal type:
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

Full file context:
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

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	Padding(1, 2)

type model struct {
	term   string
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
	m := model{
		term:   pty.Term,
		width:  pty.Window.Width,
		height: pty.Window.Height,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
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

## Part 7: Creating the Wish Server (32:00 - 36:00)

### Exercise 26
**File:** `main.go`  
**Expected:** Compiles (but main still has old code)

**Question:** Replace your `main` function to create a Wish server.

**Why:** We're transitioning from a local TUI to an SSH-served TUI. `wish.NewServer()` creates an SSH server with configuration options. `WithAddress` sets the listen address using `net.JoinHostPort` (which properly handles IPv6). `WithHostKeyPath` specifies where to store/load the server's SSH host key—SSH requires servers to have a persistent identity so clients can verify they're connecting to the right server.

**Answer:**
```go
func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
	}
}
```

---

### Exercise 27
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the middleware stack to your `wish.NewServer` call.

**Why:** Wish uses a middleware pattern—each middleware wraps the next, processing the connection before passing it along. Order matters and is reverse: last added runs first. `logging.Middleware()` logs connections (runs first). `activeterm.Middleware()` ensures the client has a PTY; non-interactive connections (like scp) get rejected here. `bubbletea.Middleware(teaHandler)` finally serves your TUI app. This layered approach keeps concerns separated and composable.

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

### Exercise 28
**File:** `main.go`  
**Expected:** Compiles

**Question:** Set up signal handling for graceful shutdown.

**Why:** Production servers need to shut down cleanly when receiving interrupt signals (Ctrl+C) or termination signals (from systemd, Docker, etc.). Go's `signal.Notify` redirects OS signals to a channel. Using a buffered channel (`make(chan os.Signal, 1)`) ensures the signal isn't lost if we're not immediately ready to receive. `SIGINT` comes from Ctrl+C, `SIGTERM` from kill commands. This pattern is standard in Go network servers.

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

### Exercise 29
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add a log message indicating the server is starting.

**Why:** Observability is crucial for production services. Logging when the server starts confirms successful initialization and tells operators what address to connect to. Charm's `log` package uses structured logging—key-value pairs like `"host", host` make logs parseable by tools like Loki or Splunk. This is better than `fmt.Printf` for production systems.

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
}
```

---

### Exercise 30
**File:** `main.go`  
**Expected:** Server starts listening

**Question:** Start the server in a goroutine.

**Why:** `ListenAndServe()` blocks forever while accepting connections—if we called it directly, the rest of main wouldn't execute. Running it in a goroutine lets the main goroutine continue to handle shutdown signals. The error check excludes `ssh.ErrServerClosed` because that's expected when we call `Shutdown()`—it's not an error, just confirmation the server stopped. If an unexpected error occurs, we send to `done` to trigger shutdown.

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

### Exercise 31
**File:** `main.go`  
**Expected:** Server shuts down gracefully on Ctrl+C

**Question:** Implement graceful shutdown.

**Why:** `<-done` blocks until a signal arrives—this keeps main alive while the server runs. On receiving a signal, we log the shutdown, create a context with 30-second timeout (giving active connections time to finish), and call `Shutdown()`. The context timeout prevents hanging forever if a connection won't close. `defer cancel()` ensures the context resources are cleaned up. This graceful shutdown pattern prevents data loss and angry users.

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

## Part 8: Testing Your SSH App (36:00+)

### Exercise 32
**File:** Terminal  
**Expected:** Server starts successfully

**Question:** Run your application.

**Why:** This tests that everything compiles and the server can bind to the port. The first run will generate an SSH host key file (`.ssh/id_ed25519`). If you see the "Starting SSH server" log, the server is ready to accept connections. If you get a "bind: address already in use" error, another process is using port 23234.

**Answer:**
```bash
go run main.go
```

Output:
```
INFO Starting SSH server host=localhost port=23234
```

---

### Exercise 33
**File:** Terminal (new window)  
**Expected:** TUI appears over SSH

**Question:** SSH into your app.

**Why:** This is the moment of truth—connecting to your TUI over SSH. The `-p` flag specifies the non-standard port. SSH will negotiate encryption, authenticate (likely prompting about the unknown host key), then pass your terminal to Wish. If everything works, you'll see your centered, styled window size display. You're now running Bubble Tea over SSH!

**Answer:**
```bash
ssh localhost -p 23234
```

---

### Exercise 34
**File:** SSH session  
**Expected:** Dimensions update

**Question:** Test window resize handling.

**Why:** Resize events are forwarded over SSH via the PTY. When you resize your terminal, the SSH client sends a "window change" message to the server, which Wish converts into a `tea.WindowSizeMsg` for your app. Seeing the dimensions update confirms the full event pipeline works: SSH client → network → Wish → Bubble Tea → your Update → View → render.

**Answer:**
Resize your terminal window while connected. The dimensions should update automatically.

---

### Exercise 35
**File:** SSH session  
**Expected:** Connection closes

**Question:** Quit the app.

**Why:** Pressing "q" triggers your KeyMsg handler, which returns `tea.Quit`. Bubble Tea exits cleanly, which signals Wish that the session is done, which closes the SSH connection. The server keeps running to accept new connections. This per-session lifecycle is key to multi-user SSH apps.

**Answer:**
Press `q`. Output: `Connection to localhost closed.`

---

### Exercise 36
**File:** Server terminal  
**Expected:** Server stops

**Question:** Stop the server.

**Why:** Pressing Ctrl+C sends SIGINT, which our signal handler receives on the `done` channel. The graceful shutdown code runs, closing the listener and waiting for active connections. You should see the "Stopping SSH server" log. Any connected users will be disconnected. The clean exit with no errors confirms our shutdown logic works.

**Answer:**
Press Ctrl+C in the server terminal.
Output:
```
INFO Stopping SSH server
```

---

## Bonus Exercises

### Exercise 37
**File:** `main.go`  
**Expected:** Styled output works over SSH

**Question:** Fix styles for SSH by using per-session renderer.

**Why:** Here's a subtle bug: Lip Gloss defaults to detecting terminal capabilities from `os.Stdout`. But over SSH, each session has its *own* output stream—not stdout! A package-level style doesn't know about SSH session capabilities. The fix is `bubbletea.MakeRenderer(sess)`, which creates a renderer aware of the specific session's terminal. Without this, colors might not work for some SSH clients, or you might get garbled output.

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

### Exercise 38
**File:** `~/.ssh/config`  
**Expected:** No known_hosts warnings

**Question:** Configure SSH to skip host key checking for localhost.

**Why:** During development, you'll regenerate the server's host key often (by deleting `.ssh/id_ed25519`). Each regeneration causes SSH to warn about a "changed host key"—a legitimate security feature that prevents man-in-the-middle attacks. For localhost development only, disabling this check saves frustration. Never do this for remote servers! The `UserKnownHostsFile /dev/null` prevents saving keys, and `StrictHostKeyChecking no` skips verification.

**Answer:**
```
Host localhost
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
```

---

## Complete Final Code

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

---

## Key Concepts Learned

1. **Elm Architecture**: Model holds state, Update handles events, View renders UI. This separation makes code predictable and testable.

2. **Bubble Tea Messages**: `tea.WindowSizeMsg` for resize events, `tea.KeyMsg` for keyboard input, `tea.Cmd` for side effects. Messages flow through Update in a single direction.

3. **Lip Gloss**: Terminal styling with a fluent API. Colors, borders, padding, and layout without manual ANSI codes. Per-session renderers are essential for SSH.

4. **Wish**: SSH server framework that integrates with Bubble Tea. Each connection gets its own program instance.

5. **Middleware Pattern**: Composable handlers (logging → activeterm → bubbletea). Order matters—last added runs first.

6. **Per-Session State**: Each SSH connection gets its own model and renderer. No shared mutable state between users.

7. **Graceful Shutdown**: Signal handling with context timeout. Production servers need to close cleanly without dropping connections.

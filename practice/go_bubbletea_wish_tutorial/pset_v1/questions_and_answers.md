# Go TUI App Served Over SSH - Questions & Answers

Based on the tutorial: "Go Tutorial: Basic TUI app served over SSH"

---

## Section 1: Project Setup

### Exercise 1
**File:** Terminal  
**Expected:** New Go module initialized

**Question:** Initialize a new Go module called `ssh-tui-app`.

**Answer:**
```bash
mkdir ssh-tui-app
cd ssh-tui-app
go mod init ssh-tui-app
```

---

### Exercise 2
**File:** Terminal  
**Expected:** Dependencies downloaded

**Question:** Install the required Charmbracelet dependencies.

**Answer:**
```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/wish
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/log
go get github.com/charmbracelet/ssh
```

Or all at once:
```bash
go get github.com/charmbracelet/bubbletea github.com/charmbracelet/wish github.com/charmbracelet/lipgloss github.com/charmbracelet/log github.com/charmbracelet/ssh
```

---

### Exercise 3
**File:** `main.go` (create new)  
**Expected:** File compiles with no errors

**Question:** Create `main.go` with the package declaration and all required imports.

**Answer:**
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

func main() {
}
```

---

## Section 2: Understanding the Elm Architecture

### Exercise 4
**File:** `main.go`  
**Expected:** Compiles (struct defined but unused)

**Question:** Create a `model` struct with fields for terminal type, width, and height.

**Answer:**
```go
type model struct {
	term   string
	width  int
	height int
}
```

---

### Exercise 5
**File:** `main.go`  
**Expected:** Compiles (methods defined)

**Question:** Implement the three required Bubble Tea methods as stubs: `Init()`, `Update()`, and `View()`.

**Answer:**
```go
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return ""
}
```

---

### Exercise 6
**File:** `main.go`  
**Expected:** Compiles, View returns a string

**Question:** Implement the `View()` method to display the terminal information.

**Answer:**
```go
func (m model) View() string {
	s := fmt.Sprintf("Your term is %s\r\n", m.term)
	s += fmt.Sprintf("Your window size is %d x %d\r\n", m.width, m.height)
	return s
}
```

---

## Section 3: Handling Messages in Update

### Exercise 7
**File:** `main.go`  
**Expected:** Compiles

**Question:** Implement `Update` to handle `tea.WindowSizeMsg` and update the model's dimensions.

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

---

### Exercise 8
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add handling for `tea.KeyMsg` to quit when user presses "q".

**Answer:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}
```

---

## Section 4: Creating the Wish SSH Server

### Exercise 9
**File:** `main.go`  
**Expected:** Compiles

**Question:** Define constants for your SSH server host and port.

**Answer:**
```go
const (
	host = "localhost"
	port = "23234"
)
```

---

### Exercise 10
**File:** `main.go`  
**Expected:** Compiles

**Question:** Create a `teaHandler` function that creates Bubble Tea programs for each SSH session.

**Answer:**
```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()
	m := model{
		term: pty.Term,
	}
	return m, nil
}
```

---

### Exercise 11
**File:** `main.go`  
**Expected:** Compiles but won't run (main is empty)

**Question:** In `main()`, create a new Wish server with address and host key options.

**Answer:**
```go
func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}
}
```

---

### Exercise 12
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the three middlewares to your Wish server.

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
		log.Error("Could not start server", "error", err)
	}
}
```

---

## Section 5: Graceful Shutdown

### Exercise 13
**File:** `main.go`  
**Expected:** Compiles

**Question:** Set up signal handling for graceful shutdown.

**Answer:**
```go
done := make(chan os.Signal, 1)
signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
```

---

### Exercise 14
**File:** `main.go`  
**Expected:** Compiles

**Question:** Log that the server is starting.

**Answer:**
```go
log.Info("Starting SSH server", "host", host, "port", port)
```

---

### Exercise 15
**File:** `main.go`  
**Expected:** Server starts and listens

**Question:** Start the SSH server in a goroutine.

**Answer:**
```go
go func() {
	if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
		done <- nil
	}
}()
```

---

### Exercise 16
**File:** `main.go`  
**Expected:** Server shuts down gracefully on Ctrl+C

**Question:** Implement graceful shutdown by blocking on the done channel and shutting down with a timeout.

**Answer:**
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

## Section 6: Testing the Application

### Exercise 17
**File:** Terminal  
**Expected:** Server starts, shows "Starting SSH server"

**Question:** Run your application.

**Answer:**
```bash
go run main.go
```

Output should show:
```
INFO Starting SSH server host=localhost port=23234
```

---

### Exercise 18
**File:** Terminal (new window)  
**Expected:** TUI displays in terminal

**Question:** SSH into your app from another terminal.

**Answer:**
```bash
ssh localhost -p 23234
```

---

### Exercise 19
**File:** SSH session  
**Expected:** TUI updates with new dimensions

**Question:** Test window resize handling.

**Answer:**
Resize your terminal window while connected. The displayed dimensions should update automatically.

---

### Exercise 20
**File:** SSH session  
**Expected:** SSH connection closes cleanly

**Question:** Quit the application.

**Answer:**
Press `q` to quit. The message "Connection to localhost closed." should appear.

---

## Section 7: Enhancements with Lip Gloss

### Exercise 21
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add a `renderer` field to your model struct.

**Answer:**
```go
type model struct {
	term     string
	width    int
	height   int
	renderer *lipgloss.Renderer
}
```

---

### Exercise 22
**File:** `main.go`  
**Expected:** Compiles

**Question:** Update `teaHandler` to create a renderer for the session.

**Answer:**
```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()
	renderer := bubbletea.MakeRenderer(sess)
	m := model{
		term:     pty.Term,
		renderer: renderer,
	}
	return m, nil
}
```

---

### Exercise 23
**File:** `main.go`  
**Expected:** TUI displays with styled text

**Question:** Use the renderer to create styled output in `View()`.

**Answer:**
```go
func (m model) View() string {
	titleStyle := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	s := titleStyle.Render("SSH TUI App") + "\r\n\r\n"
	s += fmt.Sprintf("Your term is %s\r\n", m.term)
	s += fmt.Sprintf("Your window size is %d x %d\r\n", m.width, m.height)
	return s
}
```

---

### Exercise 24
**File:** `main.go`  
**Expected:** TUI displays with background color block

**Question:** Add a styled info box with background color.

**Answer:**
```go
func (m model) View() string {
	titleStyle := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	infoStyle := m.renderer.NewStyle().
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("230")).
		Padding(1, 2)

	s := titleStyle.Render("SSH TUI App") + "\r\n\r\n"
	
	info := fmt.Sprintf("Term: %s\nSize: %d x %d", m.term, m.width, m.height)
	s += infoStyle.Render(info) + "\r\n\r\n"
	
	s += "Press 'q' to quit\r\n"
	return s
}
```

---

## Section 8: Adding More Interactivity

### Exercise 25
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add a `counter` field to your model.

**Answer:**
```go
type model struct {
	term     string
	width    int
	height   int
	renderer *lipgloss.Renderer
	counter  int
}
```

---

### Exercise 26
**File:** `main.go`  
**Expected:** Counter increments on key press

**Question:** Update the `Update` method to increment counter on any key except "q".

**Answer:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
		m.counter++
	}
	return m, nil
}
```

---

### Exercise 27
**File:** `main.go`  
**Expected:** Help text displays at bottom

**Question:** Update `View()` to show counter and help text.

**Answer:**
```go
func (m model) View() string {
	titleStyle := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	infoStyle := m.renderer.NewStyle().
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("230")).
		Padding(1, 2)

	helpStyle := m.renderer.NewStyle().
		Foreground(lipgloss.Color("241"))

	s := titleStyle.Render("SSH TUI App") + "\r\n\r\n"
	
	info := fmt.Sprintf("Term: %s\nSize: %d x %d\nCounter: %d", m.term, m.width, m.height, m.counter)
	s += infoStyle.Render(info) + "\r\n\r\n"
	
	s += helpStyle.Render("Press 'q' to quit • Press any key to increment counter") + "\r\n"
	return s
}
```

---

## Section 9: SSH Configuration Tips

### Exercise 28
**File:** `~/.ssh/config`  
**Expected:** No more known_hosts warnings during development

**Question:** Configure SSH to skip host key checking for localhost.

**Answer:**
Add to `~/.ssh/config`:
```
Host localhost
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
```

---

### Exercise 29
**File:** `.ssh/` directory  
**Expected:** Host key is generated

**Question:** Verify the host key was generated.

**Answer:**
```bash
ls -la .ssh/
```

Should show `id_ed25519` file created by Wish on first run.

---

## Final Complete Code

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
	counter  int
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
		if msg.String() == "q" {
			return m, tea.Quit
		}
		m.counter++
	}
	return m, nil
}

func (m model) View() string {
	titleStyle := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	infoStyle := m.renderer.NewStyle().
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("230")).
		Padding(1, 2)

	helpStyle := m.renderer.NewStyle().
		Foreground(lipgloss.Color("241"))

	s := titleStyle.Render("SSH TUI App") + "\r\n\r\n"

	info := fmt.Sprintf("Term: %s\nSize: %d x %d\nCounter: %d", m.term, m.width, m.height, m.counter)
	s += infoStyle.Render(info) + "\r\n\r\n"

	s += helpStyle.Render("Press 'q' to quit • Press any key to increment counter") + "\r\n"
	return s
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()
	renderer := bubbletea.MakeRenderer(sess)
	m := model{
		term:     pty.Term,
		renderer: renderer,
	}
	return m, nil
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
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
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

## Key Concepts Summary

### The Elm Architecture (Model-Update-View)
- **Model**: Holds all application state
- **Update**: Handles messages/events and returns new state
- **View**: Renders the UI as a string based on current state

### Bubble Tea Message Types
- `tea.WindowSizeMsg`: Terminal window resize events
- `tea.KeyMsg`: Keyboard input events
- `tea.Cmd`: Commands for side effects (tea.Quit, etc.)

### Wish Middleware Stack
```
Request → logging.Middleware() → activeterm.Middleware() → bubbletea.Middleware() → Your App
```
Middlewares execute in reverse order (last added = first executed).

### Per-Session State
Each SSH connection gets its own:
- `model` instance (state)
- `tea.Program` (Bubble Tea runtime)
- `lipgloss.Renderer` (for consistent styling)

### Graceful Shutdown Pattern
1. Set up signal channel
2. Start server in goroutine
3. Block on signal channel
4. Call `Shutdown()` with context timeout

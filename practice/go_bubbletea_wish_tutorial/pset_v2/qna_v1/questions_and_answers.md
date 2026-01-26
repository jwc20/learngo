# Go Tutorial: Basic TUI App Served Over SSH - Questions & Answers

Based on the video tutorial progression.

---

## Part 1: Project Setup (0:00 - 4:00)

### Exercise 1
**File:** Terminal  
**Expected:** New directory and Go module created

**Question:** Create a new directory for your project and initialize a Go module.

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

**Answer:**
```bash
go get github.com/charmbracelet/bubbletea
```

---

### Exercise 3
**File:** `main.go` (create new)  
**Expected:** File compiles with empty main

**Question:** Create `main.go` with the package declaration, import for bubbletea (aliased as `tea`), and an empty main function.

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

**Answer:**
```go
type model struct{}
```

---

### Exercise 5
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the `Init` method stub to your model.

**Answer:**
```go
func (m model) Init() tea.Cmd {
	return nil
}
```

---

### Exercise 6
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the `Update` method stub to your model.

**Answer:**
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
```

---

### Exercise 7
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the `View` method stub to your model.

**Answer:**
```go
func (m model) View() string {
	return "Hello Bubble Tea!"
}
```

---

### Exercise 8
**File:** `main.go`  
**Expected:** App runs and displays "Hello Bubble Tea!"

**Question:** In your `main` function, create a new Bubble Tea program and run it.

**Answer:**
```go
func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
```

Add imports:
```go
import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)
```

---

## Part 3: Handling Window Size (8:00 - 12:00)

### Exercise 9
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add `width` and `height` fields to your model struct.

**Answer:**
```go
type model struct {
	width  int
	height int
}
```

---

### Exercise 10
**File:** `main.go`  
**Expected:** Compiles

**Question:** Update your `Update` method to handle `tea.WindowSizeMsg`.

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

### Exercise 11
**File:** `main.go`  
**Expected:** App displays window dimensions

**Question:** Update your `View` method to display the current window size.

**Answer:**
```go
func (m model) View() string {
	return fmt.Sprintf("Window size: %d x %d", m.width, m.height)
}
```

---

## Part 4: Handling Key Input (12:00 - 16:00)

### Exercise 12
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add handling for `tea.KeyMsg` in your `Update` method to quit on "q".

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

---

### Exercise 13
**File:** `main.go`  
**Expected:** App shows help text

**Question:** Update your `View` to include a help line.

**Answer:**
```go
func (m model) View() string {
	s := fmt.Sprintf("Window size: %d x %d\n\n", m.width, m.height)
	s += "Press q to quit"
	return s
}
```

---

### Exercise 14
**File:** `main.go`  
**Expected:** App runs in fullscreen mode

**Question:** Modify your `tea.NewProgram()` call to use alternate screen mode.

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

---

## Part 5: Adding Style with Lip Gloss (16:00 - 24:00)

### Exercise 15
**File:** Terminal  
**Expected:** Package downloaded

**Question:** Install the Lip Gloss styling library.

**Answer:**
```bash
go get github.com/charmbracelet/lipgloss
```

---

### Exercise 16
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add the lipgloss import to your file.

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

**Answer:**
```go
var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```

---

### Exercise 18
**File:** `main.go`  
**Expected:** App displays styled text

**Question:** Use your style in the `View` method.

**Answer:**
```go
func (m model) View() string {
	s := style.Render(fmt.Sprintf("Window size: %d x %d", m.width, m.height))
	s += "\n\nPress q to quit"
	return s
}
```

---

### Exercise 19
**File:** `main.go`  
**Expected:** App displays centered, bordered content

**Question:** Enhance your style with bold, border, and padding.

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

### Exercise 20
**File:** `main.go`  
**Expected:** Content is centered in the terminal

**Question:** Use lipgloss's `Place` function to center your content.

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

---

## Part 6: Setting Up the SSH Server with Wish (24:00 - 32:00)

### Exercise 21
**File:** Terminal  
**Expected:** Packages downloaded

**Question:** Install the Wish SSH server library and its dependencies.

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

**Answer:**
```go
const (
	host = "localhost"
	port = "23234"
)
```

---

### Exercise 24
**File:** `main.go`  
**Expected:** Compiles

**Question:** Create a `teaHandler` function.

**Answer:**
```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
```

---

### Exercise 25
**File:** `main.go`  
**Expected:** Compiles

**Question:** Enhance `teaHandler` to get terminal info and add a `term` field to model.

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

Update View to show term:
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

## Part 7: Creating the Wish Server (32:00 - 36:00)

### Exercise 26
**File:** `main.go`  
**Expected:** Compiles (but main still has old code)

**Question:** Replace your `main` function to create a Wish server.

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

**Answer:**
```go
done := make(chan os.Signal, 1)
signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
```

---

### Exercise 29
**File:** `main.go`  
**Expected:** Compiles

**Question:** Add a log message indicating the server is starting.

**Answer:**
```go
log.Info("Starting SSH server", "host", host, "port", port)
```

---

### Exercise 30
**File:** `main.go`  
**Expected:** Server starts listening

**Question:** Start the server in a goroutine.

**Answer:**
```go
go func() {
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
		done <- nil
	}
}()
```

---

### Exercise 31
**File:** `main.go`  
**Expected:** Server shuts down gracefully on Ctrl+C

**Question:** Implement graceful shutdown.

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

## Part 8: Testing Your SSH App (36:00+)

### Exercise 32
**File:** Terminal  
**Expected:** Server starts successfully

**Question:** Run your application.

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

**Answer:**
```bash
ssh localhost -p 23234
```

---

### Exercise 34
**File:** SSH session  
**Expected:** Dimensions update

**Question:** Test window resize handling.

**Answer:**
Resize your terminal window while connected. The dimensions should update automatically.

---

### Exercise 35
**File:** SSH session  
**Expected:** Connection closes

**Question:** Quit the app.

**Answer:**
Press `q`. Output: `Connection to localhost closed.`

---

### Exercise 36
**File:** Server terminal  
**Expected:** Server stops

**Question:** Stop the server.

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

1. **Elm Architecture**: Model holds state, Update handles events, View renders UI
2. **Bubble Tea Messages**: `tea.WindowSizeMsg`, `tea.KeyMsg`, `tea.Cmd`
3. **Lip Gloss**: Styling terminal output with colors, borders, padding
4. **Wish**: SSH server framework that integrates with Bubble Tea
5. **Middleware Pattern**: Composable handlers (logging → activeterm → bubbletea)
6. **Per-Session State**: Each SSH connection gets its own model and renderer
7. **Graceful Shutdown**: Signal handling with context timeout

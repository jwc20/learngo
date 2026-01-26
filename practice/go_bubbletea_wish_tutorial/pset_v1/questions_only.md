# Go TUI App Served Over SSH - Practice Questions

Based on the tutorial: "Go Tutorial: Basic TUI app served over SSH"

---

## Section 1: Project Setup

### Exercise 1
**File:** Terminal  
**Expected:** New Go module initialized

Initialize a new Go module called `ssh-tui-app`.

---

### Exercise 2
**File:** Terminal  
**Expected:** Dependencies downloaded

Install the required Charmbracelet dependencies:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/wish` - SSH server framework
- `github.com/charmbracelet/lipgloss` - Styling library
- `github.com/charmbracelet/log` - Logging library
- `github.com/charmbracelet/ssh` - SSH types

---

### Exercise 3
**File:** `main.go` (create new)  
**Expected:** File compiles with no errors

Create `main.go` with the package declaration and all required imports for a Wish/Bubble Tea application. You'll need imports for:
- Context, errors, fmt, net, os, os/signal, syscall, time
- All the Charmbracelet packages
- The wish middlewares: activeterm, bubbletea, logging

---

## Section 2: Understanding the Elm Architecture

### Exercise 4
**File:** `main.go`  
**Expected:** Compiles (struct defined but unused)

The Elm Architecture requires a **model** that holds application state. Create a `model` struct with the following fields:
- `term` (string) - the terminal type
- `width` (int) - terminal width
- `height` (int) - terminal height

---

### Exercise 5
**File:** `main.go`  
**Expected:** Compiles (methods defined)

Bubble Tea models must implement three methods:
1. `Init() tea.Cmd` - returns initial command (return `nil` for now)
2. `Update(msg tea.Msg) (tea.Model, tea.Cmd)` - handles messages (just return model and nil for now)
3. `View() string` - renders UI (return empty string for now)

Implement these three methods on your `model` struct as stubs.

---

### Exercise 6
**File:** `main.go`  
**Expected:** Compiles, View returns a string

Implement the `View()` method to display the terminal information. Use `fmt.Sprintf` to return a string showing:
- "Your term is [term]"
- "Your window size is [width] x [height]"

Each line should end with `\r\n` for proper terminal display.

---

## Section 3: Handling Messages in Update

### Exercise 7
**File:** `main.go`  
**Expected:** Compiles

The `Update` method receives messages via a type switch. Implement `Update` to handle `tea.WindowSizeMsg`:
- When you receive a `tea.WindowSizeMsg`, update the model's `width` and `height` fields
- Return the updated model and `nil` for the command

Use a type switch: `switch msg := msg.(type) { ... }`

---

### Exercise 8
**File:** `main.go`  
**Expected:** Compiles

Add handling for `tea.KeyMsg` in your `Update` method. When the user presses "q", return the model along with `tea.Quit` to exit the application.

Check for the key using: `case tea.KeyMsg:` and `if msg.String() == "q"`

---

## Section 4: Creating the Wish SSH Server

### Exercise 9
**File:** `main.go`  
**Expected:** Compiles

Define constants for your SSH server:
- `host` = "localhost"
- `port` = "23234"

---

### Exercise 10
**File:** `main.go`  
**Expected:** Compiles

Create a `teaHandler` function that Wish will use to create Bubble Tea programs for each SSH session. The function signature should be:

```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption)
```

Inside this function:
1. Get the PTY information using `sess.Pty()`
2. Create a new `model` with the terminal type from `pty.Term`
3. Return the model and `nil` for options

---

### Exercise 11
**File:** `main.go`  
**Expected:** Compiles but won't run (main is empty)

In the `main()` function, create a new Wish server using `wish.NewServer()` with the following options:
- `wish.WithAddress(net.JoinHostPort(host, port))` - set the address
- `wish.WithHostKeyPath(".ssh/id_ed25519")` - set the host key path
- `wish.WithMiddleware(...)` - add middlewares (next exercise)

Handle the error if server creation fails using `log.Error()`.

---

### Exercise 12
**File:** `main.go`  
**Expected:** Compiles

Add three middlewares to your Wish server (order matters - last one executes first):
1. `bubbletea.Middleware(teaHandler)` - serves the Bubble Tea app
2. `activeterm.Middleware()` - ensures active terminal (PTY) connection
3. `logging.Middleware()` - logs connections

Note: Middlewares are composed from first to last, meaning the last one listed executes first.

---

## Section 5: Graceful Shutdown

### Exercise 13
**File:** `main.go`  
**Expected:** Compiles

Set up signal handling for graceful shutdown:
1. Create a channel: `done := make(chan os.Signal, 1)`
2. Register for interrupt signals: `signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)`

---

### Exercise 14
**File:** `main.go`  
**Expected:** Compiles

Log that the server is starting using `log.Info()` with the host and port information.

---

### Exercise 15
**File:** `main.go`  
**Expected:** Server starts and listens

Start the SSH server in a goroutine:
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

After starting the server goroutine:
1. Block on the `done` channel: `<-done`
2. Log that the server is stopping
3. Create a context with timeout: `ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)`
4. Defer `cancel()`
5. Call `s.Shutdown(ctx)` to gracefully shut down

---

## Section 6: Testing the Application

### Exercise 17
**File:** Terminal  
**Expected:** Server starts, shows "Starting SSH server"

Run your application:
```bash
go run main.go
```

---

### Exercise 18
**File:** Terminal (new window)  
**Expected:** TUI displays in terminal, shows term type and window size

In a separate terminal, SSH into your app:
```bash
ssh localhost -p 23234
```

You should see your terminal type and window size displayed.

---

### Exercise 19
**File:** SSH session  
**Expected:** TUI updates with new dimensions

While connected via SSH, resize your terminal window. The displayed width and height should update automatically (this is handled by `tea.WindowSizeMsg`).

---

### Exercise 20
**File:** SSH session  
**Expected:** SSH connection closes cleanly

Press `q` to quit the application. The SSH connection should close.

---

## Section 7: Enhancements with Lip Gloss

### Exercise 21
**File:** `main.go`  
**Expected:** Compiles

Add a `renderer` field of type `*lipgloss.Renderer` to your `model` struct. This will be used for styling output per-session.

---

### Exercise 22
**File:** `main.go`  
**Expected:** Compiles

Update your `teaHandler` function to create a renderer for the session:
```go
renderer := bubbletea.MakeRenderer(sess)
```

Pass this renderer to your model when creating it.

---

### Exercise 23
**File:** `main.go`  
**Expected:** TUI displays with styled text

In your `View()` method, use the renderer to create a styled output. Create a style with:
- Bold text
- A foreground color (e.g., lipgloss color "212" for pink)

Apply this style to your output text.

---

### Exercise 24
**File:** `main.go`  
**Expected:** TUI displays with background color block

Add a background color block to your view using lipgloss:
```go
style := m.renderer.NewStyle().
    Background(lipgloss.Color("63")).
    Padding(1, 2)
```

---

## Section 8: Adding More Interactivity

### Exercise 25
**File:** `main.go`  
**Expected:** Compiles

Add a `counter` field (int) to your model to track how many times a key has been pressed.

---

### Exercise 26
**File:** `main.go`  
**Expected:** Counter increments on key press

Update your `Update` method to increment the counter on any key press (except "q"). Display the counter in your `View()`.

---

### Exercise 27
**File:** `main.go`  
**Expected:** Help text displays at bottom

Add help text to your `View()` that shows available commands:
- "Press 'q' to quit"
- "Press any other key to increment counter"

---

## Section 9: SSH Configuration Tips

### Exercise 28
**File:** `~/.ssh/config`  
**Expected:** No more known_hosts warnings during development

When developing locally, you may get annoying "host key changed" warnings. Add this to your SSH config to avoid clearing `known_hosts`:

```
Host localhost
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
```

---

### Exercise 29
**File:** `.ssh/` directory  
**Expected:** Host key is generated

The first time you run your server, Wish will generate an SSH host key at `.ssh/id_ed25519`. Verify this file exists after running your server once.

---

## Final Project Structure

After completing all exercises, you should have:

```
ssh-tui-app/
├── go.mod
├── go.sum
├── main.go
└── .ssh/
    └── id_ed25519
```

Your `main.go` should contain:
- A `model` struct with term, width, height, renderer, and counter fields
- `Init()`, `Update()`, and `View()` methods implementing `tea.Model`
- A `teaHandler` function for creating per-session Bubble Tea programs
- A `main()` function that creates and runs a Wish SSH server with proper middleware and graceful shutdown

---

## Key Concepts Covered

1. **The Elm Architecture**: Model-Update-View pattern for TUI applications
2. **Bubble Tea**: Go framework for building TUIs based on Elm Architecture
3. **Wish**: SSH server framework that integrates with Bubble Tea
4. **Middleware**: Composable handlers for SSH connections (logging, terminal validation, Bubble Tea serving)
5. **Lip Gloss**: Terminal styling library for colors, padding, and formatting
6. **Graceful Shutdown**: Signal handling for clean server termination
7. **Per-Session State**: Each SSH connection gets its own model instance

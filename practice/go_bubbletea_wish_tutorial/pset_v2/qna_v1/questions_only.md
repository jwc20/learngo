# Go Tutorial: Basic TUI App Served Over SSH - Practice Questions

Based on the video tutorial progression.

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

Run with `go run main.go` to test.

---

## Part 3: Handling Window Size (8:00 - 12:00)

### Exercise 9
**File:** `main.go`  
**Expected:** Compiles

Add `width` and `height` fields (both `int`) to your model struct to store the terminal dimensions.

---

### Exercise 10
**File:** `main.go`  
**Expected:** Compiles

Update your `Update` method to handle `tea.WindowSizeMsg`. Use a type switch:
```go
switch msg := msg.(type) {
case tea.WindowSizeMsg:
    // update m.width and m.height
}
```
The message has `Width` and `Height` fields.

---

### Exercise 11
**File:** `main.go`  
**Expected:** App displays window dimensions

Update your `View` method to display the current window size:
```
"Window size: [width] x [height]"
```

---

## Part 4: Handling Key Input (12:00 - 16:00)

### Exercise 12
**File:** `main.go`  
**Expected:** Compiles

Add handling for `tea.KeyMsg` in your `Update` method. When the user presses "q", return `tea.Quit` as the command to exit the program.

Hint: Check `msg.String() == "q"` or use `msg.Type == tea.KeyCtrlC` for Ctrl+C.

---

### Exercise 13
**File:** `main.go`  
**Expected:** App shows help text

Update your `View` to include a help line at the bottom:
```
"Press q to quit"
```

---

### Exercise 14
**File:** `main.go`  
**Expected:** App runs in fullscreen mode

Modify your `tea.NewProgram()` call to use alternate screen (fullscreen) mode:
```go
tea.NewProgram(m, tea.WithAltScreen())
```

---

## Part 5: Adding Style with Lip Gloss (16:00 - 24:00)

### Exercise 15
**File:** Terminal  
**Expected:** Package downloaded

Install the Lip Gloss styling library:
```
go get github.com/charmbracelet/lipgloss
```

---

### Exercise 16
**File:** `main.go`  
**Expected:** Compiles

Add the lipgloss import to your file.

---

### Exercise 17
**File:** `main.go`  
**Expected:** Compiles

Create a style variable at the package level for styling text:
```go
var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```
Color "205" is a nice pink color.

---

### Exercise 18
**File:** `main.go`  
**Expected:** App displays styled text

Use your style in the `View` method to render the window size text:
```go
style.Render("your text here")
```

---

### Exercise 19
**File:** `main.go`  
**Expected:** App displays centered, bordered content

Enhance your style to add more visual appeal:
- Add `Bold(true)`
- Add `Border(lipgloss.RoundedBorder())`
- Add `Padding(1, 2)`

---

### Exercise 20
**File:** `main.go`  
**Expected:** Content is centered in the terminal

Use lipgloss's `Place` function to center your content:
```go
lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
```

---

## Part 6: Setting Up the SSH Server with Wish (24:00 - 32:00)

### Exercise 21
**File:** Terminal  
**Expected:** Packages downloaded

Install the Wish SSH server library and its dependencies:
```
go get github.com/charmbracelet/wish
go get github.com/charmbracelet/ssh
go get github.com/charmbracelet/log
```

---

### Exercise 22
**File:** `main.go`  
**Expected:** Compiles

Add the new imports to your file:
- `github.com/charmbracelet/wish`
- `github.com/charmbracelet/wish/bubbletea`
- `github.com/charmbracelet/wish/activeterm`
- `github.com/charmbracelet/wish/logging`
- `github.com/charmbracelet/ssh`
- `github.com/charmbracelet/log`

Also add standard library imports: `net`, `os`, `os/signal`, `syscall`, `context`, `time`, `errors`

---

### Exercise 23
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

### Exercise 24
**File:** `main.go`  
**Expected:** Compiles

Create a `teaHandler` function. This is what Wish calls to create a Bubble Tea app for each SSH connection:
```go
func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption)
```

Inside, create and return a new model instance. Return `nil` for the options.

---

### Exercise 25
**File:** `main.go`  
**Expected:** Compiles

Enhance your `teaHandler` to get terminal info from the SSH session:
```go
pty, _, _ := sess.Pty()
```

You can access `pty.Term` for the terminal type. Consider adding a `term` field to your model to store this.

---

## Part 7: Creating the Wish Server (32:00 - 36:00)

### Exercise 26
**File:** `main.go`  
**Expected:** Compiles (but main still has old code)

Replace your `main` function contents. Start by creating a new Wish server:
```go
s, err := wish.NewServer(
    wish.WithAddress(net.JoinHostPort(host, port)),
    wish.WithHostKeyPath(".ssh/id_ed25519"),
)
if err != nil {
    log.Error("Could not create server", "error", err)
}
```

---

### Exercise 27
**File:** `main.go`  
**Expected:** Compiles

Add the middleware stack to your `wish.NewServer` call. Add `wish.WithMiddleware()` with three middlewares:
1. `bubbletea.Middleware(teaHandler)` - serves your Bubble Tea app
2. `activeterm.Middleware()` - ensures a PTY is connected
3. `logging.Middleware()` - logs connections

Remember: middlewares execute in reverse order (last added = first executed).

---

### Exercise 28
**File:** `main.go`  
**Expected:** Compiles

Set up signal handling for graceful shutdown:
```go
done := make(chan os.Signal, 1)
signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
```

---

### Exercise 29
**File:** `main.go`  
**Expected:** Compiles

Add a log message indicating the server is starting:
```go
log.Info("Starting SSH server", "host", host, "port", port)
```

---

### Exercise 30
**File:** `main.go`  
**Expected:** Server starts listening

Start the server in a goroutine:
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

After the goroutine, block on the done channel and implement graceful shutdown:
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

Run your application:
```bash
go run main.go
```

You should see: `INFO Starting SSH server host=localhost port=23234`

---

### Exercise 33
**File:** Terminal (new window)  
**Expected:** TUI appears over SSH

In a new terminal, SSH into your app:
```bash
ssh localhost -p 23234
```

Your Bubble Tea app should appear!

---

### Exercise 34
**File:** SSH session  
**Expected:** Dimensions update

Resize your terminal window while connected. The displayed dimensions should update in real-time.

---

### Exercise 35
**File:** SSH session  
**Expected:** Connection closes

Press `q` to quit. The SSH connection should close gracefully.

---

### Exercise 36
**File:** Server terminal  
**Expected:** Server stops

In the server terminal, press Ctrl+C. The server should log "Stopping SSH server" and exit cleanly.

---

## Bonus Exercises

### Exercise 37
**File:** `main.go`  
**Expected:** Styled output works over SSH

The styles might look different or broken over SSH because each session needs its own renderer. Update your model to store a `*lipgloss.Renderer` and use `bubbletea.MakeRenderer(sess)` in your `teaHandler`.

---

### Exercise 38
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
2. Handles window resize events
3. Handles keyboard input
4. Uses Lip Gloss for styling
5. Runs as an SSH server using Wish
6. Supports multiple simultaneous connections
7. Shuts down gracefully

Each SSH connection gets its own independent Bubble Tea program!

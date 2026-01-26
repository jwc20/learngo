package main

import (
	_ "context"
	_ "errors"
	"fmt"
	_ "fmt"
	_ "net"
	_ "os"
	_ "os/signal"
	_ "syscall"
	_ "time"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/charmbracelet/lipgloss"
	_ "github.com/charmbracelet/log"
	_ "github.com/charmbracelet/ssh"
	_ "github.com/charmbracelet/wish"
	_ "github.com/charmbracelet/wish/activeterm"
	_ "github.com/charmbracelet/wish/bubbletea"
	_ "github.com/charmbracelet/wish/logging"
)

type model struct {
	term   string
	width  int
	height int
}

func main() {
	p := tea.NewProgram(model, opts ...tea.ProgramOption)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Your term is %s\r\n", m.term)
	s += fmt.Sprintf("Your window size is %d x %d\r\n", m.width, m.height)
	return s
}

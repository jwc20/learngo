package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()

	// New Program with inital model and program options
	p := tea.NewProgram(m)

	// Run
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Model: app state
type Model struct {
	title string
}

// NewModel: inital model
func NewModel() Model {
	return Model{
		title: "Hello world",
	}
}

// Init: kick off the event loop
func (m Model) Init() tea.Cmd {
	return nil
}

// Update: handle Msgs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View: render the string based on the state of our model
func (m Model) View() string {
	return m.title
}

// Cmd

// Msg

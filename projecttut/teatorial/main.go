package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()

	// New Program with inital model and program options
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Model: app state
type Model struct {
	title     string
	textinput textinput.Model
	terms     Terms
	err       error
}

// NewModel: inital model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter search term"
	ti.Focus()

	return Model{
		title:     "Hello world",
		textinput: ti,
	}
}

// Init: kick off the event loop
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update: handle Msgs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			v := m.textinput.Value()
			return m, handleQuerySearch(v)
		case "ctrl+c", "q":
			return m, tea.Quit // Exit the program
		}

	case TermsResponseMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		m.terms = msg.terms
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

type Terms struct {
	List []struct {
		Definition  string    `json:"definition"`
		Permalink   string    `json:"permalink"`
		ThumbsUp    int       `json:"thumbs_up"`
		Author      string    `json:"author"`
		Word        string    `json:"word"`
		Defid       int       `json:"defid"`
		CurrentVote string    `json:"current_vote"`
		WrittenOn   time.Time `json:"written_on"`
		Example     string    `json:"example"`
		ThumbsDown  int       `json:"thumbs_down"`
	} `json:"list"`
}

// View: render the string based on the state of our model
func (m Model) View() string {
	s := m.textinput.View() + "\n\n"

	if len(m.terms.List) > 0 {
		s += m.terms.List[0].Definition + "\n\n"
	}

	return s
}

// Cmd
func handleQuerySearch(q string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.urbandictionary.com/v0/define?term=%s", url.QueryEscape(q))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return TermsResponseMsg{
				err: err,
			}
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return TermsResponseMsg{
				err: err,
			}
		}

		defer res.Body.Close()

		// decode the json
		var terms Terms
		err = json.NewDecoder(res.Body).Decode(&terms)
		if err != nil {
			return TermsResponseMsg{
				err: err,
			}
		}

		return TermsResponseMsg{
			terms: terms,
		}
	}
}

// Msg

type TermsResponseMsg struct {
	terms Terms
	err   error
}

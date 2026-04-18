package main

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentView string
}

func initialModel() model {
	return model{currentView: "Stats"}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case switchViewMsg:
		m.currentView = string(msg)
	}
	return m, nil
}

func (m model) View() string {
	s := "--- PI-DASHBOARD (1080p) ---\n\n"
	s += fmt.Sprintf("Current View: %s\n\n", m.currentView)
	s += "Press 'q' to exit\n"
	return s
}

type switchViewMsg string

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	// Start a simple HTTP server to listen for Stream Deck events
	go func() {
		http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
			view := r.URL.Query().Get("view")
			if view != "" {
				p.Send(switchViewMsg(view))
				fmt.Fprintf(w, "Switched to %s", view)
			}
		})
		http.ListenAndServe(":8080", nil)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

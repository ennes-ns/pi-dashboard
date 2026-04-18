package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF")).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(1).
			Width(80).
			Align(lipgloss.Center)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1).
			Margin(1).
			Width(38)

	styleHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
)

type model struct {
	currentView string
	status      string
}

func initialModel() model {
	return model{
		currentView: "HOME",
		status:      "SYSTEM ONLINE",
	}
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
		m.currentView = strings.ToUpper(string(msg))
	}
	return m, nil
}

func (m model) View() string {
	asciiArt := `
  _____ _____       _____                _____ _    _ ____   ____          _____  _____  
 |  __ \_   _|     |  __ \     /\       / ____| |  | |  _ \ / __ \   /\    |  __ \|  __ \ 
 | |__) || | ______ | |  | |   /  \     | (___ | |__| | |_) | |  | | /  \   | |__) | |  | |
 |  ___/ | ||______| |  | |  / /\ \     \___ \|  __  |  _ <| |  | |/ /\ \  |  _  /| |  | |
 | |    _| |_      | |__| | / ____ \    ____) | |  | | |_) | |__| / ____ \ | | \ \| |__| |
 |_|   |_____|     |_____(_/_/    \_\  |_____/|_|  |_|____/ \____/_/    \_\_|  \_\_____/ 
`
	title := styleTitle.Render(asciiArt + "\n\n1080p SENTINEL INTERFACE v1.0")

	// Create sub-windows
	leftBox := styleBox.Render(fmt.Sprintf("VIEW: %s\n\n- CPU: 12%%\n- RAM: 1.2GB / 8GB\n- TEMP: 45°C", styleHighlight.Render(m.currentView)))
	rightBox := styleBox.Render(fmt.Sprintf("NETWORK: ONLINE\n\n- IP: 10.0.0.5\n- TAILSCALE: CONNECTED\n- TRAEFIK: ACTIVE"))

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	footer := "\n\n" + styleHighlight.Render(" [STREAM DECK CONTROLLED] ") + " | Press 'q' to exit"

	// Center everything
	return lipgloss.Place(1920, 1080, lipgloss.Center, lipgloss.Center, 
		lipgloss.JoinVertical(lipgloss.Center, title, content, footer))
}

type switchViewMsg string

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	go func() {
		http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
			view := r.URL.Query().Get("view")
			if view != "" {
				p.Send(switchViewMsg(view))
				fmt.Fprintf(w, "OK: Switched to %s", view)
			}
		})
		http.ListenAndServe(":8080", nil)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

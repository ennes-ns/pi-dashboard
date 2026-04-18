package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// --- Styles ---
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginRight(2).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#5A56E0")).
			Bold(true)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(1, 1).
			MarginRight(1).
			Width(45)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			MarginTop(1)
)

// --- Model ---
type model struct {
	currentView string
	cpuUsage    float64
	memUsage    float64
	uptime      uint64
	loadAvg     *load.AvgStat
	err         error
}

type statsMsg struct {
	cpuUsage float64
	memUsage float64
	uptime   uint64
	loadAvg  *load.AvgStat
}

type switchViewMsg string

func (m model) Init() tea.Cmd {
	return tea.Batch(m.tickStats(), tea.EnterAltScreen)
}

func (m model) tickStats() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		c, _ := cpu.Percent(0, false)
		vm, _ := mem.VirtualMemory()
		h, _ := host.Info()
		l, _ := load.Avg()
		return statsMsg{
			cpuUsage: c[0],
			memUsage: vm.UsedPercent,
			uptime:   h.Uptime,
			loadAvg:  l,
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case statsMsg:
		m.cpuUsage = msg.cpuUsage
		m.memUsage = msg.memUsage
		m.uptime = msg.uptime
		m.loadAvg = msg.loadAvg
		return m, m.tickStats()
	case switchViewMsg:
		m.currentView = strings.ToUpper(string(msg))
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	// Header
	header := titleStyle.Render(" PI-DASHBOARD v2.0 - SENTINEL NODE ")
	
	// Stats Column
	uptimeHours := m.uptime / 3600
	statsContent := fmt.Sprintf(
		"SYSTEM STATS\n\n"+
			"CPU:   %s %.1f%%\n"+
			"MEM:   %s %.1f%%\n"+
			"LOAD:  %.2f %.2f %.2f\n"+
			"UPTIME: %d hours",
		m.getBar(m.cpuUsage), m.cpuUsage,
		m.getBar(m.memUsage), m.memUsage,
		m.loadAvg.Load1, m.loadAvg.Load5, m.loadAvg.Load15,
		uptimeHours,
	)
	statsCol := columnStyle.Render(statsContent)

	// View Column
	viewContent := fmt.Sprintf(
		"ACTIVE VIEW: %s\n\n"+
			"INTERFACE: HDMI-A-1\n"+
			"RESOLUTION: 1920x1080\n"+
			"CONTROL: STREAM DECK\n"+
			"STATUS: %s",
		highlight.Render(m.currentView),
		special.Render("ONLINE"),
	)
	viewCol := columnStyle.Render(viewContent)

	// Assemble UI
	body := lipgloss.JoinHorizontal(lipgloss.Top, statsCol, viewCol)
	
	asciiArt := `
    _   _   _   _   _   _   _   _   _   _   _   _   _   _   _  
   / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ / \ 
  ( P | I | - | D | A | S | H | B | O | A | R | D | - | v | 2 )
   \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ 
`
	
	ui := lipgloss.JoinVertical(lipgloss.Center, header, asciiArt, body, statusStyle.Render("Press 'q' to exit | API Listening on localhost:8080"))

	return lipgloss.Place(1920, 1080, lipgloss.Center, lipgloss.Center, ui)
}

func (m model) getBar(percent float64) string {
	bars := int(percent / 10)
	return "[" + strings.Repeat("=", bars) + strings.Repeat(" ", 10-bars) + "]"
}

func main() {
	m := model{currentView: "HOME", loadAvg: &load.AvgStat{}}
	p := tea.NewProgram(m, tea.WithAltScreen())

	// HTTP API
	go func() {
		http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
			view := r.URL.Query().Get("view")
			if view != "" {
				p.Send(switchViewMsg(view))
				fmt.Fprintf(w, "OK")
			}
		})
		http.ListenAndServe("127.0.0.1:8080", nil)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

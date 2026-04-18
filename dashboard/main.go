package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00")).
			Padding(1, 4).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#00FF00"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Padding(2).
			Width(60).
			Height(15).
			Align(lipgloss.Center)

	highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true)
)

type hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

type model struct {
	currentView string
	cpuUsage    float64
	memUsage    float64
	uptime      uint64
	loadAvg     *load.AvgStat
	hub         *hub
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
		return statsMsg{cpuUsage: c[0], memUsage: vm.UsedPercent, uptime: h.Uptime, loadAvg: l}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsMsg:
		m.cpuUsage, m.memUsage, m.uptime, m.loadAvg = msg.cpuUsage, msg.memUsage, msg.uptime, msg.loadAvg
		return m, m.tickStats()
	case switchViewMsg:
		m.currentView = strings.ToUpper(string(msg))
	}
	return m, nil
}

func (m model) View() string {
	ascii := `
 ██╗  ██╗███████╗██╗     ██╗      ██████╗     ██╗    ██╗ ██████╗ ██████╗ ██╗     ██████╗ 
 ██║  ██║██╔════╝██║     ██║     ██╔═══██╗    ██║    ██║██╔═══██╗██╔══██╗██║     ██╔══██╗
 ███████║█████╗  ██║     ██║     ██║   ██║    ██║ █╗ ██║██║   ██║██████╔╝██║     ██║  ██║
 ██╔══██║██╔══╝  ██║     ██║     ██║   ██║    ██║███╗██║██║   ██║██╔══██╗██║     ██║  ██║
 ██║  ██║███████╗███████╗███████╗╚██████╔╝    ╚███╔███╔╝╚██████╔╝██║  ██║███████╗██████╔╝
 ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝ ╚═════╝      ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═════╝ 
`
	header := titleStyle.Render(ascii)
	
	statsContent := fmt.Sprintf(
		"\n%s\n\nCPU LOAD: %.1f%%\nMEMORY:   %.1f%%\nUPTIME:   %d HRS\nLOAD:     %.2f",
		highlight.Render("SYSTEM MONITOR"),
		m.cpuUsage, m.memUsage, m.uptime/3600, m.loadAvg.Load1,
	)
	
	viewContent := fmt.Sprintf(
		"\n%s\n\nACTIVE VIEW: %s\nSTATUS:      %s\nCONTROL:     STREAMDECK",
		highlight.Render("INTERFACE STATE"),
		m.currentView,
		"SECURE",
	)

	leftBox := boxStyle.Render(statsContent)
	rightBox := boxStyle.Render(viewContent)
	
	mainContent := lipgloss.JoinHorizontal(lipgloss.Center, leftBox, rightBox)
	
	ui := lipgloss.JoinVertical(lipgloss.Center, header, "\n", mainContent)
	
	return lipgloss.Place(1920, 1080, lipgloss.Center, lipgloss.Center, ui)
}

func main() {
	// Simple Hub for WebSocket
	h := &hub{clients: make(map[*websocket.Conn]bool), broadcast: make(chan []byte)}
	
	m := model{currentView: "HOME", hub: h, loadAvg: &load.AvgStat{}}
	p := tea.NewProgram(m, tea.WithAltScreen())

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		conn, _ := upgrader.Upgrade(w, r, nil)
		h.mu.Lock()
		h.clients[conn] = true
		h.mu.Unlock()
	})

	http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
		view := r.URL.Query().Get("view")
		if view != "" {
			p.Send(switchViewMsg(view))
			fmt.Fprintf(w, "OK")
		}
	})

	go http.ListenAndServe("0.0.0.0:8080", nil)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

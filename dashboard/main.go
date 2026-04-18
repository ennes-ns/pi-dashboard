package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// --- WebSocket Hub ---
type hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

func (h *hub) run() {
	for {
		select {
		case msg := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				client.WriteMessage(websocket.TextMessage, msg)
			}
			h.mu.Unlock()
		}
	}
}

// --- Dashboard Model ---
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
	// Full Screen Styles
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Padding(1, 2).BorderStyle(lipgloss.DoubleBorder())
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 4).Width(80).Height(20)

	ascii := `
  _    _ ______ _      _       ____   __          ______  _____  _      _____  
 | |  | |  ____| |    | |     / __ \  \ \        / / __ \|  __ \| |    |  __ \ 
 | |__| | |__  | |    | |    | |  | |  \ \  /\  / / |  | | |__) | |    | |  | |
 |  __  |  __| | |    | |    | |  | |   \ \/  \/ /| |  | |  _  /| |    | |  | |
 | |  | | |____| |____| |____| |__| |    \  /\  / | |__| | | \ \| |____| |__| |
 |_|  |_|______|______|______|____/      \/  \/   \____/|_|  \_\______|_____/ 
`
	header := titleStyle.Render(ascii)
	
	content := fmt.Sprintf(
		"\n%s\n\nCPU:  %.1f%%\nMEM:  %.1f%%\nLOAD: %.2f %.2f %.2f\nUPTIME: %d HRS\nVIEW: %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("SYSTEM SENTINEL v2.0"),
		m.cpuUsage, m.memUsage, m.loadAvg.Load1, m.loadAvg.Load5, m.loadAvg.Load15,
		m.uptime/3600,
		m.currentView,
	)

	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Top, header, boxStyle.Render(content)))
}

var docStyle = lipgloss.NewStyle().Margin(2, 4)

type switchViewMsg string

func main() {
	h := &hub{clients: make(map[*websocket.Conn]bool), broadcast: make(chan []byte)}
	go h.run()

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

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

// --- WebSocket Support ---
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func newHub() *hub {
	return &hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(h.clients, client)
				}
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

		// Real-time alerting via WebSocket
		if c[0] > 90 {
			m.hub.broadcast <- []byte("ALERT: HIGH CPU USAGE")
		}

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
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00D7FF")).Border(lipgloss.RoundedBorder())
	header := titleStyle.Render(" PI-DASHBOARD v2.1 (WEBSOCKET ENABLED) ")

	body := fmt.Sprintf("VIEW: %s\nCPU: %.1f%%\nMEM: %.1f%%", m.currentView, m.cpuUsage, m.memUsage)
	
	return lipgloss.Place(1920, 1080, lipgloss.Center, lipgloss.Center, 
		lipgloss.JoinVertical(lipgloss.Center, header, body))
}

func main() {
	h := newHub()
	go h.run()

	m := model{currentView: "HOME", hub: h, loadAvg: &load.AvgStat{}}
	p := tea.NewProgram(m, tea.WithAltScreen())

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		h.register <- conn
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

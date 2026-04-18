package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF")).
			Padding(0, 1).
			Underline(true)

	containerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	stoppedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

func getDockerStats() string {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}|{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return "Error fetching docker stats"
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var b strings.Builder

	b.WriteString(headerStyle.Render("DOCKER CONTAINERS") + "\n\n")

	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		status := parts[1]

		statusIndicator := runningStyle.Render("●")
		if !strings.Contains(strings.ToLower(status), "up") {
			statusIndicator = stoppedStyle.Render("○")
		}

		b.WriteString(fmt.Sprintf("%s %-20s %s\n", statusIndicator, containerStyle.Render(name), lipgloss.NewStyle().Faint(true).Render(status)))
	}

	return b.String()
}

func main() {
	// Simple loop to print to terminal
	for {
		fmt.Print("\033[H\033[J") // Clear screen
		fmt.Println(getDockerStats())
		time.Sleep(2 * time.Second)
	}
}

// splash screen animation copied from https://www.youtube.com/watch?v=YCJgQI71jEE
package ui

import (
	"math"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var sentences = []string{
	"/dream Lorem ipsum dolor sit amet, consectetur adipiscing elit, Ut enim ad minim veniam, quis nostrud exercitation ullamco",
	"/dream Ut enim ad minim veniam, quis nostrud exercitation ullamco Duis aute irure dolor in reprehenderit in voluptate velit",
	"/dream Duis aute irure dolor in reprehenderit in voluptate velit Excepteur sint occaecat cupidatat non proident, sunt in culpa",
	"/dream Excepteur sint occaecat cupidatat non proident, sunt in culpa At vero eos et accusamus et iusto odio dignissimos ducimus",
	"/dream At vero eos et accusamus et iusto odio dignissimos ducimus Lorem ipsum dolor sit amet, consectetur adipiscing elit,",
}

var dialogStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 1, 1, 0).
	Bold(true).
	Align(lipgloss.Center)

type SplashModel struct {
	width, height int
	rate          float64
	startTime     time.Time
}

func (m SplashModel) Init() tea.Cmd {
	return tick
}

func (m SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m, func() tea.Msg { return BackToLobbyMsg{} }
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		if m.startTime.IsZero() {
			m.startTime = time.Now()
		}
		return m, tick
	}
	return m, nil
}

func (m SplashModel) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.renderFrame())
	return v
}

func (m SplashModel) renderFrame() string {
	bg := m.renderVortex()

	dialog := dialogStyle.Render("                                                                            \n             d888888o.  `8.`888b           ,8' 8888888 8888888888           \n           .`8888:' `88. `8.`888b         ,8'        8 8888                 \n           8.`8888.   Y8  `8.`888b       ,8'         8 8888                 \n           `8.`8888.       `8.`888b     ,8'          8 8888                 \n            `8.`8888.       `8.`888b   ,8'           8 8888                 \n             `8.`8888.       `8.`888b ,8'            8 8888                 \n              `8.`8888.       `8.`888b8'             8 8888                 \n          8b   `8.`8888.       `8.`888'              8 8888                 \n          `8b.  ;8.`8888        `8.`8'               8 8888                 \n           `Y8888P ,88P'         `8.`                8 8888                 ")

	return m.overlayCentered(bg, dialog)
}

func (m SplashModel) overlayCentered(bg string, overlay string) string {
	bgLines := strings.Split(bg, "\n")

	bgW := m.width
	bgH := len(bgLines)

	ow := lipgloss.Width(overlay)
	oh := lipgloss.Height(overlay)

	startX := (bgW - ow) / 2
	startY := (bgH - oh) / 2

	overlayLines := strings.Split(overlay, "\n")

	for y := 0; y < oh; y++ {
		if startY+y < 0 || startY+y >= len(bgLines) {
			continue
		}

		line := []rune(bgLines[startY+y])
		for x := 0; x < ow; x++ {
			if startX+x < 0 || startX+x >= len(line) {
				continue
			}

			r := rune(overlayLines[y][x])
			line[startX+x] = r
		}
		bgLines[startY+y] = string(line)
	}

	return strings.Join(bgLines, "\n")
}

func (m SplashModel) renderVortex() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	// Calculate animation time
	// Adjust rate (e.g., 1.0 or 2.0) to change spin speed
	if m.rate == 0 {
		m.rate = 1.0
	}
	t := time.Since(m.startTime).Seconds() * m.rate

	// Initialize the character grid
	grid := make([][]rune, m.height)
	for i := range grid {
		grid[i] = make([]rune, m.width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	cx, cy := 0.5, 0.5

	// 1. Map source characters to transformed coordinates
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			char := m.getCharAt(x, y)

			// Transform logic: (x, y) -> (nx, ny)
			nx, ny := m.transform(
				float64(x)/float64(m.width),
				float64(y)/float64(m.height),
				cx, cy, t,
			)

			tx := int(nx * float64(m.width))
			ty := int(ny * float64(m.height))

			// Bounds check and placement
			if tx >= 0 && tx < m.width && ty >= 0 && ty < m.height {
				grid[ty][tx] = char
			}
		}
	}

	// 2. Build the final string from the grid
	var output strings.Builder
	for y := 0; y < m.height; y++ {
		output.WriteString(string(grid[y]))
		if y < m.height-1 {
			output.WriteByte('\n')
		}
	}

	return output.String()
}

func (m SplashModel) transform(x, y, cx, cy, t float64) (float64, float64) {
	dx := x - cx
	dy := y - cy

	dist := math.Sqrt(dx*dx + dy*dy)
	angle := math.Atan2(dy, dx)

	// The Vortex Equation
	// The math.Pow(dist, 0.5) makes the center spin differently than the edges
	newAngle := angle - math.Pow(dist, 0.5)*t*2.0

	nx := cx + math.Cos(newAngle)*dist
	ny := cy + math.Sin(newAngle)*dist

	return nx, ny
}

func (m SplashModel) getCharAt(x, y int) rune {
	si := y % len(sentences)
	line := sentences[si]
	if len(line) == 0 {
		return ' '
	}
	ci := x % len(line)
	return rune(line[ci])
}

type tickMsg time.Time

func tick() tea.Msg {
	return tickMsg(time.Now())
}

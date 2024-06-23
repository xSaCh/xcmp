package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xSaCh/xcmp/image"
	"golang.org/x/term"
)

func main() {
	physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	physicalWidth = min(physicalWidth, 70)
	// physicalHeight = min(physicalHeight, 15)

	docStyle := lipgloss.NewStyle().
		Padding(1, 2, 1, 2).
		Width(physicalWidth - 2).   // border
		Height(physicalHeight - 2). // border
		MaxHeight(physicalHeight).MaxWidth(physicalWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD"))

	bS := float32(docStyle.GetHeight() - docStyle.GetPaddingTop() - docStyle.GetPaddingBottom())
	box := lipgloss.NewStyle().
		Width(int(bS * 2)).Height(int(bS)).
		Background(lipgloss.Color("#FF0000"))

	pb := func(s string, b int) string {
		return lipgloss.NewStyle().PaddingBottom(b).Render(s)
	}
	img, _ := image.Display("a.jpg", image.KittyOpts{DisplayHeight: uint32(bS), DisplayWidth: uint32(bS * 2)})

	prog := progress.New(progress.WithSolidFill("#874BFD"))
	prog.Width = 24
	prog.Empty = '-'
	prog.Full = '▬'
	// prog. = lipgloss.NewStyle().Background(lipgloss.Color("#FF00FF"))

	final := docStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Top,
		box.Render(img),
		lipgloss.NewStyle().PaddingLeft(4).Render(lipgloss.JoinVertical(
			lipgloss.Top,
			pb("Imagation", 0),
			pb("Charlie Puth", 1),
			pb("Delux Edition", 0),
			lipgloss.NewStyle().Height(int(bS)-5-0).Render(),
			prog.ViewAs(0.2),
		)),
	))

	p := tea.NewProgram(playerModel(final), tea.WithAltScreen())
	_ = p
	// fmt.Printf("physicalWidth: %v\n", physicalWidth)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type playerModel string

func (m playerModel) Init() tea.Cmd {
	return nil
}

func (m playerModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		}

	}

	return m, nil
}

func (m playerModel) View() string {
	// go image.DisplayA("file.jpg", image.IcatOpts{Dx: 3, Dy: 2, W: uint32(28), MaintainAspectRatio: true})
	// a, _ := image.Display("file.jpg", image.KittyOpts{DisplayHeight: 14, DisplayWidth: 28})
	return string(m)
}

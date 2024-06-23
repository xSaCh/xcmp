package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xSaCh/xcmp/image"
	"golang.org/x/term"
)

type SongInfo struct {
	Title    string
	Artist   string
	Album    string
	Duration float32 // In Seconds
	AlbumArt string
}

// New Duration from starting point
type PlaybackUpdate float32

type playerUIModel struct {
	MainBox         lipgloss.Style
	ProgressBar     progress.Model
	AlbumArt        lipgloss.Style
	ArtistInfoStyle lipgloss.Style
	MetaData        SongInfo

	currentDuration float32
	imgStr          string
	progressPadding string
}

func (m playerUIModel) Init() tea.Cmd {
	physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))

	m.updateSize(physicalWidth, physicalHeight)
	return nil
}

func (m playerUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case PlaybackUpdate:
		m.currentDuration = float32(msg)
		return m, nil

	case tea.WindowSizeMsg:
		m.updateSize(msg.Width, msg.Height)
		return m, nil

	}

	return m, nil
}

func (m playerUIModel) View() string {
	prgPercent := 1 - ((m.MetaData.Duration - m.currentDuration) / m.MetaData.Duration)
	return m.MainBox.Render(lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.AlbumArt.Render(m.imgStr),
		lipgloss.NewStyle().PaddingLeft(4).Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.ArtistInfoStyle.Render(),
				"",
				lipgloss.JoinHorizontal(lipgloss.Center, m.ProgressBar.ViewAs(float64(prgPercent)), " -", getRemainingDuratrion(m.currentDuration, m.MetaData.Duration)),
			),
		),
	))
}

func (m *playerUIModel) SetInfo(newMetaData SongInfo) {
	m.MetaData = newMetaData

}

func (m *playerUIModel) updateSize(width, height int) {
	m.MainBox = m.MainBox.Width(width - 2).Height(height - 2).MaxHeight(height).MaxWidth(width)
	contentHeight := float32(m.MainBox.GetHeight() - m.MainBox.GetPaddingTop() - m.MainBox.GetPaddingBottom())

	m.AlbumArt = m.AlbumArt.Width(int(contentHeight * 2)).Height(int(contentHeight))
	m.imgStr, _ = image.Display(m.MetaData.AlbumArt, image.KittyOpts{DisplayHeight: uint32(contentHeight), DisplayWidth: uint32(contentHeight * 2)})

	m.progressPadding = lipgloss.NewStyle().Height(int(contentHeight - float32(lipgloss.Height(m.ArtistInfoStyle.Render())) - 0)).Render()

	m.ArtistInfoStyle = m.ArtistInfoStyle.SetString(lipgloss.JoinVertical(
		lipgloss.Top,

		lipgloss.NewStyle().MaxWidth(width-2-lipgloss.Width(m.AlbumArt.Render())-8).Foreground(lipgloss.Color("#FFFFFF")).PaddingBottom(1).Render(m.MetaData.Title),
		lipgloss.NewStyle().MaxWidth(width-2-lipgloss.Width(m.AlbumArt.Render())-8).Foreground(lipgloss.Color("#777777")).Render(m.MetaData.Artist),
		lipgloss.NewStyle().MaxWidth(width-2-lipgloss.Width(m.AlbumArt.Render())-8).Foreground(lipgloss.Color("#809bd1")).Render(m.MetaData.Album),
	))

	fmt.Print(image.ClearImage()) // Clear image when resize, kitty render image over prev Size
}

// Default Player UI
func DefaultPlayerUi(songinfo SongInfo) playerUIModel {
	mainBox := lipgloss.NewStyle().
		Padding(1, 2, 1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD"))

	info := lipgloss.NewStyle()

	prog := progress.New(progress.WithScaledGradient("#874BFD", "#4b57fd"))
	prog.Width = 24
	prog.Empty = '━'
	prog.Full = '━'
	prog.ShowPercentage = false

	model := playerUIModel{
		MainBox:         mainBox,
		ProgressBar:     prog,
		ArtistInfoStyle: info,
		AlbumArt:        lipgloss.NewStyle(),
	}
	model.SetInfo(songinfo)
	return model
}

func getRemainingDuratrion(to, from float32) string {
	rem := from - to

	min := int(rem / 60)
	sec := int(rem) % 60

	return fmt.Sprintf("%02d:%02d", min, sec)
}

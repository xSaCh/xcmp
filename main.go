package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xSaCh/xcmp/tui"
	"github.com/xSaCh/xcmp/util"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./xcmp songFilePath")
		return

	}

	info, err := util.GetSongInfo(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	albumFile, err := util.GetAlbumArtPath(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer os.Remove(albumFile)

	duration := (func() float32 { a, _ := strconv.ParseFloat(info.Format.Duration, 32); return float32(a) }())
	p := tea.NewProgram(tui.DefaultPlayerUi(tui.SongInfo{
		Title:    info.Format.Tags.Title,
		Artist:   info.Format.Tags.Artist,
		Album:    info.Format.Tags.Album,
		Duration: duration,
		AlbumArt: albumFile,
	}),
		tea.WithAltScreen())

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		remaningDuration := duration
		for range ticker.C {
			remaningDuration -= 1
			p.Send(tui.PlaybackUpdate(duration - remaningDuration))
			if remaningDuration <= 0 {
				break
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

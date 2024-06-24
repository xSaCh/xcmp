package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/leberKleber/go-mpris"
	"github.com/xSaCh/xcmp/models"
	"github.com/xSaCh/xcmp/tui"
	"github.com/xSaCh/xcmp/util"
)

func ffmpegMain() {

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
	p := tea.NewProgram(tui.DefaultPlayerUi(models.SongInfo{
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

func main() {

	players, _ := util.GetMprisPlayers()

	if len(players) == 0 {
		fmt.Println("Empty Players")
		return
	}

	pl, _ := mpris.NewPlayer(players[1])

	info := util.GetSongInfoFromMprisPlayer(&pl)

	p := tea.NewProgram(tui.DefaultPlayerUi(info),
		tea.WithAltScreen())

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		prevMeta, _ := pl.Metadata()
		prevTitle, _ := prevMeta.XESAMTitle()
		for range ticker.C {
			pos, _ := pl.Position()
			curMeta, _ := pl.Metadata()
			curTitle, _ := curMeta.XESAMTitle()

			if curTitle != prevTitle {
				prevTitle = curTitle
				prevMeta = curMeta

				newInfo := util.GetSongInfoFromMprisPlayer(&pl)

				p.Send(tui.ChangeSong(newInfo))
			}

			p.Send(tui.PlaybackUpdate(float32(pos) / float32(time.Millisecond)))
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/leberKleber/go-mpris"
	"github.com/xSaCh/xcmp/models"
)

func GetMprisPlayers() ([]string, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	var names []string

	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %v", err)
	}

	var filteredName []string
	const prefixLen = len("org.mpris.MediaPlayer2.")
	for _, n := range names {
		if strings.HasPrefix(n, "org.mpris.MediaPlayer2.") {
			// filteredName = append(filteredName, n[prefixLen:])
			filteredName = append(filteredName, n)
		}
	}

	return filteredName, nil
}

func GetSongInfoFromMprisPlayer(p *mpris.Player) models.SongInfo {

	metaData, _ := p.Metadata()

	// fmt.Printf("AD: %#v %#v", metaData, e)

	title, _ := metaData.XESAMTitle()
	artist, _ := metaData.XESAMArtist()
	album, _ := metaData.XESAMAlbum()
	vl := metaData["mpris:length"].Value()
	url, _ := metaData.MPRISArtURL()

	var duration int64
	var file = ""
	artistName := ""

	if len(artist) > 0 {
		artistName = artist[0]
	}

	if vl == nil {
		duration = -1
	} else if fmt.Sprintf("%T", vl) == "uint64" {
		duration = int64(vl.(uint64)) / int64(time.Millisecond)
	} else {
		duration = int64(vl.(int64)) / int64(time.Millisecond)
	}

	_ = url
	//Assuming url is http url
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {
		data, _ := io.ReadAll(resp.Body)

		tmpFilePath := fmt.Sprintf("/tmp/%s.jpg", strings.ReplaceAll(title, " ", "_"))
		ferr := os.WriteFile(tmpFilePath, data, 0644)

		if ferr == nil {
			file = tmpFilePath
		}
	}

	return models.SongInfo{
		Title:    title,
		Artist:   artistName,
		Album:    album,
		Duration: float32(duration),
		AlbumArt: file,
	}
}

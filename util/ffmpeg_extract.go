package util

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type FFMpegOutput struct {
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			Title  string `json:"title"`
			Artist string `json:"artist"`
			Album  string `json:"album"`
		} `json:"tags"`
	} `json:"format"`
}

func GetSongInfo(filePath string) (FFMpegOutput, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", filePath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running ffprobe: %v\n", err)
		return FFMpegOutput{}, fmt.Errorf("error running ffprobe: %v", err)
	}

	var ffprobeOutput FFMpegOutput
	err = json.Unmarshal(output, &ffprobeOutput)
	if err != nil {
		return FFMpegOutput{}, fmt.Errorf("error parsing JSON: %v", err)
	}

	return ffprobeOutput, nil
}

// NOTE: Remove file after use
func GetAlbumArtPath(filePath string) (string, error) {
	checkCmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=index", "-of", "csv=p=0", filePath)
	output, err := checkCmd.Output()
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(string(output)) == "" {
		return "", fmt.Errorf("no album art found in the MP3 file")
	}

	tempFile := "album_art.jpg"
	extractCmd := exec.Command("ffmpeg", "-i", filePath, "-an", "-vcodec", "copy", tempFile)
	err = extractCmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running ffmpeg to extract album art")
	}

	return tempFile, nil

}

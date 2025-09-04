package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// ffprobeData matches the JSON structure from ffprobe
type ffprobeData struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

// getVideoAspectRatio uses ffprobe to determine the aspect ratio of a video file.
// It returns "landscape", "portrait", or "other".
func getVideoAspectRatio(filePath string) (string, error) {
	// Prepare the ffprobe command
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	// Create a buffer to capture the command's stdout
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running ffprobe: %w", err)
	}

	// Unmarshal the JSON output
	var data ffprobeData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		return "", fmt.Errorf("error unmarshaling ffprobe output: %w", err)
	}

	// Check if we got any stream data
	if len(data.Streams) == 0 {
		return "other", nil
	}

	// Use the first video stream to determine aspect ratio
	stream := data.Streams[0]
	if stream.Width == 0 || stream.Height == 0 {
		return "other", nil
	}

	// Use floating-point division to get an accurate ratio
	ratio := float64(stream.Width) / float64(stream.Height)

	if ratio > 1 {
		return "landscape", nil // e.g., 16:9
	}
	if ratio < 1 {
		return "portrait", nil // e.g., 9:16
	}

	return "other", nil
}


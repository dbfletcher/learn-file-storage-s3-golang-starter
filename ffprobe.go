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
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running ffprobe: %w", err)
	}
	var data ffprobeData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		return "", fmt.Errorf("error unmarshaling ffprobe output: %w", err)
	}
	if len(data.Streams) == 0 {
		return "other", nil
	}
	stream := data.Streams[0]
	if stream.Width == 0 || stream.Height == 0 {
		return "other", nil
	}
	ratio := float64(stream.Width) / float64(stream.Height)
	if ratio > 1 {
		return "landscape", nil
	}
	if ratio < 1 {
		return "portrait", nil
	}
	return "other", nil
}

// processVideoForFastStart uses ffmpeg to move the moov atom to the beginning of the file.
func processVideoForFastStart(filePath string) (string, error) {
	outputFilePath := filePath + ".processing"
	cmd := exec.Command("ffmpeg",
		"-y", // Overwrite output file if it exists
		"-i", filePath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		outputFilePath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %w, details: %s", err, stderr.String())
	}
	return outputFilePath, nil
}


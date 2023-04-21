package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func main() {
	var url string
	fmt.Print("Enter YouTube video URL: ")
	fmt.Scanln(&url)

	// Create new YouTube client
	client := youtube.Client{}

	// Get video information
	video, err := client.GetVideo(url)
	if err != nil {
		fmt.Println("Error getting video info:", err)
		return
	}

	// Create output file
	title := strings.Replace(video.Title, " ", "_", -1)
	title = strings.Replace(title, "|", "", -1)
	title = strings.Replace(title, "?", "", -1)
	title = strings.Replace(title, ":", "", -1)
	output, err := os.Create(fmt.Sprintf("%s.m4a", title[:min(len(title), 50)]))
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	// Get highest quality video stream
	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		fmt.Println("Error getting video stream:", err)
		return
	}

	// Download video stream
	_, err = io.Copy(output, stream)
	if err != nil {
		fmt.Println("Error downloading video:", err)
		return
	}

	fmt.Println("Video downloaded successfully.")
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

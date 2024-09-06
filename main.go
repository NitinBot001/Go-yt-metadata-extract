package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lithdew/youtube"
)

type VideoResponse struct {
	Title        string `json:"title"`
	VideoURL     string `json:"video_url"`
	AudioURL     string `json:"audio_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

func getVideoInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	client := youtube.Client{}
	results, err := client.Search(query)
	if err != nil || len(results) == 0 {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	video := results[0]

	// Get video details
	formats := video.StreamingData.AdaptiveFormats
	var videoURL, audioURL string
	for _, format := range formats {
		if format.MimeType == "audio/mp4" {
			audioURL = format.URL
		} else if format.MimeType == "video/mp4" {
			videoURL = format.URL
		}
	}

	resp := VideoResponse{
		Title:        video.VideoDetails.Title,
		VideoURL:     videoURL,
		AudioURL:     audioURL,
		ThumbnailURL: video.VideoDetails.Thumbnail[0].URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	http.HandleFunc("/video", getVideoInfo)
	fmt.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

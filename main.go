package main

import (
	"encoding/json"
	"log"
	"net/http"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoEditRequest struct {
	VideoURL  string `json:"video_url"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func editVideo(videoURL, startTime, endTime string) (string, error) {
	// call ffmpeg to edit video
	// return edited video URL
	outputURL := "output.mp4"
	err := ffmpeg.Input(videoURL).Output(outputURL, ffmpeg.KwArgs{"ss": startTime, "to": endTime}).OverWriteOutput().Run()
	if err != nil {
		return "", err
	}

	return outputURL, nil
}

func handleVideoEdit(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req VideoEditRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	videoURL, err := editVideo(req.VideoURL, req.StartTime, req.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return upload URL
	w.Write([]byte(videoURL))
}

func main() {
	http.HandleFunc("/edit-video", handleVideoEdit)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

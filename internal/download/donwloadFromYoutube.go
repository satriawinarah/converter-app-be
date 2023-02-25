package download

import (
	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
	"io"
	"log"
	"os"
)

type VideoInfo struct {
	Id      string        `json:"id"`
	Title   string        `json:"title"`
	Author  string        `json:"author"`
	Formats []VideoFormat `json:"formats"`
}

type VideoFormat struct {
	ItagNo       int    `json:"itagNo"`
	Quality      string `json:"quality"`
	AudioQuality string `json:"audioQuality"`
}

func FromYoutube(url string, formatId int) (string, error) {
	log.Println("Video Url: ", url)

	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		panic(err)
	}

	format := video.Formats.FindByItag(formatId)

	stream, _, err := client.GetStream(video, format)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	fileName := uuid.New().String()
	file, err := os.Create("/tmp/" + fileName + ".mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	return fileName, nil
}

func GetYoutubeVideoInfo(url string) VideoInfo {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		panic(err)
	}

	// set the lowest format
	var format youtube.Format
	for _, f := range video.Formats {
		if f.AudioQuality != "" && f.Quality != "" {
			if &format == nil || f.ItagNo < format.ItagNo {
				format = f
			}
		}
	}

	if &format == nil {
		errMsg := "could not find a video format with both audio and video"
		panic(errMsg)
	}

	var videoFormats []VideoFormat
	for _, videoFormat := range video.Formats {
		videoFormats = append(videoFormats, VideoFormat{
			ItagNo:       videoFormat.ItagNo,
			Quality:      videoFormat.Quality,
			AudioQuality: videoFormat.AudioQuality,
		})
	}

	return VideoInfo{
		Id:      video.ID,
		Title:   video.Title,
		Author:  video.Author,
		Formats: videoFormats,
	}
}

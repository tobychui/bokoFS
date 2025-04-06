package renderer

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
	This package is used to extract meta data from files like mp3 and mp4
	Also support image caching

*/

type RenderHandler struct {
	renderingFiles  sync.Map
	renderingFolder sync.Map
}

// Create a new RenderHandler
func NewRenderHandler() *RenderHandler {
	return &RenderHandler{
		renderingFiles:  sync.Map{},
		renderingFolder: sync.Map{},
	}
}

// RenderThumbnail generates a thumbnail for the given input file and saves it to the output folder
func (rh *RenderHandler) RenderThumbnail(inputFile string, outputFolder string) error {
	if rh.fileIsBusy(inputFile) {
		return errors.New("file is rendering")
	}

	// Check if the cache file exists and is newer than the input file
	cacheFile := filepath.Join(outputFolder, filepath.Base(inputFile)+".jpg")
	cacheInfo, err := os.Stat(cacheFile)
	if err == nil {
		inputInfo, err := os.Stat(inputFile)
		if err != nil {
			// File not found, return error
			return err
		}
		if cacheInfo.ModTime().After(inputInfo.ModTime()) {
			// Cache file is newer, return the base64 encoded image
			return nil
		}
	}

	//Cache image not exists. Set this file to busy
	rh.renderingFiles.Store(inputFile, "busy")
	inputFileExt := strings.ToLower(filepath.Ext(inputFile))
	//That object not exists. Generate cache image
	//Audio formats that might contains id4 thumbnail
	id4Formats := []string{".mp3", ".ogg", ".flac"}
	if stringInSlice(inputFileExt, id4Formats) {
		err := generateThumbnailForAudio(inputFile, outputFolder)
		rh.renderingFiles.Delete(inputFileExt)
		return err
	}

	//Generate resized image for images
	imageFormats := []string{".png", ".jpeg", ".jpg"}
	if stringInSlice(inputFileExt, imageFormats) {
		err := generateThumbnailForImage(inputFile, outputFolder)
		rh.renderingFiles.Delete(inputFileExt)
		return err
	}

	//Video formats, extract from the 5 sec mark
	vidFormats := []string{".mkv", ".mp4", ".webm", ".ogv", ".avi", ".rmvb"}
	if stringInSlice(inputFileExt, vidFormats) {
		err := generateThumbnailForVideo(inputFile, outputFolder)
		rh.renderingFiles.Delete(inputFileExt)
		return err
	}

	//3D Model Formats
	modelFormats := []string{".stl", ".obj"}
	if stringInSlice(inputFileExt, modelFormats) {
		err := generateThumbnailForModel(inputFile, outputFolder)
		rh.renderingFiles.Delete(inputFileExt)
		return err
	}

	//Photoshop file
	if inputFileExt == ".psd" {
		err := generateThumbnailForPSD(inputFile, outputFolder)
		rh.renderingFiles.Delete(inputFileExt)
		return err
	}

	//Other filters
	rh.renderingFiles.Delete(inputFileExt)
	return errors.New("No supported format")
}

func stringInSlice(path string, slice []string) bool {
	for _, item := range slice {
		if item == path {
			return true
		}
	}
	return false
}

func (rh *RenderHandler) fileIsBusy(path string) bool {
	if rh == nil {
		log.Println("RenderHandler is null!")
		return true
	}
	_, ok := rh.renderingFiles.Load(path)
	if !ok {
		//File path is not being process by another process
		return false
	} else {
		return true
	}
}

// Get the image as base64 string
func getImageAsBase64(inputFile string) (string, error) {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return string(encoded), nil
}

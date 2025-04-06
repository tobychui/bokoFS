package renderer

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/oliamb/cutter"
)

func generateThumbnailForVideo(inputFile string, outputFolder string) error {
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		// The user removed this file before the thumbnail is finished
		return errors.New("source not exists")
	}

	outputFile := filepath.Join(outputFolder, filepath.Base(inputFile)+".jpg")

	absInputFile, err := filepath.Abs(inputFile)
	if err != nil {
		return errors.New("failed to get absolute path of input file")
	}

	absOutputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return errors.New("failed to get absolute path of output file")
	}

	//Get the first thumbnail using ffmpeg
	cmd := exec.Command("ffmpeg", "-i", absInputFile, "-ss", "00:00:05.000", "-vframes", "1", "-vf", "scale=-1:480", absOutputFile)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	//Resize and crop the output image
	if _, err := os.Stat(outputFile); err == nil {
		imageBytes, err := os.ReadFile(outputFile)
		if err != nil {
			return err
		}
		os.Remove(outputFile)
		img, _, err := image.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return err
		} else {
			//Crop out the center
			croppedImg, err := cutter.Crop(img, cutter.Config{
				Width:  480,
				Height: 480,
				Mode:   cutter.Centered,
			})

			if err == nil {
				//Write it back to the original file
				out, _ := os.Create(outputFile)
				jpeg.Encode(out, croppedImg, nil)
				out.Close()

			} else {
				//log.Println(err)
			}
		}

	}
	return nil

}

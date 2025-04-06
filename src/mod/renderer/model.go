package renderer

import (
	"errors"
	"image/jpeg"
	"os"
	"path/filepath"
)

// Generate thumbnail for 3D model. Output file will have the same name as the input file with .jpg extension
func generateThumbnailForModel(inputFile string, outputFolder string) error {
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return errors.New("input file does not exist")
	}

	//Generate a render of the 3d model
	outputFile := filepath.Join(outputFolder, filepath.Base(inputFile)+".jpg")
	r := New3DRenderer(RenderOption{
		Color:           "#f2f542",
		BackgroundColor: "#ffffff",
		Width:           480,
		Height:          480,
	})

	img, err := r.RenderModel(inputFile)
	if err != nil {
		return err
	}

	opt := jpeg.Options{
		Quality: 90,
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	err = jpeg.Encode(f, img, &opt)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

package renderer

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// Generate thumbnail for audio. Output file will have the same name as the input file with .jpg extension
func generateThumbnailForAudio(inputFile string, outputFolder string) error {
	//This extension is supported by id4. Call to library
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	m, err := tag.ReadFrom(f)
	if err != nil {
		return err
	}

	if m.Picture() != nil {
		//Convert the picture bytecode to image object
		img, _, err := image.Decode(bytes.NewReader(m.Picture().Data))
		if err != nil {
			//Fail to convert this image. Continue next one
			return err
		}

		//Create an empty file
		outputFilename := filepath.Join(outputFolder, filepath.Base(inputFile))
		out, err := os.Create(outputFilename + ".jpg")
		if err != nil {
			return err
		}
		defer out.Close()

		b := img.Bounds()
		imgWidth := b.Max.X
		imgHeight := b.Max.Y

		//Resize the albumn image
		var m image.Image
		if imgWidth > imgHeight {
			m = resize.Resize(0, 480, img, resize.Lanczos3)
		} else {
			m = resize.Resize(480, 0, img, resize.Lanczos3)
		}

		//Crop out the center
		croppedImg, _ := cutter.Crop(m, cutter.Config{
			Width:  480,
			Height: 480,
			Mode:   cutter.Centered,
		})

		//Write the cache image to disk
		jpeg.Encode(out, croppedImg, nil)

	}
	return errors.New("no image found")
}

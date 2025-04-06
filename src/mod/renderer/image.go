package renderer

import (
	"errors"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// Generate thumbnail for image. Output file will have the same name as the input file with .jpg extension
func generateThumbnailForImage(inputFile string, outputFolder string) error {

	var img image.Image
	var err error

	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return errors.New("failed to get file info: " + err.Error())
	}
	if fileInfo.Size() > (25 << 20) {
		// Maximum image size to be converted is 25MB
		return errors.New("image file too large")
	}

	srcImage, err := os.OpenFile(inputFile, os.O_RDONLY, 0775)
	if err != nil {
		return err
	}
	defer srcImage.Close()
	img, _, err = image.Decode(srcImage)
	if err != nil {
		return err
	}

	//Resize to desiered width
	//Check boundary to decide resize mode
	b := img.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	var m image.Image
	if imgWidth > imgHeight {
		m = resize.Resize(0, 480, img, resize.Lanczos3)
	} else {
		m = resize.Resize(480, 0, img, resize.Lanczos3)
	}

	//Crop out the center
	croppedImg, err := cutter.Crop(m, cutter.Config{
		Width:  480,
		Height: 480,
		Mode:   cutter.Centered,
	})

	//Create the thumbnail
	outputFilename := filepath.Join(outputFolder, filepath.Base(inputFile))
	out, err := os.Create(outputFilename + ".jpg")
	if err != nil {
		return err
	}

	// write new image to file
	jpeg.Encode(out, croppedImg, nil)
	out.Close()

	return nil
}

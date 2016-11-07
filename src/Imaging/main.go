package Imaging

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math"
	"os"
	"sync"

	"github.com/nfnt/resize"
)

// ResizeImage resize image
func ResizeImage(imagePath string, dest string, size int, wg *sync.WaitGroup) {
	err := resizeImageWithoutErrorHandling(imagePath, dest, size)
	if err != nil {
		fmt.Println("E " + imagePath)
	}

	wg.Done()
}

func resizeImageWithoutErrorHandling(imagePath string, dest string, size int) error {
	imgData, err := getImageFromFilePath(imagePath)
	if err != nil {
		return err
	}

	if size < 1 || size > 40 {
		size = 16
	}
	// Find width from resize
	resizeImageWidth, isSmaller := getWidthForResize(imgData, size)

	if !isSmaller {
		resizedImg := resize.Resize(resizeImageWidth, 0, imgData, resize.Lanczos2)

		err := createJpg(resizedImg, dest)
		if err != nil {
			return err
		}
		fmt.Println("+ " + imagePath)
	} else {
		err := copyFile(imagePath, dest)
		if err != nil {
			return err
		}
		fmt.Println("- " + imagePath)
	}

	return nil
}

func getWidthForResize(img image.Image, megaPixels int) (width uint, isSmaller bool) {
	imageBounds := img.Bounds()
	imageX := imageBounds.Dx()
	imageY := imageBounds.Dy()

	imagePixels := imageX * imageY
	resizePixels := megaPixels * 1000000

	resizeRatio := math.Sqrt(float64(imagePixels) / float64(resizePixels))
	resizeWidth := float64(imageX) / resizeRatio
	resizeWidthInt := uint(resizeWidth)

	smaller := uint(imageX) < resizeWidthInt

	return resizeWidthInt, smaller
}

func getImageFromFilePath(path string) (image.Image, error) {
	// open the image file
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// read the image data and build the image.Image object
	imgData, err := jpeg.Decode(reader)
	if err != nil {
		return nil, err
	}
	// close the image file
	reader.Close()

	return imgData, nil
}

func copyFile(source string, destination string) error {
	sourceReader, err := os.Open(source)
	if err != nil {
		return err
	}

	destinationReader, err := os.Create(destination)
	if err != nil {
		errOnClose := sourceReader.Close()
		if errOnClose != nil {
			return fmt.Errorf("error creating destination file ("+destination+") - (%v); error closing source file ("+source+") - (%v); ", err, errOnClose)
		}
		return err
	}

	io.Copy(destinationReader, sourceReader)

	srcCloseErr := sourceReader.Close()
	destCloseErr := destinationReader.Close()

	if srcCloseErr == nil && destCloseErr == nil {
		return nil
	}
	if srcCloseErr != nil && destCloseErr != nil {
		return fmt.Errorf("error closing destination file ("+destination+") - (%v); error closing source file ("+source+") - (%v); ", destCloseErr, srcCloseErr)
	} else if destCloseErr != nil {
		return destCloseErr
	} else if destCloseErr != nil {
		return srcCloseErr
	}

	return nil
}

func createJpg(img image.Image, destination string) error {
	file, err := os.Create(destination)
	if err != nil {
		return err
	}

	var imgOptions jpeg.Options
	imgOptions.Quality = 80
	errImageEncode := jpeg.Encode(file, img, &imgOptions)
	errFileClose := file.Close()

	if errImageEncode != nil && errFileClose != nil {
		return fmt.Errorf("error encoding and closing file: "+destination+" (%v); (%v)", errImageEncode, errFileClose)
	} else if errImageEncode != nil {
		return errImageEncode
	} else if errFileClose != nil {
		return errFileClose
	}

	return nil
}

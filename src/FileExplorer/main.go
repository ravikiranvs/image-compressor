package FileExplorer

import (
	"io/ioutil"
	"os"
	"strings"
)

// GetFiles Gets list of files from directory
func GetFiles(directory string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetJPEGFiles ...
func GetJPEGFiles(directory string) ([]os.FileInfo, error) {
	files, err := GetFiles(directory)

	if err != nil {
		return nil, err
	}

	var imageFiles []os.FileInfo
	for _, file := range files {
		fileName := strings.ToLower(file.Name())
		isImageFile := strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg")
		if isImageFile {
			imageFiles = append(imageFiles, file)
		}
	}

	return imageFiles, nil
}

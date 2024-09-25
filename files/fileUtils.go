package files

import (
	"bytes"
	"image"
	"os"
)

func LoadImgFromFile(fileName string) (image.Image, error) {
	normalFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	normalImg, _, err := image.Decode(bytes.NewReader(normalFile))
	if err != nil {
		return nil, err
	}

	return normalImg, nil
}

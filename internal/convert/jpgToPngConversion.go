package convert

import (
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func JpgToPngConversion(fileName string) error {
	jpegFile, err := os.Open("/tmp/" + fileName + ".jpg")
	if err != nil {
		log.Println(err)
		return err
	}
	defer jpegFile.Close()

	jpegImage, err := jpeg.Decode(jpegFile)
	if err != nil {
		log.Println(err)
		return err
	}

	pngFile, err := os.Create("/tmp/" + fileName + ".png")
	if err != nil {
		log.Println(err)
		return err
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, jpegImage)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

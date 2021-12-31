package compress

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type OptionCompress struct {
	Speed int
	Quality	int
	Debug bool
}

func GetFileTypeByName(path string) (string, error) {
	// Read the entire file into a byte slice
	b, err := ioutil.ReadFile("a.png")
	if err != nil {
		log.Fatal(err)
		return "Unknown", err
	}

	// Determine the content type of the image file
	return http.DetectContentType(b), nil
}

func GetFileType(file multipart.File) (string, error){
	// why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		return "Unknown", err
	}

	return http.DetectContentType(buff), nil
}

func Do(file multipart.File, fileType string) (output []byte, err error) {
	// check file type
	switch fileType {
		case "image/jpeg":
			input, _ := jpeg.Decode(file)
			compressed, err := JpgC(input, &jpeg.Options{Quality:30})
			if err != nil {
				return []byte{}, err
			}
			var w bytes.Buffer
			err = jpeg.Encode(&w, compressed, nil)
			if err != nil {
				return []byte{}, err
			}
			b := w.Bytes()

			return b, nil

		case "image/jpg":
			input, _ := jpeg.Decode(file)
			compressed, err := Jpg(input, OptionCompress{
				Speed:   40,
				Quality: 50,
				Debug:   false,
			})
			if err != nil {
				return []byte{}, err
			}
			var w bytes.Buffer
			err = jpeg.Encode(&w, compressed, nil)
			if err != nil {
				return []byte{}, err
			}
			b := w.Bytes()
			return b, nil

		case "image/png":
			input, _ := png.Decode(file)
			compressed, err := Png(input, OptionCompress{
				Speed:   2,
				Quality: 50,
				Debug:   false,
			})
			if err != nil {
				return []byte{}, err
			}
			var w bytes.Buffer
			err = png.Encode(&w, compressed)
			if err != nil {
				return []byte{}, err
			}
			b := w.Bytes()
			return b, nil

		case "image/web":
		case "image/gif":
		case "image/svg+xml":

		default:
			fmt.Println("Unsupported this file type")
	}

	return []byte{}, nil
}

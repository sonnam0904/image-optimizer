package image_optimizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDo(t *testing.T) {
	testData := []string{
		filepath.Base("img_test/i1.png"),
		filepath.Base("img_test/i2.jpg"),
	}

	for _, name := range testData {
		file, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fileType, err := GetFileTypeByName(name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ext := ""
		switch fileType {
		case "image/jpeg", "image/jpg":
			ext = "jpg"

		case "image/png":
			ext = "png"

		default:
			fmt.Println("Unsupported this file type")
		}

		zipData, err := Do(file, ext)
		if err != nil {
		}
		//imgC, err := os.Create("time.png")
		imgC, err := ioutil.TempFile("tmp", "vccloud-*."+ext)
		if err != nil {
			fmt.Println(err)
		}

		if _, err := imgC.Write(zipData); err != nil {
			fmt.Println(err)
		}

		fmt.Println("successfully compressed", file)
	}
}

func BenchmarkDo(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

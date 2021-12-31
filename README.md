# Image Optimizer
This package useful for compress image (png/jpg)

`
go get github.com/sonnam0904/image-optimizer
`


#Use
Install

`
apt-get install pngquant jpegoptim libjpeg62-turbo-dev
`
```
import (
    "fmt"
	"os"
	"github.com/sonnam0904/image-optimizer"
)

func main() {
	file, err := os.Open("img_imagepng1640922726.png")
	zip, err := compress.Do(file, "image/png")

	if err != nil {
		return
	}
	imgC, err := os.Create("upload.png")

	defer imgC.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := imgC.Write(zip); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("successfully compressed", file)
	return
}

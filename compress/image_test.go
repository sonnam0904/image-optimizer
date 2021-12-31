package compress

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
	//for i := 0; i < b.N; i++ {
	//	testData := []struct{
	//		input, expected string
	//	}{
	//		{"10/ 08- 2020. 09872  33333 111 22 day ne", ""},
	//		{"8888888888888 097.999 99.     09 xxxxxx xxx", ""},
	//		{"8888888888888 097.999.9999. 09 xxxxxx xxx", ""},
	//		{"0987 6567 xxxx xxx 87 xxxx 098 8987 765 xxxx", ""},
	//		{"0987 6567 xxxx xxx 87 xxxx 098 8987 765 x. xxx 0988888 888", ""},
	//		{"0987 6567 xxxx xxx 87 xxxx 888876765. 4543234 xxxx. xxx .0988888 888", ""},
	//		{"097.999.9998 097.999.9999. 09 xxxxxx xxx", ""},
	//		{"Zalo dùng số điện thoại 0386500999", ""},
	//		{"10/ 08- 2020. 09872  33333 111 22 day ne", ""},
	//	}
	//	//for _,i := range testData {
	//	//	_ = GetPhone(i.input)
	//	//	//fmt.Println(phone)
	//	//}
	//}
}

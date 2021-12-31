package image_optimizer
// #cgo LDFLAGS: -ljpeg
// #cgo darwin LDFLAGS: -L/opt/local/lib
// #cgo darwin CFLAGS: -I/opt/local/include
// #cgo freebsd LDFLAGS: -L/usr/local/lib
// #cgo freebsd CFLAGS: -I/usr/local/include
// #include <stdlib.h>
// extern int optimizeJPEG(unsigned char *inputbuffer, unsigned long inputsize, unsigned char **outputbuffer, unsigned long *outputsize, int quality);
// extern int encodeJPEG(unsigned char *inputbuffer, int width, int height, unsigned char **outputbuffer, unsigned long *outputsize, int quality);
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os/exec"
	"strings"
	"unsafe"
)
// need setup `apt-get install jpegoptim libjpeg62-turbo-dev`

// Run jpegoptim -h help
/**
jpegoptim v1.4.6  Copyright (C) 1996-2018, Timo Kokkonen
Usage: jpegoptim [options] <filenames>

  -d<path>, --dest=<path>
                    specify alternative destination directory for
                    optimized files (default is to overwrite originals)
  -f, --force       force optimization
  -h, --help        display this help and exit
  -m<quality>, --max=<quality>
                    set maximum image quality factor (disables lossless
                    optimization mode, which is by default on)
                    Valid quality values: 0 - 100
  -n, --noaction    don't really optimize files, just print results
  -S<size>, --size=<size>
                    Try to optimize file to given size (disables lossless
                    optimization mode). Target size is specified either in
                    kilo bytes (1 - n) or as percentage (1% - 99%)
  -T<threshold>, --threshold=<threshold>
                    keep old file if the gain is below a threshold (%)
  -b, --csv         print progress info in CSV format
  -o, --overwrite   overwrite target file even if it exists (meaningful
                    only when used with -d, --dest option)
  -p, --preserve    preserve file timestamps
  -P, --preserve-perms
                    preserve original file permissions by overwriting it
  -q, --quiet       quiet mode
  -t, --totals      print totals after processing all files
  -v, --verbose     enable verbose mode (positively chatty)
  -V, --version     print program version

  -s, --strip-all   strip all markers from output file
  --strip-none      do not strip any markers
  --strip-com       strip Comment markers from output file
  --strip-exif      strip Exif markers from output file
  --strip-iptc      strip IPTC/Photoshop (APP13) markers from output file
  --strip-icc       strip ICC profile markers from output file
  --strip-xmp       strip XMP markers markers from output file

  --all-normal      force all output files to be non-progressive
  --all-progressive force all output files to be progressive
  --stdout          send output to standard output (instead of a file)
  --stdin           read input from standard input (instead of a file)
**/
func Jpg(input image.Image, opt OptionCompress) (output image.Image, err error) {
	var oc OptionCompress
	if opt==oc {
		// set default value
		opt.Speed = 40
		opt.Quality = 50
	}

	var w bytes.Buffer
	err = jpeg.Encode(&w, input, &jpeg.Options{Quality:70})
	if err != nil {
		return input, err
	}
	args := []string{
		"-",
		"-f",
		"-o",
		"--strip-none",
		fmt.Sprintf("-m%v", opt.Quality),
		fmt.Sprintf("-S%v%%", opt.Speed),
	}
	// debug mode
	if opt.Debug {
		args = append(args,  "-v", "-t")
	}

	b := w.Bytes()
	compressed, err := JpgByte(b, args)

	if err != nil {
		return input, err
	}

	output, err = jpeg.Decode(bytes.NewReader(compressed))
	return output, err
}

// compress by jpeglib.h (C++)
func JpgC(input image.Image, opt *jpeg.Options) (output image.Image, err error) {

	var w bytes.Buffer
	err = jpeg.Encode(&w, input, opt)
	if err != nil {
		return input, err
	}
	b := w.Bytes()
	compressed, err := EncodeBytesOptimized(b, opt)
	if err != nil {
		return input, err
	}

	output, err = jpeg.Decode(bytes.NewReader(compressed))
	return output, err
}

// Optimize a JPEG bytes array if quality is -1, Optimize & Recompress if quality is between [0 - 100]
func EncodeBytesOptimized(srcBytes []byte, o *jpeg.Options) (outBytes []byte, err error) {
	if len(srcBytes) == 0 {
		err = errors.New("Image source is empty")
		return
	}
	// Clip quality to [-1, 100].
	if o != nil {
		quality := o.Quality
		if quality < -1 {
			quality = -1
		} else if quality > 100 {
			quality = 100
		}
	}

	csrcimg := (*C.uchar)(unsafe.Pointer(&srcBytes[0]))
	cinputsize := C.ulong(len(srcBytes))
	var coutimg *C.uchar
	var coutsize C.ulong
	code := C.optimizeJPEG(csrcimg, cinputsize, &coutimg, &coutsize, C.int(o.Quality))
	if code != 0 || coutsize == 0 {
		err = errors.New("Optimize failed")
		return
	}
	outBytes = C.GoBytes(unsafe.Pointer(coutimg), C.int(coutsize))
	C.free(unsafe.Pointer(coutimg))
	return
}

func JpgByte(input []byte, args []string) (output []byte, err error) {
	// "-", "--speed", speed
	cmd := exec.Command("jpegoptim", args...)
	cmd.Stdin = strings.NewReader(string(input))

	var o bytes.Buffer
	cmd.Stdout = &o
	err = cmd.Run()

	if err != nil {
		return input, err
	}
	output = o.Bytes()
	return output, nil
}




package image_optimizer

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"strings"
	"errors"
)
// need setup `apt-get install pngquant`
// @see https://salsa.debian.org/debian-phototools-team/pngquant or Run pngquant -help command
/**
usage:  pngquant [options] [ncolors] -- pngfile [pngfile ...]
        pngquant [options] [ncolors] - >stdout <stdin

options:
  --force           overwrite existing output files (synonym: -f)
  --skip-if-larger  only save converted files if they're smaller than original
  --output file     destination file path to use instead of --ext (synonym: -o)
  --ext new.png     set custom suffix/extension for output filenames
  --quality min-max don't save below min, use fewer colors below max (0-100)
  --speed N         speed/quality trade-off. 1=slow, 3=default, 11=fast & rough
  --nofs            disable Floyd-Steinberg dithering
  --posterize N     output lower-precision color (e.g. for ARGB4444 output)
  --strip           remove optional metadata (default on Mac)
  --verbose         print status messages (synonym: -v)

Quantizes one or more 32-bit RGBA PNGs to 8-bit (or smaller) RGBA-palette.
The output filename is the same as the input name except that
it ends in "-fs8.png", "-or8.png" or your custom extension (unless the
input is stdin, in which case the quantized image will go to stdout).
If you pass the special output path "-" and a single input file, that file
will be processed and the quantized image will go to stdout.
The default behavior if the output file exists is to skip the conversion;
use --force to overwrite. See man page for full list of options.
**/

func Png(input image.Image, opt OptionCompress) (output image.Image, err error) {
	var oc OptionCompress
	if opt==oc {
		// set default value
		opt.Speed = 2
		opt.Quality = 91
	}

	var w bytes.Buffer
	err = png.Encode(&w, input)
	if err != nil {
		return input, err
	}
	args := []string{
		"-",
		"--force",
		fmt.Sprintf("--speed=%v", opt.Speed),
		fmt.Sprintf("--quality=%v", opt.Quality),
	}
	// debug mode
	if opt.Debug {
		args = append(args,  "--verbose")
	}

	b := w.Bytes()
	compressed, err := PngByte(b, args)
	if err != nil {
		return input, err
	}

	output, err = png.Decode(bytes.NewReader(compressed))
	return output, err
}

func PngByte(input []byte, args []string) (output []byte, err error) {
	// "-", "--speed", speed
	cmd := exec.Command("pngquant", args...)
	cmd.Stdin = strings.NewReader(string(input))

	var o bytes.Buffer
	cmd.Stdout = &o
	err = cmd.Run()
	if err != nil {
		return input, errors.New("Can not compress this image")
	}

	output = o.Bytes()
	return output, nil
}


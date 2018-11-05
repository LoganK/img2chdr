package img2chdr

import (
	"fmt"
	"image"
	"io"
)

func cBytes(img image.Image) []uint8 {
	var res []uint8

	var buffer uint8
	count := 0
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			buffer <<= 1
			pixel, _, _, _ := img.At(x, y).RGBA()
			if pixel > 0 {
				buffer |= 1
			}

			count++
			if count == 8 {
				res = append(res, buffer)
				buffer = 0
				count = 0
			}
		}
	}

	if count > 0 {
		buffer <<= uint8(8 - count)
	}
	res = append(res, buffer)

	return res
}

func WriteHeader(img image.Image, imgName string, w io.Writer) error {
	charArray := cBytes(img)

	_, err := fmt.Fprintf(w,
`#if defined(__AVR)
  #include <avr/pgmspace.h>
#else
  #define PROGMEM
#endif

const unsigned char %s[] PROGMEM = {`, imgName)
	if err != nil {
		return fmt.Errorf("WriteHeader: %s", err)
	}

	for i, c := range(charArray) {
		if i % 16 == 0 {
			fmt.Fprint(w, "\n ")
		}

		fmt.Fprintf(w, " 0x%02x,", c)
	}

	fmt.Fprint(w, "\n};\n")

	return nil
}

package img2chdr

import (
	"github.com/esimov/dithergo"
	"github.com/disintegration/imaging"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

const (
	DITHER_ERROR = 1.18
)

type Converter struct {
	MaxX int // Maximum x resolution.
	MaxY int // Maximum y resolution.
}

// Load the given image file and write the converted form out.
func (c *Converter) ImageAsGrayscale(r io.Reader) (image.Image, error) {
	srcImage, _, err := image.Decode(r)
	if err != nil {
		return &image.Gray{}, fmt.Errorf("failed to decode image file: %s", err)
	}

	x, y := c.dimensions(srcImage)
	dstImage := imaging.Resize(srcImage, x, y, imaging.CatmullRom)

	/*
	dither.Dither{
		"Sierra-Lite",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 2.0 / 4.0 },
				[]float32{ 1.0 / 4.0, 1.0 / 4.0, 0.0 },
				[]float32{ 0.0, 0.0, 0.0 },
			},
		},
	},
	*/
	dither := dither.Dither{
		"Stucki",
		dither.Settings{
			[][]float32{
			[]float32{ 0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0 },
			[]float32{ 2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0 },
			[]float32{ 1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0 },
			},
		},
	}

	return dither.Monochrome(dstImage, DITHER_ERROR), nil
}

// Calculates the appropriate parameters for imaging.Resize to preserve the aspect ratio.
func (c *Converter) dimensions(image image.Image) (int, int) {
	bounds := image.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	srcAspect := float64(height) / float64(width)
	dstAspect := float64(c.MaxY) / float64(c.MaxX)

	if srcAspect < dstAspect {
		return 0, c.MaxY
	}
	return c.MaxX, 0
}


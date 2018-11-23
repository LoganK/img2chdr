package img2chdr

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/esimov/dithergo"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"sort"
)

var (
	DITHERERS = [...]dither.Dither{
		dither.Dither{
			"Null",
			dither.Settings{
				[][]float32{
					[]float32{0.0, 0.0, 0.0, 0.0, 0.0},
					[]float32{0.0, 0.0, 0.0, 0.0, 0.0},
					[]float32{0.0, 0.0, 0.0, 0.0, 0.0},
				},
			},
		},
		dither.Dither{
			"Burkes",
			dither.Settings{
				[][]float32{
					[]float32{0.0, 0.0, 0.0, 8.0 / 32.0, 4.0 / 32.0},
					[]float32{2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
					[]float32{0.0, 0.0, 0.0, 0.0, 0.0},
				},
			},
		},
		dither.Dither{
			"FloydSteinberg",
			dither.Settings{
				[][]float32{
					[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
					[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
					[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
				},
			},
		},
		dither.Dither{
			"Sierra-Lite",
			dither.Settings{
				[][]float32{
					[]float32{0.0, 0.0, 2.0 / 4.0},
					[]float32{1.0 / 4.0, 1.0 / 4.0, 0.0},
					[]float32{0.0, 0.0, 0.0},
				},
			},
		},
		dither.Dither{
			"Stucki",
			dither.Settings{
				[][]float32{
					[]float32{0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0},
					[]float32{2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0},
					[]float32{1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0},
				},
			},
		},
	}
	DitherError     float32 = 0.8
	Ditherer        int     = 1
	RangePercentile int     = 2
)

type Converter struct {
	MaxX int // Maximum x resolution.
	MaxY int // Maximum y resolution.
}

// Adds support for an alpha channel (as white) to color.Gray. minY/maxY represent the range of luminance and the result will
// be scaled accordingly.
func alphaGrayModel(c color.Color, minY, maxY uint8) color.Gray {
	gray := color.Gray16Model.Convert(c).(color.Gray16)
	var y16 uint32 = uint32(gray.Y) // "16" suffix indicates a 16-bit range, 32-bit value similar to RGBA()

	// Add alpha
	_, _, _, a := c.RGBA()
	y16 += (math.MaxUint16 - a)

	// Scale to range.
	minY16 := uint32(minY)<<8 | uint32(minY)
	maxY16 := uint32(maxY)<<8 | uint32(maxY)
	if y16 > minY16 {
		y16 -= minY16
	} else {
		y16 = 0
	}
	y16 = y16 * math.MaxUint16 / (maxY16 - minY16)

	// Round and scale to 8 bits
	y16 += 0x80
	if y16 > math.MaxUint16 {
		y16 = math.MaxUint16
	}
	y := uint8(y16 >> 8)

	if y > maxY {
		y = maxY
	}

	return color.Gray{y}
}

type SortableGrays []color.Gray

func (a SortableGrays) Len() int           { return len(a) }
func (a SortableGrays) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortableGrays) Less(i, j int) bool { return a[i].Y < a[j].Y }

// Determine the range of the image so we can scale pixels accordingly.
func colorRange(ys []color.Gray) (minY, maxY uint8) {
	sort.Sort(SortableGrays(ys))

	offset := len(ys) * RangePercentile / 100
	minY = ys[offset].Y
	maxY = ys[len(ys)-1-offset].Y

	// Look out for a range that's too narrow.
	if minY >= maxY {
		maxY = ys[len(ys)-1].Y
	}
	if minY >= maxY {
		minY = ys[0].Y
	}

	return
}

// Load the given image file and write the converted form out.
func (c *Converter) ImageAsGrayscale(r io.Reader) (image.Image, error) {
	srcImage, _, err := image.Decode(r)
	if err != nil {
		return &image.Gray{}, fmt.Errorf("failed to decode image file: %s", err)
	}

	// First pass: attempt to find a range to reduce the need for dithering.
	var minLum uint8 = 0
	var maxLum uint8 = math.MaxUint8
	bounds := srcImage.Bounds()
	rawLums := make([]color.Gray, 0, bounds.Dx()*bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rawLums = append(rawLums, alphaGrayModel(srcImage.At(x, y), minLum, maxLum))
		}
	}
	minLum, maxLum = colorRange(rawLums)

	// Manually convert to grayscale to use our custom model.
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.SetGray(x, y, alphaGrayModel(srcImage.At(x, y), minLum, maxLum))
		}
	}

	x, y := c.dimensions(gray)
	dstImage := imaging.Resize(gray, x, y, imaging.CatmullRom)

	return DITHERERS[Ditherer].Monochrome(dstImage, DitherError), nil
}

// Calculates the appropriate parameters for imaging.Resize to preserve the aspect ratio.
func (c *Converter) dimensions(image image.Image) (int, int) {
	bounds := image.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	srcAspect := float64(height) / float64(width)
	dstAspect := float64(c.MaxY) / float64(c.MaxX)

	if srcAspect > dstAspect {
		return 0, c.MaxY
	}
	return c.MaxX, 0
}

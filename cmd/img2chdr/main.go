package main

import (
	"flag"
	"github.com/logank/img2chdr"
	"image/png"
	"log"
	"os"
	"path"
	"strings"
)

var (
	inFileName = flag.String("in", "", "The input file name")
	imgName = flag.String("name", "", "The name of the image. Defaults based on --in")
	maxX = flag.Int("max_x", 296, "The maximum x resolution")
	maxY = flag.Int("max_y", 128, "The maximum y resolution")
	outFileName = flag.String("out", "", "The output file name")
	outImage = flag.Bool("out_image", false, "If given, write the output image for easy viewing")
)

func handleFlags() {
	flag.Parse()

	if len(*inFileName) == 0 && flag.NArg() > 0 {
		*inFileName = flag.Arg(0)
	}
	if len(*inFileName) == 0 {
		flag.Usage()
		log.Fatal("must provide an input file name")
	}

	if len(*imgName) == 0 {
		*imgName = strings.TrimSuffix(*inFileName, path.Ext(*inFileName))
	}

	if len(*outFileName) == 0 {
		if flag.NArg() > 1 {
			*outFileName = flag.Arg(1)
		} else {
			*outFileName = *imgName + ".h"
		}
	}
}

func main() {
	handleFlags()

	inReader, err := os.Open(*inFileName)
	if err != nil {
		log.Fatalf("failed to open input file '%s': %s", os.Args[1], err)
	}
	defer inReader.Close()

	converter := img2chdr.Converter {
		MaxX: *maxX,
		MaxY: *maxY,
	}
	grayImage, err := converter.ImageAsGrayscale(inReader)
	if err != nil {
		log.Fatalf("Failed to convert file: %s", err)
	}

	outWriter, err := os.Create(*outFileName)
	if err != nil {
		log.Fatalf("failed to open output file '%s': %s", os.Args[2], err)
	}
	defer outWriter.Close()

	err = img2chdr.WriteHeader(grayImage, *imgName, outWriter)
	if err != nil {
		log.Fatalf("failed to write outpfile file: %s", err)
	}

	if *outImage {
		name := *imgName + ".badgy.png"
		outImgFile, err := os.Create(name)
		if err != nil {
			log.Fatalf("failed to open output image file '%s': %s", name, err)
		}
		defer outImgFile.Close()

		// Encode the grayscale image to the output file
		png.Encode(outImgFile, grayImage)
	}

}

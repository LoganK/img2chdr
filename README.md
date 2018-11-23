# IMG2CHDR

Converts an image file to a monochrome bitmap written as a C header. Originally
intended for use wth (https://github.com/ZinggJM/GxEPD).

Supports dynamic-range scaling, resizing, and dithering.

Usage:

```bash
go get -u github.com/logank/img2chdr/cmd/...
img2chdr infile.png outfile.h
```

```
Usage of img2chdr:
  -dither
        Use dithering on the output image (default true)
  -in string
        The input file name
  -max_x int
        The maximum x resolution (default 296)
  -max_y int
        The maximum y resolution (default 128)
  -name string
        The name of the image. Defaults based on --in
  -out string
        The output file name
  -out_image
        If given, write the output image for easy viewing
  -range_percentile int
        Scale the image intensity range excluding this percentile of pixels (default 2)
```


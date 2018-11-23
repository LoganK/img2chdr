# IMG2CHDR

Converts an image file to a monochrome bitmap written as a C header. Originally
intended for use wth github.com/ZinggJM/GxEPD.

Supports dynamic-range scaling, resizing, and dithering.

Usage:

```bash
go get -u github.com/logank/img2chdr
img2chdr infile.png outfile.h
```


// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transforms

import (
	"image"
)

// Rgb2Gray function converts RGB to a gray scale array.
func Rgb2Gray(colorImg image.NRGBA) [][]float64 {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pixels := make([][]float64, h)

	for i := range pixels {
		pixels[i] = make([]float64, w)
		for j := range pixels[i] {
			color := colorImg.At(j, i)
			r, g, b, _ := color.RGBA()
			//lu := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
			lu := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
			pixels[i][j] = float64(lu)
		}
	}
	return pixels
}

// FlattenPixels function flattens 2d array into 1d array.
func FlattenPixels(pixels [][]float64, x int, y int) []float64 {
	flattens := make([]float64, x*y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[i][j]
		}
	}
	return flattens
}

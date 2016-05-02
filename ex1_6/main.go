// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	// "fmt"
)

//!-main
// Packages not needed by version in book.
import (
	"log"
	"net/http"
	"time"
)

//!+main

var palette []color.Color

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 128 // next color in palette
	twoThirdsPi = (math.Pi * 2) / 3
	fourThirdsPi = (math.Pi * 4) / 3
)

func main() {
	//!-main
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w)
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	//!+main
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 128    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	palette = make([]color.Color, 256, 256)
	palette[0] = color.Black

	for i := 1; i<256; i++ {
		var rads = ((float64(i)-1.0)/255.0) * 2.0 * math.Pi
		var r = uint8((math.Sin(rads + 0) + 1) * 128)
		var g = uint8((math.Sin(rads + twoThirdsPi) + 1) * 128)
		var b = uint8((math.Sin(rads + fourThirdsPi) + 1) * 128)

		palette[i] = color.RGBA{r, g, b, 0xff}
		// fmt.Println(rads, r, g, b)
	}


	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		colorIndex := 0
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			//colorIndex = (colorIndex+1)%255
			//colorIndex = 1 + int((t / 2*math.Pi) * 255)
			colorIndex = int(1 + ((x+1) * 255))/2
			// fmt.Println(colorIndex)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(colorIndex + 1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main

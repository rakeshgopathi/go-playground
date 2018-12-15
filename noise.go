package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	noiseImage *image.RGBA
)

type rand struct {
	x, y, z, w uint32
}

func (r *rand) next() uint32 {
	// math/rand is too slow to keep 60 FPS on web browsers.
	// Use Xorshift instead: http://en.wikipedia.org/wiki/Xorshift
	t := r.x ^ (r.x << 11)
	r.x, r.y, r.z = r.y, r.z, r.w
	r.w = (r.w ^ (r.w >> 19)) ^ (t ^ (t >> 8))
	return r.w
}

var theRand = &rand{12345678, 4185243, 776511, 45411}
var isInitialized bool = false
var sliceMap map[int][]int = make(map[int][]int)

func update(screen *ebiten.Image) error {
	const noiseRectCount = 10
	if (!isInitialized) {
		const rectSizeHorizontal =  screenWidth / noiseRectCount
		const rectSizeVertical =  screenHeight / noiseRectCount
		const sliceSize = rectSizeHorizontal * rectSizeVertical

		pixel := 0
		for i := 0; i < noiseRectCount; i++ {
			slice := make([]int, sliceSize)
			for j := 0; j < sliceSize; j++ {
				pixel = j
				startIndex := (i * rectSizeVertical) + (j * rectSizeHorizontal)
				for k := 0; k < rectSizeHorizontal; k++ {
					if pixel >= sliceSize {
						j = pixel
						break
					}
					slice[pixel] =  startIndex + k
				}
			}
			sliceMap[i] = slice
			fmt.Println(slice)
		}
		isInitialized = true
	}

	for i := 0; i < len(sliceMap); i++ {
		slice := sliceMap[i]
		x := theRand.next()
		randomColor := uint8(x)
		if i%2 == 0 {
			randomColor = 0x00
		} 
		for j := 0; j < len(slice); j++ {
			// Generate the noise with random RGB values.
			noiseImage.Pix[slice[j]] = randomColor
		}
		fmt.Println(slice)
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	screen.ReplacePixels(noiseImage.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	return nil
}

func main() {
	noiseImage = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Noise (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}

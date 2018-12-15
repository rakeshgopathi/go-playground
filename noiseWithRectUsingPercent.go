package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
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

func update(screen *ebiten.Image) error {
	// Percentage of smaller rectagles W.r.t the screen measurements
	const noiseRectPercent = 10
	const rectSizeHorizontal =  (screenWidth / 100) * noiseRectPercent
	const rectSizeVertical =  (screenHeight / 100) * noiseRectPercent
	// Number of rectagles
	const rectCountHorizontal = screenWidth / rectSizeHorizontal
	const rectCountVertical = screenHeight / rectSizeVertical
	for i := 0; i < rectCountVertical; i++ {	
		for j := 0; j < rectCountHorizontal; j++ {
			// Generate the noise with random RGB values.
			x := theRand.next()
			randomColor := uint8(x >> 24)
			randColor2 := uint8(x >> 16)
			randColor3 := uint8(x >> 8)
			randColor4 := uint8(x >> 4)//0x00

			ebitenutil.DrawRect(screen, float64(j * rectSizeHorizontal), float64(i * rectSizeVertical), rectSizeHorizontal, rectSizeVertical, color.RGBA{randomColor, randColor2, randColor3, randColor4})
		}
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Noise (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}


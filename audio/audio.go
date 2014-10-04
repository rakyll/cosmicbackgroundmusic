package audio

import (
	"image"
	"io"
	"log"
	"sync"

	_ "image/png"

	"github.com/campoy/audio/audio"
)

var (
	img image.Image

	mu      sync.Mutex
	playing *sample
)

type sample struct {
	sinFreqs []int64
	sqFreqs  []int64
	sawFreqs []int64
}

func (s *sample) Play() error {
	// create 3 go routines
	panic("not yet implemented")
}

func (s *sample) Stop() {
	panic("not yet implemented")
}

func Initialize(r io.Reader) error {
	err := audio.Initialize()
	if err != nil {
		return err
	}
	img, _, err = image.Decode(r)
	return err
}

func Play(x, y, d int) error {
	log.Printf("Playing noise for %v, %v, %v", x, y, d)

	mu.Lock()
	defer mu.Unlock()

	size := img.Bounds().Size()
	// scan the area to generate a sample
	// filter out the transparent pixels
	rT, gT, bT := 0.0, 0.0, 0.0

	for i := x - d; i < x+d; i++ {
		if i < 0 || i > size.X {
			continue
		}
		for j := y - d; j < y+d; y++ {
			if j < 0 || j > size.Y {
				continue
			}
			color := img.At(i, j)
			// if transparent, skip
			r, g, b, a := color.RGBA()
			if a == 0 {
				continue
			}
			// b is for sin, g is for square, r is for saw
			rT += float64(r)
			gT += float64(g)
			bT += float64(b)
		}
	}

	log.Println(rT, gT, bT)
	// determine the instruments depending on the microwave intensity
	// cold: sin wave
	// med: square wave
	// hot: saw wave

	// determine a frequencies for each instrumentes

	// play the sample
	s := &sample{}
	if playing != nil {
		playing.Stop()
	}
	playing = s
	return playing.Play()
}

func Terminate() error {
	mu.Lock()
	defer mu.Unlock()

	playing = nil
	return audio.Terminate()
}

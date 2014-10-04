package audio

import (
	"image"
	"io"
	"log"
	"sync"
	"time"

	_ "image/png"

	"github.com/campoy/audio/audio"
)

const (
	sampleRate = 44100
)

var (
	img image.Image

	mu      sync.Mutex
	playing *sample
)

type sample struct {
	sine   []float64
	square []float64
	saw    []float64

	sineDone   chan bool
	squareDone chan bool
	sawDone    chan bool
}

func (s *sample) Play() error {
	s.sineDone = make(chan bool, 1)
	s.squareDone = make(chan bool, 1)
	s.sawDone = make(chan bool, 1)

	play("sine", time.Second, s.sine, s.sineDone)
	play("square", 500*time.Millisecond, s.square, s.squareDone)
	play("saw", 250*time.Millisecond, s.saw, s.sawDone)
	return nil
}

func play(wave string, dur time.Duration, sample []float64, done chan bool) {
	if len(sample) == 0 {
		return
	}
	go func() {
		i := 0
		var inst audio.Instrument
		for {
			select {
			case <-done:
				inst.Stop()
				return
			default:
				if inst != nil {
					inst.Stop()
				}
				switch wave {
				case "sine":
					inst = audio.NewSine(sample[i], sampleRate)
				case "square":
					inst = audio.NewSquare(sample[i], sampleRate)
				case "saw":
					inst = audio.NewSaw(sample[i], sampleRate)
				}
				go func() {
					inst.Play()
				}()
				i = (i + 1) % len(sample)
				<-time.After(dur)
			}
		}
	}()
}

func (s *sample) Stop() {
	s.sineDone <- true
	s.squareDone <- true
	s.sawDone <- true
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
	num := float64(0)
	// scan the area to generate a sample
	// filter out the transparent pixels
	rT, gT, bT := 0.0, 0.0, 0.0

	for i := x - d; i < x+d; i++ {
		if i < 0 || i > size.X {
			continue
		}
		for j := y - d; j < y+d; j++ {
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
			num++
		}
	}

	// determine the instruments depending on the microwave intensity
	// cold: sin wave
	// med: square wave
	// hot: saw wave

	// determine a frequencies for each instrumentes
	avgR := rT / num
	avgG := gT / num
	avgB := bT / num

	log.Printf("Average intensity at [%d, %d, %d] is [%f, %f, %f]", x, y, d, avgR, avgG, avgB)

	s := &sample{
		sine:   []float64{1000, 1200, 1500},
		square: []float64{500, 300},
		saw:    []float64{200, 900, 200},
	}
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

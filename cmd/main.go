package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rakyll/cosmicmusic/audio"
)

func main() {
	r, err := os.Open("./data/bg.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := audio.Initialize(r); err != nil {
		log.Fatal(err)
	}
	log.Println("Initiated portaudio.")
	defer audio.Terminate()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "render the image here")
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		x, _ := strconv.Atoi(r.URL.Query().Get("x"))
		y, _ := strconv.Atoi(r.URL.Query().Get("y"))
		err := audio.Play(x, y, 50)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, "ok")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

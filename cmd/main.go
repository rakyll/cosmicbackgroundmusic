package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rakyll/cosmicbackgroundmusic/audio"
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

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r, err := os.Open("./index.html")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		io.Copy(w, r)
	})

	http.HandleFunc("/img", func(w http.ResponseWriter, req *http.Request) {
		r, err := os.Open("./data/bg.png")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		io.Copy(w, r)
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		x, _ := strconv.ParseFloat(r.URL.Query().Get("x"), 64)
		y, _ := strconv.ParseFloat(r.URL.Query().Get("y"), 64)
		err := audio.Play(int(x), int(y), 32)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, "ok")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

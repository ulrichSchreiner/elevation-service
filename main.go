package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ulrichSchreiner/go-elevations/geoelevations"
)

var (
	client = http.DefaultClient
	srtm   *geoelevations.Srtm
)

func main() {
	cachedir := os.Getenv("ELEVATION_CACHE")
	s, err := geoelevations.NewSrtmWithCustomCacheDir(client, cachedir)
	if err != nil {
		panic(err.Error())
	}
	srtm = s

	http.HandleFunc("/", elevation)
	http.ListenAndServe(":8000", nil)
}

func elevation(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	slat := r.Form.Get("lat")
	slng := r.Form.Get("lng")
	lat, err := strconv.ParseFloat(slat, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lng, err := strconv.ParseFloat(slng, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	elevation, err := srtm.GetElevation(client, lat, lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%f\n", elevation)
}

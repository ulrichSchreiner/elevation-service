package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ulrichSchreiner/go-elevations/geoelevations"
)

type elevationservice struct {
	client *http.Client
	srtm   *geoelevations.Srtm
}

func newElevationService() (*elevationservice, error) {
	cachedir := os.Getenv("ELEVATION_CACHE")
	s, err := geoelevations.NewSrtmWithCustomCacheDir(http.DefaultClient, cachedir)
	if err != nil {
		return nil, err
	}
	return &elevationservice{
		client: http.DefaultClient,
		srtm:   s,
	}, nil
}

func (es *elevationservice) elevation(w http.ResponseWriter, r *http.Request) {
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
	elevation, err := es.srtm.GetElevation(es.client, lat, lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%f\n", elevation)
}

func main() {
	es, err := newElevationService()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/elevation", es.elevation)
	http.ListenAndServe(":8000", nil)
}

package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/andewx/go-iio/plot"
	"github.com/andewx/go-iio/sdr"
)

func TestFFT(t *testing.T) {
	// Test FFT
	// Test FFTToReal
	// Generate Test FFT Data on Size 16 FFT

	sdr, err := sdr.ConnectNetwork("192.168.2.1", "ad9361")
	if err != nil {
		t.Errorf("Failed to connect to SDR")
	}

	data := make([]complex64, 16)

	for i := 0; i < sdr.Params.FourierPoint; i++ {
		rl := rand.Float64()
		im := rand.Float64()
		data[i] = complex(float32(rl), float32(im))
	}

	grapher := plot.NewFFTPlotServer(sdr, data)
	fmt.Printf("Starting FFT Plot Server\n")
	fmt.Printf("...Listening on locahost:8080\n")
	grapher.Start()
}

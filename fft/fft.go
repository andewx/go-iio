package fft

import (
	"fmt"
	"math"
	"math/cmplx"
)

// FFT performs Fast Fourier Transform on complex input
// Input length must be a power of 2
func FFT(x []complex128) ([]complex128, error) {
	n := len(x)
	if n&(n-1) != 0 {
		return nil, fmt.Errorf("input length must be power of 2, got %d", n)
	}

	// Make a copy to avoid modifying input
	out := make([]complex128, n)
	copy(out, x)

	// Bit-reverse permutation
	for i := 0; i < n; i++ {
		j := bitReverse(i, n)
		if i < j {
			out[i], out[j] = out[j], out[i]
		}
	}

	// Cooley-Tukey FFT
	for size := 2; size <= n; size *= 2 {
		halfsize := size / 2
		//tablestep := n / size

		for i := 0; i < n; i += size {
			for j := i; j < i+halfsize; j++ {
				k := j + halfsize
				temp := out[j]

				// twiddle := cexp(-2π * i * j / n)
				theta := -2.0 * math.Pi * float64(j-i) / float64(size)
				twiddle := cmplx.Rect(1, theta)

				out[j] = temp + twiddle*out[k]
				out[k] = temp - twiddle*out[k]
			}
		}
	}

	return out, nil
}

// IFFT performs Inverse Fast Fourier Transform
func IFFT(x []complex128) ([]complex128, error) {
	n := len(x)

	// Make a copy and conjugate
	input := make([]complex128, n)
	for i := 0; i < n; i++ {
		input[i] = cmplx.Conj(x[i])
	}

	// Perform forward FFT
	result, err := FFT(input)
	if err != nil {
		return nil, err
	}

	// Conjugate and scale
	scale := complex(1/float64(n), 0)
	for i := 0; i < n; i++ {
		result[i] = cmplx.Conj(result[i]) * scale
	}

	return result, nil
}

// bitReverse reverses the bits of x up to width n
func bitReverse(x, n int) int {
	result := 0
	width := uint(math.Log2(float64(n)))

	for i := uint(0); i < width; i++ {
		if x&(1<<i) != 0 {
			result |= 1 << (width - 1 - i)
		}
	}

	return result
}

// Helper functions for common operations
type Window int

const (
	Rectangular Window = iota
	Hanning
	Hamming
	Blackman
)

// ApplyWindow applies the specified window function to the input data
func ApplyWindow(x []complex128, window Window) []complex128 {
	n := len(x)
	result := make([]complex128, n)

	for i := 0; i < n; i++ {
		var w float64
		switch window {
		case Rectangular:
			w = 1.0
		case Hanning:
			w = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(n-1)))
		case Hamming:
			w = 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
		case Blackman:
			w = 0.42 - 0.5*math.Cos(2*math.Pi*float64(i)/float64(n-1)) +
				0.08*math.Cos(4*math.Pi*float64(i)/float64(n-1))
		}
		result[i] = complex(real(x[i])*w, imag(x[i])*w)
	}

	return result
}

// PowerSpectrum computes the power spectrum of the FFT result
func PowerSpectrum(x []complex128) []float64 {
	n := len(x)
	result := make([]float64, n)

	for i := 0; i < n; i++ {
		result[i] = cmplx.Abs(x[i]) * cmplx.Abs(x[i])
	}

	return result
}

// Convenience function for real input
func RealFFT(x []float64) ([]complex128, error) {
	n := len(x)
	input := make([]complex128, n)
	for i, v := range x {
		input[i] = complex(v, 0)
	}
	return FFT(input)
}

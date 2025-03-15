package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// ReadFir reads a FIR filter from a file
// and returns the coefficients as a slice of float64
func ReadFir(filename string, filterLen int) ([]int16, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fir []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		coefficient, err := strconv.ParseFloat(scanner.Text(), 64)
		if err == nil {
			fir = append(fir, coefficient)
		}
	}

	if len(fir) != filterLen {
		return ConvertFir(fir), fmt.Errorf("Expected %d coefficients, got %d", filterLen, len(fir))
	}
	return ConvertFir(fir), nil

}

// ConvertFir converts a slice of float64 to a slice of int16 2's complement data
func ConvertFir(fir []float64) []int16 {
	var firInt []int16
	for _, val := range fir {
		firInt = append(firInt, int16(val*(1<<15)))
	}
	return firInt
}

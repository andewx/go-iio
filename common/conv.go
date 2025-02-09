package common

func MhZ(x float64) float64 {
	return x*(1000000) + 0.5
}

func GhZ(x float64) float64 {
	return x*(1000000000) + 0.5
}

func KhZ(x float64) float64 {
	return x * (1000)
}

func Rate(x float64) float64 {
	return 1 / x
}

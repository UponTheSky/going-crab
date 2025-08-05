package calculate

func SumInt(numbers []int) int {
	total := 0

	for _, n := range numbers {
		total += n
	}

	return total
}

func SumFloat64(numbers []float64) float64 {
	total := 0.0

	for _, n := range numbers {
		total += n
	}

	return total
}

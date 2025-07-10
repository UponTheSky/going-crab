package calculate

func Sum(numbers []int) int {
	total := 0

	for _, n := range numbers {
		total += n
	}

	return total
}

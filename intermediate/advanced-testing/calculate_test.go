package calculate

import "testing"

func TestSum(t *testing.T) {
	// Asset
	input := []int{1, 2, 3, 4, 5}
	expected := 15

	// Act
	total := Sum(input)

	// Assert
	if total != expected {
		t.Errorf("test unsuccessful - got: %v, want: %v", total, expected)
	}
}

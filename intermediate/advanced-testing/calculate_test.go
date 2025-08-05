package calculate

import "testing"

func TestSum(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		// Asset
		input := []int{1, 2, 3, 4, 5}
		expected := 15

		// Act
		total := SumInt(input)

		// Assert
		if total != expected {
			t.Errorf("test unsuccessful - got: %v, want: %v", total, expected)
		}
	})

	t.Run("Float64", func(t *testing.T) {
		// Asset
		input := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		expected := 15.0

		// Act
		total := SumFloat64(input)

		// Assert
		if total != expected {
			t.Errorf("test unsuccessful - got: %v, want: %v", total, expected)
		}
	})
}

func TestSumWithDB(t *testing.T) {
	// setup
	db := make(map[string][]int)
	db["odd"] = []int{1, 3, 5}
	db["even"] = []int{2, 4, 6}

	testCases := []struct {
		name     string
		key      string
		expected int
	}{
		{name: "sum odd integers", key: "odd", expected: 9},
		{name: "sum even integers", key: "even", expected: 12},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Asset
			input := db[tc.key]

			// Act
			total := SumInt(input)

			// Assert
			if total != tc.expected {
				t.Errorf("test unsuccessful - got: %v, want: %v", total, tc.expected)
			}
		})
	}

	// teardown
	delete(db, "odd")
	delete(db, "even")
}

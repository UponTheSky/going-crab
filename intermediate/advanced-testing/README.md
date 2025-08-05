# Test Revisited - Advanced Testing
As we have more features being added, we need more test functions to guarantee that the entire application works as expected. In one of the previous chapters, [Test Your Server Automatically](../../basic/hello-test/README.md), we covered how to test Go functions in simple ways, using [`testing`](https://pkg.go.dev/testing) package. In this chapter, we will look into more details about the package and how to organize your test code for a large-scaled application.

Of course, nothing is better than reading through the whole documentation page, but 

## Recap - Writing a big `TestXxx(t *testing.T)` function
As you remember, every test in Go using `testing` package starts from writing a big `TestXxx(t *testing.T)` function, inside a file with `_test` suffix. Inside the function, normally you follow the following steps:

1. Asset - Prepare input, expected output, and other required setups for testing.
2. Act - Call the function to be tested, and record the result.
3. Assert - Compare the result from Act and the expected result from Asset. Here methods like `t.Errorf` or `t.Fatalf`(or any similar function) are used to mark failures of comparisons. 

The following code example summarizes the recap:
```go
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
```

## Organize a Test Using Subtests
Now, what if an application scales up and we need more complex tests? Suppose we test various cases of summations, such as summing not only integers, but also floats, etc. Of course we could separately implement test functions(`TestIntegerSum`, `TestFloatSum`). However, it would be better if we organize those test functions in a same category - `TestSum`. 

In `testing` package, we have something called "subtests" that enable us to tidy up multiple tests into several categories. 

```go
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

```

You would recognize the exact same pattern of testing inside each `t.Run()`as the original `TestSum`. By doing so, we can write several tests that are related to summing numbers under the same `TestSum` big function. For your information, each `t.Run()` runs separate tests in separate goroutines. 

### Setup and Teardown
One advantage of arranging several tests with similar interests is that we can share the same *setup* and *teardown* mechanism. For those who are not familiar with these concepts, *setup* is the same as preparing *assets* phase but connotating resources of big scales, such as running a database docker container dedicated to testing. Since such resources usually require to be freed after using, we need the tidying up phase, *teardown*.

To demonstrate the idea simple, let's *setup* a mock database using `map` data structure, and *teardown* it after all the tests are run.

```go
func TestSumWithDB(t *testing.T) {
	// setup
	db := make(map[string][]int)
	db["odd"] = []int{1, 3, 5}
	db["even"] = []int{2, 4, 6}

	t.Run("odd", func(t *testing.T) {
		// Asset
		input := db["odd"]
		expected := 9

		// Act
		total := SumInt(input)

		// Assert
		if total != expected {
			t.Errorf("test unsuccessful - got: %v, want: %v", total, expected)
		}
	})

	t.Run("even", func(t *testing.T) {
		// Asset
		input := db["even"]
		expected := 12

		// Act
		total := SumInt(input)

		// Assert
		if total != expected {
			t.Errorf("test unsuccessful - got: %v, want: %v", total, expected)
		}
	})

	// teardown
	delete(db, "odd")
	delete(db, "even")
}
```

The logic is very straightforward hence we won't talk about the code example in details, but you will see how the *setup* and *teardown* phases work here.

## Refactoring Tests - Table-driven Strategy
When the application gets bigger, there would be repeats in tests against the core logics. If you see the code examples above, all the subtests are having the exact same pattern. How can we reduce these unnecessary duplications?

In this case, we can think of the [table-driven test technique](https://go.dev/wiki/TableDrivenTests), where the test inputs and exected outputs are provided(other information as metdata could also be provided). 

Let us refactor the `TestSumWithDB` first.

```go
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
```

As you see, it is very useful for testing one single function against multiple input cases. Note that when defining the `testCases` variable, we simply define a nameless `struct` type on the fly, which has `name`, `key`, and `expected` fields. 

## Conclusion
Of course, this couldn't be the end of the story about testing. We still have several more topics, including [helper functions](https://pkg.go.dev/testing#T.Helper), [parallel testing](https://pkg.go.dev/testing#T.Parallel), etc. But don't worry - once you understand the overall picture, these are all the techniques that you can pick quickly when you actually need them.

## Exercise
No exercises today! Sometimes teachers also need to take some rest...

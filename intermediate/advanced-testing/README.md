# Test Revisited - Advanced Testing
As we have more features being added, we need more test functions to guarantee that the entire application works as expected. In one of the previous chapters, [Test Your Server Automatically](../../basic/hello-test/README.md), we covered how to test Go functions in simple ways, using [`testing`](https://pkg.go.dev/testing) package. In this chapter, we will look into more details about the package and how to organize your test code for a large-scaled application.

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

## Setup and Teardown

## Table-Driven Tests

## Conclusion

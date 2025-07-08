package app

import "testing"

func TestEmpty(t *testing.T) {
	input := ""
	want := ""

	got := Capitalize(input)

	if got != want {
		t.Errorf("Test Empty - got=%v, want=%v", got, want)
	}
}

func TestText(t *testing.T) {
	input := "Ken Jeong!"
	want := "KEN JEONG!"

	got := Capitalize(input)

	if got != want {
		t.Errorf("Test Empty - got=%v, want=%v", got, want)
	}
}

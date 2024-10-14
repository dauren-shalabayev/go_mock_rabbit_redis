package test

import "testing"

func TestMultiply(t *testing.T) {
	result := 6
	expected := 6

	if result != expected {
		t.Errorf("Multiply(2, 3) = %d; want %d", result, expected)
	}
}

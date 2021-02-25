package main

import "testing"

func TestIsPrime(t *testing.T) {
	tests := []struct {
		n       uint64
		isPrime bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{11, true},
		{21, false},
		{83, true},
		{89, true},
		{97, true},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			if test.isPrime != isPrime(test.n) {
				t.Errorf("something wrong with %d", test.n)
			}
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		x, y, expected int
	}{
		{0, 0, 1},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual := pow(test.x, test.y)
			if test.expected != actual {
				t.Errorf("according to your code %d^%d is %d", test.x, test.y, actual)
			}
		})
	}
}

func TestLargestPandigitalPrime(t *testing.T) {
	tests := []struct {
		n          int
		pandigital uint64
	}{
		{1, 1},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			ch := make(chan uint64)
			go largestPandigitalPrime(ch, test.n)
			actual := <-ch
			if test.pandigital != actual {
				t.Errorf("expected: %d, actual: %d", test.pandigital, actual)
			}
		})
	}
}

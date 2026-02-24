package arrays

import (
	"reflect"
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want int) {
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	}

	t.Run("collection of size 5", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		got := Sum(numbers)
		want := 15
		assertCorrectMessage(t, got, want)
	})

	t.Run("collction of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7}
		got := Sum(numbers)
		want := 28
		assertCorrectMessage(t, got, want)
	})
}

func TestSumAll(t *testing.T) {
	t.Run("make the sums of slices", func(t *testing.T) {
		got := SumAll([]int{1, 2}, []int{0, 9})
		want := []int{3, 9}
		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestSumAllTails(t *testing.T) {
	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make the sums of slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})
	t.Run("safely sum empty slies", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})
}

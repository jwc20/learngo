package iterations

import "testing"

func TestRepeat(t *testing.T) {
	t.Run("repeats a character", func(t *testing.T) {
		repeated := Repeat("a", 5)
		expected := "aaaaa"

		assertCorrectMessage(t, repeated, expected)
	})

	t.Run("repeats a character using range", func(t *testing.T) {
		repeated := Repeat2("a", 10)
		expected := "aaaaaaaaaa"

		assertCorrectMessage(t, repeated, expected)
	})
}

func assertCorrectMessage(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("expected %q but got %q", got, want)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a", 1000000)
	}
}

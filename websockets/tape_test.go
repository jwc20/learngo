package poker_test

import (
	"io"
	"testing"

	poker "github.com/jwc20/learngowithtests/websockets"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	// init pointer to the Tape
	tape := &poker.Tape{File: file}

	// write to the Tape
	tape.Write([]byte("abc"))

	// seek to the start
	file.Seek(0, io.SeekStart)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

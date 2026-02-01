package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

func NewCLI(playerStore PlayerStore, in io.Reader) *CLI {
	return &CLI{playerStore, bufio.NewScanner(in)}
}

func (cli *CLI) PlayPoker() {
	userInput := readLine(cli.in)
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func readLine(reader *bufio.Scanner) string {
	reader.Scan()
	return reader.Text()
}

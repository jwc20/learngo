package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type CLI struct {
	game Game
	in   *bufio.Scanner
	out  io.Writer
}

func NewCLI(game Game, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		game: game,
		in:   bufio.NewScanner(in),
		out:  out,
	}
}

const PlayerPrompt = "Please enter the number of players: "

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, _ := strconv.Atoi(cli.readLine())

	cli.scheduleBlindAlerts(numberOfPlayers)
	userInput := cli.readLine()
	winner := extractWinner(userInput)
	cli.game.Finish(winner)
	// cli.game.store.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts(numberOfPlayers int) {
	cli.game.Start(numberOfPlayers)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/jwc20/learngowithtests/practice/01_httpserver_practices/httpserver"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()

}

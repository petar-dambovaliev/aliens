package main

import (
	"flag"
	"os"
)

func main() {
	num := flag.Int("n", 2, "# of aliens")
	f := flag.String("f", "", "file containing the cities map")
	flag.Parse()

	n := *num
	citiesMapFile := *f
	citiesMapFile = "world.test"

	//if citiesMapFile == "" {
	//	panic("-f; file is required")
	//}

	readFile, err := os.Open(citiesMapFile)

	if err != nil {
		panic(err)
	}
	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			panic(err)
		}
	}(readFile)

	world, err := WorldFromReader(readFile)
	if err != nil {
		panic(err)
	}

	aliens := GenerateAliens(world, n)
	game := NewGame(world, aliens)
	game.Start(maxMoves, os.Stdout)
	err = game.export(os.Stdout)
	if err != nil {
		panic(err)
	}
}

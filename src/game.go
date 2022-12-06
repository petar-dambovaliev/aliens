package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

const maxMoves = 10000

func GenerateAliens(w *World, n int) []*Alien {
	aliens := make([]*Alien, n)

	rand.Seed(time.Now().UnixNano())

	for i := range aliens {
		city := w.cities[rand.Intn(len(w.cities))]
		aliens[i] = &Alien{i: i, Loc: city}
		city.Aliens = append(city.Aliens, aliens[i])
	}
	return aliens
}

type Game struct {
	World  *World
	Aliens []*Alien
}

func (g *Game) export(out io.Writer) error {
	for _, city := range g.World.cities {
		_, err := out.Write([]byte(city.Name))
		if err != nil {
			return err
		}

		_, err = out.Write([]byte(" "))
		if err != nil {
			return err
		}

		for dir, c := range city.Directions {
			if c == nil {
				continue
			}
			_, err = out.Write([]byte(dir))
			if err != nil {
				return err
			}

			_, err = out.Write([]byte("="))
			if err != nil {
				return err
			}

			_, err = out.Write([]byte(c.Name))
			if err != nil {
				return err
			}

			_, err = out.Write([]byte(" "))
			if err != nil {
				return err
			}
		}
		_, err = out.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func NewGame(w *World, aliens []*Alien) *Game {
	return &Game{World: w, Aliens: aliens}
}

func (g *Game) Start(moves int, out io.Writer) {
	g.World.removeDestroyedCities(out)

	for i := 0; i < moves; i++ {
		for _, alien := range g.Aliens {
			if alien.Loc.isDestroyed() {
				g.World.removeDestroyedCity(alien.Loc, out)
			} else {
				for len(g.World.cities) > 0 {
					moved := alien.move(randLoc())
					if moved || alien.dead {
						break
					}
				}
			}
		}
		g.removeDeadAliens()
	}
}

func (g *Game) removeDeadAliens() {
	if g == nil {
		return
	}
	aliens := make([]*Alien, 0)

	for i, alien := range g.Aliens {
		if !alien.dead {
			aliens = append(aliens, g.Aliens[i])
		} else {
			fmt.Printf("%v is dead\n", alien.i)
		}
	}
	g.Aliens = aliens
}

func randLoc() string {
	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(4)
	switch n {
	case 0:
		return east
	case 1:
		return west
	case 2:
		return north
	case 3:
		return south
	default:
		panic("can't happen")
	}
}

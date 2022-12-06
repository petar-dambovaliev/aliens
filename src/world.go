package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func contains(elems []string, v string) (bool, int) {
	for i, s := range elems {
		if v == s {
			return true, i
		}
	}
	return false, -1
}

func containsCity(elems []*City, v string) (bool, int) {
	for i, s := range elems {
		if v == s.Name {
			return true, i
		}
	}
	return false, -1
}

var directions = [...]string{east, north, south, west}

const east = "east"
const north = "north"
const south = "south"
const west = "west"

type City struct {
	Name       string
	Directions map[string]*City
	Aliens     []*Alien
}

func (c *City) print(out io.Writer) {
	out.Write([]byte(
		fmt.Sprintf(
			"%s has been destroyed by alien %v and alien %v!\n",
			c.Name,
			c.Aliens[0].i,
			c.Aliens[1].i),
	),
	)
}

func (c *City) isDestroyed() bool {
	if c != nil && len(c.Aliens) > 1 {
		for _, alien := range c.Aliens {
			alien.dead = true
		}
		return true
	}
	return false
}

type World struct {
	cities    []*City
	cityIndex map[string]int
}

func (w *World) removeDestroyedCities(out io.Writer) {
	if w == nil {
		return
	}
	cities := make([]*City, 0)

	for i := range w.cities {
		if !w.cities[i].isDestroyed() {
			cities = append(cities, w.cities[i])
		} else {
			w.cities[i].print(out)
		}
	}
	w.cities = cities
}

func (w *World) removeDestroyedCity(city *City, out io.Writer) {
	if w == nil {
		return
	}

	i, ok := w.cityIndex[city.Name]

	if !ok {
		return
	}
	city.print(out)
	w.cities = append(w.cities[:i], w.cities[i+1:]...)
}

func WorldFromReader(r io.Reader) (*World, error) {
	buf := bufio.NewReader(r)
	cities := make([]*City, 0)
	alienIndex := make(map[string]int)
	var line int

	for {
		str, err := buf.ReadString('\n')
		str = strings.TrimSpace(str)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if len(str) == 0 && err == io.EOF {
			break
		}

		strs := strings.Split(str, " ")

		if len(strs) == 0 {
			return nil, fmt.Errorf("invalid input missing city name on line %v\n", line)
		}
		name := strings.TrimSpace(strs[0])
		var city *City
		c, i := containsCity(cities, name)

		if c {
			city = cities[i]
		} else {
			city = &City{
				Name:       name,
				Directions: make(map[string]*City),
			}
			alienIndex[city.Name] = len(cities)
			cities = append(cities, city)
		}

		links_len := len(strs[1:])

		if links_len == 0 {
			continue
		}

		if links_len > 4 {
			return nil, fmt.Errorf("invalid city link directions. 4 are allowed directions are `east`, `north`, `south`, `west` on line %v\n", line)
		}

		for _, s := range strs[1:] {
			s = strings.TrimSpace(s)
			if len(s) < 3 {
				continue
			}
			link := strings.Split(s, "=")
			if len(link) != 2 {
				return nil, fmt.Errorf("invalid city link `%s`. expected format example `east=bar`", s)
			}

			dir := strings.TrimSpace(link[0])
			dir_name := strings.TrimSpace(link[1])

			if c, _ := contains(directions[:], dir); !c {
				return nil, fmt.Errorf("invalid city link direction `%s`. allowed directions are `east`, `north`, `south`, `west`", dir)
			}

			var link_city *City
			c, i := containsCity(cities, dir_name)

			if !c {
				link_city = &City{
					Name:       dir_name,
					Directions: make(map[string]*City),
				}
				alienIndex[city.Name] = len(cities)
				cities = append(cities, link_city)
			} else {
				link_city = cities[i]
			}

			err := allowDirection(city, link_city, dir)

			if err != nil {
				return nil, err
			}

			city.Directions[dir] = link_city
		}

		line += 1

		if err == io.EOF {
			break
		}
	}

	return &World{cities: cities}, nil
}

func oppositeDirection(dir string) string {
	switch dir {
	case east:
		return west
	case west:
		return east
	case north:
		return south
	case south:
		return north
	default:
		panic("invalid direction")
	}
}

func allowDirection(city *City, other *City, dir string) error {
	oppositeDir := oppositeDirection(dir)
	copyDirections := make(map[string]*City)

	for k, v := range city.Directions {
		copyDirections[k] = v
	}

	copyOppositeDirections := make(map[string]*City)

	for k, v := range other.Directions {
		copyOppositeDirections[k] = v
	}

	if copyDirections[dir] != nil && copyOppositeDirections[oppositeDir] != nil && copyDirections[dir].Name != copyOppositeDirections[oppositeDir].Name {
		return fmt.Errorf("nonsensical directions: %s points %s to %s", city.Name, dir, copyDirections[dir].Name)
	}

	delete(copyDirections, dir)
	delete(copyOppositeDirections, oppositeDir)

	for _, n := range copyDirections {
		if n.Name == city.Name {
			return fmt.Errorf("a city cannot have a point to itself")
		}
	}

	for _, n := range copyOppositeDirections {
		if n.Name == other.Name {
			return fmt.Errorf("cannot have more than one reference to another city")
		}
	}

	return nil
}

package main

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func TestWorldFromReaderUndeclaredLinks(t *testing.T) {
	input := `Foo north=Bar west=Baz south=Qu-ux`
	reader := strings.NewReader(input)
	_, err := WorldFromReader(reader)

	assert.Nil(t, err)
}

func TestWorldFromReaderDuplicateCities(t *testing.T) {
	input := `Foo north=Bar west=Baz south=Qu-ux
				Foo north=Bar west=Baz south=Qu-ux`
	reader := strings.NewReader(input)
	_, err := WorldFromReader(reader)

	assert.Nil(t, err)

}

func TestWorldFromReaderItWorks(t *testing.T) {
	input := `Foo north=Bar west=Baz south=Qu-ux
			    Bar
				Baz
				Qu-ux`
	reader := strings.NewReader(input)
	world, err := WorldFromReader(reader)

	if err != nil {
		t.Errorf("should not have errors: %+v\n", err)
	}

	sort.Slice(world.cities, func(i, j int) bool {
		return world.cities[i].Name < world.cities[j].Name
	})

	bar := &City{Name: "Bar", Directions: make(map[string]*City)}
	baz := &City{Name: "Baz", Directions: make(map[string]*City)}
	q := &City{Name: "Qu-ux", Directions: make(map[string]*City)}
	dir := map[string]*City{
		"north": bar, "west": baz, "south": q,
	}
	foo := &City{Name: "Foo", Directions: dir}

	assert.Equal(t, world.cities[2], foo)
	assert.Equal(t, world.cities[0], bar)
	assert.Equal(t, world.cities[1], baz)
	assert.Equal(t, world.cities[3], q)
}

func TestWorldFromReaderInadequatePositions(t *testing.T) {
	input := `Foo north=Bar
				Bar north=Foo`
	reader := strings.NewReader(input)
	_, err := WorldFromReader(reader)

	assert.Nil(t, err)
}

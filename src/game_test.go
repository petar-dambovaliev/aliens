package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGame_Start(t *testing.T) {
	input := `Foo north=Bar
				Bar north=Foo`
	reader := strings.NewReader(input)
	world, err := WorldFromReader(reader)
	assert.Nil(t, err)
	aliens := GenerateAliens(world, 2)
	assert.Equal(t, len(aliens), 2)

	game := NewGame(world, aliens)
	var msg strings.Builder
	game.Start(maxMoves, &msg)
	println(msg.String())

	var sb strings.Builder
	err = game.export(&sb)
	if err != nil {
		t.Errorf("%+v\n", err)
	}
}

package main

import (
	"log"

	"github.com/brantleyr/go-snake/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Set window size
	log.Println("Setting window size to", game.ScreenWidth, "x", game.ScreenHeight)
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)

	// Set window title
	log.Println("Setting window title to", game.GameTitle)
	ebiten.SetWindowTitle(game.GameTitle)

	// Set game state
	game.GameState = "intro"

	// Run the game
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}

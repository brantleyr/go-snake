package main

import (
	"log"

	"github.com/brantleyr/go-snake/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Set window size
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)

	// Set window title
	ebiten.SetWindowTitle(game.GameTitle)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Set game state
	game.GameState = "intro"
	game.GameStarted = false
	game.GamePaused = false

	// Run the game
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}

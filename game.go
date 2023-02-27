package main

import (
	//"bytes"
	//"image"
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	gameTitle = "Go Snake"
)

func init() {
	// Init stuff
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw base images

	// Draw new images

	// Get keys pressed
	keyStrs := []string{}
	for _, k := range g.keys {
		keyStrs = append(keyStrs, k.String())
	}

	// Debug key presses
	if len(keyStrs) > 0 {
		log.Println("Pressing keys", strings.Join(keyStrs, ", "))
	}
	ebitenutil.DebugPrint(screen, strings.Join(keyStrs, ", "))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Set window size
	log.Println("Setting window size to", screenWidth*2, "x", screenHeight*2)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)

	// Set window title
	log.Println("Setting window title to", gameTitle)
	ebiten.SetWindowTitle(gameTitle)

	// Run the game
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

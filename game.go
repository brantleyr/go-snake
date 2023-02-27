package main

import (
	_ "image/png"
	"image/color"
	"log"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidth  = 640
	screenHeight = 640
	dpi = 72
	baseFontSize = 24
	gameTitle = "Go Snake"
	startGameText = "Use arrow keys to guide Snake.\n    Press Enter to start."
)

var (
	tilesImage *ebiten.Image
	snakeHead *ebiten.Image
	baseFont font.Face
	gameActive bool
)

func init() {
	// Load background image
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("images/tiles.png")
	if err != nil {
		log.Fatal(err)
	}
	// Load snake head
	snakeHead, _, err = ebitenutil.NewImageFromFile("images/snake-head.png")
	if err != nil {
		log.Fatal(err)
	}
	// Load basic font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	baseFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    baseFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.DrawImage(tilesImage, nil)

	// Draw snake head
	screen.DrawImage(snakeHead, nil)

	// Draw start game text
	if ( gameActive == false ) {
		text.Draw(screen, startGameText, baseFont, (screenWidth/3)-50, (screenHeight/3)+90, color.White)
	}

	// Get keys pressed
	keyStrs := []string{}
	for _, k := range g.keys {
		keyStrs = append(keyStrs, k.String())
		if ( k.String() == "Enter" ){
			gameActive = true
		}
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
	log.Println("Setting window size to", screenWidth, "x", screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	// Set window title
	log.Println("Setting window title to", gameTitle)
	ebiten.SetWindowTitle(gameTitle)

	// Run the game
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

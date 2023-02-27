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

	// TODO: Dynamic screen widths/heights
	//		 Fonts, DPI, font settings and images will also need to scale appropriately
	screenWidth  = 640
	screenHeight = 640
	dpi = 72
	baseFontSize = 24
	gameTitle = "Go Snake"
	startGameText = "Use arrow keys to guide Snake.\n    Press Enter to start."
	bgImage = "images/tiles.png"
	snakeHead = "images/snake-head.png"
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
	tilesImage, _, err = ebitenutil.NewImageFromFile(bgImage)
	if err != nil {
		log.Fatal(err)
	}
	// Load snake head
	snakeHead, _, err = ebitenutil.NewImageFromFile(snakeHead)
	if err != nil {
		log.Fatal(err)
	}
	//TODO: Add "snake body" that can be extended
	//TODO: Add "snake tail" that will be the end

	// Load basic font
	// TODO: Investigate loading custom fonts or including our own
	//		 font with our source
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

		// TODO: Add some sort of way to detect the center of the screen
		// TODO: Add some sort of BG overlay so font is more easily readable
		text.Draw(screen, startGameText, baseFont, (screenWidth/3)-50, (screenHeight/3)+90, color.White)
	}

	// Get keys pressed
	keyStrs := []string{}
	for _, k := range g.keys {
		keyStrs = append(keyStrs, k.String())

		// TODO: Have a better way to detect game states
		// TODO: Pull out game states into its own package/library
		// 		 instead of just cramming it all in this one "Draw" function
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

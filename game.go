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
	DEBUG_MODE = true
	screenWidth  = 640
	screenHeight = 640
	dpi = 72
	baseFontSize = 24
	gameTitle = "Go Snake"
	startGameText = "Use arrow keys to guide Snake.\n    Press Enter to start."
	bgImageSrc = "images/tiles.png"
	snakeHeadImageSrc = "images/snake-head.png"
)

var (
	tilesImage *ebiten.Image
	snakeHead *ebiten.Image
	baseFont font.Face
	gameActive bool
	gameState string // intro, title, game, exit
)

func init() {
	// Load background image
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile(bgImageSrc)
	if err != nil {
		log.Fatal(err)
	}
	// Load snake head
	snakeHead, _, err = ebitenutil.NewImageFromFile(snakeHeadImageSrc)
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

func handleKeys(g* Game, screen *ebiten.Image) []string {
	keyStrs := []string{}
	for _, k := range g.keys {
		keyStrs = append(keyStrs, k.String())
		if ( k.String() == "Enter" ){
			gameActive = true
			gameState = "game"
		}
	}

	// Debug key presses
	if ( DEBUG_MODE == true ) {
		if len(keyStrs) > 0 {
			log.Println("Pressing keys", strings.Join(keyStrs, ", "))
		}
		ebitenutil.DebugPrint(screen, strings.Join(keyStrs, ", "))
	}

	return keyStrs
}

func doIntro(g* Game, screen *ebiten.Image) {

	// Show some opening credit images
	// Red Hat
	// Dr. Nick
	// Schneiders picture of choice
	// Golang picture
	// Etc

}

func doTitle(g* Game, screen *ebiten.Image) {
	// TODO: Make menu system

	// TODO: Make some snake game logo to display above menu

	// TODO: Add some sort of way to detect the center of the screen
	// TODO: Add some sort of BG overlay so font is more easily readable
	text.Draw(screen, startGameText, baseFont, (screenWidth/3)-50, (screenHeight/3)+90, color.White)

	// Handle keys
	handleKeys(g, screen)
	
}

func doGame(g* Game, screen *ebiten.Image) {
	// Draw background
	screen.DrawImage(tilesImage, nil)

	// Draw snake head
	screen.DrawImage(snakeHead, nil)

	// Handle keys
	handleKeys(g, screen)

}

func doExit(g* Game, screen *ebiten.Image) {
	// TODO: To trigger a game close / exit when chosen from the menu system
}

func handleGameState(g* Game, screen *ebiten.Image) {

	// TODO: Maybe make this cleaner
	if ( gameState == "intro" ) {
		doIntro(g, screen)
	}

	if ( gameState == "title" ) {
		doTitle(g, screen)
	}

	if ( gameState == "game" ){
		doGame(g, screen)
	}

	if ( gameState == "exit" ){
		doExit(g, screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	handleGameState(g, screen)
}

func main() {
	// Set window size
	log.Println("Setting window size to", screenWidth, "x", screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	// Set window title
	log.Println("Setting window title to", gameTitle)
	ebiten.SetWindowTitle(gameTitle)

	// Set game state
	gameState = "title"

	// Run the game
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

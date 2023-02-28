package game

import (
	"image/color"
	_ "image/png"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (

	// TODO: Dynamic screen widths/heights
	//		 Fonts, DPI, font settings and images will also need to scale appropriately
	DEBUG_MODE        = true
	ScreenWidth       = 640
	ScreenHeight      = 640
	dpi               = 72
	baseFontSize      = 24
	GameTitle         = "Go Snake"
	startGameText     = "Use arrow keys to guide Snake.\n    Press Enter to start."
	bgImageSrc        = "images/tiles.png"
	snakeHeadImageSrc = "images/snake-head.png"
	drNickImageSrc    = "images/dr-nick.png"
	schImageSrc       = "images/schneider.png"
	rhImageSrc        = "images/red-hat.png"
	ebImageSrc        = "images/ebitengine.png"
	goImageSrc        = "images/golang.png"
)

var (
	tilesImage *ebiten.Image
	snakeHead  *ebiten.Image
	drNick     *ebiten.Image
	schImage   *ebiten.Image
	rhImage    *ebiten.Image
	ebImage    *ebiten.Image
	goImage    *ebiten.Image
	baseFont   font.Face
	gameActive bool
	GameState  string // intro, title, game, exit
)

func init() {
	var err error
	// Load intro images
	drNick, _, err = ebitenutil.NewImageFromFile(drNickImageSrc)
	if err != nil {
		log.Fatal(err)
	}
	schImage, _, err = ebitenutil.NewImageFromFile(schImageSrc)
	if err != nil {
		log.Fatal(err)
	}
	rhImage, _, err = ebitenutil.NewImageFromFile(rhImageSrc)
	if err != nil {
		log.Fatal(err)
	}
	ebImage, _, err = ebitenutil.NewImageFromFile(ebImageSrc)
	if err != nil {
		log.Fatal(err)
	}
	goImage, _, err = ebitenutil.NewImageFromFile(goImageSrc)
	if err != nil {
		log.Fatal(err)
	}
	// Load background image
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
	
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if GameState == "intro" {
			log.Println("Enter key pressed in intro")
			GameState = "title"
		} else if GameState == "title" {
			log.Println("Enter key pressed in title")
			GameState = "game"
		}
	}

	return nil
}

func doIntro(g *Game, screen *ebiten.Image) {

	// TODO: Clean up these images and spacing
	//	     Try to find a way to dynamically place them instead of hard-coding

	// Show some opening credit images
	// Red Hat
	// Dr. Nick
	// Schneiders picture of choice
	// Golang picture
	// Etc

	// Show images
	// Dr. Nick
	nickOp := &ebiten.DrawImageOptions{}
	nickOp.GeoM.Scale(.25, .25)
	nickOp.GeoM.Translate(75, 125)
	screen.DrawImage(drNick, nickOp)

	// Schneider
	schOp := &ebiten.DrawImageOptions{}
	schOp.GeoM.Scale(.40, .40)
	schOp.GeoM.Translate(250, 125)
	screen.DrawImage(schImage, schOp)

	// Red Hat
	rhOp := &ebiten.DrawImageOptions{}
	rhOp.GeoM.Scale(.65, .65)
	rhOp.GeoM.Translate(425, 125)
	screen.DrawImage(rhImage, rhOp)

	// Ebitengine
	ebOp := &ebiten.DrawImageOptions{}
	ebOp.GeoM.Scale(.65, .65)
	ebOp.GeoM.Translate(125, 300)
	screen.DrawImage(ebImage, ebOp)

	// Golang
	goOp := &ebiten.DrawImageOptions{}
	goOp.GeoM.Scale(.25, .25)
	goOp.GeoM.Translate(375, 300)
	screen.DrawImage(goImage, goOp)

}

func doTitle(g *Game, screen *ebiten.Image) {
	// TODO: Make menu system

	// TODO: Make some snake game logo to display above menu

	// TODO: Add some sort of way to detect the center of the screen
	// TODO: Add some sort of BG overlay so font is more easily readable
	text.Draw(screen, startGameText, baseFont, (ScreenWidth/3)-50, (ScreenHeight/3)+90, color.White)

}

func doGame(g *Game, screen *ebiten.Image) {
	// Draw background
	screen.DrawImage(tilesImage, nil)

	// Draw snake head
	screen.DrawImage(snakeHead, nil)

}

func doExit(g *Game, screen *ebiten.Image) {
	// TODO: To trigger a game close / exit when chosen from the menu system
}

func handleGameState(g *Game, screen *ebiten.Image) {

	// TODO: Maybe make this cleaner
	if GameState == "intro" {
		doIntro(g, screen)
	}

	if GameState == "title" {
		doTitle(g, screen)
	}

	if GameState == "game" {
		doGame(g, screen)
	}

	if GameState == "exit" {
		doExit(g, screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	handleGameState(g, screen)
}

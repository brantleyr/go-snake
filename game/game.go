package game

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"os"

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
	titleFontSize     = 48
	GameTitle         = "Go Snake"
	startGameText     = "Use arrow keys to guide Snake.\n    Press Enter to start."
	bgImageSrc        = "images/tiles.png"
	snakeHeadImageSrc = "images/snake-head.png"
	drNickImageSrc    = "images/dr-nick.png"
	schImageSrc       = "images/schneider.png"
	rhImageSrc        = "images/red-hat.png"
	ebImageSrc        = "images/ebitengine.png"
	goImageSrc        = "images/golang.png"
	snakeLogoImageSrc = "images/snake-logo.png"
	baseSpeed         = 1
	cellSizeWidth     = 76
	cellSizeHeight    = 76
	gridLineSize      = 3.75
)

type snake struct {
	snakeHead *ebiten.Image
	xPos      float64
	yPos      float64
	speed     float64
}

var (
	tilesImage  *ebiten.Image
	snakeHead   *ebiten.Image
	snakePlayer snake
	drNick      *ebiten.Image
	schImage    *ebiten.Image
	rhImage     *ebiten.Image
	ebImage     *ebiten.Image
	goImage     *ebiten.Image
	snakeLogo   *ebiten.Image
	baseFont    font.Face
	titleFont   font.Face
	gameActive  bool
	GameState   string // intro, title, game, exit
	menuItem    string
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
	// Load snake logo
	snakeLogo, _, err = ebitenutil.NewImageFromFile(snakeLogoImageSrc)
	if err != nil {
		log.Fatal(err)
	}

	// Make snake
	snakePlayer = snake{snakeHead, gridLineSize, gridLineSize, baseSpeed}

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
	titleFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	menuItem = "new_game"
}

type Game struct {
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if GameState == "intro" {
			log.Println("Enter key pressed in intro")
			GameState = "title"
		} else if GameState == "title" {
			if menuItem == "new_game" {
				GameState = "game"
			} else if menuItem == "exit" {
				// Loop to the down
				GameState = "exit"
			}
		}
	} else if GameState == "title" {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			log.Println("Arrow down key pressed in title")

			if menuItem == "new_game" {
				// They just moved down
				menuItem = "exit"
			} else if menuItem == "exit" {
				// Loop to the top
				menuItem = "new_game"
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			log.Println("Arrow up key pressed in title")

			if menuItem == "new_game" {
				// They just moved up
				menuItem = "exit"
			} else if menuItem == "exit" {
				// Loop to the down
				menuItem = "new_game"
			}
		}
	}

	if GameState == "game" {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			newYpos := float64(cellSizeHeight) + float64(gridLineSize)
			if (snakePlayer.yPos - newYpos) > 0 {
				snakePlayer.yPos -= newYpos
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			newYpos := float64(cellSizeHeight) + float64(gridLineSize)
			if (snakePlayer.yPos + newYpos) < ScreenHeight {
				snakePlayer.yPos += newYpos
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			newXpos := float64(cellSizeWidth) + float64(gridLineSize)
			if (snakePlayer.xPos - newXpos) > 0 {
				snakePlayer.xPos -= newXpos
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			newXpos := float64(cellSizeWidth) + float64(gridLineSize)
			if (snakePlayer.xPos + newXpos) < ScreenHeight {
				snakePlayer.xPos += newXpos
			}
		}
	}

	return nil
}

func (g *Game) newGame() {
	GameState = "game"
}

func (g *Game) exitGame() {
	GameState = "exit"
}

func ParseHexColor(s string) (c color.RGBA) {
	c.A = 0xff
	switch len(s) {
	case 7:
		fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	}
	return
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

func drawTitle(screen *ebiten.Image) {
	snake := &ebiten.DrawImageOptions{}
	snake.GeoM.Scale(.25, .25)
	snake.GeoM.Translate((ScreenWidth/2)-130, 20)
	screen.DrawImage(snakeLogo, snake)
	text.Draw(screen, "nake", titleFont, (ScreenWidth/2)-45, 135, color.White)

	if menuItem == "new_game" {
		text.Draw(screen, "> New Game", titleFont, (ScreenWidth/3)-56, (ScreenHeight/3)+90, color.White)
		text.Draw(screen, "Exit", titleFont, (ScreenWidth/3)-10, (ScreenHeight/3)+140, ParseHexColor("#8c8c8c"))
	} else if menuItem == "exit" {
		text.Draw(screen, "New Game", titleFont, (ScreenWidth/3)-10, (ScreenHeight/3)+90, ParseHexColor("#8c8c8c"))
		text.Draw(screen, "> Exit", titleFont, (ScreenWidth/3)-56, (ScreenHeight/3)+140, color.White)
	}
}

func doTitle(g *Game, screen *ebiten.Image) {
	// TODO: Make some snake game logo to display above menu

	// TODO: Add some sort of way to detect the center of the screen
	// TODO: Add some sort of BG overlay so font is more easily readable

	drawTitle(screen)
	// Use this once we actually have a working menu
	// text.Draw(screen, startGameText, baseFont, (ScreenWidth/3)-50, (ScreenHeight/3)+90, color.White)

}

func doGame(g *Game, screen *ebiten.Image) {
	// Draw background
	screen.DrawImage(tilesImage, nil)

	// Draw snake head
	hOp := &ebiten.DrawImageOptions{}
	hOp.GeoM.Translate(snakePlayer.xPos, snakePlayer.yPos)
	screen.DrawImage(snakeHead, hOp)

	// Draw middle parts

	// Draw tail

}

func doExit(g *Game) {
	g.exitGame()
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
		doExit(g)
		os.Exit(0)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	handleGameState(g, screen)
}

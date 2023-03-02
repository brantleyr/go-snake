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
	dpi               = 72
	baseFontSize      = 36
	titleFontSize     = 72
	GameTitle         = "Go Snake"
	startGameText     = "Use arrow keys to guide Snake.\n    Press Enter to start."
	drNickImageSrc    = "images/dr-nick.png"
	schImageSrc       = "images/schneider.png"
	rhImageSrc        = "images/red-hat.png"
	ebImageSrc        = "images/ebitengine.png"
	goImageSrc        = "images/golang.png"
	snakeLogoImageSrc = "images/snake-logo.png"
	gridHeight        = 20
	gridWidth         = 20
	gridSolidColor    = "#002200"
	pieceColor        = "#00ff00"
	clockSpeed        = 30
)

type pathPair struct {
	xPos, yPos int
}

type snakeBody struct {
	xPos    int
	yPos    int
	segment int
}

type snake struct {
	snakeBody []snakeBody
	xPos      int
	yPos      int
	direction string // up, down, left, right
}

var (
	snakePlayer    snake
	drNick         *ebiten.Image
	schImage       *ebiten.Image
	rhImage        *ebiten.Image
	ebImage        *ebiten.Image
	goImage        *ebiten.Image
	snakeLogo      *ebiten.Image
	baseFont       font.Face
	titleFont      font.Face
	GameStarted    bool
	GamePaused     bool
	GameState      string // intro, title, game, exit
	menuItem       string
	ScreenWidth    = 1200
	ScreenHeight   = 1024
	gridCellHeight int
	gridCellWidth  int
	snakePath      []pathPair
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
	// Load snake logo
	snakeLogo, _, err = ebitenutil.NewImageFromFile(snakeLogoImageSrc)
	if err != nil {
		log.Fatal(err)
	}

	// Make snake
	var initialBodyPieces []snakeBody
	initialBodyPieces = []snakeBody{
		snakeBody{0, 0, 2},
		snakeBody{0, 1, 1},
		snakeBody{0, 2, 0},
	}
	snakePlayer = snake{initialBodyPieces, 0, 3, "down"} // Head

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
	clockSpeedCount int
}

func (g *Game) Update() error {
	// Handle "intro" game state key events
	if GameState == "intro" {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			GameState = "title"
		}

		// Handle "title" game state key events
	} else if GameState == "title" {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if menuItem == "new_game" {
				GameState = "game"
			} else if menuItem == "exit" {
				// Loop to the down
				GameState = "exit"
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			if menuItem == "new_game" {
				// They just moved down
				menuItem = "exit"
			} else if menuItem == "exit" {
				// Loop to the top
				menuItem = "new_game"
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			if menuItem == "new_game" {
				// They just moved up
				menuItem = "exit"
			} else if menuItem == "exit" {
				// Loop to the down
				menuItem = "new_game"
			}
		}

		// Handle "game" game state key events
	} else if GameState == "game" {
		if GameStarted == true {
			if GamePaused == false {
				if snakePlayer.direction != "up" && snakePlayer.direction != "down" {
					if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
						snakePlayer.direction = "up"
					}
					if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
						snakePlayer.direction = "down"
					}
				}
				if snakePlayer.direction != "left" && snakePlayer.direction != "right" {
					if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
						snakePlayer.direction = "left"
					}
					if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
						snakePlayer.direction = "right"
					}
				}
				if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
					GamePaused = true
				}
			} else {
				if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
					GamePaused = false
				}
			}
		} else {
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				GameStarted = true
			}
		}
	}

	return nil
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

	// Show images
	// Dr. Nick
	nickOp := &ebiten.DrawImageOptions{}
	nickOp.GeoM.Scale(.5, .5)
	nickOp.GeoM.Translate(float64(ScreenWidth)*0.1171875, float64(ScreenHeight)*0.1953125)
	screen.DrawImage(drNick, nickOp)

	// Schneider
	schOp := &ebiten.DrawImageOptions{}
	schOp.GeoM.Scale(.80, .80)
	schOp.GeoM.Translate(float64(ScreenWidth)*0.390625, float64(ScreenHeight)*0.1953125)
	screen.DrawImage(schImage, schOp)

	// Red Hat
	rhOp := &ebiten.DrawImageOptions{}
	rhOp.GeoM.Scale(1.3, 1.3)
	rhOp.GeoM.Translate(float64(ScreenWidth)*0.6640625, float64(ScreenHeight)*0.1953125)
	screen.DrawImage(rhImage, rhOp)

	// Ebitengine
	ebOp := &ebiten.DrawImageOptions{}
	ebOp.GeoM.Scale(1.3, 1.3)
	ebOp.GeoM.Translate(float64(ScreenWidth)*0.1953125, float64(ScreenHeight)*0.46875)
	screen.DrawImage(ebImage, ebOp)

	// Golang
	goOp := &ebiten.DrawImageOptions{}
	goOp.GeoM.Scale(.50, .50)
	goOp.GeoM.Translate(float64(ScreenWidth)*0.5859375, float64(ScreenHeight)*0.46875)
	screen.DrawImage(goImage, goOp)

}

func drawTitle(screen *ebiten.Image) {
	// Logo and text on top
	snake := &ebiten.DrawImageOptions{}
	snake.GeoM.Scale(.50, .50)
	snake.GeoM.Translate(float64((ScreenWidth/2))-(float64(ScreenWidth)*0.203125), float64(ScreenHeight)*0.08125)
	screen.DrawImage(snakeLogo, snake)

	if menuItem == "new_game" {
		text.Draw(screen, "> New Game", titleFont, (ScreenWidth/3)-105, (ScreenHeight/3)+90, color.White)
		text.Draw(screen, "Exit", titleFont, (ScreenWidth/3)-35, (ScreenHeight/3)+170, ParseHexColor("#8c8c8c"))
	} else if menuItem == "exit" {
		text.Draw(screen, "New Game", titleFont, (ScreenWidth/3)-35, (ScreenHeight/3)+90, ParseHexColor("#8c8c8c"))
		text.Draw(screen, "> Exit", titleFont, (ScreenWidth/3)-105, (ScreenHeight/3)+170, color.White)
	}
}

func doTitle(g *Game, screen *ebiten.Image) {

	// TODO: Add some sort of way to detect the center of the screen
	// TODO: Add some sort of BG overlay so font is more easily readable

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	drawTitle(screen)

}

func getGridColor(ix int, iy int) color.Color {
	var theColor color.Color
	if ix%2 == 0 {
		if iy%2 == 0 {
			theColor = ParseHexColor(gridSolidColor)
		} else {
			theColor = color.Black
		}
	} else {
		if iy%2 == 0 {
			theColor = color.Black
		} else {
			theColor = ParseHexColor(gridSolidColor)
		}
	}
	return theColor
}

func buildGrid(screen *ebiten.Image) {

	gridCellWidth = ScreenWidth / gridWidth
	gridCellHeight = ScreenHeight / gridHeight
	for ix := 0; ix < gridCellWidth; ix++ {
		for iy := 0; iy < gridCellHeight; iy++ {
			ebitenutil.DrawRect(screen, float64(ix*gridCellWidth), float64(iy*gridCellHeight), float64(gridCellWidth), float64(gridCellHeight), getGridColor(ix, iy))
		}
	}

}

func drawPiece(screen *ebiten.Image, ix int, iy int, theColor color.Color, shapeType string) {
	if shapeType == "rect" {
		ebitenutil.DrawRect(screen, float64(ix*gridCellWidth), float64(iy*gridCellHeight), float64(gridCellWidth), float64(gridCellHeight), theColor)
	}
	if shapeType == "smallcircle" {
		radius := float64((float64(gridCellWidth/5) + float64(gridCellHeight/5)) / 2)
		ebitenutil.DrawCircle(screen, float64(ix*gridCellWidth)+(radius*2.5), float64(iy*gridCellHeight)+(radius*2.5), radius, theColor)
	}
	if shapeType == "circle" {
		radius := float64((float64(gridCellWidth/3) + float64(gridCellHeight/3)) / 2)
		ebitenutil.DrawCircle(screen, float64(ix*gridCellWidth)+(radius*1.5), float64(iy*gridCellHeight)+(radius*1.5), radius, theColor)
	}
	if shapeType == "triangle" {
		//TODO: For tail, for now its a smaller cicle
	}

}

func doGame(g *Game, screen *ebiten.Image) {
	// Draw background
	buildGrid(screen)

	// Handle game started vs paused
	if GameStarted == true && GamePaused == false {
		var moveCounter int

		if g.clockSpeedCount == 0 {
			moveCounter = 1
		} else {
			moveCounter = 0
		}

		// Update snake
		if snakePlayer.yPos > 0 && snakePlayer.direction == "up" {
			snakePlayer.yPos -= moveCounter
		}
		if snakePlayer.yPos < 19 && snakePlayer.direction == "down" {
			snakePlayer.yPos += moveCounter
		}
		if snakePlayer.xPos > 0 && snakePlayer.direction == "left" {
			snakePlayer.xPos -= moveCounter
		}
		if snakePlayer.xPos < 19 && snakePlayer.direction == "right" {
			snakePlayer.xPos += moveCounter
		}
	}

	// Handle game started vs paused
	if GameStarted == true && GamePaused == true {
		text.Draw(screen, "Game Paused. Escape to resume.", baseFont, (ScreenWidth/3)-56, (ScreenHeight/3)+90, color.White)
	} else if GameStarted == false {
		// Do not update snake
		// Show start text
		text.Draw(screen, "Arrow keys move snake\nEnter starts game", baseFont, (ScreenWidth/3)+20, (ScreenHeight/3)+180, color.White)
	}

	// Draw head
	drawPiece(screen, snakePlayer.xPos, snakePlayer.yPos, ParseHexColor(pieceColor), "rect")

	// Draw pieces
	for idx, snakePiece := range snakePlayer.snakeBody {
		// Update snake pieces
		// TODO: More directions than just right, and turn with the snake
		// The pieces start at a negative index which is cutting off the snake as well
		snakePiece.yPos = snakePlayer.yPos
		snakePiece.xPos = snakePlayer.xPos - idx
		// snakePiece.segment = len(snakePlayer.snakeBody)

		log.Println("Snake Piece:", snakePiece.xPos, snakePiece.yPos, idx)
		log.Println("Player:     ", snakePlayer.xPos, snakePlayer.yPos)
		log.Println("Direction:  ", snakePlayer.direction)
		log.Println("Segment:     ", snakePiece.segment)
		log.Println("Body Length:", len(snakePlayer.snakeBody))

		// TODO: Get the tail drawing again
		if snakePiece.segment == (len(snakePlayer.snakeBody) - len(snakePath)) {
			// Tail
			drawPiece(screen, snakePiece.xPos, snakePiece.yPos, ParseHexColor(pieceColor), "smallcircle")
		} else {
			// Other pieces
			drawPiece(screen, snakePiece.xPos, snakePiece.yPos, ParseHexColor(pieceColor), "circle")
		}
	}

	// Update clock speed count
	// TODO: Is there some way to control game fps or clock speed or ticks in ebitengine?
	g.clockSpeedCount += 1
	if g.clockSpeedCount > clockSpeed {
		if len(snakePath) > len(snakePlayer.snakeBody) {
			snakePath = append([]pathPair{pathPair{snakePlayer.xPos, snakePlayer.yPos}}, snakePath[len(snakePath)-len(snakePlayer.snakeBody):]...)
		} else {
			snakePath = append([]pathPair{pathPair{snakePlayer.xPos, snakePlayer.yPos}}, snakePath...)

		}

		g.clockSpeedCount = 0
	}
}

func doExit(g *Game) {
	GameState = "exit"
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
		os.Exit(0)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ScreenWidth = outsideWidth
	ScreenHeight = outsideHeight

	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	handleGameState(g, screen)
}

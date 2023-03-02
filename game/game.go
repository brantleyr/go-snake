package game

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
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
	gridWidth         = 25
	gridSolidColor    = "#002200"
	pieceColor        = "#00ff00"
	nomColor          = "#ff0000"
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
	GameStarted    = false
	GamePaused     = false
	GameOver       = false
	GameState      = "title" // intro, title, game, exit
	menuItem       string
	ScreenWidth    = 1024
	ScreenHeight   = 768
	gridCellHeight int
	gridCellWidth  int
	snakePath      []pathPair
	nomActive      = false
	currentNom     pathPair
	clockSpeed     = 20
)

func setupInitialSnake() {
	var initialBodyPieces = []snakeBody{
		snakeBody{0, 0, 2},
		snakeBody{0, 1, 1},
		snakeBody{0, 2, 0},
	}
	snakePlayer = snake{initialBodyPieces, 0, 3, "down"}

	// Initial Path
	snakePath = []pathPair{
		pathPair{0, 2},
		pathPair{0, 1},
		pathPair{0, 0},
	}
}

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
	setupInitialSnake()

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
		if GameStarted && !GameOver {
			if !GamePaused {
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
		if GameOver {
			if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				GameStarted = true
				GameOver = false
				setupInitialSnake()
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
	nickOp.GeoM.Scale(.4, .4)
	nickOp.GeoM.Translate(float64(ScreenWidth)*0.0571875, float64(ScreenHeight)*0.1953125)
	screen.DrawImage(drNick, nickOp)

	// Schneider
	schOp := &ebiten.DrawImageOptions{}
	schOp.GeoM.Scale(.70, .70)
	schOp.GeoM.Translate(float64(ScreenWidth)*0.370625, float64(ScreenHeight)*0.1853125)
	screen.DrawImage(schImage, schOp)

	// Red Hat
	rhOp := &ebiten.DrawImageOptions{}
	rhOp.GeoM.Scale(1.2, 1.2)
	rhOp.GeoM.Translate(float64(ScreenWidth)*0.6940625, float64(ScreenHeight)*0.1653125)
	screen.DrawImage(rhImage, rhOp)

	// Ebitengine
	ebOp := &ebiten.DrawImageOptions{}
	ebOp.GeoM.Scale(1.2, 1.2)
	ebOp.GeoM.Translate(float64(ScreenWidth)*0.0571875, float64(ScreenHeight)*0.56875)
	screen.DrawImage(ebImage, ebOp)

	// Golang
	goOp := &ebiten.DrawImageOptions{}
	goOp.GeoM.Scale(.40, .40)
	goOp.GeoM.Translate(float64(ScreenWidth)*0.5859375, float64(ScreenHeight)*0.46875)
	screen.DrawImage(goImage, goOp)

}

func drawTitle(screen *ebiten.Image) {
	// Logo and text on top
	snake := &ebiten.DrawImageOptions{}
	snake.GeoM.Scale(.50, .50)
	snake.GeoM.Translate(float64((ScreenWidth/2))-(float64(ScreenWidth)*0.203125), float64(ScreenHeight)*0.06125)
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

	for ix := 0; ix < gridWidth; ix++ {
		for iy := 0; iy < gridHeight; iy++ {
			//log.Println("grid cell",ix,",",iy)
			ebitenutil.DrawRect(screen, float64(ix*gridCellWidth), float64(iy*gridCellHeight), float64(gridCellWidth), float64(gridCellHeight), getGridColor(ix, iy))
		}
	}

}

func drawGridPiece(screen *ebiten.Image, ix int, iy int, theColor color.Color, shapeType string) {
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

func doNoms(g *Game, screen *ebiten.Image) {

	// Only generate a new nom if there isn't currently one
	if !nomActive {
		notValidNom := true
		var randX, randY int

		// Generate random x,y pairs until it is not found in the existing snake path
		for notValidNom {
			rand.Seed(time.Now().UnixNano())
			randX = rand.Intn(gridWidth - 1)
			randY = rand.Intn(gridHeight - 1)
			notValidNom = false
			for _, pathPair := range snakePath {
				if pathPair.xPos == randX && pathPair.yPos == randY {
					notValidNom = true
					break
				}
			}
		}
		// Set the new nom and draw it on the screen
		nomActive = true
		currentNom = pathPair{randX, randY}
		drawGridPiece(screen, randX, randY, ParseHexColor(nomColor), "smallcircle")
	}

	// If theres already a nom, check to see if it intersects with the snake head or draw it
	if nomActive {
		if snakePlayer.xPos == currentNom.xPos && snakePlayer.yPos == currentNom.yPos {
			nomActive = false

			// Add new snake piece
			newSnakeBodyPiece := snakeBody{snakePlayer.xPos, snakePlayer.yPos, len(snakePlayer.snakeBody)}
			snakePlayer.snakeBody = append([]snakeBody{newSnakeBodyPiece}, snakePlayer.snakeBody...)

			// TODO: Make clock speed faster when you acquire more noms
			// clockSpeed -= someNumber and floor() it so it only gets so fast

		} else {
			drawGridPiece(screen, currentNom.xPos, currentNom.yPos, ParseHexColor(nomColor), "smallcircle")
		}
	}
}

func doGame(g *Game, screen *ebiten.Image) {
	// Draw background
	buildGrid(screen)

	// Handle game started vs paused
	if GameStarted && !GamePaused && !GameOver {
		var moveCounter int
		if g.clockSpeedCount == 0 {
			moveCounter = 1
		} else {
			moveCounter = 0
		}

		// Update snake
		if snakePlayer.direction == "up" {
			snakePlayer.yPos -= moveCounter
		}
		if snakePlayer.direction == "down" {
			snakePlayer.yPos += moveCounter
		}
		if snakePlayer.direction == "left" {
			snakePlayer.xPos -= moveCounter
		}
		if snakePlayer.direction == "right" {
			snakePlayer.xPos += moveCounter
		}
	}

	// Handle game started vs paused
	if GameStarted && GamePaused {
		text.Draw(screen, "Game Paused. Escape to resume.", baseFont, (ScreenWidth/3)-56, (ScreenHeight/3)+90, color.White)
	} else if !GameStarted && !GameOver {
		// Do not update snake
		// Show start text
		text.Draw(screen, "Arrow keys move snake\nEnter starts game", baseFont, (ScreenWidth/3)+20, (ScreenHeight/3)+180, color.White)
	}

	// Draw head
	drawGridPiece(screen, snakePlayer.xPos, snakePlayer.yPos, ParseHexColor(pieceColor), "rect")

	// Draw pieces
	for idx, snakePiece := range snakePlayer.snakeBody {
		snakePlayer.snakeBody[idx].xPos = snakePath[snakePiece.segment].xPos
		snakePlayer.snakeBody[idx].yPos = snakePath[snakePiece.segment].yPos
		if snakePiece.segment == (len(snakePlayer.snakeBody) - 1) {
			// Tail
			drawGridPiece(screen, snakePiece.xPos, snakePiece.yPos, ParseHexColor(pieceColor), "smallcircle")
		} else {
			// Other pieces
			drawGridPiece(screen, snakePiece.xPos, snakePiece.yPos, ParseHexColor(pieceColor), "circle")
		}
	}

	// Draw noms
	if GameStarted {
		doNoms(g, screen)
	}

	// Update clock speed count
	// TODO: Is there some way to control game fps or clock speed or ticks in ebitengine?

	g.clockSpeedCount += 1
	if g.clockSpeedCount > clockSpeed {
		if GameStarted && !GamePaused {
			snakePath = append([]pathPair{pathPair{snakePlayer.xPos, snakePlayer.yPos}}, snakePath[0:len(snakePlayer.snakeBody)]...)
		}
		g.clockSpeedCount = 0
	}

	// Check collision paths to end game
	for idx, snakePathPair := range snakePath {
		// Check if the head collided with the body
		if idx != 0 {
			// Did the snake collide with itself?
			if snakePlayer.xPos == snakePathPair.xPos && snakePlayer.yPos == snakePathPair.yPos {
				GameStarted = false
				GameOver = true
			}
		}
		// Check if the head collided with a wall
		if idx == 0 {
			if snakePlayer.xPos >= gridWidth || snakePlayer.xPos < 0 || snakePlayer.yPos >= gridHeight || snakePlayer.yPos < 0 {
				GameStarted = false
				GameOver = true
			}
		}
	}

	// Show Game Over
	if GameOver {
		text.Draw(screen, "Womp womp. Game over.\nPress Enter for New Game", baseFont, (ScreenWidth/2)-200, (ScreenHeight / 2), color.White)
	}

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

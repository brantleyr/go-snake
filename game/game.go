package game

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	DEBUG_MODE        = true
	dpi               = 72
	baseFontSize      = 36
	titleFontSize     = 72
	scoreFontSize     = 24
	timerFontSize     = 24
	GameTitle         = "Go Snake"
	startGameText     = "Use arrow keys to guide Snake.\n    Press Enter to start."
	drNickImageSrc    = "images/dr-nick.png"
	schImageSrc       = "images/schneider.png"
	rhImageSrc        = "images/red-hat.png"
	ebImageSrc        = "images/ebitengine.png"
	goImageSrc        = "images/golang.png"
	snakeLogoImageSrc = "images/snake-logo.png"
	globBgImageSrc    = "images/green-bg.png"
	appleImageSrc     = "images/apple.png"
	greenGridImageSrc = "images/green-grid.png"
	gridHeight        = 20
	gridWidth         = 25
	gridSolidColor    = "#002200"
	gridAltColor      = "#000000"
	gridCellOpacity   = 0xaf
	gridBorderColor   = "#005500"
	gridBorderSize    = 3
	pieceColor        = "#00ff00"
	nomColor          = "#ff0000"
	borderTop         = 50
	borderBottom      = 10
	borderLeft        = 10
	borderRight       = 10
	sampleRate        = 22050
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
	snakePlayer        snake
	drNick             *ebiten.Image
	schImage           *ebiten.Image
	rhImage            *ebiten.Image
	ebImage            *ebiten.Image
	goImage            *ebiten.Image
	snakeLogo          *ebiten.Image
	globBg             *ebiten.Image
	apple              *ebiten.Image
	greenGrid          *ebiten.Image
	baseFont           font.Face
	titleFont          font.Face
	scoreFont          font.Face
	timerFont          font.Face
	GameStarted        = false
	GamePaused         = false
	GameOver           = false
	GameState          = "title" // intro, title, game, exit
	menuItem           string
	ScreenWidth        = 1024
	ScreenHeight       = 768
	gridCellHeight     int
	gridCellWidth      int
	snakePath          []pathPair
	nomActive          = false
	currentNom         pathPair
	clockSpeed         = 20
	currScore          = 0
	globBgRot          = 0.75
	zoomingBg          = true
	introOpacity       = 0.0
	fadingOutIntro     = false
	timeElapsed        = 0
	timerDone          = make(chan bool)
	timerTicker        = time.NewTicker(1 * time.Second)
	appleScale         = .1
	zoomingApple       = true
	emptyImage         = ebiten.NewImage(3, 3)
	emptySubImage      = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	gameOverSnd        *audio.Player
	gameOverFile       fs.File
	GameJustEnded      = false
	GameOverSndPlaying = true
)

func setupInitialSnake() {
	var initialBodyPieces = []snakeBody{
		{0, 0, 2},
		{0, 1, 1},
		{0, 2, 0},
	}
	snakePlayer = snake{initialBodyPieces, 0, 3, "down"}

	// Initial Path
	snakePath = []pathPair{
		{0, 2},
		{0, 1},
		{0, 0},
	}
}

func openFile(path string) fs.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func decodeMP3(ctx *audio.Context, src io.ReadSeeker) *audio.Player {
	s, err := mp3.Decode(ctx, src)
	if err != nil {
		log.Fatal(err)
	}
	player, err := audio.NewPlayer(ctx, s)
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func init() {
	var err error

	// Fill the subimage
	// Used for DrawLine
	emptyImage.Fill(color.White)

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

	// Load global bg
	globBg, _, err = ebitenutil.NewImageFromFile(globBgImageSrc)
	if err != nil {
		log.Fatal(err)
	}

	// Load snake logo
	snakeLogo, _, err = ebitenutil.NewImageFromFile(snakeLogoImageSrc)
	if err != nil {
		log.Fatal(err)
	}

	// Load apple image
	apple, _, err = ebitenutil.NewImageFromFile(appleImageSrc)
	if err != nil {
		log.Fatal(err)
	}

	// Load green grid tile
	greenGrid, _, err = ebitenutil.NewImageFromFile(greenGridImageSrc)
	if err != nil {
		log.Fatal(err)
	}

	// Make snake
	setupInitialSnake()

	// Load basic font
	externalFont, err := os.ReadFile("fonts/JungleAdventurer.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(externalFont)
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
	scoreFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    scoreFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	timerFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    timerFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	menuItem = "new_game"

	// Initialize sounds
	ctx := audio.NewContext(sampleRate)
	gameOverFile = openFile("sounds/game-over.mp3")
	gameOverSnd = decodeMP3(ctx, gameOverFile.(io.ReadSeeker))
}

type Game struct {
	clockSpeedCount int
}

func (g *Game) Update() error {
	// Handle "intro" game state key events
	if GameState == "intro" {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			fadingOutIntro = true
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
				// Loop to the bottom
				menuItem = "new_game"
			}
		}

		// Handle "game" game state key events
	} else if GameState == "game" {
		if GameStarted && !GameOver {
			if !GamePaused {
				if snakePlayer.direction != "up" && snakePlayer.direction != "down" {
					if inpututil.IsKeyJustPressed(ebiten.KeyUp) ||
						inpututil.IsKeyJustPressed(ebiten.KeyW) {
						snakePlayer.direction = "up"
					}
					if inpututil.IsKeyJustPressed(ebiten.KeyDown) ||
						inpututil.IsKeyJustPressed(ebiten.KeyS) {
						snakePlayer.direction = "down"
					}
				}
				if snakePlayer.direction != "left" && snakePlayer.direction != "right" {
					if inpututil.IsKeyJustPressed(ebiten.KeyLeft) ||
						inpututil.IsKeyJustPressed(ebiten.KeyA) {
						snakePlayer.direction = "left"
					}
					if inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
						inpututil.IsKeyJustPressed(ebiten.KeyD) {
						snakePlayer.direction = "right"
					}
				}
				if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
					GamePaused = true
				}
			} else {
				if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
					GamePaused = false
				} else if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
					GameState = "exit"
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
				GameOverSndPlaying = false
				GameJustEnded = false
				currScore = 0
				timeElapsed = 1 // it starts slower than the first timer for some reason
				timerTicker.Reset(1 * time.Second)

				setupInitialSnake()
			} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
				GameState = "exit"
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

func ParseHexColorAlpha(s string, a uint8) (c color.RGBA) {
	c.A = a
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

	// Show images
	// Dr. Nick
	nickOp := &ebiten.DrawImageOptions{}
	nickOp.GeoM.Scale(.4, .4)
	nickOp.GeoM.Translate(float64(ScreenWidth)*0.0571875, float64(ScreenHeight)*0.1953125)
	nickOp.ColorM.Scale(1, 1, 1, introOpacity)
	screen.DrawImage(drNick, nickOp)

	// Schneider
	schOp := &ebiten.DrawImageOptions{}
	schOp.GeoM.Scale(.70, .70)
	schOp.GeoM.Translate(float64(ScreenWidth)*0.370625, float64(ScreenHeight)*0.1853125)
	schOp.ColorM.Scale(1, 1, 1, introOpacity)
	screen.DrawImage(schImage, schOp)

	// Red Hat
	rhOp := &ebiten.DrawImageOptions{}
	rhOp.GeoM.Scale(1.2, 1.2)
	rhOp.GeoM.Translate(float64(ScreenWidth)*0.6940625, float64(ScreenHeight)*0.1653125)
	rhOp.ColorM.Scale(1, 1, 1, introOpacity)
	screen.DrawImage(rhImage, rhOp)

	// Ebitengine
	ebOp := &ebiten.DrawImageOptions{}
	ebOp.GeoM.Scale(1.2, 1.2)
	ebOp.GeoM.Translate(float64(ScreenWidth)*0.0571875, float64(ScreenHeight)*0.56875)
	ebOp.ColorM.Scale(1, 1, 1, introOpacity)
	screen.DrawImage(ebImage, ebOp)

	// Golang
	goOp := &ebiten.DrawImageOptions{}
	goOp.GeoM.Scale(.40, .40)
	goOp.GeoM.Translate(float64(ScreenWidth)*0.5859375, float64(ScreenHeight)*0.46875)
	goOp.ColorM.Scale(1, 1, 1, introOpacity)
	screen.DrawImage(goImage, goOp)

	// Increment opacity (fade in)
	if fadingOutIntro {
		introOpacity -= .01
		if introOpacity <= 0 {
			GameState = "title"
			fadingOutIntro = false
		}
	} else {
		introOpacity += .01
		if introOpacity >= 1 {
			introOpacity = 1
		}
	}

}

func drawBg(screen *ebiten.Image) {
	globBgOp := &ebiten.DrawImageOptions{}
	globBgOp.GeoM.Scale(globBgRot, globBgRot)
	screen.DrawImage(globBg, globBgOp)

	if zoomingBg {
		globBgRot += .0001
	} else {
		globBgRot -= .0001
	}

	if globBgRot >= 1.25 {
		zoomingBg = false
	}
	if globBgRot <= .75 {
		zoomingBg = true
	}
}

func drawTitle(screen *ebiten.Image) {
	// Background
	drawBg(screen)

	// Logo and text on top
	snake := &ebiten.DrawImageOptions{}
	snake.GeoM.Scale(.50, .50)
	snake.GeoM.Translate(float64((ScreenWidth/2))-(float64(ScreenWidth)*0.17), float64(ScreenHeight)*0.06125)
	screen.DrawImage(snakeLogo, snake)

	// Handle Menu
	if menuItem == "new_game" {
		text.Draw(screen, "> New Game", titleFont, (ScreenWidth/3)-30, (ScreenHeight/3)+190, color.White)
		text.Draw(screen, "Exit", titleFont, (ScreenWidth/3)+20, (ScreenHeight/3)+270, ParseHexColor("#8c8c8c"))
	} else if menuItem == "exit" {
		text.Draw(screen, "New Game", titleFont, (ScreenWidth/3)+20, (ScreenHeight/3)+190, ParseHexColor("#8c8c8c"))
		text.Draw(screen, "> Exit", titleFont, (ScreenWidth/3)-30, (ScreenHeight/3)+270, color.White)
	}
}

func doTitle(g *Game, screen *ebiten.Image) {
	// TODO: Add some sort of way to detect the center of the screen
	drawTitle(screen)

}

func getGridCellColor(ix int, iy int) color.Color {

	var theColor color.Color

	if ix%2 == 0 {
		if iy%2 == 0 {
			theColor = ParseHexColorAlpha(gridSolidColor, gridCellOpacity)
		} else {
			theColor = ParseHexColorAlpha(gridAltColor, gridCellOpacity)
		}
	} else {
		if iy%2 == 0 {
			theColor = ParseHexColorAlpha(gridAltColor, gridCellOpacity)
		} else {
			theColor = ParseHexColorAlpha(gridSolidColor, gridCellOpacity)
		}
	}

	return theColor
}

func DrawLine(dst *ebiten.Image, x1, y1, x2, y2, width float64, clr color.Color) {
	length := math.Hypot(x2-x1, y2-y1)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(length, width)
	op.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	op.GeoM.Translate(x1, y1)
	op.ColorM.ScaleWithColor(clr)
	// Filter must be 'nearest' filter (default).
	// Linear filtering would make edges blurred.
	dst.DrawImage(emptySubImage, op)
}

func buildGrid(screen *ebiten.Image) {

	// Draw BG
	drawBg(screen)

	// Draw Grid Border
	// Top
	DrawLine(screen, float64(borderLeft-gridBorderSize), float64(borderTop-gridBorderSize), float64(ScreenWidth-borderRight), float64(borderTop-gridBorderSize), gridBorderSize, ParseHexColor(gridBorderColor))
	// Bottom
	DrawLine(screen, float64(borderLeft-gridBorderSize), float64(ScreenHeight-borderBottom)-(float64(gridBorderSize)*3), float64(ScreenWidth-borderRight), float64(ScreenHeight-borderBottom)-(float64(gridBorderSize)*3), gridBorderSize, ParseHexColor(gridBorderColor))
	// Left
	DrawLine(screen, float64(borderLeft), float64(borderTop), float64(borderLeft), float64(ScreenHeight-borderBottom)-(float64(gridBorderSize)*2.5), gridBorderSize, ParseHexColor(gridBorderColor))
	// Right
	DrawLine(screen, float64(ScreenWidth-borderRight), float64(borderTop), float64(ScreenWidth-borderRight), float64(ScreenHeight-borderBottom)-(float64(gridBorderSize)*2.5), gridBorderSize, ParseHexColor(gridBorderColor))

	// Calculate Grid width using borders
	gridCellWidth = (ScreenWidth - borderLeft - borderRight) / gridWidth
	gridCellHeight = (ScreenHeight - borderTop - borderBottom) / gridHeight

	// Draw grid
	for ix := 0; ix < gridWidth; ix++ {
		for iy := 0; iy < gridHeight; iy++ {
			ebitenutil.DrawRect(screen, float64(ix*gridCellWidth)+float64(borderLeft), float64(iy*gridCellHeight)+float64(borderTop), float64(gridCellWidth), float64(gridCellHeight), getGridCellColor(ix, iy))
		}
	}

}

func drawGridPiece(screen *ebiten.Image, ix int, iy int, theColor color.Color, shapeType string) {
	if shapeType == "rect" {
		ebitenutil.DrawRect(screen, float64(ix*gridCellWidth)+float64(borderLeft), float64(iy*gridCellHeight)+float64(borderTop), float64(gridCellWidth), float64(gridCellHeight), theColor)
	}
	if shapeType == "smallcircle" {
		radius := float64((float64(gridCellWidth/5) + float64(gridCellHeight/5)) / 2)
		ebitenutil.DrawCircle(screen, float64(ix*gridCellWidth)+(radius*2.5)+float64(borderLeft), float64(iy*gridCellHeight)+(radius*2.5)+float64(borderTop), radius, theColor)
	}
	if shapeType == "circle" {
		radius := float64((float64(gridCellWidth/3) + float64(gridCellHeight/3)) / 2)
		ebitenutil.DrawCircle(screen, float64(ix*gridCellWidth)+(radius*1.5)+float64(borderLeft), float64(iy*gridCellHeight)+(radius*1.5)+float64(borderTop), radius, theColor)
	}
	if shapeType == "triangle" {
		//TODO: For tail, for now its a smaller cicle
	}
	if shapeType == "apple" {
		a := &ebiten.DrawImageOptions{}
		a.GeoM.Scale(appleScale, appleScale)
		a.GeoM.Translate(float64(ix*gridCellWidth)+5+float64(borderLeft), float64(iy*gridCellHeight)+2+float64(borderTop))
		screen.DrawImage(apple, a)
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
		drawGridPiece(screen, randX, randY, ParseHexColor(nomColor), "apple")
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

			// Increment score
			currScore += 1

		} else {
			drawGridPiece(screen, currentNom.xPos, currentNom.yPos, ParseHexColor(nomColor), "apple")
		}
	}
}

func showScore(screen *ebiten.Image) {

	diam := (float64(gridCellWidth/3) + float64(gridCellHeight/3))
	radius := diam / 2

	// Draw the apple
	a := &ebiten.DrawImageOptions{}
	a.GeoM.Scale(appleScale, appleScale)
	a.GeoM.Translate(float64(ScreenWidth/2)+105+float64(borderLeft), float64(borderTop/2)-18)
	screen.DrawImage(apple, a)

	// Score
	text.Draw(screen, strconv.Itoa(currScore), scoreFont, (ScreenWidth/2)+int(radius)+140, (borderTop/2)+(int(radius)/2)+3, color.White)

}

func doTimer() {
	go func() {
		for {
			select {
			case <-timerDone:
				return
			case <-timerTicker.C:
				timeElapsed += 1
			}
		}
	}()
}

func doAppleScale() {
	if GameStarted && !GamePaused {
		if zoomingApple {
			appleScale += .0005
		} else {
			appleScale -= .0005
		}

		if appleScale >= .1 {
			zoomingApple = false
		}
		if appleScale <= .09 {
			zoomingApple = true
		}
	}
}

func doGame(g *Game, screen *ebiten.Image) {
	// Draw background
	buildGrid(screen)

	// FX for apple
	doAppleScale()

	// Show score count
	showScore(screen)

	// Handle game started vs paused
	if GameStarted && !GamePaused && !GameOver {
		var moveCounter int
		if g.clockSpeedCount == 0 {
			moveCounter = 1
			// Show timer
			doTimer()
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
		text.Draw(screen, "Game Paused. Escape to resume\nor Q to quit.", baseFont, (ScreenWidth/3)-56, (ScreenHeight/3)+90, color.White)
	} else if !GameStarted && !GameOver {
		// Do not update snake
		// Show start text
		text.Draw(screen, "Arrow keys or WASD keys move snake\nEnter starts game", baseFont, (ScreenWidth/3)-70, (ScreenHeight/3)+130, color.White)
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
			snakePath = append([]pathPair{{snakePlayer.xPos, snakePlayer.yPos}}, snakePath[0:len(snakePlayer.snakeBody)]...)
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
				GameJustEnded = true
			}
		}
		// Check if the head collided with a wall
		if idx == 0 {
			if snakePlayer.xPos >= gridWidth || snakePlayer.xPos < 0 || snakePlayer.yPos >= gridHeight || snakePlayer.yPos < 0 {
				GameStarted = false
				GameOver = true
				GameJustEnded = true
			}
		}
	}

	text.Draw(screen, "Seconds Survived: "+strconv.Itoa(timeElapsed), timerFont, (ScreenWidth/3)-40, (int(math.Round(borderTop / 1.5))), color.White)
	// Show Game Over
	if GameOver {
		text.Draw(screen, "Womp womp. Game over.\n\nPress Enter for New Game\nor Escape to quit", baseFont, (ScreenWidth/2)-200, (ScreenHeight/2)-50, color.White)
		timerTicker.Stop()
	}

	// Handle game over sound
	if GameOver && GameJustEnded && !GameOverSndPlaying {
		GameOverSndPlaying = true
		gameOverSnd.Seek(0)
		gameOverSnd.Play()
	}
	if GameOver && GameOverSndPlaying {
		gameOverSnd.Play()
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

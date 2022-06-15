package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
	"time"
)

const (
	screenWidth  = 300
	screenHeight = 300

	ballRadius                        = 15
	ballAccelerationConstant          = float64(0.000000015)
	ballAccelerationSpeedUpMultiplier = float64(2)
	ballResistance                    = float64(0.975)
)

var (
	ballPositionX     = float64(screenWidth) / 2
	ballPositionY     = float64(screenHeight) / 2
	ballMovementX     = float64(0)
	ballMovementY     = float64(0)
	ballAccelerationX = float64(0)
	ballAccelerationY = float64(0)
	prevUpdateTime    = time.Now() // Acts as delta time would in a game engine
)

type Game struct {
	pressedKeys []ebiten.Key
}

func (g *Game) Update() error {
	timeDelta := float64(time.Since(prevUpdateTime))
	prevUpdateTime = time.Now()

	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	ballAccelerationX = 0
	ballAccelerationY = 0

	acc := ballAccelerationConstant

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "Space":
			acc *= ballAccelerationSpeedUpMultiplier
		}
	}

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "S":
			ballAccelerationY = acc
		case "W":
			ballAccelerationY = -acc
		case "D":
			ballAccelerationX = acc
		case "A":
			ballAccelerationX = -acc
		}
	}

	ballMovementY += ballAccelerationY
	ballMovementX += ballAccelerationX

	ballMovementX *= ballResistance
	ballMovementY *= ballResistance

	ballPositionX += ballMovementX * timeDelta
	ballPositionY += ballMovementY * timeDelta

	const minX = ballRadius
	const minY = ballRadius
	const maxX = screenWidth - ballRadius
	const maxY = screenHeight - ballRadius

	if ballPositionX >= maxX || ballPositionX <= minX {
		if ballPositionX > maxX {
			ballPositionX = maxX
		} else if ballPositionX < minX {
			ballPositionX = minX
		}
		ballMovementX *= -1
	}
	if ballPositionY >= maxY || ballPositionY <= minY {
		if ballPositionY > maxY {
			ballPositionY = maxY
		} else if ballPositionY < minY {
			ballPositionY = minY
		}
		ballMovementY *= -1
	}
	return nil
}

var simpleShader *ebiten.Shader

func init() {
	var err error
	simpleShader, err = ebiten.NewShader([]byte(`
	package main

	func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
		return color
	}
	`))
	if err != nil {
		panic(err)
	}
}

func (g *Game) DrawSquare(screen *ebiten.Image, x, y int, color color.Color) {
	for width := 100; width < x; width++ {
		for height := 100; height < y; height++ {
			screen.Set(width, height, color)
		}
	}
}

func (g *Game) DrawCircle(screen *ebiten.Image, x, y, radius float32, color color.RGBA) {
	var path vector.Path

	path.MoveTo(x, y)
	path.Arc(x, y, radius, 0, math.Pi*2, vector.Clockwise)

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	redScaled := float32(color.R) / 255
	greenScaled := float32(color.G) / 255
	blueScaled := float32(color.B) / 255
	alphaScaled := float32(color.A) / 255

	for i := range vertices {
		v := &vertices[i]
		v.ColorR = redScaled
		v.ColorG = greenScaled
		v.ColorB = blueScaled
		v.ColorA = alphaScaled
	}

	screen.DrawTrianglesShader(vertices, indices, simpleShader, &ebiten.DrawTrianglesShaderOptions{
		FillRule: ebiten.EvenOdd,
	})
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Game render logic
	purpleCol := color.RGBA{255, 0, 255, 255}
	redCol := color.RGBA{255, 0, 0, 255}

	g.DrawSquare(screen, 200, 200, purpleCol)
	g.DrawCircle(screen, float32(ballPositionX), float32(ballPositionY), float32(ballRadius), redCol)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Game layout, make dynamic later, currently fixed
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Test Game")
	// Call and catch error
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

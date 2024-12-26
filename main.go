package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Point struct {
	X, Y int
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

type Game struct {
	snake     []Point
	direction Point
	food      Point
	score     int
	gameOver  bool
	moveTime  float64
	moveSpeed float64
}

const (
	gridSize   = 20
	gridWidth  = 32
	gridHeight = 24
)

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{
		snake: []Point{
			// {X: gridWidth/2 - 3, Y: gridHeight / 2},
			// {X: gridWidth/2 - 2, Y: gridHeight / 2},
			// {X: gridWidth/2 - 1, Y: gridHeight / 2},
			{X: gridWidth / 2, Y: gridHeight / 2},
		},
		direction: Point{X: 1, Y: 0},
		moveSpeed: 0.15,
	}
	g.placeFood()
	fmt.Println("New game: ", g)
	return g
}

func (g *Game) placeFood() {
	g.food = Point{
		X: rand.Intn(gridWidth),
		Y: rand.Intn(gridHeight),
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			*g = *NewGame()
		}
		return nil
	}

	// Handle input
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction.X == 0 {
		g.direction = Point{X: -1, Y: 0}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction.X == 0 {
		g.direction = Point{X: 1, Y: 0}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction.Y == 0 {
		g.direction = Point{X: 0, Y: -1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction.Y == 0 {
		g.direction = Point{X: 0, Y: 1}
	}

	// Update movement timer
	g.moveTime += 1.0 / 60.0 // Assuming 60 FPS

	// Move snake when timer exceeds moveSpeed
	if g.moveTime >= g.moveSpeed {
		g.moveTime = 0

		head := g.snake[0]
		newHead := Point{
			X: head.X + g.direction.X,
			Y: head.Y + g.direction.Y,
		}

		// Wrap around logic
		if newHead.X < 0 {
			newHead.X = gridWidth - 1 // Wrap to right edge
		} else if newHead.X >= gridWidth {
			newHead.X = 0 // Wrap to left edge
		}

		if newHead.Y < 0 {
			newHead.Y = gridHeight - 1 // Wrap to bottom
		} else if newHead.Y >= gridHeight {
			newHead.Y = 0 // Wrap to top
		}

		// Check wall collision
		// if newHead.X < 0 || newHead.X >= gridWidth ||
		// 	newHead.Y < 0 || newHead.Y >= gridHeight {
		// 	fmt.Println("Wall collision")
		// 	g.gameOver = true
		// 	return nil
		// }

		// Check self collision - check if new head position hits any part of the snake
		for i := 0; i < len(g.snake); i++ {
			if newHead.X == g.snake[i].X && newHead.Y == g.snake[i].Y {
				fmt.Println("Self collision")
				g.gameOver = true
				return nil
			}
		}

		// fmt.Printf("Before move - Snake segments: %v\n", g.snake)
		// fmt.Printf("New head position: {X:%d, Y:%d}\n", newHead.X, newHead.Y)

		if newHead.Equals(g.food) {
			// Instead of append with the whole snake, just add the current tail
			g.snake = append([]Point{newHead}, g.snake...)
			// fmt.Printf("After eating - Snake segments: %v\n", g.snake)
			g.score++
			g.placeFood()
		} else {
			g.snake = append([]Point{newHead}, g.snake[:len(g.snake)-1]...)
			// fmt.Printf("After move - Snake segments: %v\n", g.snake)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw snake
	for _, p := range g.snake {
		ebitenutil.DrawRect(screen,
			float64(p.X*gridSize),
			float64(p.Y*gridSize),
			gridSize-1,
			gridSize-1,
			color.RGBA{0, 255, 0, 255})
	}

	// Draw food
	ebitenutil.DrawRect(screen,
		float64(g.food.X*gridSize),
		float64(g.food.Y*gridSize),
		gridSize-1,
		gridSize-1,
		color.RGBA{255, 0, 0, 255})

	// Draw score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrint(screen, scoreText)

	if g.gameOver {
		msg := "Game Over!\nPress ENTER to restart"
		fmt.Println("Score: ", g.score)

		// Get the bounds of the screen
		bounds := screen.Bounds()
		x := bounds.Dx()/2 - len(msg)*3 // Approximate center position
		y := bounds.Dy() / 2

		// Draw game over text
		ebitenutil.DebugPrint(screen, msg)

		// Draw centered text at x, y position
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(NewGame()); err != nil {
		fmt.Println("Fatal error: ", err)
		panic(err)
	}
}

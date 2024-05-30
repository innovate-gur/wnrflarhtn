package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  int32 = 800
	screenHeight int32 = 800
)

type Game struct {
	framecount    int32
	player        Player
	arrows        []Arrow
	arrowConstant float32
	screen        string
	alive         bool
	button        Button
}
type Player struct {
	pos   rl.Vector2
	color rl.Color
}
type Arrow struct {
	pos           rl.Vector2
	angle         float64
	angleToplayer float64
	slope         float32
	active        bool
}
type Button struct {
	frame rl.Rectangle
	text  string
}

func main() {

	g := Game{}
	g.Init()

	rl.InitWindow(screenWidth, screenHeight, "죽림고수")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		switch g.screen {
		case "PLAY", "DEAD":
			g.Update()
			g.DrawPlay()
		case "START":
			g.DrawStart()
		}
	}
}

func (g *Game) Init() {
	g.framecount = 0
	g.player = Player{rl.NewVector2(400, 400), rl.Black}
	g.arrows = make([]Arrow, 150)
	g.screen = "START"
	g.alive = true
	g.arrowConstant = 3

}

func (g *Game) Update() {
	if g.screen == "PLAY" {
		g.framecount += 1
		var dt float32 = rl.GetFrameTime() * float32(rl.GetFPS())

		if rl.IsKeyDown(rl.KeyDown) {
			g.player.pos.Y += dt * 5
		}
		if rl.IsKeyDown(rl.KeyUp) {
			g.player.pos.Y -= dt * 5
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			g.player.pos.X -= dt * 5
		}
		if rl.IsKeyDown(rl.KeyRight) {
			g.player.pos.X += dt * 5
		}

		count := 0

		if g.framecount%120 == 0 {

			for i := 0; i < len(g.arrows); i++ {
				if count == 9 {
					break
				}
				if !g.arrows[i].active {
					g.arrows[i].InitPosition(g.player.pos)
					count += 1
				}
			}
		}

		for i := 0; i < len(g.arrows); i++ {
			if g.arrows[i].active {
				g.arrows[i].pos.X += float32(math.Cos(g.arrows[i].angleToplayer)) * g.arrowConstant * dt
				g.arrows[i].pos.Y += float32(math.Sin(g.arrows[i].angleToplayer)) * g.arrowConstant * dt
				if math.Hypot(float64(400-g.arrows[i].pos.X), float64(400-g.arrows[i].pos.Y)) >= 399 {
					g.arrows[i].active = false
				}

			}

			if rl.CheckCollisionCircleRec(g.arrows[i].pos, 5, rl.NewRectangle(g.player.pos.X, g.player.pos.Y, 25, 25)) {
				g.screen = "DEAD"
			}
		}
	}

}

func (g *Game) DrawPlay() {

	rl.BeginDrawing()

	//Initialize Window
	rl.ClearBackground(rl.RayWhite)
	rl.DrawCircleLines(400, 400, 400, rl.Black)

	//score
	score := float64(g.framecount) / float64(60)

	rl.DrawText(fmt.Sprintf("%.2f", score), 150, 300, 300, rl.LightGray)

	//player
	rl.DrawRectangleV(g.player.pos, rl.NewVector2(25, 25), g.player.color)

	//arrow
	for _, arrow := range g.arrows {
		if arrow.active {
			rl.DrawCircle(int32(arrow.pos.X), int32(arrow.pos.Y), 5, rl.Red)
			// rl.DrawLine(int32(arrow.pos.X), int32(arrow.pos.Y), int32(p.pos.X), int32(p.pos.Y), rl.Black)
		}
	}
	if g.screen == "DEAD" {
		rl.DrawRectangle(250, 250, 300, 300, rl.Green)
		rl.DrawText(fmt.Sprintf("Score: %.2f", score), 250, 250, 50, rl.Black)
		btnToStart := Button{rl.NewRectangle(200, 400, 200, 50), "To Start"}
		btnPlayAgain := Button{rl.NewRectangle(400, 400, 200, 50), "Play Again"}
		btnToStart.Render()
		btnPlayAgain.Render()

		if btnToStart.Hover() {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				g.Init()
				g.screen = "START"
			}
		}
		if btnPlayAgain.Hover() {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				g.Init()
				g.screen = "PLAY"
			}
		}
	}
	rl.EndDrawing()
}

func (g *Game) DrawStart() {

	rl.BeginDrawing()

	rl.ClearBackground(rl.White)

	btnGameStart := Button{rl.NewRectangle(300, 300, 200, 50), "game start"}
	btnHowToPlay := Button{rl.NewRectangle(300, 400, 200, 50), "how to play"}

	btnGameStart.Render()
	btnHowToPlay.Render()

	if btnGameStart.Hover() {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			g.screen = "PLAY"
		}
	}
	if btnHowToPlay.Hover() {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			fmt.Println("Click")
		}
	}

	rl.EndDrawing()

}

func (a *Arrow) InitPosition(playerPosition rl.Vector2) {
	a.angle = rand.Float64() * 2 * math.Pi
	a.pos.X = float32(400*(math.Cos(a.angle))) + 400
	a.pos.Y = float32(400*(math.Sin(a.angle))) + 400
	a.slope = (playerPosition.Y - a.pos.Y) / (playerPosition.X - a.pos.X)
	a.active = true
	a.angleToplayer = math.Atan2(float64(playerPosition.Y-a.pos.Y), float64(playerPosition.X-a.pos.X))
}

func (b Button) Render() {
	rl.DrawRectangle(b.frame.ToInt32().X, b.frame.ToInt32().Y, b.frame.ToInt32().Width, b.frame.ToInt32().Height, rl.Black)
	rl.DrawText(b.text, b.frame.ToInt32().X, b.frame.ToInt32().Y, 30, rl.White)
}

func (b Button) Hover() bool {
	return rl.CheckCollisionPointRec(rl.GetMousePosition(), b.frame)
}

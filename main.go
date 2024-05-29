package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth, screenHeight int32 = 800, 800
)

type player struct {
	pos   rl.Vector2
	color rl.Color
}
type arrow struct {
	pos   rl.Vector2
	slope float32
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Title Here")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	p := player{rl.NewVector2(0, 0), rl.Black}

	a := arrow{}

	for !rl.WindowShouldClose() {

		var dt float32 = rl.GetFrameTime() * float32(rl.GetFPS())

		if rl.IsKeyDown(rl.KeyS) {
			p.pos.Y += dt * 5
		}
		if rl.IsKeyDown(rl.KeyW) {
			p.pos.Y -= dt * 5
		}
		if rl.IsKeyDown(rl.KeyA) {
			p.pos.X -= dt * 5
		}
		if rl.IsKeyDown(rl.KeyD) {
			p.pos.X += dt * 5
		}
		if rl.IsKeyPressed(rl.KeyR) {
			a.InitPosition(p.pos)
		}

		rl.BeginDrawing()

		//Initialize Window
		rl.ClearBackground(rl.RayWhite)
		rl.DrawCircleLines(400, 400, 400, rl.Black)

		//player
		rl.DrawRectangleV(p.pos, rl.NewVector2(25, 25), p.color)

		//arrow
		rl.DrawCircle(int32(a.pos.X), int32(a.pos.Y), 10, rl.Black)
		rl.DrawLine(int32(a.pos.X), int32(a.pos.Y), int32(p.pos.X), int32(p.pos.Y), rl.Black)

		rl.EndDrawing()
	}
}

func (a *arrow) InitPosition(playerPosition rl.Vector2) {
	var angle float64 = rand.Float64() * 2 * math.Pi
	a.pos.X = float32(400*(math.Cos(angle))) + 400
	a.pos.Y = float32(400*(math.Sin(angle))) + 400
	a.slope = (playerPosition.Y - a.pos.Y) / (playerPosition.X - a.pos.X)
}

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
	framecount2   int32
	player        Player
	arrows        []Arrow
	arrowConstant float32
	screen        string
	alive         bool
	wave          int32
	timer         Timer

	texture_arrow  rl.Texture2D
	texture_player rl.Texture2D
}
type Player struct {
	pos   rl.Vector2
	color rl.Color
}
type Arrow struct {
	pos    rl.Vector2
	angle  float64
	slope  float32
	active bool
}
type Button struct {
	frame rl.Rectangle
	text  string
}
type Timer struct {
	startTime float64
}

func (t *Timer) StartTime() {
	t.startTime = rl.GetTime()
}
func (t Timer) ElapsedTime() float64 {
	return (rl.GetTime() - t.startTime)
}

func main() {

	g := Game{}

	rl.InitWindow(screenWidth, screenHeight, "죽림고수")
	g.Init()
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()
	defer rl.UnloadTexture(g.texture_arrow)
	defer rl.UnloadTexture(g.texture_player)

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
	g.wave = 0
	g.Load()
}

func (g *Game) Update() {
	if g.screen == "PLAY" {
		g.framecount += 1
		g.framecount2 += 1
		var dt float32 = rl.GetFrameTime() * float32(rl.GetFPS())

		if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
			g.player.pos.Y += dt * 4
		}
		if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
			g.player.pos.Y -= dt * 4
		}
		if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
			g.player.pos.X -= dt * 4
		}
		if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
			g.player.pos.X += dt * 4
		}

		switch g.wave {

		case 0:
			if g.framecount%60 == 0 {
				count := 0
				i := 0
				for count < 9 {
					if !g.arrows[i].active {
						g.arrows[i].angle = rand.Float64() * 2 * math.Pi
						g.arrows[i].pos.X = float32(400*(math.Cos(g.arrows[i].angle))) + 400
						g.arrows[i].pos.Y = float32(400*(math.Sin(g.arrows[i].angle))) + 400
						g.arrows[i].slope = (g.player.pos.Y - g.arrows[i].pos.Y) / (g.player.pos.X - g.arrows[i].pos.X)
						g.arrows[i].angle = math.Atan2(float64(g.player.pos.Y-g.arrows[i].pos.Y), float64(g.player.pos.X-g.arrows[i].pos.X))
						g.arrows[i].active = true
						count += 1

					}
					i += 1
				}
			}

		case 1:
			fmt.Println(g.framecount)
			if g.framecount%200 == 0 {
				fmt.Println("Temp Called")
				go g.Temp()
			}

		case 2:
			if g.framecount%120 == 0 {
				count := 0
				i := 0
				for count < 9 {
					i += 1
					if !g.arrows[i].active {
						g.arrows[i].angle = 0
						g.arrows[i].pos.X = float32(400*(math.Cos(g.arrows[i].angle))) + 400
						g.arrows[i].pos.Y = float32(400*(math.Sin(g.arrows[i].angle))) + 400
						g.arrows[i].slope = (g.player.pos.Y - g.arrows[i].pos.Y) / (g.player.pos.X - g.arrows[i].pos.X)
						g.arrows[i].angle = math.Atan2(float64(g.player.pos.Y-g.arrows[i].pos.Y), float64(g.player.pos.X-g.arrows[i].pos.X))
						g.arrows[i].active = true
						count += 1
					}
				}
			}

		}

		for i := 0; i < len(g.arrows); i++ {
			if g.arrows[i].active {
				g.arrows[i].pos.X += float32(math.Cos(g.arrows[i].angle)) * g.arrowConstant * dt
				g.arrows[i].pos.Y += float32(math.Sin(g.arrows[i].angle)) * g.arrowConstant * dt
				if math.Hypot(float64(400-g.arrows[i].pos.X), float64(400-g.arrows[i].pos.Y)) >= 399 {
					g.arrows[i].active = false
				}
			}

			if rl.CheckCollisionCircles(g.arrows[i].pos, 3.5, rl.NewVector2(g.player.pos.X, g.player.pos.Y), 10) {
				g.screen = "DEAD"
			}
		}

	}
}

func (g *Game) DrawPlay() {

	rl.BeginDrawing()

	//Initialize Window
	rl.ClearBackground(rl.RayWhite)
	rl.DrawCircle(400, 400, 400, rl.Black)

	//score
	score := float64(g.framecount) / float64(60)

	if g.screen == "DEAD" {
		rl.DrawText(fmt.Sprintf("%.2f", score), 150, 300, 300, rl.LightGray)
	} else {
		rl.DrawText(fmt.Sprintf("%.2f", g.timer.ElapsedTime()), 150, 300, 300, rl.LightGray)

	}

	//player
	rl.DrawCircleV(g.player.pos, 10, g.player.color)
	rl.DrawTextureEx(g.texture_player, rl.NewVector2(g.player.pos.X-12, g.player.pos.Y-13), 0, 1, rl.White)

	//arrow
	for _, arrow := range g.arrows {
		if arrow.active {
			// rl.DrawCircle(int32(arrow.pos.X), int32(arrow.pos.Y), 3.5, rl.Red)
			// rl.DrawLine(int32(arrow.pos.X), int32(arrow.pos.Y), int32(p.pos.X), int32(p.pos.Y), rl.Black)
			dest := rl.NewRectangle(arrow.pos.X, arrow.pos.Y, 30, 30)
			rl.DrawTexturePro(g.texture_arrow, rl.NewRectangle(0, 0, 25, 25), dest, rl.NewVector2(dest.Width-6, dest.Height/2+1), (float32(arrow.angle)*180/math.Pi + math.Pi), rl.White)
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
				g.timer.StartTime()
			}
		}
	}
	rl.EndDrawing()
}

func (g *Game) Load() {
	g.texture_arrow = rl.LoadTexture("assets/img/arrow.png")
	g.texture_player = rl.LoadTexture("assets/img/player.png")
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
			g.timer.StartTime()
		}
	}
	if btnHowToPlay.Hover() {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			fmt.Println("Click")
		}
	}

	rl.EndDrawing()

}

func (a *Arrow) FixPosition(playerPosition rl.Vector2, angle float64) {
	a.angle = angle
	a.pos.X = float32(400*(math.Cos(a.angle))) + 400
	a.pos.Y = float32(400*(math.Sin(a.angle))) + 400
	a.slope = (playerPosition.Y - a.pos.Y) / (playerPosition.X - a.pos.X)
	// a.angleToplayer = math.Atan2(float64(playerPosition.Y-a.pos.Y), float64(playerPosition.X-a.pos.X))
	a.active = true
}

func (b Button) Render() {
	rl.DrawRectangle(b.frame.ToInt32().X, b.frame.ToInt32().Y, b.frame.ToInt32().Width, b.frame.ToInt32().Height, rl.Black)
	rl.DrawText(b.text, b.frame.ToInt32().X, b.frame.ToInt32().Y, 30, rl.White)
}

func (b Button) Hover() bool {
	return rl.CheckCollisionPointRec(rl.GetMousePosition(), b.frame)
}

func (g *Game) Temp() {
	count := 0
	firstAngle := rand.Float64() * 2 * math.Pi
	i := 0

	for count < 9 {
		i += 1

		if !g.arrows[i].active {
			g.arrows[i].angle = firstAngle - (math.Pi/float64(6))*float64(i)
			g.arrows[i].pos.X = float32(400*(math.Cos(g.arrows[i].angle))) + 400
			g.arrows[i].pos.Y = float32(400*(math.Sin(g.arrows[i].angle))) + 400
			g.arrows[i].slope = (g.player.pos.Y - g.arrows[i].pos.Y) / (g.player.pos.X - g.arrows[i].pos.X)
			g.arrows[i].angle = math.Atan2(float64(g.player.pos.Y-g.arrows[i].pos.Y), float64(g.player.pos.X-g.arrows[i].pos.X))
			g.arrows[i].active = true
			count += 1
			rl.WaitTime(0.2)
		}
	}
	fmt.Println("temp End")
}

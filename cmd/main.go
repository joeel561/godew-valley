package main

import (
	"encoding/json"
	"fmt"
	"godew-valley/pkg/player"
	"io/ioutil"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

type JsonMap struct {
	Layers    []Layer `json:"layers"`
	MapHeight int     `json:"mapHeight"`
	MapWidth  int     `json:"mapWidth"`
	TileSize  int     `json:"tileSize"`
}

type Layer struct {
	Name  string `json:"name"`
	Tiles []Tile `json:"tiles"`
}

type Tile struct {
	Id string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

var (
	running = true
	bgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite    rl.Texture2D
	waterSprite    rl.Texture2D
	tex            rl.Texture2D
	playerSprite   rl.Texture2D
	spritesheetMap rl.Texture2D
	wall           rl.Rectangle
	oldX, oldY     float32

	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int
	playerHitBox                                  rl.Rectangle
	playerHitBoxYOffset                           float32 = 3

	frameCount int

	tileDest rl.Rectangle
	tileSrc  rl.Rectangle
	jsonMap  JsonMap

	playerSpeed float32 = 1.4

	musicPaused bool
	music       rl.Music
	printDebug  bool

	cam rl.Camera2D
)

func rectToString(rec rl.Rectangle) string {
	return fmt.Sprintf("X:%v, Y:%v, H:%v, W:%v", rec.X, rec.Y, rec.Height, rec.Width)
}

func vec2ToString(vec rl.Vector2) string {
	return fmt.Sprintf("X:%v, Y:%v", vec.X, vec.Y)
}

func debugText() []string {
	return []string{
		fmt.Sprintf("FPS: %v", rl.GetFPS()),
		fmt.Sprintf("Cam Target %v", vec2ToString(cam.Target)),
		fmt.Sprintf("Player Direction: %v   U:%v, D:%v, L:%v, R:%v", playerDir, playerUp, playerDown, playerLeft, playerRight),
		fmt.Sprintf("Player Speed: %v", playerSpeed),
		fmt.Sprintf("Player Frame: %v", playerFrame),
		fmt.Sprintf("Player Moving: %v", playerMoving),
		fmt.Sprintf("Player Src %v", rectToString(playerSrc)),
		fmt.Sprintf("Player Dest %v", rectToString(playerDest)),
		fmt.Sprintf("Music Paused: %v", musicPaused),
	}
}

func drawDebug(debugText []string) {
	textSize := 10
	lineSpace := 15

	offsetX := 10
	offsetY := 10

	for i, line := range debugText {
		rl.DrawText(line, int32(offsetX), int32(offsetY+lineSpace*i), int32(textSize), rl.Black)
	}
}

func drawScene() {
	var groundTiles []Tile
	var waterTiles []Tile

	for i := 0; i < len(jsonMap.Layers); i++ {
		if jsonMap.Layers[i].Name == "Background" {
			groundTiles = jsonMap.Layers[i].Tiles
		}

		if jsonMap.Layers[i].Name == "Water" {
			waterTiles = jsonMap.Layers[i].Tiles
		}
	}

	for i := 0; i < len(waterTiles); i++ {
		s, _ := strconv.ParseInt(waterTiles[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(jsonMap.TileSize)
		tileSrc.X = float32(jsonMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(jsonMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(waterTiles[i].X * jsonMap.TileSize)
		tileDest.Y = float32(waterTiles[i].Y * jsonMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	for i := 0; i < len(groundTiles); i++ {
		s, _ := strconv.ParseInt(groundTiles[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(jsonMap.TileSize)
		tileSrc.X = float32(jsonMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(jsonMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(groundTiles[i].X * jsonMap.TileSize)
		tileDest.Y = float32(groundTiles[i].Y * jsonMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	rl.DrawRectangle(int32(wall.X), int32(wall.Y), int32(wall.Width), int32(wall.Height), rl.Red)

	if printDebug {
		// Draw cetner map cross
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(-20, 0), 1, rl.Gray)
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(20, 0), 1, rl.Red)
		rl.DrawTriangle(rl.NewVector2(16, 2), rl.NewVector2(20, 0), rl.NewVector2(16, -2), rl.Red)
		rl.DrawText("X", int32(22), int32(-5), int32(10), rl.Black)
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(0, -20), 1, rl.Gray)
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(0, 20), 1, rl.Blue)
		rl.DrawTriangle(rl.NewVector2(-2, 16), rl.NewVector2(0, 20), rl.NewVector2(2, 16), rl.Blue)
		rl.DrawText("Y", int32(-2), int32(22), int32(10), rl.Black)

		// Draw collision rectangle
		rl.DrawRectangleLinesEx(playerHitBox, 1, rl.Green)
		rl.DrawRectangleLinesEx(playerDest, 1, rl.Purple)
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(0, 0), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyF10) {
		display := rl.GetCurrentMonitor()
		if rl.IsWindowFullscreen() {
			rl.SetWindowSize(screenWidth, screenHeight)
		} else {
			rl.SetWindowSize(rl.GetMonitorWidth(display), rl.GetMonitorHeight(display))
		}

		rl.ToggleFullscreen()
	}

	player.PlayerInput()
	if rl.IsKeyPressed(rl.KeyF3) {
		printDebug = !printDebug
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}
func update() {
	running = !rl.WindowShouldClose()

	oldX, oldY := playerDest.X, playerDest.Y

	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed

			if playerSpeed == 2 {
				playerDir = 9
			}
		}
		if playerDown {
			playerDest.Y += playerSpeed

			if playerSpeed == 2 {
				playerDir = 8
			}
		}
		if playerLeft {
			playerDest.X -= playerSpeed

			if playerSpeed == 2 {
				playerDir = 11
			}

		}
		if playerRight {
			playerDest.X += playerSpeed

			if playerSpeed == 2 {
				playerDir = 10
			}
		}

		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++

	}

	frameCount++
	if playerFrame >= 8 {
		playerFrame = 0
	}

	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.Y = playerSrc.Height * float32(playerDir)
	playerSrc.X = playerSrc.Width * float32(playerFrame)

	playerHitBox.X = playerDest.X + (playerDest.Width / 2) - playerHitBox.Width/2
	playerHitBox.Y = playerDest.Y + (playerDest.Height / 2) + playerHitBoxYOffset

	fmt.Println(oldX, oldY, "oldX, oldY")
	fmt.Println(playerDest.X, playerDest.Y, "playerDest.X, playerDest.Y")

	if rl.CheckCollisionRecs(playerHitBox, wall) {
		playerDest.X = oldX
		playerDest.Y = oldY
	}

	var waterTiles []Tile

	for i := 0; i < len(jsonMap.Layers); i++ {
		if jsonMap.Layers[i].Name == "Water" {
			waterTiles = jsonMap.Layers[i].Tiles
		}
	}

	for i := 0; i < len(waterTiles); i++ {
		if playerHitBox.X < float32(waterTiles[i].X*jsonMap.TileSize+jsonMap.TileSize) &&
			playerHitBox.X+playerHitBox.Width > float32(waterTiles[i].X*jsonMap.TileSize) &&
			playerHitBox.Y < float32(waterTiles[i].Y*jsonMap.TileSize+jsonMap.TileSize) &&
			playerHitBox.Y+playerHitBox.Height > float32(waterTiles[i].Y*jsonMap.TileSize) {

			playerDest.X = oldX
			playerDest.Y = oldY
		}

	}

	groundTilesRect := rl.NewRectangle(0, 0, float32(jsonMap.MapWidth*jsonMap.TileSize), float32(jsonMap.MapHeight*jsonMap.TileSize))
	if !rl.CheckCollisionRecs(playerHitBox, groundTilesRect) {
		playerDest.X = oldX
		playerDest.Y = oldY
	}

	//rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
}

func update() {
	running = !rl.WindowShouldClose()

	PlayerMoving()

	//rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bgColor)
	rl.BeginMode2D(cam)

	drawScene()
	rl.EndMode2D()

	if printDebug {
		drawDebug(debugText())
	}

	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := os.Open(mapFile)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &jsonMap)
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "godew valley - a game by joeel56")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	spritesheetMap = rl.LoadTexture("assets/spritesheet.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	playerSprite = rl.LoadTexture("assets/Characters/CharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	wall = rl.NewRectangle(50, 100, 200, 100)
	playerDest = rl.NewRectangle(100, 50, 60, 60)
	playerHitBox = rl.NewRectangle(0, 0, 10, 10)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/bgmusic.mp3")

	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2))), 0, 3)

	printDebug = false

	loadMap("map.json")

}
func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(waterSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadTexture(spritesheetMap)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 1080
)

type JsonMap struct {
	Layers    []Layer `json:"layers"`
	MapHeight int     `json:"mapHeight"`
	MapWidth  int     `json:"mapWidth"`
	TileMap   int     `json:"TileMap"`
}

type Layer struct {
	Name  string `json:"name"`
	Tiles []Tile `json:"tiles"`
}

type Tile struct {
	Id int `json:"id"`
	X  int `json:"x"`
	Y  int `json:"y"`
}

var (
	running = true
	bgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite    rl.Texture2D
	waterSprite    rl.Texture2D
	tex            rl.Texture2D
	playerSprite   rl.Texture2D
	spritesheetMap rl.Texture2D

	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int

	frameCount int

	tileDest            rl.Rectangle
	tileSrc             rl.Rectangle
	tileMap             []int
	srcMap              []string
	mapWidth, mapHeight int

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
		fmt.Sprintf("srcMap %v", srcMap),
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
	var jsonMap JsonMap

	for i := 0; i < len(jsonMap.Layers[0].Tiles); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapWidth)
			tileDest.Y = tileDest.Height * float32(i/mapWidth)

			/* 			if srcMap[i] == "g" {
			   				tex = grassSprite
			   			}

			   			if srcMap[i] == "w" {
			   				tex = waterSprite
			   			} */

			tex = spritesheetMap

			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))

			rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
		}
	}
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
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

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 5
		playerUp = true
	}

	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 4
		playerDown = true
	}

	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 7
		playerLeft = true
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 6
		playerRight = true
	}

	if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
		playerSpeed = 2
	} else {
		playerSpeed = 1.4
	}

	if rl.IsKeyPressed(rl.KeyF3) {
		printDebug = !printDebug
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}
func update() {
	running = !rl.WindowShouldClose()

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

	var jsonMap JsonMap

	json.Unmarshal(byteValue, &jsonMap)

	for i := 0; i < len(jsonMap.Layers); i++ {
		fmt.Println(jsonMap.Layers[i].Name)
		fmt.Println(jsonMap.MapHeight)
		os.Exit(1)
	}

	/* 	var result map[string]interface{}
	   	json.Unmarshal([]byte(byteValue), &result)

	   	remNewLine := strings.Replace(mapFile, "\r\n", " ", -1)
	   	sliced := strings.Split(remNewLine, " ") */
	mapWidth = -1
	mapHeight = -1

	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)

		if mapWidth == -1 {
			mapWidth = m
		} else if mapHeight == -1 {
			mapHeight = m
		} else if i < mapWidth*mapHeight+2 {
			tileMap = append(tileMap, m)

		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}

	if len(tileMap) > mapWidth*mapHeight {
		tileMap = tileMap[:len(tileMap)-1]
	}
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "godew valley - a game by joeel56")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("res/Tilesets/ground-tiles/New-tiles/Grass_tiles_v2.png")
	waterSprite = rl.LoadTexture("res/Tilesets/ground-tiles/water-frames/Water_1.png")

	spritesheetMap = rl.LoadTexture("res/spritesheet.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	playerSprite = rl.LoadTexture("res/Characters/CharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 60, 60)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/bgmusic.mp3")

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

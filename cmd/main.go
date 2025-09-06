package main

import (
	"godew-valley/pkg/debug"
	"godew-valley/pkg/items"
	"godew-valley/pkg/player"
	"godew-valley/pkg/save"
	"godew-valley/pkg/userinterface"
	"godew-valley/pkg/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

var (
	running = true
	bgColor = rl.NewColor(147, 211, 196, 255)

	musicPaused bool
	music       rl.Music
	printDebug  bool
)

func drawScene() {
	world.DrawWorld(player.Cam, screenWidth, screenHeight)

	items.DrawItems()

	world.DrawDoors()

	if printDebug {
		debug.DrawPlayerOutlines()
	}

	player.DrawPlayerTexture()
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "godew valley - a game by joeel56")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	world.InitWorld()
	world.InitDoors()
	player.InitPlayer()
	userinterface.InitUserInterface()

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/bgmusic.mp3")

	musicPaused = false
	rl.PlayMusicStream(music)

	printDebug = false

	world.LoadMap("pkg/world/world.json")
	userinterface.LoadUserInterfaceMap("pkg/userinterface/userinterface.json")
	items.InitItemTextures()
	save.LoadGame()
	items.InitItems()
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

	if rl.IsKeyPressed(rl.KeyF5) {
		save.SaveGame()
	}

	if rl.IsKeyPressed(rl.KeyF9) {
		save.LoadGame()
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		running = false
	}

	userinterface.ItemBarInput()

	items.InputHoe()
}

func update() {
	running = !rl.WindowShouldClose()
	player.PlayerMoving()
	items.UpdateItems()

	//rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}
}

func render() {
	var cam = player.Cam

	rl.BeginDrawing()
	rl.ClearBackground(bgColor)
	rl.BeginMode2D(cam)

	drawScene()
	rl.EndMode2D()

	userinterface.DrawUserInterface()
	if printDebug {
		debug.DrawDebug(debug.DebugText())
	}

	rl.EndDrawing()
}

func quit() {
	save.SaveGame()
	player.UnloadPlayerTexture()
	world.UnloadWorldTexture()
	world.UnloadDoors()
	items.UnloadItems()
	userinterface.UnloadUserInterface()
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

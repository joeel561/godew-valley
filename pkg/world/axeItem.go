package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	AxeSprite rl.Texture2D
	AxeSrc    rl.Rectangle
	AxeDest   rl.Rectangle
)

func InitAxe() {
	AxeSprite = rl.LoadTexture("assets/Objects/Items/tools-n-meterial-items.png")
	AxeSrc = rl.NewRectangle(16, 0, 16, 16)
	AxeDest = rl.NewRectangle(0, 0, 16, 16)
}

func DrawAxe() {
	AxeDest.X = 400
	AxeDest.Y = 430
	AxeDest.Width = 16
	AxeDest.Height = 16

	rl.DrawTexturePro(AxeSprite, AxeSrc, AxeDest, rl.NewVector2(0, 0), 0, rl.White)
}

func UnloadAxe() {
	rl.UnloadTexture(AxeSprite)
}

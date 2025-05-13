package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	axeSprite rl.Texture2D
	axeSrc    rl.Rectangle
	axeDest   rl.Rectangle
)

func InitAxe() {
	axeSprite = rl.LoadTexture("assets/Objects/Items/tools-n-meterial-items.png")
	axeSrc = rl.NewRectangle(16, 0, 16, 16)
	axeDest = rl.NewRectangle(0, 0, 16, 16)
}

func DrawAxe() {
	axeDest.X = 400
	axeDest.Y = 430
	axeDest.Width = 16
	axeDest.Height = 16

	rl.DrawTexturePro(axeSprite, axeSrc, axeDest, rl.NewVector2(0, 0), 0, rl.White)
}

func UnloadAxe() {
	rl.UnloadTexture(axeSprite)
}

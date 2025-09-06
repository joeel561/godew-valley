package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	houseFrameCount int
	houseFrame      int
	barnFrameCount  int
	barnFrame       int

	doorsSprite   rl.Texture2D
	HouseDoorSrc  rl.Rectangle
	HouseDoorDest rl.Rectangle
	DoorsMaxFrame int = 5

	BarnDoorSrc  rl.Rectangle
	BarnDoorDest rl.Rectangle

	houseBaseX float32
	barnBaseX  float32
)

func InitDoors() {
	doorsSprite = rl.LoadTexture("assets/Tilesets/building-parts/dooranimationsprites.png")
	HouseDoorSrc = rl.NewRectangle(80, 0, 16, 16)
	HouseDoorDest = rl.NewRectangle(528, 352, 16, 16)

	BarnDoorSrc = rl.NewRectangle(240, 16, 48, 16)
	BarnDoorDest = rl.NewRectangle(886, 448, 48, 16)

	houseBaseX = HouseDoorSrc.X
	barnBaseX = BarnDoorSrc.X
}

func OpenHouseDoor() {
	houseFrameCount++

	if houseFrameCount >= DoorsMaxFrame {
		houseFrameCount = 0
		houseFrame++
	}

	houseFrame = houseFrame % DoorsMaxFrame
	HouseDoorSrc.X = houseBaseX + float32(houseFrame)*HouseDoorSrc.Width
}

func OpenBarnDoor() {
	barnFrameCount++

	if barnFrameCount >= DoorsMaxFrame {
		barnFrameCount = 0
		barnFrame++
	}

	barnFrame = barnFrame % DoorsMaxFrame
	BarnDoorSrc.X = barnBaseX + float32(barnFrame)*BarnDoorSrc.Width
}

func DrawDoors() {
	rl.DrawTexturePro(doorsSprite, HouseDoorSrc, HouseDoorDest, rl.NewVector2(0, 0), 0, rl.White)
	rl.DrawTexturePro(doorsSprite, BarnDoorSrc, BarnDoorDest, rl.NewVector2(0, 0), 0, rl.White)
}

func UnloadDoors() {
	rl.UnloadTexture(doorsSprite)
}

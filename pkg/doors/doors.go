package doors

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	frameCountDoor int
	frameDoor      int = 0
	doorsSprite    rl.Texture2D
	HouseDoorSrc   rl.Rectangle
	HouseDoorDest  rl.Rectangle
	DoorsMaxFrame  int = 5

	BarnDoorSrc  rl.Rectangle
	BarnDoorDest rl.Rectangle
)

func InitDoors() {
	doorsSprite = rl.LoadTexture("assets/Tilesets/building-parts/dooranimationsprites.png")
	HouseDoorSrc = rl.NewRectangle(80, 0, 16, 16)
	HouseDoorDest = rl.NewRectangle(528, 352, 16, 16)

	BarnDoorSrc = rl.NewRectangle(240, 16, 48, 16)
	BarnDoorDest = rl.NewRectangle(886, 448, 48, 16)
}

func OpenHouseDoor() {
	frameCountDoor++

	if frameCountDoor >= DoorsMaxFrame {
		frameCountDoor = 0
		frameDoor++
	}

	HouseDoorSrc.X = 16

	frameDoor = frameDoor % DoorsMaxFrame

}

func OpenBarnDoor() {
	frameCountDoor++

	if frameCountDoor >= DoorsMaxFrame {
		frameCountDoor = 0
		frameDoor++
	}

	BarnDoorSrc.X = 48

	frameDoor = frameDoor % DoorsMaxFrame

}

func DrawDoors() {
	rl.DrawTexturePro(doorsSprite, HouseDoorSrc, HouseDoorDest, rl.NewVector2(0, 0), 0, rl.White)
	rl.DrawTexturePro(doorsSprite, BarnDoorSrc, BarnDoorDest, rl.NewVector2(0, 0), 0, rl.White)
}

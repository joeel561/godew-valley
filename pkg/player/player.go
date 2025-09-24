package player

import (
	"fmt"
	"godew-valley/pkg/userinterface"
	"godew-valley/pkg/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

var (
	playerSprite rl.Texture2D
	oldX, oldY   float32

	playerSrc                                     rl.Rectangle
	PlayerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int
	PlayerHitBox                                  rl.Rectangle
	playerHitBoxYOffset                           float32 = 3
	playerHoe                                     bool
	playerMoveTool                                bool
	playerDirection                               int
	playerAxe                                     bool
	playerWateringCan                             bool
	PlayerToolHitBox                              rl.Rectangle
	playerToolFrame                               int
	wateringSpriteSheet                           rl.Texture2D
	wateringTileSrc                               rl.Rectangle
	WateringTileDest                              rl.Rectangle
	wateringDir                                   int

	frameCount int

	playerSpeed float32 = 1.4

	Cam rl.Camera2D
)

func InitPlayer() {
	playerSprite = rl.LoadTexture("assets/Characters/CharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)

	PlayerDest = rl.NewRectangle(370, 270, 60, 60)
	PlayerHitBox = rl.NewRectangle(0, 0, 6, 6)
	PlayerToolHitBox = rl.NewRectangle(0, 0, 1, 1)

	wateringSpriteSheet = rl.LoadTexture("assets/Characters/watercanframes.png")
	WateringTileDest = rl.NewRectangle(0, 0, 16, 16)
	wateringTileSrc = rl.NewRectangle(0, 0, 16, 16)

	Cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(PlayerDest.X-(PlayerDest.Width/2)), float32(PlayerDest.Y-(PlayerDest.Height/2))), 0, 2)
}

func DrawPlayerTexture() {
	rl.DrawTexturePro(playerSprite, playerSrc, PlayerDest, rl.NewVector2(0, 0), 0, rl.White)
}

func DrawWateringCan() {
	rl.DrawTexturePro(wateringSpriteSheet, wateringTileSrc, WateringTileDest, rl.NewVector2(0, 0), 0, rl.White)
	//fmt.Println("draw watering can", playerDirection, WateringTileDest)
}

func PlayerInput() {
	activeItem := userinterface.PlayerActiveItem

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		if !playerMoveTool {
			playerMoving = true
			playerDir = 5
			playerUp = true
		}
	}

	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		if !playerMoveTool {
			playerMoving = true
			playerDir = 4
			playerDown = true
		}
	}

	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {

		if !playerMoveTool {
			playerMoving = true
			playerDir = 7
			playerLeft = true
		}
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		if !playerMoveTool {
			playerMoving = true
			playerDir = 6
			playerRight = true
		}
	}

	if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
		playerSpeed = 2
	} else {
		playerSpeed = 1.4
	}

	if activeItem.Name == "Hoe" && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		playerHoe = true
		playerMoveTool = true
	}

	if activeItem.Name == "Axe" && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		playerMoveTool = true
		playerAxe = true
	}

	if activeItem.Name == "Watering Can" && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		playerMoveTool = true
		playerWateringCan = true
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		playerMoveTool = false
		playerHoe = false
		playerAxe = false
		playerWateringCan = false
	}
}

func PlayerUseTools() {
	checkCollision := PlayerToolCollision()

	if playerMoveTool {
		if playerDirection == 8 {
			if playerHoe {
				playerDir = 12
			}

			if playerAxe {
				playerDir = 16
			}

			if playerWateringCan {
				playerDir = 20
				wateringDir = 1
			}
		}

		if playerDirection == 9 {
			if playerHoe {
				playerDir = 13
			}

			if playerAxe {
				playerDir = 17
			}

			if playerWateringCan {
				playerDir = 21
				wateringDir = 1
			}
		}

		if playerDirection == 11 {
			if playerHoe {
				playerDir = 15
			}

			if playerAxe {
				playerDir = 19
			}

			if playerWateringCan {
				playerDir = 23
				wateringDir = 1
			}
		}
		if playerDirection == 10 {
			if playerHoe {
				playerDir = 14
			}

			if playerAxe {
				playerDir = 18
			}

			if playerWateringCan {
				playerDir = 22
				wateringDir = 2
			}
		}

		if frameCount%8 == 1 {
			playerFrame++
		}

		wateringTileSrc.X = wateringTileSrc.Width * float32(playerFrame)
		wateringTileSrc.Y = wateringTileSrc.Height * float32(wateringDir)

		fmt.Println(playerFrame)

		if playerFrame >= 4 {
			if playerHoe && checkCollision {
				PlowTile(int(PlayerToolHitBox.X/16), int(PlayerToolHitBox.Y/16))
			}

			if playerFrame >= 8 {
				playerFrame = 0
				playerMoveTool = false
			}
		}
	}
}

func PlayerMoving() {
	oldX, oldY = PlayerDest.X, PlayerDest.Y
	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {
			playerDirection = 9
			PlayerDest.Y -= playerSpeed
			if playerSpeed == 2 {
				playerDir = 9
			}
		}
		if playerDown {
			playerDirection = 8
			PlayerDest.Y += playerSpeed
			if playerSpeed == 2 {
				playerDir = 8
			}
		}
		if playerLeft {
			playerDirection = 11
			PlayerDest.X -= playerSpeed
			if playerSpeed == 2 {
				playerDir = 11
			}

		}
		if playerRight {
			playerDirection = 10
			PlayerDest.X += playerSpeed
			if playerSpeed == 2 {
				playerDir = 10
			}
		}

		if frameCount%8 == 1 {
			playerFrame++
		}

		PlayerOpenHouseDoor()
		PlayerOpenBarnDoor()
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	frameCount++
	if playerFrame >= 8 {
		playerFrame = 0
	}

	if !playerMoveTool && !playerMoving && playerFrame > 1 {
		playerFrame = 0
		playerDir = playerDirection
	}

	playerSrc.Y = playerSrc.Height * float32(playerDir)
	playerSrc.X = playerSrc.Width * float32(playerFrame)

	PlayerHitBox.X = PlayerDest.X + (PlayerDest.Width / 2) - PlayerHitBox.Width/2
	PlayerHitBox.Y = PlayerDest.Y + (PlayerDest.Height / 2) + playerHitBoxYOffset

	PlayerToolHitBox.X = PlayerDest.X + (PlayerDest.Width / 2) - PlayerToolHitBox.Width/2
	PlayerToolHitBox.Y = PlayerDest.Y + (PlayerDest.Height / 2) + 2

	WateringTileDest.X = PlayerDest.X + (PlayerDest.Width / 2) - WateringTileDest.Width/2
	WateringTileDest.Y = PlayerDest.Y + (PlayerDest.Height / 2) + 2

	switch playerDirection {
	case 8: // Down
		PlayerToolHitBox.Y = PlayerToolHitBox.Y + float32(world.WorldMap.TileSize)
	case 9: // Up
		PlayerToolHitBox.Y = PlayerToolHitBox.Y - float32(world.WorldMap.TileSize)
	case 10: // Right
		PlayerToolHitBox.X = PlayerToolHitBox.X + float32(world.WorldMap.TileSize)

	case 11: // Left
		PlayerToolHitBox.X = PlayerToolHitBox.X - float32(world.WorldMap.TileSize)
	}

	//	PlayerCollision(world.WaterTiles)
	PlayerCollision(world.Structures)
	PlayerCollision(world.Furniture)

	Cam.Target = rl.NewVector2(float32(PlayerDest.X-(PlayerDest.Width/2)), float32(PlayerDest.Y-(PlayerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
}

func PlayerCollision(tiles []world.Tile) {
	var jsonMap = world.WorldMap

	for i := 0; i < len(tiles); i++ {
		if PlayerHitBox.X < float32(tiles[i].X*jsonMap.TileSize+jsonMap.TileSize) &&
			PlayerHitBox.X+PlayerHitBox.Width > float32(tiles[i].X*jsonMap.TileSize) &&
			PlayerHitBox.Y < float32(tiles[i].Y*jsonMap.TileSize+jsonMap.TileSize) &&
			PlayerHitBox.Y+PlayerHitBox.Height > float32(tiles[i].Y*jsonMap.TileSize) {

			PlayerDest.X = oldX
			PlayerDest.Y = oldY
		}
	}
}

func PlayerToolCollision() bool {
	var jsonMap = world.WorldMap

	var tiles = world.Dirt

	for i := 0; i < len(tiles); i++ {
		if PlayerToolHitBox.X < float32(tiles[i].X*jsonMap.TileSize+jsonMap.TileSize) &&
			PlayerToolHitBox.X+PlayerToolHitBox.Width > float32(tiles[i].X*jsonMap.TileSize) &&
			PlayerToolHitBox.Y < float32(tiles[i].Y*jsonMap.TileSize+jsonMap.TileSize) &&
			PlayerToolHitBox.Y+PlayerToolHitBox.Height > float32(tiles[i].Y*jsonMap.TileSize) {

			return true
		}
	}

	return false
}

func PlayerOpenHouseDoor() {
	world.HouseDoorSrc.X = 80

	if PlayerHitBox.X < float32(world.HouseDoorDest.X+world.HouseDoorDest.Width) &&
		PlayerHitBox.X+PlayerHitBox.Width > float32(world.HouseDoorDest.X) &&
		PlayerHitBox.Y < float32(world.HouseDoorDest.Y+world.HouseDoorDest.Height) &&
		PlayerHitBox.Y+PlayerHitBox.Height > float32(world.HouseDoorDest.Y) {

		world.OpenHouseDoor()
	}
}

func PlayerOpenBarnDoor() {
	world.BarnDoorSrc.X = 240

	if PlayerHitBox.X < float32(world.BarnDoorDest.X+world.BarnDoorDest.Width) &&
		PlayerHitBox.X+PlayerHitBox.Width > float32(world.BarnDoorDest.X) &&
		PlayerHitBox.Y < float32(world.BarnDoorDest.Y+world.BarnDoorDest.Height) &&
		PlayerHitBox.Y+PlayerHitBox.Height > float32(world.BarnDoorDest.Y) {

		world.OpenBarnDoor()
	}
}

func UnloadPlayerTexture() {
	rl.UnloadTexture(playerSprite)
}

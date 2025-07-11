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

	frameCount int

	playerSpeed float32 = 1.4

	Cam rl.Camera2D
)

func InitPlayer() {
	playerSprite = rl.LoadTexture("assets/Characters/CharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)

	PlayerDest = rl.NewRectangle(370, 270, 60, 60)
	PlayerHitBox = rl.NewRectangle(0, 0, 10, 10)

	Cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(PlayerDest.X-(PlayerDest.Width/2)), float32(PlayerDest.Y-(PlayerDest.Height/2))), 0, 2)
}

func DrawPlayerTexture() {
	rl.DrawTexturePro(playerSprite, playerSrc, PlayerDest, rl.NewVector2(0, 0), 0, rl.White)
}

func PlayerInput() {

	activeItem := userinterface.PlayerActiveItem

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

	if activeItem.Name == "Hoe" && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		playerHoe = true
		playerMoveTool = true

	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		playerMoveTool = false
		playerHoe = false
	}
}

func PlayerUseTools() {
	if playerMoveTool {
		fmt.Println("playermovetool", playerDirection)

		if playerDirection == 8 {
			if playerHoe {
				fmt.Println(playerDir, "playerHoe")
				playerDir = 12
			}
		}

		if playerDirection == 9 {
			if playerHoe {
				playerDir = 13
			}
		}

		if playerDirection == 11 {
			if playerHoe {
				playerDir = 15
			}

		}
		if playerDirection == 10 {
			if playerHoe {
				playerDir = 14
			}
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		playerMoveTool = false
		playerHoe = false
	}

	playerMoveTool = false
	playerHoe = false
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

	/* 	if !playerMoveTool {
		playerDir = playerDirection // Reset player direction to the current direction
	} */

	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.Y = playerSrc.Height * float32(playerDir)
	playerSrc.X = playerSrc.Width * float32(playerFrame)

	PlayerHitBox.X = PlayerDest.X + (PlayerDest.Width / 2) - PlayerHitBox.Width/2
	PlayerHitBox.Y = PlayerDest.Y + (PlayerDest.Height / 2) + playerHitBoxYOffset

	PlayerCollision(world.WaterTiles)
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

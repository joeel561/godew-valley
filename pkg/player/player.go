package player

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	playerSprite rl.Texture2D
	oldX, oldY   float32

	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int
	playerHitBox                                  rl.Rectangle
	playerHitBoxYOffset                           float32 = 3

	frameCount int

	playerSpeed float32 = 1.4

	cam rl.Camera2D
)

func PlayerInput() {
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
}

func PlayerMoving() {
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

	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
}

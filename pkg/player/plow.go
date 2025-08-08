package player

import (
	"fmt"
	"godew-valley/pkg/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var plowGrid [][]world.Tile

var (
	tileSrc         rl.Rectangle
	tileDest        rl.Rectangle
	dirtSpriteSheet rl.Texture2D
)

func InitPlowGrid() {
	tileSrc = rl.NewRectangle(0, 80, 16, 16)
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	plowGrid = make([][]world.Tile, world.WorldMap.MapHeight)
	dirtSpriteSheet = rl.LoadTexture("assets/Tilesets/ground-tiles/Old-tiles/Tilled_Dirt.png")

	for i := range plowGrid {
		plowGrid[i] = make([]world.Tile, world.WorldMap.MapWidth)
	}

	for _, tile := range world.GroundTiles {
		plowGrid[tile.Y][tile.X] = tile
		fmt.Println("Plowed tile at:", tile.X, tile.Y)
	}
}

func PlowTile(x, y int) {
	tile := &plowGrid[y][x]

	if tile.State == world.TileDefault {
		tile.State = world.TilePlowed
	}
}

func DrawPlowGrid() {
	for y := 0; y < world.WorldMap.MapHeight; y++ {
		for x := 0; x < world.WorldMap.MapWidth; x++ {
			tile := plowGrid[y][x]
			if tile.State == world.TilePlowed {

				tileSrc = rl.NewRectangle(0, 80, 16, 16) // Assuming plowed tile is at (32, 128) in the spritesheet
				tileDest = rl.NewRectangle(float32(x*world.WorldMap.TileSize), float32(y*world.WorldMap.TileSize), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize))
				fmt.Println(tileDest, "Plowed tile at:")
				rl.DrawTexturePro(dirtSpriteSheet, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
			} else {
				rl.DrawTexturePro(world.SpritesheetMap, tileSrc, rl.NewRectangle(float32(tile.X), float32(tile.Y), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize)), rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

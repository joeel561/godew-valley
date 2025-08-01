package player

import (
	"fmt"
	"godew-valley/pkg/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var plowGrid [][]world.Tile

var (
	tileSrc  rl.Rectangle
	tileDest rl.Rectangle
)

func InitPlowGrid() {
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	plowGrid = make([][]world.Tile, len(world.GroundTiles))

	fmt.Println(len(plowGrid), world.WorldMap.MapHeight, world.WorldMap.MapWidth)

	for y := range plowGrid {
		plowGrid[y] = make([]world.Tile, world.WorldMap.MapWidth)
		for x := range plowGrid[y] {
			plowGrid[y][x] = world.GroundTiles[y*world.WorldMap.MapWidth+x]
		}
	}

	fmt.Println(plowGrid)
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
				fmt.Println("Plowed tile at:", x, y)
				rl.DrawTexturePro(world.SpritesheetMap, tileSrc, rl.NewRectangle(float32(tile.X), float32(tile.Y), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize)), rl.NewVector2(0, 0), 0, rl.White)
			} else {
				rl.DrawTexturePro(world.SpritesheetMap, tileSrc, rl.NewRectangle(float32(tile.X), float32(tile.Y), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize)), rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

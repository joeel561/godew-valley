package player

import (
	"godew-valley/pkg/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var plowGrid [][]world.Tile

var (
	tileSrc         rl.Rectangle
	tileDest        rl.Rectangle
	dirtSpriteSheet rl.Texture2D
	blobMapping     map[int]rl.Rectangle
)

func InitPlowGrid() {
	tileSrc = rl.NewRectangle(0, 80, 16, 16)
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	plowGrid = make([][]world.Tile, world.WorldMap.MapHeight)
	dirtSpriteSheet = rl.LoadTexture("assets/Tilesets/ground-tiles/Old-tiles/Tilled_Dirt_Wide.png")

	for i := range plowGrid {
		plowGrid[i] = make([]world.Tile, world.WorldMap.MapWidth)
	}

	for _, tile := range world.GroundTiles {
		plowGrid[tile.Y][tile.X] = tile
	}
}

func InitBlobMapping() {
	blobMapping = map[int]rl.Rectangle{
		// Einzelner Block
		0: {48, 48, 16, 16},

		// Voller Block (alle Nachbarn)
		255: {16, 16, 16, 16},

		// Gerade Linien
		1 | 16: {32, 0, 16, 16}, // N + S (vertikale Linie)
		4 | 64: {48, 0, 16, 16}, // E + W (horizontale Linie)

		// Ecken
		1 | 4:   {64, 0, 16, 16},  // N + E
		4 | 16:  {80, 0, 16, 16},  // E + S
		16 | 64: {96, 0, 16, 16},  // S + W
		64 | 1:  {112, 0, 16, 16}, // W + N

		// T-Stücke
		1 | 4 | 64:  {0, 16, 16, 16},  // N, E, W
		4 | 16 | 1:  {16, 16, 16, 16}, // E, S, N
		16 | 64 | 4: {32, 16, 16, 16}, // S, W, E
		64 | 1 | 16: {48, 16, 16, 16}, // W, N, S

		// Kreuzung
		1 | 4 | 16 | 64: {64, 16, 16, 16},

		// Endstücke (nur eine Seite)
		1:  {80, 16, 16, 16},  // N
		4:  {96, 16, 16, 16},  // E
		16: {112, 16, 16, 16}, // S
		64: {0, 32, 16, 16},   // W

		// Ecken außen (Diagonalen)
		1 | 4 | 16 | 64 | 128: {16, 32, 16, 16}, // NW Ecke gefüllt
		1 | 4 | 16 | 64 | 2:   {32, 32, 16, 16}, // NE Ecke gefüllt
		1 | 4 | 16 | 64 | 8:   {48, 32, 16, 16}, // SE Ecke gefüllt
		1 | 4 | 16 | 64 | 32:  {64, 32, 16, 16}, // SW Ecke gefüllt
	}
}

func PlowTile(x, y int) {
	tile := &plowGrid[y][x]

	if tile.State == world.TileDefault {
		tile.State = world.TilePlowed
	}
}

func GetBlobMask(plowGrid [][]world.Tile, x, y int) int {
	mask := 0

	check := func(dx, dy int, bit int) {
		nx, ny := x+dx, y+dy
		if ny >= 0 && ny < len(plowGrid) && nx >= 0 && nx < len(plowGrid[0]) {
			if plowGrid[ny][nx].State == world.TilePlowed {
				mask |= bit
			}
		}
	}

	check(-1, 0, 1)   // left
	check(1, 0, 2)    // right
	check(0, -1, 4)   // up
	check(0, 1, 8)    // down
	check(-1, -1, 16) // top-left
	check(1, -1, 32)  // top-right
	check(-1, 1, 64)  // bottom-left
	check(1, 1, 128)  // bottom-right
	check(0, 2, 256)  // down

	return mask
}

func DrawPlowGrid() {
	for y := 0; y < world.WorldMap.MapHeight; y++ {
		for x := 0; x < world.WorldMap.MapWidth; x++ {
			tile := plowGrid[y][x]
			if tile.State == world.TilePlowed {
				mask := GetBlobMask(plowGrid, x, y)
				tileSrc = blobMapping[mask]
				tileDest = rl.NewRectangle(float32(x*world.WorldMap.TileSize), float32(y*world.WorldMap.TileSize), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize))
				rl.DrawTexturePro(dirtSpriteSheet, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
			} else {
				rl.DrawTexturePro(world.SpritesheetMap, tileSrc, rl.NewRectangle(float32(tile.X), float32(tile.Y), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize)), rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

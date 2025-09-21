package player

import (
	"godew-valley/pkg/world"
	"math/bits"

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

	for _, tile := range world.Dirt {
		plowGrid[tile.Y][tile.X] = tile
	}
}

func InitBlobMapping() {
	blobMapping = map[int]rl.Rectangle{
		// Einzelner Block

		0:      {48, 48, 16, 16}, // kein nachbar
		4:      {0, 48, 16, 16},  // rechter nachbar
		64:     {32, 48, 16, 16}, // linker nachbar
		4 | 64: {16, 48, 16, 16}, // linker und rechter nachbar
		128:    {0, 0, 16, 16},   // obern links nachbar
		1:      {48, 32, 16, 16}, // obern mitte nachbar
		2:      {0, 32, 16, 16},  // obern rechts nachbar
		8:      {32, 32, 16, 16}, // unten rechts nachbar
		16:     {48, 0, 16, 16},  // unten mitte nachbar
		1 | 16: {48, 16, 16, 16}, // N + S (vertikale Linie)

		4 | 8 | 16:                         {0, 0, 16, 16},
		4 | 8 | 16 | 32 | 64:               {16, 0, 16, 16},
		64 | 32 | 16:                       {32, 0, 16, 16},
		128 | 1 | 64 | 32 | 16:             {32, 16, 16, 16},
		128 | 1 | 64:                       {32, 32, 16, 16},
		128 | 1 | 2 | 64 | 4:               {16, 32, 16, 16},
		1 | 2 | 4:                          {0, 32, 16, 16},
		1 | 2 | 4 | 8 | 16:                 {0, 16, 16, 16},
		128 | 1 | 2 | 64 | 32 | 16 | 8 | 4: {16, 16, 16, 16},

		4 | 16:                             {64, 0, 16, 16},
		64 | 4 | 32 | 16:                   {80, 0, 16, 16},
		64 | 4 | 16 | 8:                    {96, 0, 16, 16},
		64 | 16:                            {112, 0, 16, 16},
		64 | 4 | 16:                        {128, 0, 16, 16},
		128 | 1 | 64 | 4 | 16 | 8:          {144, 0, 16, 16},
		1 | 2 | 4 | 16:                     {64, 16, 16, 16},
		128 | 1 | 2 | 64 | 32 | 16:         {80, 16, 16, 16},
		128 | 1 | 2 | 64 | 4 | 16 | 8:      {96, 16, 16, 16},
		128 | 1 | 64 | 16:                  {112, 16, 16, 16},
		128 | 1 | 2 | 64 | 4 | 16:          {128, 16, 16, 16},
		1 | 2 | 64 | 4 | 32 | 16:           {144, 16, 16, 16},
		1 | 4 | 16 | 8:                     {64, 32, 16, 16},
		128 | 1 | 4 | 64 | 4 | 32 | 16 | 8: {80, 32, 16, 16},
		1 | 2 | 4 | 64 | 32 | 16 | 8:       {96, 32, 16, 16},
		1 | 64 | 32 | 16:                   {112, 32, 16, 16},
		1 | 64 | 4 | 32 | 16 | 8:           {128, 32, 16, 16},
		1 | 64 | 4 | 16 | 8:                {144, 32, 16, 16},
		1 | 64 | 4 | 32 | 16:               {160, 32, 16, 16},
		1 | 4:                              {64, 48, 16, 16},
		128 | 1 | 4 | 64:                   {80, 48, 16, 16},
		1 | 2 | 64 | 4:                     {96, 48, 16, 16},
		1 | 64:                             {112, 48, 16, 16},
		1 | 64 | 4:                         {128, 48, 16, 16},
		1 | 2 | 4 | 64 | 16:                {144, 48, 16, 16},
		128 | 1 | 4 | 64 | 16:              {160, 48, 16, 16},
		1 | 4 | 16:                         {64, 64, 16, 16},
		128 | 1 | 4 | 64 | 32 | 16:         {80, 64, 16, 16},
		1 | 2 | 64 | 4 | 16 | 8:            {96, 64, 16, 16},
		1 | 64 | 16:                        {112, 64, 16, 16},
		1 | 4 | 16 | 64:                    {128, 64, 16, 16},
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

	check(0, -1, 1)    //Norden      (Oben)
	check(1, -1, 2)    //Nord-Ost    (Oben-Rechts)
	check(1, 0, 4)     //Osten       (Rechts)
	check(1, 1, 8)     //Süd-Osten   (Unten-Rechts)
	check(0, 1, 16)    //Süden       (Unten)
	check(-1, 1, 32)   //Süd-Westen  (Unten-Links)
	check(-1, 0, 64)   //Westen      (Links)
	check(-1, -1, 128) //Nord-Westen (Links-Oben)

	return mask
}

func getClosestPlowedTile(mask int, blobMapping map[int]rl.Rectangle) rl.Rectangle {
	if rect, exists := blobMapping[mask]; exists {
		return rect
	}

	bestKey := 0
	bestScore := 0
	found := false

	// Suche nach dem ähnlichsten Schlüssel
	for key := range blobMapping {
		score := 8 * bits.OnesCount(uint(key)&uint(mask)&uint(0b01010101))
		score += 4 * bits.OnesCount(uint(key)&uint(mask)&uint(0b10101010))
		score -= 6 * bits.OnesCount(uint(key)&uint(^mask)&uint(0b01010101))
		score -= 2 * bits.OnesCount(uint(key)&uint(^mask)&uint(0b10101010))
		if score > bestScore {
			bestScore = score
			bestKey = key
			found = true
		}
	}

	if found {
		return blobMapping[bestKey]
	}

	return rl.NewRectangle(16, 16, 16, 16)
}

func DrawPlowGrid() {
	for y := 0; y < world.WorldMap.MapHeight; y++ {
		for x := 0; x < world.WorldMap.MapWidth; x++ {
			tile := plowGrid[y][x]
			if tile.State == world.TilePlowed {
				mask := GetBlobMask(plowGrid, x, y)
				if _, exists := blobMapping[mask]; !exists {
					mask = GetBlobMask(plowGrid, x, y)
					blobMapping[mask] = getClosestPlowedTile(mask, blobMapping)
				}

				tileSrc = blobMapping[mask]
				tileDest = rl.NewRectangle(float32(x*world.WorldMap.TileSize), float32(y*world.WorldMap.TileSize), float32(world.WorldMap.TileSize), float32(world.WorldMap.TileSize))
				rl.DrawTexturePro(dirtSpriteSheet, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

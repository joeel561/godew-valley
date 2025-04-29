package itembar

import (
	"fmt"
	"godew-valley/pkg/world"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	tileDest rl.Rectangle
	tileSrc  rl.Rectangle
	tex      rl.Texture2D
)

func DrawItemBar() {
	renderItemBarLayer(world.ItemBarTiles)
}

func renderItemBarLayer(Layer []world.Tile) {
	for i := 0; i < len(Layer); i++ {
		s, _ := strconv.ParseInt(Layer[i].Id, 10, 64)
		tileId := int(s)
		tex = world.SpritesheetMap

		texColumns := tex.Width / int32(world.WorldMap.TileSize)
		tileSrc.X = float32(world.WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(world.WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(Layer[i].X * world.WorldMap.TileSize)
		tileDest.Y = float32(Layer[i].Y * world.WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)

		fmt.Println("Drawing item bar tile ID:", tileId, "at position:", tileDest.X, tileDest.Y)
	}
}

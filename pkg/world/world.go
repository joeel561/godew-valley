package world

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	tileDest        rl.Rectangle
	tileSrc         rl.Rectangle
	WorldMap        JsonMap
	SpritesheetMap  rl.Texture2D
	tex             rl.Texture2D
	doorSprite      rl.Texture2D
	DoorSrc         rl.Rectangle
	DoorDest        rl.Rectangle
	WaterTiles      []Tile
	Structures      []Tile
	Furniture       []Tile
	WalkableWater   []Tile
	Paths           []Tile
	ItemBarTiles    []Tile
	BackgroundTiles []Tile
)

type JsonMap struct {
	Layers    []Layer `json:"layers"`
	MapHeight int     `json:"mapHeight"`
	MapWidth  int     `json:"mapWidth"`
	TileSize  int     `json:"tileSize"`
}

type Layer struct {
	Name  string `json:"name"`
	Tiles []Tile `json:"tiles"`
}

type Tile struct {
	Id string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func LoadMap(mapFile string) {
	file, err := os.Open(mapFile)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &WorldMap)

	// Cache layer
	BackgroundTiles = nil
	WaterTiles = nil
	Structures = nil
	Furniture = nil
	WalkableWater = nil
	Paths = nil

	for i := 0; i < len(WorldMap.Layers); i++ {
		switch WorldMap.Layers[i].Name {
		case "Background":
			BackgroundTiles = WorldMap.Layers[i].Tiles
		case "Water":
			WaterTiles = WorldMap.Layers[i].Tiles
		case "Structures":
			Structures = WorldMap.Layers[i].Tiles
		case "Furniture":
			Furniture = WorldMap.Layers[i].Tiles
		case "WalkableWater":
			WalkableWater = WorldMap.Layers[i].Tiles
		case "Paths":
			Paths = WorldMap.Layers[i].Tiles
		}
	}
}

func InitWorld() {
	SpritesheetMap = rl.LoadTexture("assets/spritesheet.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
}

func DrawWorld(cam rl.Camera2D, screenW, screenH int) {
	RenderLayer(WaterTiles, cam, screenW, screenH)
	RenderLayer(WalkableWater, cam, screenW, screenH)
	RenderLayer(BackgroundTiles, cam, screenW, screenH)
	RenderLayer(Structures, cam, screenW, screenH)
	RenderLayer(Paths, cam, screenW, screenH)
	RenderLayer(Furniture, cam, screenW, screenH)
}

func RenderLayer(layer []Tile, cam rl.Camera2D, screenW, screenH int) {
	if len(layer) == 0 {
		return
	}

	// Tile-Culling: compute view bounds in tile coordinates
	halfW := float32(screenW) / (2.0 * cam.Zoom)
	halfH := float32(screenH) / (2.0 * cam.Zoom)
	minX := cam.Target.X - halfW
	maxX := cam.Target.X + halfW
	minY := cam.Target.Y - halfH
	maxY := cam.Target.Y + halfH

	tsize := float32(WorldMap.TileSize)
	minTileX := int(minX/tsize) - 1
	maxTileX := int(maxX/tsize) + 1
	minTileY := int(minY/tsize) - 1
	maxTileY := int(maxY/tsize) + 1

	if minTileX < 0 {
		minTileX = 0
	}
	if minTileY < 0 {
		minTileY = 0
	}

	tex = SpritesheetMap
	texColumns := tex.Width / int32(WorldMap.TileSize)

	for i := 0; i < len(layer); i++ {
		tx := layer[i].X
		ty := layer[i].Y

		if tx < minTileX || tx > maxTileX || ty < minTileY || ty > maxTileY {
			continue
		}

		s, _ := strconv.ParseInt(layer[i].Id, 10, 64)
		tileId := int(s)

		tileSrc.X = float32(WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(tx * WorldMap.TileSize)
		tileDest.Y = float32(ty * WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func UnloadWorldTexture() {
	rl.UnloadTexture(SpritesheetMap)
}

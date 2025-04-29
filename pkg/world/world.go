package world

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	tileDest       rl.Rectangle
	tileSrc        rl.Rectangle
	WorldMap       JsonMap
	SpritesheetMap rl.Texture2D
	tex            rl.Texture2D
	doorSprite     rl.Texture2D
	DoorSrc        rl.Rectangle
	DoorDest       rl.Rectangle
	WaterTiles     []Tile
	Structures     []Tile
	Furniture      []Tile
	WalkableWater  []Tile
	Paths          []Tile
	ItemBarTiles   []Tile
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
}

func InitWorld() {
	SpritesheetMap = rl.LoadTexture("assets/spritesheet.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
}

func DrawWorld() {
	var groundTiles []Tile

	for i := 0; i < len(WorldMap.Layers); i++ {
		if WorldMap.Layers[i].Name == "Background" {
			groundTiles = WorldMap.Layers[i].Tiles
		}

		if WorldMap.Layers[i].Name == "Water" {
			WaterTiles = WorldMap.Layers[i].Tiles
		}

		if WorldMap.Layers[i].Name == "Structures" {
			Structures = WorldMap.Layers[i].Tiles
		}

		if WorldMap.Layers[i].Name == "Furniture" {
			Furniture = WorldMap.Layers[i].Tiles
		}

		if WorldMap.Layers[i].Name == "WalkableWater" {
			WalkableWater = WorldMap.Layers[i].Tiles
		}

		if WorldMap.Layers[i].Name == "Paths" {
			Paths = WorldMap.Layers[i].Tiles
		}

		if WorldMap.Layers[i].Name == "ItemBar" {
			ItemBarTiles = WorldMap.Layers[i].Tiles
		}
	}

	rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)

	RenderLayer(WaterTiles)
	RenderLayer(WalkableWater)
	RenderLayer(groundTiles)
	RenderLayer(Structures)
	RenderLayer(Paths)
	RenderLayer(Furniture)
}

func RenderLayer(Layer []Tile) {
	for i := 0; i < len(Layer); i++ {
		s, _ := strconv.ParseInt(Layer[i].Id, 10, 64)
		tileId := int(s)
		tex = SpritesheetMap

		texColumns := tex.Width / int32(WorldMap.TileSize)
		tileSrc.X = float32(WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(Layer[i].X * WorldMap.TileSize)
		tileDest.Y = float32(Layer[i].Y * WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func UnloadWorldTexture() {
	rl.UnloadTexture(SpritesheetMap)
}

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
	spritesheetMap rl.Texture2D
	tex            rl.Texture2D
	WaterTiles     []Tile
	Structures     []Tile
	Furniture      []Tile
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
	spritesheetMap = rl.LoadTexture("assets/spritesheet.png")
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
	}

	for i := 0; i < len(WaterTiles); i++ {
		s, _ := strconv.ParseInt(WaterTiles[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(WorldMap.TileSize)
		tileSrc.X = float32(WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(WaterTiles[i].X * WorldMap.TileSize)
		tileDest.Y = float32(WaterTiles[i].Y * WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	for i := 0; i < len(groundTiles); i++ {
		s, _ := strconv.ParseInt(groundTiles[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(WorldMap.TileSize)
		tileSrc.X = float32(WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(groundTiles[i].X * WorldMap.TileSize)
		tileDest.Y = float32(groundTiles[i].Y * WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	for i := 0; i < len(Structures); i++ {
		s, _ := strconv.ParseInt(Structures[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(WorldMap.TileSize)
		tileSrc.X = float32(WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(Structures[i].X * WorldMap.TileSize)
		tileDest.Y = float32(Structures[i].Y * WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	for i := 0; i < len(Furniture); i++ {
		s, _ := strconv.ParseInt(Furniture[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(WorldMap.TileSize)
		tileSrc.X = float32(WorldMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(WorldMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(Furniture[i].X * WorldMap.TileSize)
		tileDest.Y = float32(Furniture[i].Y * WorldMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func UnloadWorldTexture() {
	rl.UnloadTexture(spritesheetMap)
}

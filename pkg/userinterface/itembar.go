package userinterface

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

var (
	UserInterface jsonMap
	spritesheet   rl.Texture2D
	tileDest      rl.Rectangle
	tileSrc       rl.Rectangle
	tex           rl.Texture2D
	texColumns    int32
	buttonSprite  rl.Texture2D
)

type jsonMap struct {
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

func InitUserInterface() {
	spritesheet = rl.LoadTexture("assets/userinterface/userinterfacespritesheet.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	buttonSprite = rl.LoadTexture("assets/userinterface/squarebutton.png")
}

func LoadUserInterfaceMap(mapFile string) {
	file, err := os.Open(mapFile)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &UserInterface)
}

func DrawUserInterface() {
	var itembar []Tile

	for i := 0; i < len(UserInterface.Layers); i++ {
		if UserInterface.Layers[i].Name == "itembar" {
			itembar = UserInterface.Layers[i].Tiles
		}
	}

	renderItemBarLayer(itembar)

	DrawItemBar()
}

type Item struct {
	Name string
	Icon rl.Texture2D
}

type Hotbar struct {
	Slots         []Item
	SelectedIndex int
}

func DrawItemBar() {
	hotbar := Hotbar{
		Slots:         make([]Item, 5),
		SelectedIndex: 0,
	}

	buttonSrc := rl.NewRectangle(0, 0, 32, 32)

	buttonDest := rl.NewRectangle(0, 0, 32, 32)

	for i := 0; i < len(hotbar.Slots); i++ {
		x := int32(screenWidth/2 - 160 + (i * 32))
		y := int32(screenHeight - UserInterface.MapHeight*UserInterface.TileSize + 11)
		buttonDest.X = float32(x)
		buttonDest.Y = float32(y)

		if i == hotbar.SelectedIndex {
			rl.DrawTexturePro(buttonSprite, buttonSrc, buttonDest, rl.NewVector2(0, 0), 0, rl.White)
		} else {
			rl.DrawTexturePro(buttonSprite, buttonSrc, buttonDest, rl.NewVector2(0, 0), 0, rl.White)
		}
	}
}

func renderItemBarLayer(Layer []Tile) {
	for i := 0; i < len(Layer); i++ {
		s, _ := strconv.ParseInt(Layer[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheet

		texColumns := tex.Width / int32(UserInterface.TileSize)
		tileSrc.X = float32(UserInterface.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(UserInterface.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(Layer[i].X*UserInterface.TileSize) + float32(screenWidth/2-UserInterface.MapWidth*UserInterface.TileSize/2)
		tileDest.Y = float32(Layer[i].Y*UserInterface.TileSize) + float32(screenHeight-UserInterface.MapHeight*UserInterface.TileSize-5)
		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 1, rl.White)
	}
}

func UnloadUserInterface() {
	rl.UnloadTexture(spritesheet)
}

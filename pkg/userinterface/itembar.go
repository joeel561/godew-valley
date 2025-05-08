package userinterface

import (
	"encoding/json"
	"fmt"
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
	hotbar        Hotbar
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

type Item struct {
	Name string
	Icon rl.Texture2D
}

type Hotbar struct {
	Slots         []Item
	SelectedIndex int
}

func InitUserInterface() {
	spritesheet = rl.LoadTexture("assets/userinterface/userinterfacespritesheet.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	buttonSprite = rl.LoadTexture("assets/userinterface/Inventory_Spritesheet.png")

	hotbar = Hotbar{
		Slots:         make([]Item, 10),
		SelectedIndex: 0,
	}
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

func DrawItemBar() {

	buttonSrc := rl.NewRectangle(224, 112, 48, 48)
	buttonDest := rl.NewRectangle(0, 0, 48, 48)

	buttonSelected := rl.NewRectangle(272, 112, 48, 48)
	buttonSelectedDest := rl.NewRectangle(0, 0, 48, 48)

	buttonActive := rl.NewRectangle(224, 0, 48, 48)
	buttonActiveDest := rl.NewRectangle(0, 0, 48, 48)

	for i := 0; i < len(hotbar.Slots); i++ {
		x := int32(screenWidth/2 - 182 + (i * 35))
		y := int32(screenHeight - UserInterface.MapHeight*UserInterface.TileSize + 2)
		buttonDest.X = float32(x)
		buttonDest.Y = float32(y)

		buttonSelectedDest.X = float32(x)
		buttonSelectedDest.Y = float32(y)

		buttonActiveDest.X = float32(x)
		buttonActiveDest.Y = float32(y)

		if i == hotbar.SelectedIndex {
			rl.DrawTexturePro(buttonSprite, buttonSelected, buttonSelectedDest, rl.NewVector2(0, 0), 0, rl.White)
			rl.DrawTexturePro(buttonSprite, buttonActive, buttonActiveDest, rl.NewVector2(0, 0), 0, rl.White)
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

func ItemBarInput() {
	key := rl.GetKeyPressed()
	if key >= rl.KeyOne && key <= rl.KeyNine {
		hotbar.SelectedIndex = int(key) - rl.KeyOne
	}

	if key == rl.KeyZero {
		hotbar.SelectedIndex = 9
	}

	fmt.Println("Selected Index:", hotbar.SelectedIndex)
}

func UnloadUserInterface() {
	rl.UnloadTexture(spritesheet)
}

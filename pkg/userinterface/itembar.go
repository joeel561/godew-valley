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
	UserInterface   jsonMap
	spritesheet     rl.Texture2D
	tileDest        rl.Rectangle
	tileSrc         rl.Rectangle
	tex             rl.Texture2D
	texColumns      int32
	buttonSprite    rl.Texture2D
	PlayerHotbar    Hotbar
	openInventory   bool
	PlayerInventory Hotbar
	Dragging        DraggedItem
	cachedItem      Item
	maxQuantity     int = 64
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
	Name     string
	Icon     rl.Texture2D
	IconSrc  rl.Rectangle
	Quantity int
	X        int32
	Y        int32
	Active   bool
}

type Hotbar struct {
	Slots         []Item
	SelectedIndex int
}

type DraggedItem struct {
	Item       Item
	Source     int
	SourceType string // "hotbar" or "inventory"
	Position   rl.Vector2
	Drag       bool
}

func InitUserInterface() {
	spritesheet = rl.LoadTexture("assets/userinterface/userinterfacespritesheet.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	buttonSprite = rl.LoadTexture("assets/userinterface/Inventory_Spritesheet.png")

	PlayerHotbar = Hotbar{
		Slots:         make([]Item, 10),
		SelectedIndex: 0,
	}

	PlayerInventory = Hotbar{
		Slots:         make([]Item, 27),
		SelectedIndex: 0,
	}

	Dragging = DraggedItem{
		Item:       Item{},
		Source:     -1,
		SourceType: "",
		Position:   rl.NewVector2(0, 0),
		Drag:       false,
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
	var inventory []Tile

	for i := 0; i < len(UserInterface.Layers); i++ {
		if UserInterface.Layers[i].Name == "itembar" {
			itembar = UserInterface.Layers[i].Tiles
		}

		if UserInterface.Layers[i].Name == "inventory" {
			inventory = UserInterface.Layers[i].Tiles
		}
	}

	renderItemBarLayer(itembar)

	if openInventory {
		renderItemBarLayer(inventory)
		DrawInventorySlots()
	}

	DrawItemBar()
}

func DrawItemBar() {
	buttonSrc := rl.NewRectangle(224, 112, 48, 48)
	buttonDest := rl.NewRectangle(0, 0, 48, 48)

	buttonSelected := rl.NewRectangle(272, 112, 48, 48)
	buttonSelectedDest := rl.NewRectangle(0, 0, 48, 48)

	buttonActive := rl.NewRectangle(224, 0, 48, 48)
	buttonActiveDest := rl.NewRectangle(0, 0, 48, 48)

	mousePosition := rl.GetMousePosition()
	isValidMousePosition := mousePosition.X >= 0 && mousePosition.Y >= 0 && mousePosition.X <= screenWidth && mousePosition.Y <= screenHeight

	for i, item := range PlayerHotbar.Slots {
		if i >= len(PlayerHotbar.Slots) {
			break
		}
		PlayerHotbar.Slots[i].X = int32(screenWidth/2 - 182 + (i * 35))
		PlayerHotbar.Slots[i].Y = int32(screenHeight - UserInterface.MapHeight*UserInterface.TileSize + 194)
		buttonDest.X = float32(PlayerHotbar.Slots[i].X)
		buttonDest.Y = float32(PlayerHotbar.Slots[i].Y)

		buttonSelectedDest.X = float32(PlayerHotbar.Slots[i].X)
		buttonSelectedDest.Y = float32(PlayerHotbar.Slots[i].Y)

		buttonActiveDest.X = float32(PlayerHotbar.Slots[i].X)
		buttonActiveDest.Y = float32(PlayerHotbar.Slots[i].Y + 4)

		if i == PlayerHotbar.SelectedIndex {
			rl.DrawTexturePro(buttonSprite, buttonSelected, buttonSelectedDest, rl.NewVector2(0, 0), 0, rl.White)
			rl.DrawTexturePro(buttonSprite, buttonActive, buttonActiveDest, rl.NewVector2(0, 0), 0, rl.White)
		} else {
			rl.DrawTexturePro(buttonSprite, buttonSrc, buttonDest, rl.NewVector2(0, 0), 0, rl.White)
		}

		item = PlayerHotbar.Slots[i]

		// left click on hotbar slot
		if isValidMousePosition && rl.CheckCollisionPointRec(mousePosition, buttonDest) && rl.IsMouseButtonPressed(rl.MouseLeftButton) && !Dragging.Drag {
			PlayerHotbar.SelectedIndex = i
			item = PlayerHotbar.Slots[i]

			if item.Name != "" {
				// ctrl+click = quick move to inventory
				if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
					if PlayerInventory.AddItemToHotbar(item) {
						PlayerHotbar.Slots[i] = Item{} // remove from hotbar
					} else {
						// inventory full, do normal drag instead
						Dragging.Drag = true
						Dragging.Item = item
						Dragging.Source = i
						Dragging.SourceType = "hotbar"
						PlayerHotbar.Slots[i] = Item{}
					}
				} else {
					// normal drag
					Dragging.Drag = true
					Dragging.Item = item
					Dragging.Source = i
					Dragging.SourceType = "hotbar"
					PlayerHotbar.Slots[i] = Item{}
				}
			}
		}

		if isValidMousePosition && rl.CheckCollisionPointRec(mousePosition, buttonDest) && rl.IsMouseButtonPressed(rl.MouseRightButton) && !Dragging.Drag {
			splitItems(&PlayerHotbar.Slots[i], "hotbar")
		}

		if isValidMousePosition && rl.CheckCollisionPointRec(mousePosition, buttonDest) && rl.IsMouseButtonReleased(rl.MouseLeftButton) && Dragging.Drag {
			if i == Dragging.Source && Dragging.SourceType == "hotbar" {
				placeItemInSlot(&PlayerHotbar.Slots[i], "hotbar")
			} else if PlayerHotbar.Slots[i].Name != "" {
				swapItems(&PlayerHotbar.Slots[i], "hotbar")
			} else {
				placeItemInSlot(&PlayerHotbar.Slots[i], "hotbar")
			}
		}

		if item.Name != "" {
			rl.DrawTexturePro(item.Icon, item.IconSrc, ScaleItemDest(buttonDest, -10), rl.NewVector2(0, 0), 0, rl.White)
			rl.DrawText(fmt.Sprintf("%d", item.Quantity), int32(buttonDest.X+25), int32(buttonDest.Y+30), 0, rl.White)
		}

		if Dragging.Drag {
			mouse := rl.GetMousePosition()

			item = Dragging.Item
			rl.DrawTexturePro(item.Icon, item.IconSrc, rl.NewRectangle(mouse.X-24, mouse.Y-24, 32, 32), rl.NewVector2(0, 0), 0, rl.White)
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
	if key >= rl.KeyOne && key <= rl.KeyNine && len(PlayerHotbar.Slots) > int(key)-rl.KeyOne {
		PlayerHotbar.SelectedIndex = int(key) - rl.KeyOne
	}

	if key == rl.KeyZero && len(PlayerHotbar.Slots) > 9 {
		PlayerHotbar.SelectedIndex = 9
	}

	mousePosition := rl.GetMousePosition()
	if mousePosition.X >= 0 && mousePosition.Y >= 0 && mousePosition.X <= screenWidth && mousePosition.Y <= screenHeight {
		scrollPosition := rl.GetMouseWheelMove()

		if scrollPosition > 0 {
			PlayerHotbar.SelectedIndex--
			if PlayerHotbar.SelectedIndex < 0 {
				PlayerHotbar.SelectedIndex = 9
			}
		}

		if scrollPosition < 0 {
			PlayerHotbar.SelectedIndex++
			if PlayerHotbar.SelectedIndex > 9 {
				PlayerHotbar.SelectedIndex = 0
			}
		}
	}

	if rl.IsKeyPressed(rl.KeyE) {
		openInventory = !openInventory
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		openInventory = false
	}

}

func DrawInventorySlots() {
	buttonSrc := rl.NewRectangle(224, 112, 48, 48)
	buttonDest := rl.NewRectangle(0, 0, 48, 48)
	mousePosition := rl.GetMousePosition()
	isValidMousePosition := mousePosition.X >= 0 && mousePosition.Y >= 0 && mousePosition.X <= screenWidth && mousePosition.Y <= screenHeight

	for i, item := range PlayerInventory.Slots {
		if i >= len(PlayerInventory.Slots) {
			break
		}
		PlayerInventory.Slots[i].X = int32(screenWidth/2 - 165 + (i % 9 * 35))
		PlayerInventory.Slots[i].Y = int32(screenHeight/2 + 170 + (i / 9 * 40))
		buttonDest.X = float32(PlayerInventory.Slots[i].X)
		buttonDest.Y = float32(PlayerInventory.Slots[i].Y)

		rl.DrawTexturePro(buttonSprite, buttonSrc, buttonDest, rl.NewVector2(0, 0), 0, rl.White)

		item = PlayerInventory.Slots[i]

		// left click on inventory slot
		if isValidMousePosition && rl.CheckCollisionPointRec(mousePosition, buttonDest) && rl.IsMouseButtonPressed(rl.MouseLeftButton) && !Dragging.Drag {
			PlayerInventory.SelectedIndex = i
			item = PlayerInventory.Slots[i]

			if item.Name != "" {
				// ctrl+click = quick move to hotbar
				if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
					if PlayerHotbar.AddItemToHotbar(item) {
						PlayerInventory.Slots[i] = Item{} // remove from inv
					} else {
						// hotbar full, do normal drag instead
						Dragging.Drag = true
						Dragging.Item = item
						Dragging.Source = i
						Dragging.SourceType = "inventory"
						PlayerInventory.Slots[i] = Item{}
					}
				} else {
					// normal drag
					Dragging.Drag = true
					Dragging.Item = item
					Dragging.Source = i
					Dragging.SourceType = "inventory"
					PlayerInventory.Slots[i] = Item{}
				}
			}
		}

		if isValidMousePosition && rl.CheckCollisionPointRec(mousePosition, buttonDest) && rl.IsMouseButtonPressed(rl.MouseRightButton) && !Dragging.Drag {
			splitItems(&PlayerInventory.Slots[i], "inventory")
		}

		if isValidMousePosition && rl.CheckCollisionPointRec(mousePosition, buttonDest) && rl.IsMouseButtonReleased(rl.MouseLeftButton) && Dragging.Drag {
			if i == Dragging.Source && Dragging.SourceType == "inventory" {
				placeItemInSlot(&PlayerInventory.Slots[i], "inventory")
			} else if PlayerInventory.Slots[i].Name != "" {
				swapItems(&PlayerInventory.Slots[i], "inventory")
			} else {
				placeItemInSlot(&PlayerInventory.Slots[i], "inventory")
			}
		}

		if item.Name != "" {
			rl.DrawTexturePro(item.Icon, item.IconSrc, ScaleItemDest(buttonDest, -10), rl.NewVector2(0, 0), 0, rl.White)
			rl.DrawText(fmt.Sprintf("%d", item.Quantity), int32(buttonDest.X+25), int32(buttonDest.Y+30), 0, rl.White)
		}

		if Dragging.Drag {
			mouse := rl.GetMousePosition()

			item = Dragging.Item
			rl.DrawTexturePro(item.Icon, item.IconSrc, rl.NewRectangle(mouse.X-24, mouse.Y-24, 32, 32), rl.NewVector2(0, 0), 0, rl.White)
		}
	}
}

func resetDragState() {
	Dragging.Drag = false
	Dragging.SourceType = ""
}

func placeItemInSlot(slot *Item, targetSourceType string) {
	*slot = Dragging.Item
	resetDragState()
}

func splitItems(slot *Item, targetSourceType string) {
	if slot.Name != "" {
		if slot.Quantity == maxQuantity {
			Dragging.Drag = true
			Dragging.Item = *slot
			Dragging.Item.Quantity = maxQuantity / 2
			slot.Quantity -= maxQuantity / 2
			Dragging.SourceType = targetSourceType
		} else if slot.Quantity > 1 {
			Dragging.Drag = true
			Dragging.Item = *slot
			Dragging.Item.Quantity = 1
			slot.Quantity -= 1
			Dragging.SourceType = targetSourceType
		} else {
			Dragging.Drag = true
			Dragging.Item = *slot
			Dragging.SourceType = targetSourceType
			*slot = Item{}
		}
	}
}

func swapItems(slot *Item, targetSourceType string) {
	if slot.Name == Dragging.Item.Name && slot.Quantity+Dragging.Item.Quantity <= maxQuantity {
		slot.Quantity += Dragging.Item.Quantity
		resetDragState()
	} else {
		cachedItem = *slot
		*slot = Dragging.Item
		Dragging.Item = cachedItem
		Dragging.SourceType = targetSourceType
	}
}

func ScaleItemDest(i rl.Rectangle, s float32) rl.Rectangle {
	return rl.NewRectangle(i.X-s, i.Y-s, i.Width+s*2, i.Height+s*2)
}

func (h *Hotbar) AddItemToHotbar(newItem Item) bool {
	if len(h.Slots) == 0 {
		return false
	}
	
	for i := range h.Slots {
		if h.Slots[i].Name == newItem.Name && h.Slots[i].Quantity+newItem.Quantity <= maxQuantity {
			h.Slots[i].Quantity += newItem.Quantity
			return true
		}
	}

	for i := range h.Slots {
		if h.Slots[i].Name == "" {
			h.Slots[i] = newItem
			return true
		}
	}

	return false
}

func UnloadUserInterface() {
	rl.UnloadTexture(spritesheet)
}

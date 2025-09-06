package items

import (
	"godew-valley/pkg/player"
	"godew-valley/pkg/userinterface"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldItem struct {
	Position rl.Vector2
	Item     userinterface.Item
	Active   bool
}

var (
	ItemsSprite    rl.Texture2D
	AxeSrc         rl.Rectangle
	AxeDest        rl.Rectangle
	HoeSrc         rl.Rectangle
	GrassSrc       rl.Rectangle
	StickSrc       rl.Rectangle
	BranchSrc      rl.Rectangle
	WateringCanSrc rl.Rectangle
	worldItems     []WorldItem
)

func InitItemTextures() {
	ItemsSprite = rl.LoadTexture("assets/Objects/Items/tools-n-meterial-items.png")
	AxeSrc = rl.NewRectangle(16, 0, 16, 16)
	WateringCanSrc = rl.NewRectangle(0, 0, 16, 16)
	HoeSrc = rl.NewRectangle(32, 0, 16, 16)
	GrassSrc = rl.NewRectangle(48, 0, 16, 16)
	StickSrc = rl.NewRectangle(16, 16, 16, 16)
	BranchSrc = rl.NewRectangle(32, 32, 16, 16)
}

func InitItems() {
	if !isItemInInventory("Axe") {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(380, 450),
			Item: userinterface.Item{
				Name:     "Axe",
				Icon:     ItemsSprite,
				IconSrc:  AxeSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}

	if !isItemInInventory("Watering Can") {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(430, 450),
			Item: userinterface.Item{
				Name:     "Watering Can",
				Icon:     ItemsSprite,
				IconSrc:  WateringCanSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}

	if !isItemInInventory("Hoe") {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(480, 450),
			Item: userinterface.Item{
				Name:     "Hoe",
				Icon:     ItemsSprite,
				IconSrc:  HoeSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}

	if !isItemInInventory("Grass") {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(530, 450),
			Item: userinterface.Item{
				Name:     "Grass",
				Icon:     ItemsSprite,
				IconSrc:  GrassSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}

	if !isItemInInventory("Branch") {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(580, 450),
			Item: userinterface.Item{
				Name:     "Branch",
				Icon:     ItemsSprite,
				IconSrc:  BranchSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}

	if !isItemInInventory("Stick") {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(630, 450),
			Item: userinterface.Item{
				Name:     "Stick",
				Icon:     ItemsSprite,
				IconSrc:  StickSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}
}

func isItemInInventory(itemName string) bool {
	for _, slot := range userinterface.PlayerHotbar.Slots {
		if slot.Name == itemName && slot.Active && slot.Quantity > 0 {
			return true
		}
	}

	for _, slot := range userinterface.PlayerInventory.Slots {
		if slot.Name == itemName && slot.Active && slot.Quantity > 0 {
			return true
		}
	}

	return false
}

func InputHoe() {
	if rl.IsKeyPressed(rl.KeyH) {
		worldItems = append(worldItems, WorldItem{
			Position: rl.NewVector2(player.PlayerDest.X+50, player.PlayerDest.Y+50),
			Item: userinterface.Item{
				Name:     "Hoe",
				Icon:     ItemsSprite,
				IconSrc:  HoeSrc,
				Quantity: 1,
			},
			Active: true,
		})
	}
}

func DrawItems() {
	for _, item := range worldItems {
		if !item.Active {
			continue
		}

		itemRect := rl.NewRectangle(item.Position.X, item.Position.Y, 16, 16)
		rl.DrawTexturePro(item.Item.Icon, item.Item.IconSrc, itemRect, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func UpdateItems() {
	for i := range worldItems {
		item := &worldItems[i]

		if !item.Active {
			continue
		}

		itemRect := rl.NewRectangle(item.Position.X, item.Position.Y, 32, 32)

		if rl.CheckCollisionRecs(player.PlayerHitBox, itemRect) {
			successHotBar := userinterface.PlayerHotbar.AddItemToHotbar(item.Item)

			if successHotBar {
				item.Active = false
			} else {
				successInventory := userinterface.PlayerInventory.AddItemToHotbar(item.Item)

				if successInventory {
					item.Active = false
				}
			}
		}
	}
}

func UnloadItems() {
	rl.UnloadTexture(ItemsSprite)
}

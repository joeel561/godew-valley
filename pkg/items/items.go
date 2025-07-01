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
	ItemsSprite rl.Texture2D
	AxeSrc      rl.Rectangle
	AxeDest     rl.Rectangle
	HoeSrc      rl.Rectangle
	worldItems  []WorldItem
)

func InitItems() {
	ItemsSprite = rl.LoadTexture("assets/Objects/Items/tools-n-meterial-items.png")
	AxeSrc = rl.NewRectangle(16, 0, 16, 16)
	WateringCanSrc := rl.NewRectangle(0, 0, 16, 16)
	HoeSrc = rl.NewRectangle(32, 0, 16, 16)

	worldItems = append(worldItems, WorldItem{
		Position: rl.NewVector2(400, 430),
		Item: userinterface.Item{
			Name:     "Axe",
			Icon:     ItemsSprite,
			IconSrc:  AxeSrc,
			Quantity: 1,
		},
		Active: true,
	})

	worldItems = append(worldItems, WorldItem{
		Position: rl.NewVector2(430, 460),
		Item: userinterface.Item{
			Name:     "Watering Can",
			Icon:     ItemsSprite,
			IconSrc:  WateringCanSrc,
			Quantity: 1,
		},
		Active: true,
	})

	worldItems = append(worldItems, WorldItem{
		Position: rl.NewVector2(480, 465),
		Item: userinterface.Item{
			Name:     "Hoe",
			Icon:     ItemsSprite,
			IconSrc:  HoeSrc,
			Quantity: 1,
		},
		Active: true,
	})
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

func UnloadAxe() {
	rl.UnloadTexture(ItemsSprite)
}

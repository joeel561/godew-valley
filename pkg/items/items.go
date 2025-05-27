package items

import (
	"fmt"
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
	worldItems  []WorldItem
)

func InitItems() {
	ItemsSprite = rl.LoadTexture("assets/Objects/Items/tools-n-meterial-items.png")
	AxeSrc = rl.NewRectangle(16, 0, 16, 16)
	WateringCanSrc := rl.NewRectangle(0, 0, 16, 16)

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
		Position: rl.NewVector2(500, 550),
		Item: userinterface.Item{
			Name:     "Watering Can",
			Icon:     ItemsSprite,
			IconSrc:  WateringCanSrc,
			Quantity: 1,
		},
		Active: true,
	})
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
			success := userinterface.PlayerHotbar.AddItemToHotbar(item.Item)

			if success {
				item.Active = false
			} else {
				fmt.Println("Item not added to hotbar")
			}
		}

	}

}

func UnloadAxe() {
	rl.UnloadTexture(ItemsSprite)
}

func loadMap(mapFile string) {
	file, err := os.Open(mapFile)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &jsonMap)
}

func drawScene() {
	var groundTiles []Tile
	var waterTiles []Tile

	for i := 0; i < len(jsonMap.Layers); i++ {
		if jsonMap.Layers[i].Name == "Background" {
			groundTiles = jsonMap.Layers[i].Tiles
		}

		if jsonMap.Layers[i].Name == "Water" {
			waterTiles = jsonMap.Layers[i].Tiles
		}
	}

	for i := 0; i < len(waterTiles); i++ {
		s, _ := strconv.ParseInt(waterTiles[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(jsonMap.TileSize)
		tileSrc.X = float32(jsonMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(jsonMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(waterTiles[i].X * jsonMap.TileSize)
		tileDest.Y = float32(waterTiles[i].Y * jsonMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	for i := 0; i < len(groundTiles); i++ {
		s, _ := strconv.ParseInt(groundTiles[i].Id, 10, 64)
		tileId := int(s)
		tex = spritesheetMap

		texColumns := tex.Width / int32(jsonMap.TileSize)
		tileSrc.X = float32(jsonMap.TileSize) * float32((tileId)%int(texColumns))
		tileSrc.Y = float32(jsonMap.TileSize) * float32((tileId)/int(texColumns))

		tileDest.X = float32(groundTiles[i].X * jsonMap.TileSize)
		tileDest.Y = float32(groundTiles[i].Y * jsonMap.TileSize)

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
	}

	rl.DrawRectangle(int32(wall.X), int32(wall.Y), int32(wall.Width), int32(wall.Height), rl.Red)

	if printDebug {
		// Draw cetner map cross
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(-20, 0), 1, rl.Gray)
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(20, 0), 1, rl.Red)
		rl.DrawTriangle(rl.NewVector2(16, 2), rl.NewVector2(20, 0), rl.NewVector2(16, -2), rl.Red)
		rl.DrawText("X", int32(22), int32(-5), int32(10), rl.Black)
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(0, -20), 1, rl.Gray)
		rl.DrawLineEx(rl.NewVector2(0, 0), rl.NewVector2(0, 20), 1, rl.Blue)
		rl.DrawTriangle(rl.NewVector2(-2, 16), rl.NewVector2(0, 20), rl.NewVector2(2, 16), rl.Blue)
		rl.DrawText("Y", int32(-2), int32(22), int32(10), rl.Black)

		// Draw collision rectangle
		rl.DrawRectangleLinesEx(playerHitBox, 1, rl.Green)
		rl.DrawRectangleLinesEx(playerDest, 1, rl.Purple)
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(0, 0), 0, rl.White)
}
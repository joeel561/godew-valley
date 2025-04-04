

func rectToString(rec rl.Rectangle) string {
	return fmt.Sprintf("X:%v, Y:%v, H:%v, W:%v", rec.X, rec.Y, rec.Height, rec.Width)
}

func vec2ToString(vec rl.Vector2) string {
	return fmt.Sprintf("X:%v, Y:%v", vec.X, vec.Y)
}

func debugText() []string {
	return []string{
		fmt.Sprintf("FPS: %v", rl.GetFPS()),
		fmt.Sprintf("Cam Target %v", vec2ToString(cam.Target)),
		fmt.Sprintf("Player Direction: %v   U:%v, D:%v, L:%v, R:%v", playerDir, playerUp, playerDown, playerLeft, playerRight),
		fmt.Sprintf("Player Speed: %v", playerSpeed),
		fmt.Sprintf("Player Frame: %v", playerFrame),
		fmt.Sprintf("Player Moving: %v", playerMoving),
		fmt.Sprintf("Player Src %v", rectToString(playerSrc)),
		fmt.Sprintf("Player Dest %v", rectToString(playerDest)),
		fmt.Sprintf("Music Paused: %v", musicPaused),
	}
}

func drawDebug(debugText []string) {
	textSize := 10
	lineSpace := 15

	offsetX := 10
	offsetY := 10

	for i, line := range debugText {
		rl.DrawText(line, int32(offsetX), int32(offsetY+lineSpace*i), int32(textSize), rl.Black)
	}
}

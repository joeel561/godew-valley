package save

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"godew-valley/pkg/items"
	"godew-valley/pkg/player"
	"godew-valley/pkg/userinterface"
	"io"
	"io/ioutil"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState struct {
	PlayerData    PlayerData      `json:"player"`
	InventoryData InventoryData   `json:"inventory"`
	WorldItems    []WorldItemData `json:"worldItems"`
	Timestamp     time.Time       `json:"timestamp"`
	GameTime      float64         `json:"gameTime"`
	Checksum      string          `json:"checksum"`
}

type PlayerData struct {
	Position  rl.Vector2 `json:"position"`
	Direction int        `json:"direction"`
	Frame     int        `json:"frame"`
}

type InventoryData struct {
	HotbarSlots    []ItemData `json:"hotbarSlots"`
	InventorySlots []ItemData `json:"inventorySlots"`
	SelectedIndex  int        `json:"selectedIndex"`
}

type ItemData struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Active   bool   `json:"active"`
	X        int32  `json:"x"`
	Y        int32  `json:"y"`
}

type WorldItemData struct {
	Position rl.Vector2 `json:"position"`
	Item     ItemData   `json:"item"`
	Active   bool       `json:"active"`
}

const SaveFileName = "savegame.dat"
const encryptionKey = "godew-valley-secret-key-2024"

func calculateChecksum(gameState GameState) string {
	tempState := gameState
	tempState.Checksum = ""

	data, err := json.Marshal(tempState)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func validateChecksum(gameState GameState) error {
	expectedChecksum := calculateChecksum(gameState)
	if gameState.Checksum != expectedChecksum {
		return fmt.Errorf("savegame integrity check failed: checksum mismatch")
	}
	return nil
}

func validateItemData(item ItemData) error {
	maxQuantities := map[string]int{
		"Axe":         1,
		"Watering Can": 1,
		"Hoe":         1,
		"Grass":       99,
		"Branch":      99,
		"Stick":       99,
	}

	if item.Name != "" {
		maxQty, exists := maxQuantities[item.Name]
		if !exists {
			return fmt.Errorf("unknown item: %s", item.Name)
		}

		if item.Quantity < 0 || item.Quantity > maxQty {
			return fmt.Errorf("invalid quantity for %s: %d (max: %d)", item.Name, item.Quantity, maxQty)
		}
	}

	return nil
}

func validateInventoryData(data InventoryData) error {
	for i, item := range data.HotbarSlots {
		if err := validateItemData(item); err != nil {
			return fmt.Errorf("invalid hotbar item at slot %d: %v", i, err)
		}
	}

	for i, item := range data.InventorySlots {
		if err := validateItemData(item); err != nil {
			return fmt.Errorf("invalid inventory item at slot %d: %v", i, err)
		}
	}

	return nil
}

// encrypt encrypts data using AES-GCM
func encrypt(data []byte) ([]byte, error) {
	hash := sha256.Sum256([]byte(encryptionKey))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func decrypt(data []byte) ([]byte, error) {
	hash := sha256.Sum256([]byte(encryptionKey))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func SaveGame() error {
	gameState := GameState{
		PlayerData:    getPlayerData(),
		InventoryData: getInventoryData(),
		WorldItems:    getWorldItemsData(),
		Timestamp:     time.Now(),
		GameTime:      0,
	}

	gameState.Checksum = calculateChecksum(gameState)

	data, err := json.Marshal(gameState)
	if err != nil {
		return err
	}

	encryptedData, err := encrypt(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(SaveFileName, encryptedData, 0644)
}

func LoadGame() error {
	if _, err := os.Stat(SaveFileName); os.IsNotExist(err) {
		return nil
	}

	encryptedData, err := ioutil.ReadFile(SaveFileName)
	if err != nil {
		return err
	}

	data, err := decrypt(encryptedData)
	if err != nil {
		return err
	}

	var gameState GameState
	err = json.Unmarshal(data, &gameState)
	if err != nil {
		return err
	}

	err = validateChecksum(gameState)
	if err != nil {
		return fmt.Errorf("savegame corrupted or tampered: %v", err)
	}

	err = validateInventoryData(gameState.InventoryData)
	if err != nil {
		return fmt.Errorf("invalid savegame data: %v", err)
	}

	applyPlayerData(gameState.PlayerData)
	applyInventoryData(gameState.InventoryData)
	applyWorldItemsData(gameState.WorldItems)

	return nil
}

func getPlayerData() PlayerData {
	return PlayerData{
		Position:  rl.NewVector2(player.PlayerDest.X, player.PlayerDest.Y),
		Direction: 0,
		Frame:     0,
	}
}

func getInventoryData() InventoryData {
	hotbarSlots := make([]ItemData, len(userinterface.PlayerHotbar.Slots))
	for i, slot := range userinterface.PlayerHotbar.Slots {
		hotbarSlots[i] = ItemData{
			Name:     slot.Name,
			Quantity: slot.Quantity,
			Active:   slot.Active,
			X:        slot.X,
			Y:        slot.Y,
		}
	}

	inventorySlots := make([]ItemData, len(userinterface.PlayerInventory.Slots))
	for i, slot := range userinterface.PlayerInventory.Slots {
		inventorySlots[i] = ItemData{
			Name:     slot.Name,
			Quantity: slot.Quantity,
			Active:   slot.Active,
			X:        slot.X,
			Y:        slot.Y,
		}
	}

	return InventoryData{
		HotbarSlots:    hotbarSlots,
		InventorySlots: inventorySlots,
		SelectedIndex:  userinterface.PlayerHotbar.SelectedIndex,
	}
}

func getWorldItemsData() []WorldItemData {
	return []WorldItemData{}
}

func applyPlayerData(data PlayerData) {
	player.PlayerDest.X = data.Position.X
	player.PlayerDest.Y = data.Position.Y
}

func applyInventoryData(data InventoryData) {
	for i, itemData := range data.HotbarSlots {
		if i < len(userinterface.PlayerHotbar.Slots) {
			userinterface.PlayerHotbar.Slots[i] = reconstructItem(itemData)
		}
	}

	for i, itemData := range data.InventorySlots {
		if i < len(userinterface.PlayerInventory.Slots) {
			userinterface.PlayerInventory.Slots[i] = reconstructItem(itemData)
		}
	}

	userinterface.PlayerHotbar.SelectedIndex = data.SelectedIndex
}

func reconstructItem(data ItemData) userinterface.Item {
	item := userinterface.Item{
		Name:     data.Name,
		Quantity: data.Quantity,
		Active:   data.Name != "" && data.Quantity > 0,
		X:        data.X,
		Y:        data.Y,
	}

	// Setze die korrekten Texturen basierend auf dem Item-Namen
	switch data.Name {
	case "Axe":
		item.Icon = items.ItemsSprite
		item.IconSrc = items.AxeSrc
	case "Watering Can":
		item.Icon = items.ItemsSprite
		item.IconSrc = items.WateringCanSrc
	case "Hoe":
		item.Icon = items.ItemsSprite
		item.IconSrc = items.HoeSrc
	case "Grass":
		item.Icon = items.ItemsSprite
		item.IconSrc = items.GrassSrc
	case "Stick":
		item.Icon = items.ItemsSprite
		item.IconSrc = items.StickSrc
	case "Branch":
		item.Icon = items.ItemsSprite
		item.IconSrc = items.BranchSrc
	}

	return item
}

func applyWorldItemsData(data []WorldItemData) {

}

func SaveExists() bool {
	_, err := os.Stat(SaveFileName)
	return !os.IsNotExist(err)
}
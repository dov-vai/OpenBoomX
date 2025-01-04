package services

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"obx/gui/constants"
	"os"
	"path/filepath"
)

type ColorChangeListener interface {
	OnColorListChanged(colors []color.NRGBA)
}

type ColorPresetService struct {
	configDir      string
	presetFilePath string
	colors         []color.NRGBA
	listeners      []ColorChangeListener
}

type ColorPresetData struct {
	Colors []color.NRGBA `json:"colors"`
}

var defaultColorPresets = []color.NRGBA{
	{R: 237, G: 34, B: 36, A: 255},   // Red
	{R: 243, G: 115, B: 33, A: 255},  // Orange
	{R: 249, G: 188, B: 40, A: 255},  // Yellow
	{R: 131, G: 197, B: 50, A: 255},  // Green
	{R: 73, G: 195, B: 176, A: 255},  // Teal
	{R: 75, G: 178, B: 251, A: 255},  // Light Blue
	{R: 0, G: 111, B: 249, A: 255},   // Blue
	{R: 81, G: 32, B: 223, A: 255},   // Indigo
	{R: 180, G: 44, B: 215, A: 255},  // Violet
	{R: 255, G: 255, B: 255, A: 255}, // White
}

func NewColorPresetService() *ColorPresetService {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}

	service := &ColorPresetService{
		configDir:      filepath.Join(configDir, constants.AppName),
		presetFilePath: filepath.Join(configDir, constants.AppName, "colors.json"),
		colors:         []color.NRGBA{},
		listeners:      []ColorChangeListener{},
	}

	if err := service.ensureConfigDir(); err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}

	if err := service.loadColorPresets(); err != nil {
		log.Fatalf("Error loading color presets: %v", err)
	}

	return service
}

func (service *ColorPresetService) ensureConfigDir() error {
	if _, err := os.Stat(service.configDir); os.IsNotExist(err) {
		return os.MkdirAll(service.configDir, 0755)
	}
	return nil
}

func (service *ColorPresetService) loadColorPresets() error {
	dataFile, err := os.ReadFile(service.presetFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize with default presets
			service.colors = defaultColorPresets
			return service.saveColorPresets()
		}
		return fmt.Errorf("error reading color preset file: %w", err)
	}

	var configData ColorPresetData
	err = json.Unmarshal(dataFile, &configData)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	service.colors = configData.Colors

	return nil
}

func (service *ColorPresetService) saveColorPresets() error {
	configData := ColorPresetData{
		Colors: service.colors,
	}

	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = os.WriteFile(service.presetFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing color preset file: %w", err)
	}

	return nil
}

func (service *ColorPresetService) AddColor(c color.NRGBA) error {
	for _, existingColor := range service.colors {
		if existingColor == c {
			return fmt.Errorf("color already exists in the list")
		}
	}

	service.colors = append(service.colors, c)

	if err := service.saveColorPresets(); err != nil {
		return fmt.Errorf("error saving color presets after adding: %w", err)
	}

	service.notifyListeners()
	return nil
}

func (service *ColorPresetService) DeleteColor(c color.NRGBA) error {
	found := false
	for i, existingColor := range service.colors {
		if existingColor == c {
			service.colors = append(service.colors[:i], service.colors[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("color not found in the list")
	}

	if err := service.saveColorPresets(); err != nil {
		return fmt.Errorf("error saving color presets after deletion: %w", err)
	}

	service.notifyListeners()
	return nil
}

func (service *ColorPresetService) ListColors() []color.NRGBA {
	colorsCopy := make([]color.NRGBA, len(service.colors))
	copy(colorsCopy, service.colors)
	return colorsCopy
}

func (service *ColorPresetService) RegisterListener(listener ColorChangeListener) {
	service.listeners = append(service.listeners, listener)
}

func (service *ColorPresetService) RemoveListener(listener ColorChangeListener) {
	for i, l := range service.listeners {
		if l == listener {
			service.listeners = append(service.listeners[:i], service.listeners[i+1:]...)
			break
		}
	}
}

func (service *ColorPresetService) notifyListeners() {
	for _, listener := range service.listeners {
		listener.OnColorListChanged(service.ListColors())
	}
}

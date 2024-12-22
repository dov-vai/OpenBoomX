package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type EqPresetService struct {
	configDir      string
	presetFilePath string
	activePreset   string
	presets        map[string][]float32
}

type PresetData struct {
	ActivePreset string               `json:"activePreset"`
	Presets      map[string][]float32 `json:"presets"`
}

func NewEqPresetService() *EqPresetService {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}

	service := &EqPresetService{
		configDir:      filepath.Join(configDir, "OpenBoomX"),
		presetFilePath: filepath.Join(configDir, "OpenBoomX", "presets.json"),
		presets:        make(map[string][]float32),
	}

	if err := service.ensureConfigDir(); err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}

	if err := service.loadPresets(); err != nil {
		log.Fatalf("Error loading presets: %v", err)
	}

	return service
}

func (service *EqPresetService) ensureConfigDir() error {
	if _, err := os.Stat(service.configDir); os.IsNotExist(err) {
		return os.MkdirAll(service.configDir, 0755)
	}
	return nil
}

func (service *EqPresetService) loadPresets() error {
	dataFile, err := os.ReadFile(service.presetFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return service.savePresets()
		}
		return fmt.Errorf("error reading preset file: %w", err)
	}

	var configData PresetData
	err = json.Unmarshal(dataFile, &configData)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	service.activePreset = configData.ActivePreset
	service.presets = configData.Presets

	return nil
}

func (service *EqPresetService) savePresets() error {
	configData := PresetData{
		ActivePreset: service.activePreset,
		Presets:      service.presets,
	}

	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = os.WriteFile(service.presetFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing preset file: %w", err)
	}

	return nil
}

func (service *EqPresetService) AddPreset(title string, values []float32) error {
	service.presets[title] = values
	service.activePreset = title

	if err := service.savePresets(); err != nil {
		return fmt.Errorf("error saving presets after adding: %w", err)
	}

	log.Printf("Added and activated preset: '%s' with values: %v\n", title, values)
	return nil
}

func (service *EqPresetService) DeletePreset(title string) error {
	if _, exists := service.presets[title]; !exists {
		return fmt.Errorf("preset with title '%s' not found", title)
	}

	delete(service.presets, title)

	// If the deleted preset was the active one, clear the active preset.
	if service.activePreset == title {
		service.activePreset = ""
	}

	if err := service.savePresets(); err != nil {
		return fmt.Errorf("error saving presets after deletion: %w", err)
	}

	log.Printf("Deleted preset: '%s'\n", title)
	return nil
}

func (service *EqPresetService) GetActivePreset() string {
	return service.activePreset
}

func (service *EqPresetService) SetActivePreset(title string) error {
	if _, exists := service.presets[title]; !exists {
		return fmt.Errorf("preset with title '%s' not found", title)
	}

	service.activePreset = title
	if err := service.savePresets(); err != nil {
		return fmt.Errorf("error saving presets after setting active: %w", err)
	}

	log.Printf("Set active preset to: '%s'\n", title)
	return nil
}

func (service *EqPresetService) GetPresetValues(title string) ([]float32, error) {
	values, exists := service.presets[title]
	if !exists {
		return nil, fmt.Errorf("preset with title '%s' not found", title)
	}
	return values, nil
}

func (service *EqPresetService) ListPresets() []string {
	titles := make([]string, 0, len(service.presets))
	for title := range service.presets {
		titles = append(titles, title)
	}
	return titles
}

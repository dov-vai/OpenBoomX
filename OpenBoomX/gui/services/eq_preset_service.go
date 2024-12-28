package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"golang.org/x/exp/maps"
)

type PresetChangeListener interface {
	OnPresetChanged(newPreset string, values []float32)
}

type EqPresetService struct {
	configDir      string
	presetFilePath string
	activePreset   string
	presets        map[string]PresetDetails
	listeners      []PresetChangeListener
}

type PresetDetails struct {
	Values    []float32 `json:"values"`
	Timestamp int64     `json:"timestamp"`
}

type PresetData struct {
	ActivePreset string                   `json:"activePreset"`
	Presets      map[string]PresetDetails `json:"presets"`
}

func NewEqPresetService() *EqPresetService {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}

	service := &EqPresetService{
		configDir:      filepath.Join(configDir, "OpenBoomX"),
		presetFilePath: filepath.Join(configDir, "OpenBoomX", "presets.json"),
		presets:        make(map[string]PresetDetails),
		listeners:      []PresetChangeListener{},
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
	service.presets[title] = PresetDetails{
		Values:    values,
		Timestamp: time.Now().Unix(),
	}
	service.activePreset = title

	if err := service.savePresets(); err != nil {
		return fmt.Errorf("error saving presets after adding: %w", err)
	}

	service.notifyListeners()
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
	service.notifyListeners()

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

	service.notifyListeners()
	return nil
}

func (service *EqPresetService) GetPresetValues(title string) ([]float32, error) {
	presetDetails, exists := service.presets[title]
	if !exists {
		return nil, fmt.Errorf("preset with title '%s' not found", title)
	}

	valuesCopy := make([]float32, len(presetDetails.Values))
	copy(valuesCopy, presetDetails.Values)

	return valuesCopy, nil
}

func (service *EqPresetService) ListPresets() []string {
	titles := maps.Keys(service.presets)

	sort.Slice(titles, func(i, j int) bool {
		return service.presets[titles[i]].Timestamp > service.presets[titles[j]].Timestamp
	})

	return titles
}

func (service *EqPresetService) RegisterListener(listener PresetChangeListener) {
	service.listeners = append(service.listeners, listener)
}

func (service *EqPresetService) RemoveListener(listener PresetChangeListener) {
	for i, l := range service.listeners {
		if l == listener {
			service.listeners = append(service.listeners[:i], service.listeners[i+1:]...)
			break
		}
	}
}

func (service *EqPresetService) notifyListeners() {
	activePresetValues, _ := service.GetPresetValues(service.activePreset)
	for _, listener := range service.listeners {
		listener.OnPresetChanged(service.activePreset, activePresetValues)
	}
}

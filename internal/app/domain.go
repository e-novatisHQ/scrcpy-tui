package app

type Device struct {
	Serial string
	State  string
	Model  string
}
type Preset struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Args        []string `json:"args"`
}
type Config struct {
	Presets    []Preset `json:"presets"`
	LastDevice string   `json:"last_device,omitempty"`
	LastPreset string   `json:"last_preset,omitempty"`
}

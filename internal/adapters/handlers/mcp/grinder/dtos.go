package grinder

type ListGrindersRequest struct {
	Type string `json:"type,omitempty" jsonschema:"optional filter: manual or electric"`
}

type ListGrindersResponse struct {
	Grinders []Grinder `json:"grinders"`
}

type Grinder struct {
	ID                    string     `json:"id"`
	Name                  string     `json:"name"`
	Type                  string     `json:"type"`
	Aliases               []string   `json:"aliases"`
	SettingNotation       string     `json:"setting_notation"`
	RangeUM               [2]float64 `json:"range_um"`
	EffectiveUMPerStep    float64    `json:"effective_um_per_step"`
	ManufacturerUMPerStep string     `json:"manufacturer_um_per_step,omitempty"`
}

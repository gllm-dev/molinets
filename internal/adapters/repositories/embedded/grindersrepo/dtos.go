package grindersrepo

type grinderDTO struct {
	ID                    string       `json:"id"`
	Name                  string       `json:"name"`
	Aliases               []string     `json:"aliases"`
	Type                  string       `json:"type"`
	Notation              notationDTO  `json:"notation"`
	NotationHint          string       `json:"notation_hint"`
	Calibration           []segmentDTO `json:"calibration"`
	RangeUM               [2]float64   `json:"range_um"`
	ManufacturerUMPerStep string       `json:"manufacturer_um_per_step,omitempty"`
	FitRMSEUM             *float64     `json:"fit_rmse_um,omitempty"`
	Reference             string       `json:"reference"`
}

type notationDTO struct {
	Kind              string  `json:"kind"`
	ClicksPerRotation int     `json:"clicks_per_rotation,omitempty"`
	ClicksPerNumber   int     `json:"clicks_per_number,omitempty"`
	Step              float64 `json:"step,omitempty"`
}

type segmentDTO struct {
	From      float64 `json:"from"`
	OffsetUM  float64 `json:"offset_um"`
	UMPerUnit float64 `json:"um_per_unit"`
}

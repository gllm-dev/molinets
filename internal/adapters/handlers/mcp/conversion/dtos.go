package conversion

type ConvertRequest struct {
	FromGrinder string   `json:"from_grinder" jsonschema:"source grinder id, name or alias, e.g. comandante, c3s pro, k6"`
	Setting     string   `json:"setting" jsonschema:"setting in the source grinder's own notation: 26 (clicks), 1.22 (Kingrinder rotation.clicks), 2.8.4 (1Zpresso rotation.number.click), 4.2 (Fellow Ode or OXO, 4 and 2/3), 53 (Niche second rotation). See list_grinders for each notation"`
	ToGrinders  []string `json:"to_grinders,omitempty" jsonschema:"optional target grinders; omit to convert to every supported grinder"`
}

type ConvertResponse struct {
	From            string       `json:"from"`
	Setting         string       `json:"setting"`
	ApproxMicrons   int          `json:"approx_microns"`
	MatchingMethods []string     `json:"matching_methods"`
	Equivalents     []Equivalent `json:"equivalents"`
	Note            string       `json:"note"`
}

type Equivalent struct {
	Grinder         string  `json:"grinder"`
	ID              string  `json:"id"`
	InRange         bool    `json:"in_range"`
	Setting         string  `json:"setting,omitempty"`
	ExactValue      float64 `json:"exact_value,omitempty"`
	SettingNotation string  `json:"setting_notation"`
	Note            string  `json:"note,omitempty"`
}

type MethodRequest struct {
	Grinder string `json:"grinder" jsonschema:"grinder id, name or alias"`
	Method  string `json:"method" jsonschema:"turkish, espresso, moka, v60, pour_over, aeropress, steep_release, cupping, filter_machine, french_press or cold_brew; brewer names such as chemex, origami, switch, clever or moccamaster also work"`
}

type MethodResponse struct {
	Grinder         string `json:"grinder"`
	Method          string `json:"method"`
	MicronRange     string `json:"micron_range"`
	SettingFrom     string `json:"setting_from"`
	SettingMiddle   string `json:"setting_middle"`
	SettingTo       string `json:"setting_to"`
	SettingNotation string `json:"setting_notation"`
	Note            string `json:"note"`
}

type MicronsRequest struct {
	Grinder string  `json:"grinder" jsonschema:"grinder id, name or alias"`
	Microns float64 `json:"microns" jsonschema:"target particle size in microns, e.g. 650"`
}

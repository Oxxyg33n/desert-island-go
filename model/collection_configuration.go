package model

type CollectionConfiguration struct {
	Rareness []RarenessConfiguration `json:"rareness"`
	Layers   LayersConfiguration     `json:"layers"`
}

type RarenessConfiguration struct {
	Name   Rareness `json:"name"`
	Chance float32  `json:"chance"`
}

type LayersConfiguration struct {
	SkipMultiple bool                 `json:"skip_multiple"`
	Groups       []GroupConfiguration `json:"groups"`
}

type GroupConfiguration struct {
	Name       string  `json:"name"`
	Priority   int     `json:"priority"`
	CanSkip    bool    `json:"can_skip"`
	SkipChance float32 `json:"skip_chance,omitempty"`
}

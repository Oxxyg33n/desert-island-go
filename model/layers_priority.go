package model

type LayersPriority struct {
	Layers []LayerPriority `json:"layers_priority"`
}

type LayerPriority struct {
	Name     string `json:"name"`
	Priority string `json:"priority"`
}

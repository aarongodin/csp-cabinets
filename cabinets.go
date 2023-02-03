package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type Input struct {
	Config   Config    `json:"config"`
	Cabinets []Cabinet `json:"cabinets"`
}

type Config struct {
	Kerf        float64 `json:"kerf"`
	PanelWidth  float64 `json:"panel_width"`
	PanelHeight float64 `json:"panel_height"`
}

type Cabinet struct {
	ID     string  `json:"id"`
	Kind   string  `json:"kind"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Depth  float64 `json:"depth"`
}

func (cab Cabinet) Title() string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		cab.Kind,
		strconv.FormatFloat(cab.Width, 'f', -1, 64),
		strconv.FormatFloat(cab.Height, 'f', -1, 64),
		strconv.FormatFloat(cab.Depth, 'f', -1, 64),
	)
}

type Cut struct {
	CabinetID string
	X         float64
	Y         float64
}

type Offcut struct {
	PanelID int
	X       float64
	Y       float64
	PosX    float64
	PosY    float64
	Used    bool
}

type CutFit struct {
}

type ResultSet struct {
	FailureReason string `json:"failure_reason"`
	Config        Config `json:"config"`
	Cuts          []Cut  `json:"cuts"`
}

// AddPanel will increment the number of panels to accomodate the cut.
// Offcuts are processed and added to the list.
func (r *ResultSet) AddPanel(cut Cut) {
	// not sure if I'll need this here?
}

// Waste returns the offcuts that were not used.
func (r ResultSet) Waste() []Offcut {
	// todo
	return nil
}

func (r ResultSet) WasteSqIn() float64 {
	waste := r.Waste()
	wasteSqIn := 0.0
	for _, item := range waste {
		wasteSqIn += item.X * item.Y
	}
	return wasteSqIn
}

func ProcessCabinets(rawInput []byte) (*ResultSet, error) {
	var input Input
	if err := json.Unmarshal(rawInput, &input); err != nil {
		return nil, fmt.Errorf("error reading input data: %w", err)
	}

	// Assign IDs to each cabinet from the input
	for idx, cab := range input.Cabinets {
		if cab.ID == "" {
			id, err := nanoid.New()
			if err != nil {
				return nil, fmt.Errorf("error creating cabinet ID: %w", err)
			}
			input.Cabinets[idx].ID = id
		}
	}

	res := &ResultSet{
		Config: input.Config,
		Cuts:   make([]Cut, 0),
	}

	for _, cab := range input.Cabinets {
		res.Cuts = append(res.Cuts, createCuts(cab)...)
	}

	if err := solve(res); err != nil {

	}

	return res, nil
}

func solve(res *ResultSet) error {
	// TODO
	return nil
}

func createCuts(cab Cabinet) []Cut {
	cuts := make([]Cut, 0)
	switch cab.Kind {
	case "base":
		cuts = append(
			cuts,
			// sides
			Cut{cab.ID, cab.Depth, cab.Height - 4},
			Cut{cab.ID, cab.Depth, cab.Height - 4},
			// bottom
			Cut{cab.ID, cab.Width - 1.5, cab.Depth},
			// toe kick face
			Cut{cab.ID, cab.Width, 4},
			Cut{cab.ID, cab.Width, 4},
			// toe kick side
			Cut{cab.ID, cab.Depth - 3 - 1.5, 4},
			Cut{cab.ID, cab.Depth - 3 - 1.5, 4},
			// stretchers/nailers
			Cut{cab.ID, cab.Width - 1.5, 4},
			Cut{cab.ID, cab.Width - 1.5, 4},
			Cut{cab.ID, cab.Width - 1.5, 4},
			Cut{cab.ID, cab.Width - 1.5, 4},
		)
	case "wall":
		cuts = append(
			cuts,
			// sides
			Cut{cab.ID, cab.Depth, cab.Height},
			Cut{cab.ID, cab.Depth, cab.Height},
			// bottom/top
			Cut{cab.ID, cab.Width - 1.5, cab.Depth},
			Cut{cab.ID, cab.Width - 1.5, cab.Depth},
			// nailers
			Cut{cab.ID, cab.Width - 1.5, 4},
			Cut{cab.ID, cab.Width - 1.5, 4},
		)
	}
	return cuts
}

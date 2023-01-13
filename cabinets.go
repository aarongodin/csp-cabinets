package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

const PANEL_X = 96.0
const PANEL_Y = 48.0

type Cabinet struct {
	Kind   string
	Width  float64
	Height float64
	Depth  float64
}

type ResultSet struct {
	// Required cuts for the input set of cabinets
	Cuts []Cut
	// Count of standard 4' x 8' panels required to make the desired cuts
	Panels int32
	// Set of offcuts used as additional material throughout the
	Offcuts []Offcut
}

// AddPanel will increment the number of panels to accomodate the cut.
// Offcuts are processed and added to the list.
func (r *ResultSet) AddPanel(cut Cut) {
	r.Panels = r.Panels + 1
	remainderX := PANEL_X - cut.X
	remainderY := PANEL_Y - cut.Y
	r.Offcuts = append(
		r.Offcuts,
		Offcut{remainderX, cut.Y, false},
		Offcut{PANEL_X, remainderY, false},
	)
}

// Waste returns the offcuts that were not used.
func (r ResultSet) Waste() []Offcut {
	waste := make([]Offcut, 0)
	for _, offcut := range r.Offcuts {
		if !offcut.Used {
			waste = append(waste, offcut)
		}
	}
	return waste
}

func (r ResultSet) WasteSqIn() float64 {
	waste := r.Waste()
	wasteSqIn := 0.0
	for _, item := range waste {
		wasteSqIn += item.X * item.Y
	}
	return wasteSqIn
}

type Cut struct {
	CabinetRef string
	X          float64
	Y          float64
}

type Offcut struct {
	X    float64
	Y    float64
	Used bool
}

var calculateCabinets = &cobra.Command{
	Use:  "cabinet-calc",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		raw, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("error reading file: %s", err)
		}
		var data [][]any
		if err := json.Unmarshal(raw, &data); err != nil {
			log.Fatalf("error reading json: %s", err)
		}

		cabs := make([]Cabinet, 0)
		for _, item := range data {
			cab := Cabinet{
				Kind:   item[0].(string),
				Width:  item[1].(float64),
				Height: item[2].(float64),
				Depth:  item[3].(float64),
			}
			cabs = append(cabs, cab)
		}

		cuts := createCuts(cabs)
		// Modify the cuts to ensure that the largest dimension is on the x-axis
		for idx, cut := range cuts {
			if cut.X < cut.Y {
				cuts[idx] = Cut{cut.CabinetRef, cut.Y, cut.X}
			}
		}

		// Sort the cuts according to those with the longest side first
		sort.Slice(cuts, func(i, j int) bool {
			first := cuts[i]
			second := cuts[j]
			return first.X > second.X
		})

		res := &ResultSet{
			Cuts:    cuts,
			Panels:  0,
			Offcuts: make([]Offcut, 0),
		}

		for _, cut := range res.Cuts {
			cutMadeFromOffcut := processOffcuts(res, cut)
			if !cutMadeFromOffcut {
				res.AddPanel(cut)
			}
		}

		spew.Dump(res.Cuts)
	},
}

// processOffcuts checks if any offcuts can be used for a given cut.
// If yes, the offcut is consumed with new offcuts appended.
// The return value indicates whether this operation consumed an offcut
func processOffcuts(res *ResultSet, cut Cut) bool {
	for offcutIdx, offcut := range res.Offcuts {
		if offcut.Used {
			continue
		}

		var x, y float64

		// Condition below is for the default orientation
		if offcut.X > cut.X && offcut.Y > cut.Y {
			x = offcut.X
			y = offcut.Y
		}

		// Condition below is for a 90 degree orientation
		if offcut.Y > cut.X && offcut.X > cut.Y {
			x = offcut.Y
			y = offcut.X
		}

		if x != 0 && y != 0 {
			remainderX := x - cut.X
			remainderY := y - cut.Y
			if remainderX != 0 {
				res.Offcuts = append(res.Offcuts, Offcut{remainderX, cut.Y, false})
			}
			if remainderY != 0 {
				res.Offcuts = append(res.Offcuts, Offcut{x, remainderY, false})
			}
			res.Offcuts[offcutIdx].Used = true
			return true
		}
	}

	return false
}

// cref creates a cabinet reference identifier
func cabinetRef(cab Cabinet) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		cab.Kind,
		strconv.FormatFloat(cab.Width, 'f', -1, 64),
		strconv.FormatFloat(cab.Height, 'f', -1, 64),
		strconv.FormatFloat(cab.Depth, 'f', -1, 64),
	)
}

func createCuts(cabs []Cabinet) []Cut {
	cuts := make([]Cut, 0)

	for _, cab := range cabs {
		cref := cabinetRef(cab)
		switch cab.Kind {
		case "base":
			cuts = append(
				cuts,
				// sides
				Cut{cref, cab.Depth, cab.Height - 4},
				Cut{cref, cab.Depth, cab.Height - 4},
				// bottom
				Cut{cref, cab.Width - 1.5, cab.Depth},
				// toe kick face
				Cut{cref, cab.Width, 4},
				Cut{cref, cab.Width, 4},
				// toe kick side
				Cut{cref, cab.Depth - 3 - 1.5, 4},
				Cut{cref, cab.Depth - 3 - 1.5, 4},
				// stretchers/nailers
				Cut{cref, cab.Width - 1.5, 4},
				Cut{cref, cab.Width - 1.5, 4},
				Cut{cref, cab.Width - 1.5, 4},
				Cut{cref, cab.Width - 1.5, 4},
			)
		case "wall":
			cuts = append(
				cuts,
				// sides
				Cut{cref, cab.Depth, cab.Height},
				Cut{cref, cab.Depth, cab.Height},
				// bottom/top
				Cut{cref, cab.Width - 1.5, cab.Depth},
				Cut{cref, cab.Width - 1.5, cab.Depth},
				// nailers
				Cut{cref, cab.Width - 1.5, 4},
				Cut{cref, cab.Width - 1.5, 4},
			)
		}
	}

	return cuts
}

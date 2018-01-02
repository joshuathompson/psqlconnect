package ui

import (
	"fmt"

	"github.com/joshuathompson/psqlconnect/utils"
	"github.com/jroimartin/gocui"
)

func RenderInstructions(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("instructions", -1, maxY-2, maxX, maxY)

	v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
	}

	instructions := utils.RightPaddedString("[Enter] Connect [f] Filter [r] Refresh [j] Down [k] Up [Ctrl-c] Quit/Back", maxX, 2)
	fmt.Fprintf(v, "%s", instructions)

	return nil
}

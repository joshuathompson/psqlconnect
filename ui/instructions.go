package ui

import (
	"fmt"

	"github.com/joshuathompson/pgm/utils"
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
		//v.BgColor = gocui.ColorBlue
	}

	instructions := utils.RightPaddedString("[Enter] Connect [e] Edit [c] Create [d] Delete [Ctrl-c] Quit/Back [j] Down [k] Up [pgUp/pgDn] Reorder", maxX, 2)
	fmt.Fprintf(v, "\u001b[1m%s\u001b[0m", instructions)

	return nil
}

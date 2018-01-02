package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func RenderFilterView(g *gocui.Gui, filter string) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("filter", (maxX / 4), (maxY/2)-5, (maxX/2)+(maxX/4), (maxY/2)-3)

	v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		g.Highlight = true
		g.SelFgColor = gocui.ColorBlue

		v.Frame = true
		v.Editable = true
		v.Title = "Filter?  (blank to remove filter)"

		_, err := g.SetCurrentView("filter")

		if err != nil {
			return err
		}
	}

	if len(filter) > 0 {
		fmt.Fprint(v, filter)
		v.SetCursor(len(filter), 0)
	}

	return nil
}

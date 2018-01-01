package ui

import (
	"fmt"

	"github.com/joshuathompson/pgm/pgpass"
	"github.com/jroimartin/gocui"
)

func RenderDeleteConfirmation(g *gocui.Gui, connection *pgpass.Connection) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("delete_confirmation", (maxX / 6), (maxY/2)-5, (maxX/2)+(maxX/3), (maxY/2)-3)

	v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Frame = true
		v.Title = fmt.Sprintf("Re-enter name '%s' to delete. Ctrl-c to cancel.", connection.Name)

		g.Highlight = true
		g.SelFgColor = gocui.ColorRed

		_, err = g.SetCurrentView("delete_confirmation")

		if err != nil {
			return err
		}
	}

	return nil
}

package ui

import (
	"fmt"

	"github.com/joshuathompson/pgm/pgpass"
	"github.com/joshuathompson/pgm/utils"
	"github.com/jroimartin/gocui"
)

func RenderConnectionsView(g *gocui.Gui, connections []*pgpass.Connection) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("connections", -1, 2, maxX, maxY-1)

	v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorWhite
		v.SelFgColor = gocui.ColorBlack

		_, err = g.SetCurrentView("connections")

		if err != nil {
			return err
		}
	}

	if len(connections) > 0 {
		for _, c := range connections {
			name := utils.RightPaddedString(c.Name, maxX/6, 2)
			host := utils.RightPaddedString(c.Host, maxX/3, 2)
			port := utils.RightPaddedString(c.Port, maxX/10, 2)
			db := utils.RightPaddedString(c.Database, maxX/5, 2)
			username := utils.RightPaddedString(c.Username, maxX/5, 2)

			fmt.Fprintf(v, "%s %s %s %s %s\n", name, host, port, db, username)
		}
	} else {
		fmt.Fprint(v, utils.RightPaddedString("No connections found in ~/.pgpass.", maxX, 2))
	}

	return nil
}

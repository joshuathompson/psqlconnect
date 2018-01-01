package ui

import (
	"fmt"

	"github.com/joshuathompson/pgm/utils"
	"github.com/jroimartin/gocui"
)

func RenderHeaderView(g *gocui.Gui) error {
	maxX, _ := g.Size()

	v, err := g.SetView("header", -1, -1, maxX, 2)

	v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
	}

	// TOP ROW
	title := utils.RightPaddedString("Postgres/Redshift Connections", maxX, 2)

	fmt.Fprintf(v, "\u001b[1m%s\u001b[0m\n", title)

	//BOTTOM ROW
	nameHeader := utils.RightPaddedString("NAME", maxX/6, 2)
	hostHeader := utils.RightPaddedString("HOST", maxX/3, 2)
	portHeader := utils.RightPaddedString("PORT", maxX/10, 2)
	dbHeader := utils.RightPaddedString("DATABASE", maxX/5, 2)
	usernameHeader := utils.RightPaddedString("USERNAME", maxX/5, 2)

	fmt.Fprintf(v, "\u001b[1m%s %s %s %s %s\u001b[0m", nameHeader, hostHeader, portHeader, dbHeader, usernameHeader)

	return nil
}

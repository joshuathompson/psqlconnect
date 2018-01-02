package ui

import (
	"fmt"

	"github.com/joshuathompson/psqlconnect/utils"
	"github.com/jroimartin/gocui"
)

func RenderHeaderView(g *gocui.Gui, numberOfConnections int, filter string) error {
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
	title := utils.RightPaddedString("Connections", (maxX / 2), 2)

	var infoSection string
	if len(filter) > 0 {
		infoSection = utils.LeftPaddedString(fmt.Sprintf("Total: %d  |  Active filter: %s", numberOfConnections, filter), (maxX / 2), 2)
	} else {
		infoSection = utils.LeftPaddedString(fmt.Sprintf("Total: %d", numberOfConnections), (maxX / 2), 2)
	}

	fmt.Fprintf(v, "%s %s\n", title, infoSection)

	//BOTTOM ROW
	nameHeader := utils.RightPaddedString("NAME", maxX/6, 2)
	hostHeader := utils.RightPaddedString("HOST", maxX/3, 2)
	portHeader := utils.RightPaddedString("PORT", maxX/10, 2)
	dbHeader := utils.RightPaddedString("DATABASE", maxX/5, 2)
	usernameHeader := utils.RightPaddedString("USERNAME", maxX/5, 2)

	fmt.Fprintf(v, "%s %s %s %s %s", nameHeader, hostHeader, portHeader, dbHeader, usernameHeader)

	return nil
}

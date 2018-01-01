package ui

import (
	"fmt"

	"github.com/joshuathompson/pgm/pgpass"
	"github.com/jroimartin/gocui"
)

func RenderCreateEdit(g *gocui.Gui, connection *pgpass.Connection) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("create_edit", (maxX / 6), (maxY/2)-5, (maxX/2)+(maxX/3), (maxY/2)-3)

	v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Frame = true
		v.Title = "Enter in form: name:host:port:db:user:pw. Ctrl-c to cancel."

		g.Highlight = true
		g.SelFgColor = gocui.ColorGreen

		_, err = g.SetCurrentView("create_edit")

		if err != nil {
			return err
		}
	}

	if connection != nil {
		pgpassString := fmt.Sprintf("%s:%s:%s:%s:%s:%s", connection.Name, connection.Host, connection.Port, connection.Database, connection.Username, connection.Password)
		fmt.Fprintf(v, pgpassString)
		v.SetCursor(len(pgpassString), 0)
	}

	return nil
}

package main

import (
	"log"
	"strings"

	"github.com/joshuathompson/psqlconnect/pgpass"
	"github.com/joshuathompson/psqlconnect/ui"
	"github.com/jroimartin/gocui"
)

var (
	connections         []*pgpass.Connection
	filteredConnections []*pgpass.Connection
	selectedIndex       = 0
	connectOnExit       = false
	filter              = ""
)

func main() {
	var err error

	connections, err = pgpass.LoadConnectionsFromPgpass()

	if err != nil {
		log.Fatal(err)
	}

	filteredConnections = connections

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}

	g.Cursor = true
	g.SetManagerFunc(layout)

	err = keybindings(g)

	if err != nil {
		g.Close()
		log.Fatal(err)
	}

	err = g.MainLoop()

	if err != nil && err != gocui.ErrQuit {
		g.Close()
		log.Fatal(err)
	}

	g.Close()

	if connectOnExit {
		pgpass.ConnectToDatabase(connections[selectedIndex])
	}
}

func layout(g *gocui.Gui) error {
	err := ui.RenderHeaderView(g, len(filteredConnections), filter)
	err = ui.RenderConnectionsView(g, filteredConnections)
	err = ui.RenderInstructions(g)

	if err != nil {
		return err
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})

	err = g.SetKeybinding("connections", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if selectedIndex != 0 {
			selectedIndex--
			v.MoveCursor(0, -1, false)
		}
		return nil
	})

	err = g.SetKeybinding("connections", 'k', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if selectedIndex != 0 {
			selectedIndex--
			v.MoveCursor(0, -1, false)
		}
		return nil
	})

	err = g.SetKeybinding("connections", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if selectedIndex != len(filteredConnections)-1 {
			selectedIndex++
			v.MoveCursor(0, 1, false)
		}
		return nil
	})

	err = g.SetKeybinding("connections", 'j', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if selectedIndex != len(filteredConnections)-1 {
			selectedIndex++
			v.MoveCursor(0, 1, false)
		}
		return nil
	})

	err = g.SetKeybinding("connections", 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := ui.RenderFilterView(g, filter)
		return err
	})

	err = g.SetKeybinding("connections", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		connectOnExit = true
		return gocui.ErrQuit
	})

	err = g.SetKeybinding("connections", 'r', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		connections, err = pgpass.LoadConnectionsFromPgpass()

		if err != nil {
			log.Fatal(err)
		}

		if len(filter) > 0 {
			filteredConnections = pgpass.GetFilteredConnections(connections, filter)
		} else {
			filteredConnections = connections
		}

		return nil
	})

	err = g.SetKeybinding("filter", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		buffer := v.Buffer()

		err := g.DeleteView("filter")

		if err != nil {
			return err
		}

		filter = strings.TrimSuffix(strings.TrimSpace(buffer), "\n")

		g.Highlight = false
		g.SelFgColor = gocui.ColorWhite
		cv, err := g.SetCurrentView("connections")

		if len(filter) > 0 {
			filteredConnections = pgpass.GetFilteredConnections(connections, filter)
			selectedIndex = 0
			err = cv.SetCursor(0, 0)
		} else {
			filteredConnections = connections
		}

		return err
	})

	return err
}

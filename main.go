package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joshuathompson/pgm/pgpass"
	"github.com/joshuathompson/pgm/ui"
	"github.com/jroimartin/gocui"
)

var (
	connections   []*pgpass.Connection
	selectedIndex = 0
	connectOnExit = false
	editing       = false
)

func main() {
	var err error

	connections, err = pgpass.LoadConnectionsFromPgpass()

	if err != nil {
		log.Fatal(err)
	}

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
		connectToDatabase()
	}
}

func layout(g *gocui.Gui) error {
	err := ui.RenderHeaderView(g)
	err = ui.RenderConnectionsView(g, connections)
	err = ui.RenderInstructions(g)

	if err != nil {
		return err
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("connections", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})

	err = g.SetKeybinding("connections", gocui.KeyPgup, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(connections) > 1 && selectedIndex > 0 {
			tmp := connections[selectedIndex-1]
			connections[selectedIndex-1] = connections[selectedIndex]
			connections[selectedIndex] = tmp

			err := pgpass.SaveConnectionsToPgPass(connections)

			if err != nil {
				return err
			}

			selectedIndex--
			v.MoveCursor(0, -1, false)
		}

		return nil
	})

	err = g.SetKeybinding("connections", gocui.KeyPgdn, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(connections) > 1 && selectedIndex != len(connections)-1 {
			tmp := connections[selectedIndex+1]
			connections[selectedIndex+1] = connections[selectedIndex]
			connections[selectedIndex] = tmp

			err := pgpass.SaveConnectionsToPgPass(connections)

			if err != nil {
				return err
			}

			selectedIndex++
			v.MoveCursor(0, 1, false)
		}

		return nil
	})

	err = g.SetKeybinding("connections", 'c', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		newConnection := &pgpass.Connection{
			Name:     "new",
			Host:     "/var/run/postgresql",
			Port:     "5432",
			Database: "db",
			Username: "user",
			Password: "pw",
		}

		connections = append(connections, newConnection)
		selectedIndex = len(connections) - 1
		v.SetCursor(0, selectedIndex)

		ui.RenderCreateEdit(g, newConnection)

		return nil
	})

	err = g.SetKeybinding("connections", 'e', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(connections) > 0 {
			editing = true
			ui.RenderCreateEdit(g, connections[selectedIndex])
		}

		return nil
	})

	err = g.SetKeybinding("connections", 'd', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(connections) > 0 {
			ui.RenderDeleteConfirmation(g, connections[selectedIndex])
		}

		return nil
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
		if selectedIndex != len(connections)-1 {
			selectedIndex++
			v.MoveCursor(0, 1, false)
		}
		return nil
	})

	err = g.SetKeybinding("connections", 'j', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if selectedIndex != len(connections)-1 {
			selectedIndex++
			v.MoveCursor(0, 1, false)
		}
		return nil
	})

	err = g.SetKeybinding("connections", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		connectOnExit = true
		return gocui.ErrQuit
	})

	err = g.SetKeybinding("create_edit", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err = g.DeleteView("create_edit")

		if err != nil {
			return err
		}

		g.Highlight = false
		g.SelFgColor = gocui.ColorWhite

		cv, err := g.SetCurrentView("connections")

		if err != nil {
			return err
		}

		if !editing {
			connections = append(connections[:selectedIndex], connections[selectedIndex+1:]...)
			selectedIndex--
			cv.MoveCursor(0, -1, false)
		}

		editing = false

		return nil
	})

	err = g.SetKeybinding("create_edit", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		currentConnection := connections[selectedIndex]
		buffer := strings.Split(strings.TrimSuffix(strings.TrimSpace(v.Buffer()), "\n"), ":")

		if len(buffer) == 6 {
			currentConnection.Name = buffer[0]
			currentConnection.Host = buffer[1]
			currentConnection.Port = buffer[2]
			currentConnection.Database = buffer[3]
			currentConnection.Username = buffer[4]
			currentConnection.Password = buffer[5]

			err := pgpass.SaveConnectionsToPgPass(connections)

			if err != nil {
				return err
			}
		}

		err = g.DeleteView("create_edit")

		if err != nil {
			return err
		}

		g.Highlight = false
		g.SelFgColor = gocui.ColorWhite

		_, err = g.SetCurrentView("connections")

		if err != nil {
			return err
		}

		return nil
	})

	err = g.SetKeybinding("delete_confirmation", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err = g.DeleteView("delete_confirmation")

		if err != nil {
			return err
		}

		g.Highlight = false
		g.SelFgColor = gocui.ColorWhite

		_, err = g.SetCurrentView("connections")

		if err != nil {
			return err
		}

		return nil
	})

	err = g.SetKeybinding("delete_confirmation", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		currentConnection := connections[selectedIndex]

		if currentConnection.Name == strings.TrimSuffix(strings.TrimSpace(v.Buffer()), "\n") {
			newConnections := append(connections[:selectedIndex], connections[selectedIndex+1:]...)

			err := pgpass.SaveConnectionsToPgPass(newConnections)

			if err != nil {
				return err
			}

			cv, err := g.View("connections")

			if err != nil {
				return err
			}

			if selectedIndex > 0 {
				selectedIndex--
				cv.MoveCursor(0, -1, false)
			}
			connections = newConnections
		}

		err = g.DeleteView("delete_confirmation")

		if err != nil {
			return err
		}

		g.Highlight = false
		g.SelFgColor = gocui.ColorWhite

		_, err = g.SetCurrentView("connections")

		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func connectToDatabase() {
	scanner := bufio.NewScanner(os.Stdin)
	connection := connections[selectedIndex]
	cmdString := "psql"

	var args []string

	if len(connection.Host) > 0 {
		if connection.Host == "*" {
			fmt.Print("Host was a wildcard, enter a host name: ")
			scanner.Scan()
			connection.Host = strings.TrimSpace(scanner.Text())
		}

		args = append(args, "-h")
		args = append(args, connection.Host)
		cmdString = fmt.Sprintf("%s -h %s", cmdString, connection.Host)
	}

	if len(connection.Port) > 0 {
		if connection.Port == "*" {
			fmt.Print("Port was a wildcard, enter a port number: ")
			scanner.Scan()
			connection.Port = strings.TrimSpace(scanner.Text())
		}

		args = append(args, "-p")
		args = append(args, connection.Port)
		cmdString = fmt.Sprintf("%s -p %s", cmdString, connection.Port)
	}

	if len(connection.Database) > 0 {
		if connection.Database == "*" {
			fmt.Print("Database was a wildcard, enter a database name: ")
			scanner.Scan()
			connection.Database = strings.TrimSpace(scanner.Text())
		}

		args = append(args, "-d")
		args = append(args, connection.Database)
		cmdString = fmt.Sprintf("%s -d %s", cmdString, connection.Database)
	}

	if len(connection.Username) > 0 {
		if connection.Username == "*" {
			fmt.Print("Username was a wildcard, enter a username: ")
			scanner.Scan()
			connection.Username = strings.TrimSpace(scanner.Text())
		}

		args = append(args, "-U")
		args = append(args, connection.Username)
		cmdString = fmt.Sprintf("%s -U %s", cmdString, connection.Username)
	}

	fmt.Printf("Running command %s\n", cmdString)

	cmd := exec.Command("psql", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

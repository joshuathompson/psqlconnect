package pgpass

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
)

type Connection struct {
	Name     string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func LoadConnectionsFromPgpass() (connections []*Connection, err error) {
	usr, err := user.Current()

	if err != nil {
		return connections, err
	}

	file, err := os.OpenFile(usr.HomeDir+"/.pgpass", os.O_RDONLY|os.O_CREATE, 0600)

	if err != nil {
		return connections, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	newConnection := &Connection{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			if strings.HasPrefix(line, "# NAME=") {
				newConnection.Name = strings.Replace(line, "# NAME=", "", 1)
			}
		} else {
			split := strings.Split(line, ":")

			if len(split) != 5 {
				err = fmt.Errorf("Line %s is malformed, should be in the form hostname:port:database:username:password", line)
				return connections, err
			}

			newConnection.Host = split[0]
			newConnection.Port = split[1]
			newConnection.Database = split[2]
			newConnection.Username = split[3]
			newConnection.Password = split[4]

			connections = append(connections, newConnection)

			newConnection = &Connection{}
		}
	}

	err = scanner.Err()
	if err != nil {
		return connections, err
	}

	return connections, err
}

func SaveConnectionsToPgPass(connections []*Connection) error {
	usr, err := user.Current()

	if err != nil {
		return err
	}

	file, err := os.OpenFile(usr.HomeDir+"/.pgpass", os.O_RDWR|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}
	defer file.Close()

	for _, c := range connections {
		fmt.Fprintf(file, "# NAME=%s\n", c.Name)
		fmt.Fprintf(file, "%s:%s:%s:%s:%s\n", c.Host, c.Port, c.Database, c.Username, c.Password)
	}

	return nil
}

package pgpass

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func GetFilteredConnections(connections []*Connection, filter string) (filteredConnections []*Connection) {
	filterLower := strings.ToLower(filter)

	for _, c := range connections {
		nameLower := strings.ToLower(c.Name)
		hostLower := strings.ToLower(c.Host)
		portLower := strings.ToLower(c.Port)
		databaseLower := strings.ToLower(c.Database)
		usernameLower := strings.ToLower(c.Username)

		if strings.Contains(nameLower, filterLower) ||
			strings.Contains(hostLower, filterLower) ||
			strings.Contains(portLower, filterLower) ||
			strings.Contains(databaseLower, filterLower) ||
			strings.Contains(usernameLower, filterLower) {
			filteredConnections = append(filteredConnections, c)
		}

	}

	return filteredConnections
}

func ConnectToDatabase(connection *Connection) {
	scanner := bufio.NewScanner(os.Stdin)
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

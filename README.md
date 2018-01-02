# psqlconnect
> Interface to quickly build and run a psql command from your ~/.pgpass file  

[![asciicast](https://asciinema.org/a/UDHxTcmiRSmOozpAsoL5jimIn.png)](https://asciinema.org/a/UDHxTcmiRSmOozpAsoL5jimIn)

## Install
Fetch the latest release for your platform [from the following page](https://github.com/joshuathompson/psqlconnect/releases).

## Why would I want to use this?
I built this for myself after getting tired of copy pasting information from 1Password / looking through terminal history in order to connect to client Redshift/Postgres databases.  The [~/.pgpass file](https://www.postgresql.org/docs/9.3/static/libpq-pgpass.html) simplifies the situation but it still takes time to open that file, find the connection that I wanted, and write a psql command to connect.  This tool simplifies all that to selecting a connection from a table.

## How do I add connections?
Connections are handled by a ~/.pgpass file as described by the [Postgres docs](https://www.postgresql.org/docs/9.3/static/libpq-pgpass.html).  Each entry should be in the following format:

```
# NAME=<your connection name>
hostname:port:database:username:password
```

The commented line with `NAME=` is optional but will give you an extra descriptor about the database in psqlconnect.  An example pgpass file is included.

### TUI Keybinds

Keybind              | Description
---------------------|---------------------------------------
<kbd>j</kbd>         | move the cursor down a line
<kbd>k</kbd>         | move the cursor up a line
<kbd>Enter</kbd>     | run psql with the selected connection
<kbd>f</kbd>         | add/remove a filter
<kbd>r</kbd>         | refresh connections from ~/.pgpass
<kbd>Ctrl+c</kbd>    | quit

## License
MIT

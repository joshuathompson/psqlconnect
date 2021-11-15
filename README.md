# psqlconnect
> Interface to quickly run a psql command from your ~/.pgpass file  

[![asciicast](https://asciinema.org/a/UDHxTcmiRSmOozpAsoL5jimIn.png)](https://asciinema.org/a/UDHxTcmiRSmOozpAsoL5jimIn)

## Install
Fetch the latest release for your platform [from the following page](https://github.com/joshuathompson/psqlconnect/releases).

## How do I add connections?
Connections are handled by a ~/.pgpass file as described by the [Postgres docs](https://www.postgresql.org/docs/9.3/static/libpq-pgpass.html).  Each entry should be in the following format:

```
# NAME=<your connection name>
hostname:port:database:username:password
```

The commented line with `NAME=` is optional but will give you an extra descriptor about the database in psqlconnect.  An example pgpass file is included.  Make sure that your ~/.pgpass has permissions of 0600.

## Security
The `.pgpass` file stores all information in plain text including usernames and passwords. Arguably this isn't a huge deal because if someone has access to your computer you're compromised but if you need an encrypted file/database you won't get it from this solution.

### Keybinds

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

package nodedb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // registers as sql client
)

// Connection provide access to the database
type Connection struct {
	DB      *sql.DB
	RootID  int64
	Created time.Time
	Start   time.Time
	End     time.Time
	Deleted time.Time
	logging bool
}

func newConnection(db *sql.DB, rootID int64) (*Connection, error) {
	connection := &Connection{}
	connection.DB = db
	connection.RootID = rootID
	connection.Created = time.Now()
	connection.Start = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	connection.End = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	connection.Deleted = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	connection.logging = false

	return connection, nil
}

//Open creates a new Connection Object and opens the underlying sqlite database
func Open(filename string) (*Connection, error) {

	DB, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return newConnection(DB, 1)
}

//Close closes the underlying sqlite database
func (c *Connection) Close() error {
	return c.DB.Close()
}

//Collections returs a string of collection names
func (c *Connection) Collections() (collections []string, err error) {
	q, err := c.DB.Query(SQLCollections)
	c.log(SQLCollections)
	if err != nil {
		return nil, err
	}

	defer q.Close()

	for q.Next() {
		var tableName string
		if err := q.Scan(&tableName); err != nil {
			return nil, err
		}
		collections = append(collections, tableName)
	}

	return collections, nil
}

//Collection opens (and creates) collections
func (c *Connection) Collection(name string) (*Collection, error) {
	col := &Collection{c, name}

	if name == "sqlite_sequence" {
		return nil, errors.New("reserved name")
	}

	c.log(SQLCollections, "Parameter:", name)

	rows, err := c.DB.Query(SQLCollection, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var table, tabletype string
	if rows.Next() { // only pick the first row
		rows.Scan(&table, &tabletype)
		if tabletype != "table" {
			return nil, errors.New("wrong db object type")
		}
	} else {
		_, err := c.DB.Exec(SQLString(SQLNewCollection).Collection(name).String())
		if err != nil {
			return nil, err
		}
	}

	return col, nil
}

//RemoveCollection removes collections, it is destructive, no undo possible!
func (c *Connection) RemoveCollection(name string) error {

	sql := SQLString(SQLDropCollection).Collection(name).String()
	c.log(sql)
	_, err := c.DB.Exec(sql) // TODO Should use transactions here !!!
	return err
}

//CloneCollection makes a full copy of source collection in a new clone collection
func (c *Connection) CloneCollection(sourceName string, cloneName string) error {
	cols, err := c.Collections()
	if err != nil {
		return err
	}
	_source, _clone := "", ""
	for _, colname := range cols {
		if colname == sourceName {
			_source = sourceName
		}
		if colname == cloneName {
			_clone = cloneName
		}
	}
	if _source == "" {
		return fmt.Errorf("collection '%s' does not exist", sourceName)
	}
	if _clone != "" {
		return fmt.Errorf("collection '%s' already exists", _clone)
	}

	_clone = cloneName

	_, err = c.Collection(_clone)
	if err != nil {
		return err
	}

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	sql := processString(SQLCloneCollection, map[string]string{"targetCollection": _clone, "sourceCollection": _source})
	c.log(sql)
	_, err = tx.Exec(sql)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}

// Ping the underlying data base
func (c *Connection) Ping() error {
	return c.DB.Ping()
}

// SetLogging to true or false, default is false
func (c *Connection) SetLogging(logging bool) { // LoggingEnabled returns true if logging is enabled, false otherwise.
	c.logging = logging
}

// LoggingEnabled returns true if logging is enabled, false otherwise
func (c *Connection) LoggingEnabled() bool {
	return c.logging

}

func (c *Connection) log(v ...interface{}) {
	if c.logging {
		log.Println(v...)
	}
}

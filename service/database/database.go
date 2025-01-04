/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
// all the function that I creat in db nust be declated here
type AppDatabase interface {
	CreateUser(username string) (User, error)
	GetUser(username string) (User, error)
	SetToken(user_id int, token string) (err error)
	DeleteToken(token string) (err error)
	GetMyConversations_db(token string) (conversations []Conversation, err error)	
	AddUsersToConversation(user_id int, conversation_id int) (err error)
	CreateConversation_db(isGroup bool, name string, photo string) (conversation Conversation, err error)	
	GetUserId(token string) (user UserToken, err error)
	UpdateUserName(id int, newname string) (err error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func createDatabase(db *sql.DB) error {
	tables := [5] string {
		`CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
			name VARCHAR(25) NOT NULL,
			photo VARCHAR(255)
		);`,

		`CREATE TABLE IF NOT EXISTS conversations(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
			lastconvo TIMESTAMP NOT NULL,
			is_group BOOLEAN DEFAULT FALSE,
			photo VARCHAR(255),
			name VARCHAR(255)


		);`,
		
		`CREATE TABLE IF NOT EXISTS convmembers (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
			conversation_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (conversation_id) REFERENCES conversations (id),
			FOREIGN KEY (user_id) REFERENCES users (id)
		);`,

		`CREATE TABLE IF NOT EXISTS messages(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
			datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			content  TEXT NOT NULL,
			sender	 INTEGER NOT NULL,
			conversation_id	 INTEGER NOT NULL,
			status VARCHAR(10) DEFAULT'sent',
			FOREIGN KEY(sender) REFERENCES users(id),
			FOREIGN KEY(conversation_id) REFERENCES conversation(id)



		);`,
		`CREATE TABLE IF NOT EXISTS usertokens(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
			user_id INTEGER NOT NULL,
			token VARCHAR(64),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
	}
	for t := 0; t < len(tables); t++ {
		sqlStmt := tables[t]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}

	return nil
}
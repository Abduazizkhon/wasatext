
// functions that generate the tables 

package database

import (
 "database/sql"
 "fmt"
 "github.com/gofrs/uuid"
 "log"
)

func CreateTables(db *sql.DB) error {
 if err := CreateUsersTable(db); err != nil {
  return fmt.Errorf("creating users table: %w", err)
 }
 if err := createConversationsTable(db); err != nil {
  return fmt.Errorf("creating Conversations table: %w", err)
 }
 if err := createConversationBodyTable(db); err != nil {
  return fmt.Errorf("creating conversations body table: %w", err)
 }
 return nil
}
func CreateUsersTable(db *sql.DB) error {
 query := `
 CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL
 );`
 _, err := db.Exec(query)
 if err != nil {
  log.Printf("Error creating users table: %v", err)
  return err
 }
 log.Println("Users table created successfully")
 return nil
}

func createConversationsTable(db *sql.DB) error {
 query := `
 CREATE TABLE IF NOT EXISTS conversations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  is_group BOOLEAN NOT NULL DEFAULT FALSE,
  picture TEXT NOT NULL UNIQUE,
  latest_message TEXT,
  snippet TEXT,
  uuid TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL
 );`
 _, err := db.Exec(query)
 if err != nil {
  log.Printf("Error creating conversations table: %v", err)
  return err
 }
 log.Println("Conversations table created successfully")
 return nil
}

func createConversationBodyTable(db *sql.DB) error {
 query := `
 CREATE TABLE IF NOT EXISTS conversation_body (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  conversation_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  message TEXT NOT NULL,
  FOREIGN KEY(conversation_id) REFERENCES conversations(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
 );`
 _, err := db.Exec(query)
 if err != nil {
  log.Printf("Error creating conversation_body table: %v", err)
  return err
 }
 log.Println("Conversation_body table created successfully")
 return nil
}

func InsertUser(db *sql.DB, id uuid.UUID, username string) error {
 query := INSERT INTO users (uuid, username) VALUES (?, ?)
 _, err := db.Exec(query, id.String(), username)
 return err
}

func GetUsers(db *sql.DB) ([]User, error) {
 query := SELECT uuid, username FROM users
 rows, err := db.Query(query)
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 var users []User
 for rows.Next() {
  var user User
  var uuidStr string
  if err := rows.Scan(&uuidStr, &user.UserName); err != nil {
   return nil, err
  }
  user.Id, err = uuid.FromString(uuidStr)
  if err != nil {
   return nil, err
  }
  users = append(users, user)
 }
 if err := rows.Err(); err != nil {
  return nil, err
 }
 return users, nil
}
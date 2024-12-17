package database

import (
 "fmt"
)

func (db *appdbimpl) GetUsers() ([]User, error) {
 var Users []User
 err := db.c.QueryRow(`SELECT uuid, username FROM users `).Scan(&Users)
 if err != nil || len(Users) == 0 {
  return Users, fmt.Errorf("No users exist: %w", err)
 }
 return Users, nil
}
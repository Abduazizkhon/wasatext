package database

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
)

// func (db *appdbimpl) CreateUser(username string) (user User, err error) {

// 	query := `INSERT INTO users ( name ) VALUES (?);`
// 	_, err = db.c.Exec(query, username)
// 	if err != nil {
// 		return
// 	}
// 	return db.GetUser(username)
// }

func (db *appdbimpl) CreateUser(username string) (User, error) {
	// Start a transaction
	tx, err := db.c.Begin()
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback() // Rollback the transaction if it's not committed

	// Generate a UUID
	id, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}

	// Insert the new user
	query := `INSERT INTO users (id, name) VALUES (?, ?);`
	_, err = tx.Exec(query, id.String(), username)
	if err != nil {
		return User{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return User{}, err
	}

	// Retrieve the created user
	return db.GetUser(username)
}

func (db *appdbimpl) GetUser(username string) (user User, err error) {
	query := `SELECT id, name FROM users WHERE name = ?;`
	row := db.c.QueryRow(query, username)

	// Attempt to scan the result into the User struct
	err = row.Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		// Return a clear error if no user is found
		return
	} else if err != nil {
		// Return the database error for debugging
		return
	}

	// Successfully found the user
	return user, nil
}

// func (db *appdbimpl) SetToken(user_id int, token string) (err error) {
// 	query := `DELETE FROM usertokens WHERE user_id = ?`
// 	_, err = db.c.Exec(query, user_id)
// 	if err != nil {
// 		return
// 	}
// 	query = `INSERT INTO usertokens(user_id, token) VALUES(?, ?);`
// 	_, err = db.c.Exec(query, user_id, token)
// 	if err != nil {
// 		return
// 	}
// 	return err

// }

// func (db *appdbimpl) DeleteToken(token string) (err error) {
// 	query := `DELETE FROM usertokens WHERE token = ?`
// 	_, err = db.c.Exec(query, token)
// 	if err != nil {
// 		return
// 	}
// 	return err
// }

func (db *appdbimpl) GetUserId(id string) (user User, err error) {
	// Query the database for the user by ID (UUID)
	query := `SELECT id, name, photo FROM users WHERE id = ?;`
	row := db.c.QueryRow(query, id)

	// Attempt to scan the result into the User struct
	err = row.Scan(&user.ID, &user.Username, &user.Photo)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a clear error if no user is found
			return User{}, sql.ErrNoRows
		}
		// Return other database-related errors
		return User{}, err
	}

	// Successfully found the user
	return user, nil
}

// -------------
func (db *appdbimpl) UpdateUserName(id string, newname string) (err error) {
    tx, err := db.c.Begin()
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }
    defer tx.Rollback() 

    // Check if the username already exists
    var count int
    checkQuery := `SELECT COUNT(*) FROM users WHERE name = ?;`
    err = tx.QueryRow(checkQuery, newname).Scan(&count)
    if err != nil {
        return fmt.Errorf("failed to check if username exists: %w", err)
    }
    if count > 0 {
        return fmt.Errorf("username '%s' is already taken", newname)
    }

    // Update the username in the users table
    updateUserQuery := `UPDATE users SET name = ? WHERE id = ?;`
    _, err = tx.Exec(updateUserQuery, newname, id)
    if err != nil {
        return fmt.Errorf("failed to update username in users table: %w", err)
    }

    // ✅ Update all conversation names linked to this user
    updateConversationQuery := `UPDATE conversations SET name = ? WHERE id IN (
        SELECT conversation_id FROM convmembers WHERE user_id = ?
    );`
    _, err = tx.Exec(updateConversationQuery, newname, id)
    if err != nil {
        return fmt.Errorf("failed to update conversation names: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (db *appdbimpl) UpdateUserPhoto(userID string, photoPath string) error {
	query := `UPDATE users SET photo = ? WHERE id = ?;`
	_, err := db.c.Exec(query, photoPath, userID)
	return err
}

package database

import "database/sql"


func (db *appdbimpl) CreateUser(username string) (user User, err error) {

	query := `INSERT INTO users ( name ) VALUES (?);`
	_, err = db.c.Exec(query, username)
	if err != nil {
		return
	} 
	return db.GetUser(username)
}

func (db *appdbimpl) GetUser(username string) (user User, err error) {
	query := `SELECT * FROM users WHERE name = ?`
	err = db.c.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Photo) // retirive one row QueryRow, Exec executer 
	if err != nil && err != sql.ErrNoRows{
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		return user, nil
	}
	
}

func (db *appdbimpl) SetToken(user_id int, token string) (err error) {
	query := `DELETE FROM usertokens WHERE user_id = ?`
	_, err = db.c.Exec(query, user_id)
	if err != nil {
		return
	} 
	query = `INSERT INTO usertokens(user_id, token) VALUES(?, ?);`
	_, err = db.c.Exec(query, user_id, token)
	if err != nil {
		return
	} 
	return err

}

func (db *appdbimpl) DeleteToken(token string) (err error) {
	query := `DELETE FROM usertokens WHERE token = ?`
	_, err = db.c.Exec(query, token)
	if err != nil {
		return
	} 
	return err
}

func (db *appdbimpl) GetUserId(token string) (user UserToken, err error) {
	query := `SELECT user_id FROM usertokens WHERE token = ?`
	err = db.c.QueryRow(query, token).Scan(&user.User_id) // retirive one row QueryRow, Exec executer 
	if err != nil && err != sql.ErrNoRows{
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		return user, nil
	}

}
// -------------

func (db *appdbimpl) UpdateUserName(id int, newname string) (err error) {
	query := `UPDATE users SET name = ? WHERE id = ?;`
	_, err = db.c.Exec(query, newname, id)
	if err != nil {
		return
	} 
	query1 := `UPDATE conversations SET name = ? WHERE id IN (SELECT conversation_id FROM convmembers WHERE user_id = ?);`
	_, err = db.c.Exec(query1, newname, id)
	if err != nil {
		return
	} 
	return
}












package repository

import (
    "database/sql"
    "user-service/model"
)

// Thêm user vào DB
func InsertUser(db *sql.DB, user *model.User) (int, error) {
    stmt := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
    result, err := db.Exec(stmt, user.Username, user.Password, user.Email)
    if err != nil {
        return 0, err
    }
    id, _ := result.LastInsertId()
    return int(id), nil
}

// Lấy user theo username
func GetUserByUsername(db *sql.DB, username string) (*model.User, error) {
    var user model.User
    err := db.QueryRow("SELECT id, username, password, email FROM users WHERE username = ?", username).
        Scan(&user.ID, &user.Username, &user.Password, &user.Email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Lấy user theo id
func GetUserByID(db *sql.DB, id int) (*model.User, error) {
    var user model.User
    err := db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id).
        Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
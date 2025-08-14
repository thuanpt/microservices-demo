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

// Lấy user theo ID (bao gồm created_at và updated_at)
func GetUserByID(db *sql.DB, id int) (*model.User, error) {
    var user model.User
    err := db.QueryRow(`
        SELECT id, username, password, email, created_at, updated_at 
        FROM users WHERE id = ?`, id).
        Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Lấy user theo username (cần cho login)
func GetUserByUsername(db *sql.DB, username string) (*model.User, error) {
    var user model.User
    err := db.QueryRow(`
        SELECT id, username, password, email, created_at, updated_at 
        FROM users WHERE username = ?`, username).
        Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Lấy user theo email (có thể cần sau này)
func GetUserByEmail(db *sql.DB, email string) (*model.User, error) {
    var user model.User
    err := db.QueryRow(`
        SELECT id, username, password, email, created_at, updated_at 
        FROM users WHERE email = ?`, email).
        Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

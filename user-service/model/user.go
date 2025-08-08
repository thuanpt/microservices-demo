package model

// Định nghĩa struct User cho DB và API
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password,omitempty"` // Không trả password khi response
    Email    string `json:"email"`
}
package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"

    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Load .env từ thư mục cha
    err := godotenv.Load("../.env")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file .env: %v", err)
    }

    // Kết nối MySQL (không chỉ định database cụ thể để tạo DB)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
    )

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Lỗi khi kết nối DB: %v", err)
    }
    defer db.Close()

    // Kiểm tra tham số command line
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run migrate.go [up|down]")
    }

    command := os.Args[1]
    
    switch command {
    case "up":
        migrateUp(db)
    case "down":
        migrateDown(db)
    default:
        log.Fatal("Command phải là 'up' hoặc 'down'")
    }
}

// Chạy migration up
func migrateUp(db *sql.DB) {
    files, err := filepath.Glob("../migrations/*.up.sql")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file migration: %v", err)
    }
    
    sort.Strings(files) // Sort để đảm bảo thứ tự migration
    
    for _, file := range files {
        fmt.Printf("Chạy migration: %s\n", file)
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("Lỗi khi đọc file %s: %v", file, err)
        }
        
        // Split các câu SQL (có thể có nhiều câu trong 1 file)
        statements := strings.Split(string(content), ";")
        for _, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" {
                continue
            }
            
            _, err := db.Exec(stmt)
            if err != nil {
                log.Printf("Lỗi khi thực hiện SQL: %s\nError: %v", stmt, err)
                // Không fatal để tiếp tục chạy các statement khác
            }
        }
        fmt.Printf("✅ Migration %s hoàn thành\n", filepath.Base(file))
    }
    
    fmt.Println("🎉 Tất cả migration UP đã hoàn thành!")
}

// Chạy migration down (rollback)
func migrateDown(db *sql.DB) {
    files, err := filepath.Glob("../migrations/*.down.sql")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file migration: %v", err)
    }
    
    // Sort ngược để rollback theo thứ tự ngược lại
    sort.Sort(sort.Reverse(sort.StringSlice(files)))
    
    for _, file := range files {
        fmt.Printf("Rollback migration: %s\n", file)
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("Lỗi khi đọc file %s: %v", file, err)
        }
        
        statements := strings.Split(string(content), ";")
        for _, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" {
                continue
            }
            
            _, err := db.Exec(stmt)
            if err != nil {
                log.Printf("Lỗi khi thực hiện SQL: %s\nError: %v", stmt, err)
            }
        }
        fmt.Printf("✅ Rollback %s hoàn thành\n", filepath.Base(file))
    }
    
    fmt.Println("🎉 Tất cả migration DOWN đã hoàn thành!")
}
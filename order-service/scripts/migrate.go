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
    // Load .env từ thư mục cha (order-service/.env)
    err := godotenv.Load("../.env")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file .env: %v", err)
    }

    // Kết nối MySQL
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

    if err := db.Ping(); err != nil {
        log.Fatalf("Không thể ping DB: %v", err)
    }

    // Kiểm tra command
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

func migrateUp(db *sql.DB) {
    files, err := filepath.Glob("../migrations/*.up.sql")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file migration: %v", err)
    }
    
    sort.Strings(files)
    fmt.Println("🚀 Bắt đầu chạy migration UP cho Order Service...")
    
    for _, file := range files {
        fmt.Printf("📄 Chạy migration: %s\n", filepath.Base(file))
        
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("Lỗi khi đọc file %s: %v", file, err)
        }
        
        statements := strings.Split(string(content), ";")
        for _, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" || strings.HasPrefix(stmt, "--") {
                continue
            }
            
            _, err := db.Exec(stmt)
            if err != nil {
                log.Printf("⚠️  Warning: %s\nError: %v", stmt, err)
            }
        }
        fmt.Printf("✅ Migration %s hoàn thành\n", filepath.Base(file))
    }
    
    fmt.Println("🎉 Migration UP cho Order Service hoàn thành!")
}

func migrateDown(db *sql.DB) {
    files, err := filepath.Glob("../migrations/*.down.sql")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file migration: %v", err)
    }
    
    sort.Sort(sort.Reverse(sort.StringSlice(files)))
    fmt.Println("🔄 Bắt đầu rollback migration cho Order Service...")
    
    for _, file := range files {
        fmt.Printf("📄 Rollback: %s\n", filepath.Base(file))
        
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("Lỗi khi đọc file %s: %v", file, err)
        }
        
        statements := strings.Split(string(content), ";")
        for _, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" || strings.HasPrefix(stmt, "--") {
                continue
            }
            
            _, err := db.Exec(stmt)
            if err != nil {
                log.Printf("⚠️  Warning: %s\nError: %v", stmt, err)
            }
        }
        fmt.Printf("✅ Rollback %s hoàn thành\n", filepath.Base(file))
    }
    
    fmt.Println("🎉 Migration DOWN cho Order Service hoàn thành!")
}
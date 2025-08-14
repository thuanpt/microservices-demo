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
    // Load .env từ thư mục cha (product-service/.env)
    err := godotenv.Load("../.env")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file .env: %v", err)
    }

    // Kết nối MySQL root để tạo database và user
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

    // Test kết nối
    if err := db.Ping(); err != nil {
        log.Fatalf("Không thể ping DB: %v", err)
    }

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

// Chạy migration up cho Product Service
func migrateUp(db *sql.DB) {
    // Lấy tất cả file .up.sql trong thư mục migrations
    files, err := filepath.Glob("../migrations/*.up.sql")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file migration: %v", err)
    }
    
    if len(files) == 0 {
        log.Fatal("Không tìm thấy file migration nào trong thư mục ../migrations/")
    }
    
    sort.Strings(files) // Sort để đảm bảo thứ tự migration đúng
    
    fmt.Println("🚀 Bắt đầu chạy migration UP cho Product Service...")
    
    for _, file := range files {
        fmt.Printf("📄 Chạy migration: %s\n", filepath.Base(file))
        
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("Lỗi khi đọc file %s: %v", file, err)
        }
        
        // Split các câu SQL (có thể có nhiều câu trong 1 file)
        statements := strings.Split(string(content), ";")
        for _, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" || strings.HasPrefix(stmt, "--") {
                continue // Bỏ qua câu SQL rỗng hoặc comment
            }
            
            _, err := db.Exec(stmt)
            if err != nil {
                // Log lỗi nhưng vẫn tiếp tục (có thể DB/table đã tồn tại)
                log.Printf("⚠️  Warning khi thực hiện SQL: %s\nError: %v", stmt, err)
            }
        }
        fmt.Printf("✅ Migration %s hoàn thành\n", filepath.Base(file))
    }
    
    fmt.Println("🎉 Tất cả migration UP cho Product Service đã hoàn thành!")
    fmt.Printf("🔗 Product Service có thể kết nối database: %s\n", os.Getenv("DB_NAME"))
}

// Chạy migration down (rollback) cho Product Service
func migrateDown(db *sql.DB) {
    // Lấy tất cả file .down.sql trong thư mục migrations
    files, err := filepath.Glob("../migrations/*.down.sql")
    if err != nil {
        log.Fatalf("Lỗi khi đọc file migration: %v", err)
    }
    
    if len(files) == 0 {
        fmt.Println("⚠️  Không tìm thấy file rollback migration nào")
        return
    }
    
    // Sort ngược để rollback theo thứ tự ngược lại
    sort.Sort(sort.Reverse(sort.StringSlice(files)))
    
    fmt.Println("🔄 Bắt đầu rollback migration cho Product Service...")
    
    for _, file := range files {
        fmt.Printf("📄 Rollback migration: %s\n", filepath.Base(file))
        
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
                log.Printf("⚠️  Warning khi rollback SQL: %s\nError: %v", stmt, err)
            }
        }
        fmt.Printf("✅ Rollback %s hoàn thành\n", filepath.Base(file))
    }
    
    fmt.Println("🎉 Tất cả migration DOWN cho Product Service đã hoàn thành!")
}

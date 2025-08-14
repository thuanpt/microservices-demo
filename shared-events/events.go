package events

import (
    "encoding/json"
    "time"
    "github.com/google/uuid"
)

// Base Event interface
type Event interface {
    GetType() string
    GetID() string
    GetTimestamp() time.Time
    GetData() interface{}
}

// Base Event struct
type BaseEvent struct {
    ID        string      `json:"id"`
    Type      string      `json:"type"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data"`
}

func (e BaseEvent) GetType() string        { return e.Type }
func (e BaseEvent) GetID() string          { return e.ID }
func (e BaseEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e BaseEvent) GetData() interface{}   { return e.Data }

// Helper để tạo BaseEvent mới
func NewBaseEvent(eventType string, data interface{}) BaseEvent {
    return BaseEvent{
        ID:        uuid.New().String(),
        Type:      eventType,
        Timestamp: time.Now().UTC(),
        Data:      data,
    }
}

// User Events
type UserRegisteredEvent struct {
    BaseEvent
    UserData UserEventData `json:"user_data"`
}

func NewUserRegisteredEvent(userID int, username, email string) *UserRegisteredEvent {
    userData := UserEventData{
        UserID:   userID,
        Username: username,
        Email:    email,
    }
    
    return &UserRegisteredEvent{
        BaseEvent: NewBaseEvent(EventTypeUserRegistered, userData),
        UserData:  userData,
    }
}

type UserEventData struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// Product Events
type ProductCreatedEvent struct {
    BaseEvent
    ProductData ProductEventData `json:"product_data"`
}

func NewProductCreatedEvent(productID int, name string, price float64, stock int) *ProductCreatedEvent {
    productData := ProductEventData{
        ProductID: productID,
        Name:      name,
        Price:     price,
        Stock:     stock,
    }
    
    return &ProductCreatedEvent{
        BaseEvent:   NewBaseEvent(EventTypeProductCreated, productData),
        ProductData: productData,
    }
}

type ProductStockUpdatedEvent struct {
    BaseEvent
    ProductData ProductStockData `json:"product_data"`
}

func NewProductStockUpdatedEvent(productID, oldStock, newStock int) *ProductStockUpdatedEvent {
    productData := ProductStockData{
        ProductID: productID,
        OldStock:  oldStock,
        NewStock:  newStock,
    }
    
    return &ProductStockUpdatedEvent{
        BaseEvent:   NewBaseEvent(EventTypeProductStockUpdated, productData),
        ProductData: productData,
    }
}

type ProductEventData struct {
    ProductID int     `json:"product_id"`
    Name      string  `json:"name"`
    Price     float64 `json:"price"`
    Stock     int     `json:"stock"`
}

type ProductStockData struct {
    ProductID int `json:"product_id"`
    OldStock  int `json:"old_stock"`
    NewStock  int `json:"new_stock"`
}

// Order Events
type OrderCreatedEvent struct {
    BaseEvent
    OrderData OrderEventData `json:"order_data"`
}

func NewOrderCreatedEvent(orderID, userID int, totalAmount float64, status string, items []OrderItemData) *OrderCreatedEvent {
    orderData := OrderEventData{
        OrderID:     orderID,
        UserID:      userID,
        TotalAmount: totalAmount,
        Status:      status,
        Items:       items,
    }
    
    return &OrderCreatedEvent{
        BaseEvent: NewBaseEvent(EventTypeOrderCreated, orderData),
        OrderData: orderData,
    }
}

type OrderStatusUpdatedEvent struct {
    BaseEvent
    OrderData OrderStatusData `json:"order_data"`
}

func NewOrderStatusUpdatedEvent(orderID int, oldStatus, newStatus string) *OrderStatusUpdatedEvent {
    orderData := OrderStatusData{
        OrderID:   orderID,
        OldStatus: oldStatus,
        NewStatus: newStatus,
    }
    
    return &OrderStatusUpdatedEvent{
        BaseEvent: NewBaseEvent(EventTypeOrderStatusUpdated, orderData),
        OrderData: orderData,
    }
}

type OrderEventData struct {
    OrderID     int             `json:"order_id"`
    UserID      int             `json:"user_id"`
    TotalAmount float64         `json:"total_amount"`
    Status      string          `json:"status"`
    Items       []OrderItemData `json:"items"`
}

type OrderItemData struct {
    ProductID int     `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

type OrderStatusData struct {
    OrderID   int    `json:"order_id"`
    OldStatus string `json:"old_status"`
    NewStatus string `json:"new_status"`
}

// Event Types constants
const (
    EventTypeUserRegistered      = "user.registered"
    EventTypeProductCreated      = "product.created"
    EventTypeProductStockUpdated = "product.stock_updated"
    EventTypeOrderCreated        = "order.created"
    EventTypeOrderStatusUpdated  = "order.status_updated"
)

// Helper functions
func ToJSON(event Event) ([]byte, error) {
    return json.Marshal(event)
}

func FromJSON(data []byte, event interface{}) error {
    return json.Unmarshal(data, event)
}
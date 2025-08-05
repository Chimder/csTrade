package order

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCanceled   OrderStatus = "canceled"
)

type Order struct {
	ID        string
	BuyerID   string
	SellerID  string
	ListingID string
	Price     float64
	CreatedAt time.Time
	Status    OrderStatus
}

func NewOrder() *Order {
return &Order{

	// BuyerID: ,
	Status: OrderStatusPending,
}
}

package db

import (
	"time"
)

type Order struct {
	OrderUID        string    `gorm:"primaryKey;size:50" json:"order_uid"`
	TrackNumber     string    `gorm:"size:50" json:"track_number"`
	Entry           string    `gorm:"size:20" json:"entry"`
	Locale          string    `gorm:"size:10" json:"locale"`
	InternalSig     string    `gorm:"type:text" json:"internal_signature"`
	CustomerID      string    `gorm:"size:50" json:"customer_id"`
	DeliveryService string    `gorm:"size:50" json:"delivery_service"`
	ShardKey        string    `gorm:"size:10" json:"shardkey"`
	SmID            int       `json:"sm_id"`
	DateCreated     time.Time `json:"date_created"`
	OofShard        string    `gorm:"size:10" json:"oof_shard"`

	// Associations
	Delivery Delivery `gorm:"foreignKey:OrderUID;references:OrderUID;constraint:OnDelete:CASCADE" json:"delivery"`
	Payment  Payment  `gorm:"foreignKey:OrderUID;references:OrderUID;constraint:OnDelete:CASCADE" json:"payment"`
	Items    []Item   `gorm:"foreignKey:OrderUID;references:OrderUID;constraint:OnDelete:CASCADE" json:"items"`
}

type Delivery struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	OrderUID string `gorm:"size:50;index" json:"-"`
	Name     string `gorm:"size:100" json:"name"`
	Phone    string `gorm:"size:20" json:"phone"`
	Zip      string `gorm:"size:20" json:"zip"`
	City     string `gorm:"size:50" json:"city"`
	Address  string `gorm:"size:200" json:"address"`
	Region   string `gorm:"size:50" json:"region"`
	Email    string `gorm:"size:100" json:"email"`
}

type Payment struct {
	ID           uint    `gorm:"primaryKey" json:"-"`
	OrderUID     string  `gorm:"size:50;index" json:"-"`
	Transaction  string  `gorm:"size:100" json:"transaction"`
	RequestID    string  `gorm:"size:100" json:"request_id"`
	Currency     string  `gorm:"size:10" json:"currency"`
	Provider     string  `gorm:"size:50" json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDT    int64   `json:"payment_dt"`
	Bank         string  `gorm:"size:50" json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ID          uint    `gorm:"primaryKey" json:"-"`
	OrderUID    string  `gorm:"size:50;index" json:"-"`
	ChrtID      int64   `json:"chrt_id"`
	TrackNumber string  `gorm:"size:50" json:"track_number"`
	Price       float64 `json:"price"`
	RID         string  `gorm:"size:100" json:"rid"`
	Name        string  `gorm:"size:100" json:"name"`
	Sale        float64 `json:"sale"`
	Size        string  `gorm:"size:10" json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int64   `json:"nm_id"`
	Brand       string  `gorm:"size:100" json:"brand"`
	Status      int     `json:"status"`
}

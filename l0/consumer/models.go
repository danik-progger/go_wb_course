package main

import "time"

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	PaymentDt    int64   `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   float64 `json:"goods_total" validate:"required,gte=0"`
	CustomFee    float64 `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int64   `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        float64 `json:"sale" validate:"gte=0,lte=100"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price" validate:"required,gt=0"`
	NmID        int64   `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      int     `json:"status" validate:"required"`
}

type Order struct {
	OrderUID        string    `json:"order_uid" validate:"required,uuid"`
	TrackNumber     string    `json:"track_number" validate:"required"`
	Entry           string    `json:"entry" validate:"required"`
	Delivery        Delivery  `json:"delivery" validate:"required"`
	Payment         Payment   `json:"payment" validate:"required"`
	Items           []Item    `json:"items" validate:"required,min=1,dive"`
	Locale          string    `json:"locale" validate:"required"`
	InternalSig     string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id" validate:"required"`
	DeliveryService string    `json:"delivery_service" validate:"required"`
	Shardkey        string    `json:"shardkey" validate:"required"`
	SmID            int       `json:"sm_id" validate:"required"`
	DateCreated     time.Time `json:"date_created" validate:"required"`
	OofShard        string    `json:"oof_shard" validate:"required"`
}

type OrderMessage struct {
	Action string `json:"action" validate:"required"`
	Body   Order  `json:"body" validate:"required"`
}

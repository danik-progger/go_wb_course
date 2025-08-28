package db

import (
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	prov *gorm.DB
}

func getDSN() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("🟡  .env file not found, using system environment variables")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, db_name, port,
	)

	return dsn
}

func InitDB() *DB {
	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("🔴 Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&Order{}, &Delivery{}, &Payment{}, &Item{})
	if err != nil {
		log.Fatal("🔴 Failed to migrate database: ", err)
	}

	return &DB{db}
}

func (db *DB) GetAllOrders(limit, offset int) ([]Order, error) {
	var orders []Order

	q := db.prov.Preload("Delivery").Preload("Payment").Preload("Items").Order("date_created desc")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&orders).Error; err != nil {
		fmt.Printf("🔴 DB: Failed to get all orders %s\n", err)
		return nil, err
	}

	fmt.Printf("🟢 DB: got all orders (limit=%d, offset=%d)\n", limit, offset)
	return orders, nil
}

func (db *DB) HasOrders() (bool, error) {
	var count int64
	if err := db.prov.Model(&Order{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) GetOrderById(id string) (*Order, error) {
	var order Order
	if err := db.prov.Preload("Delivery").Preload("Payment").Preload("Items").First(&order, "order_uid = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("🟡 DB: Order not found, uid: %s\n", order.OrderUID)
			return nil, nil
		}
		fmt.Printf("🔴 DB: Failed to get order, uid:%s\n %s\n", order.OrderUID, err)
		return nil, err
	}

	fmt.Printf("🟢 DB: Found order, uid: %s\n", order.OrderUID)
	return &order, nil
}

func (db *DB) AddOrder(order *Order) error {
	if order.OrderUID == "" {
		order.OrderUID = gofakeit.UUID()
	}

	// check existence first
	var existing Order
	if err := db.prov.First(&existing, "order_uid = ?", order.OrderUID).Error; err == nil {
		fmt.Printf("🟡 DB: Order already exists, uid: %s\n", order.OrderUID)
		return nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("🔴 DB: Failed checking existing order: %s\n", err)
		return err
	}

	// ensure associations have the parent's OrderUID
	order.Delivery.OrderUID = order.OrderUID
	order.Payment.OrderUID = order.OrderUID
	for i := range order.Items {
		order.Items[i].OrderUID = order.OrderUID
	}

	if err := db.prov.Session(&gorm.Session{FullSaveAssociations: true}).Create(order).Error; err != nil {
		fmt.Printf("🔴 DB: Failed to add order: %s\n", err)
		return err
	}

	fmt.Printf("🟢 DB: Added order, uid: %s\n", order.OrderUID)
	return nil
}

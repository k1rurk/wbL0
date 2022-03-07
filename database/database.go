package database

import (
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

const dsn = "user=wb_user password=12345 dbname=wb_db port=5432 sslmode=disable"

type Item struct {
	Id          uint    `json:"-" gorm:"primaryKey"`
	ChrtId      uint    `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        uint    `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

type Payment struct {
	Id           uint    `json:"-" gorm:"primaryKey"`
	Transaction  string  `json:"transaction"`
	RequestId    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       int     `json:"amount"`
	PaymentDt    int     `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   int     `json:"goods_total"`
	CustomFee    int     `json:"custom_fee"`
}

type Delivery struct {
	Id      uint   `json:"-" gorm:"primaryKey"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Order struct {
	Id                uint      `json:"-" gorm:"primaryKey"`
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" gorm:"foreignKey:Id;references:Id"`
	Payment           Payment   `json:"payment" gorm:"foreignKey:Transaction;references:OrderUid"`
	Item              []Item    `json:"items" gorm:"foreignKey:TrackNumber;references:TrackNumber"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (Order) TableName() string {
	return "Order"
}

func (Item) TableName() string {
	return "Item"
}

func (Delivery) TableName() string {
	return "Delivery"
}

func (Payment) TableName() string {
	return "Payment"
}

func Open() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func Create(order *Order, db *gorm.DB) error {
	//err := db.AutoMigrate(Order{})
	//if err != nil {
	//	log.Println(err)
	//}
	res := db.Create(order)
	if res.Error != nil {
		log.Println(res.Error)
	}
	return nil
}

func LoadDataFromDb(db *gorm.DB) []Order {
	var orders []Order
	result := db.Preload("Delivery").Preload("Payment").Preload("Item").Find(&orders)
	log.Printf("Data from database received, count: %v\n", result.RowsAffected)
	return orders
}

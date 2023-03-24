package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var firstNames = []string{
	"Alice", "Bob", "Charlie", "David", "Ivan", "Isaac", "Emma", "Frank",
	"Grace", "Henry", "Isabella", "Jack", "Kate", "Liam", "Peter", "Justus",
	"Mia", "Nora", "Oliver", "Penny", "Quinn", "Riley",
	"Sarah", "Thomas", "Ursula", "Victor", "Wendy", "Xander",
	"Yara", "Zachary"}

var lastNames = []string{
	"Adams", "Baker", "Clark", "Davis", "Saringer", "Evans", "Franklin",
	"Gonzalez", "Harris", "Irwin", "Johnson", "Klein", "Lee", "Gieskanne",
	"Miller", "Nelson", "O'Connor", "Patel", "Quinn", "Roberts", "Kanunnikov",
	"Smith", "Thompson", "Upton", "Valdez", "Walker", "Xu",
	"Young", "Zhang",
}

// structs for JSON unmarshalling
type Ordering struct {
	OrderUid          string    `json:"order_uid" validate:"max=40"`
	TrackNumber       string    `json:"track_number" validate:"max=40"`
	Entry             string    `json:"entry" validate:"max=40"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale" validate:"max=40"`
	InternalSignature string    `json:"internal_signature" validate:"max=40"`
	CustomerId        string    `json:"customer_id" validate:"max=40"`
	DeliveryService   string    `json:"delivery_service" validate:"max=40"`
	ShardKey          string    `json:"shardkey" validate:"max=40"`
	SmId              int       `json:"sm_id"`
	DateOfCreation    time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard" validate:"max=40"`
}

// alias fuction to convert the string from JSON to time type
func (o *Ordering) UnmarshalJSON(b []byte) error {
	// auxiliary Alias struct where DateOfCreation is unmarshalled into a string
	type Alias Ordering
	aux := &struct {
		DateOfCreation string `json:"date_created"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	// parse the DateString to time type
	if aux.DateOfCreation != "" {
		parsedTime, err := time.Parse(time.RFC3339, aux.DateOfCreation)
		if err != nil {
			return err
		}
		o.DateOfCreation = parsedTime
	}

	return nil
}

type Delivery struct {
	Name    string `json:"name" validate:"max=40"`
	Phone   string `json:"phone" validate:"max=40"`
	Zip     string `json:"zip" validate:"max=40"`
	City    string `json:"city" validate:"max=40"`
	Address string `json:"address" validate:"max=40"`
	Region  string `json:"region" validate:"max=40"`
	Email   string `json:"email" validate:"max=40"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number" validate:"max=40"`
	Price       int    `json:"price"`
	Rid         string `json:"rid" validate:"max=40"`
	Name        string `json:"name" validate:"max=40"`
	Sale        int    `json:"sale"`
	Size        string `json:"size" validate:"max=40"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand" validate:"max=40"`
	Status      int    `json:"status"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"max=40"`
	RequestId    string `json:"request_id" validate:"max=40"`
	Currency     string `json:"currency" validate:"max=40"`
	Provider     string `json:"provider" validate:"max=40"`
	Amount       int    `json:"amount"`
	PaymentDate  int    `json:"payment_dt"`
	Bank         string `json:"bank" validate:"max=40"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

func generateRandomItem(amount int) []Item {
	rand.Seed(time.Now().UnixNano())
	// generate i items with some random field-values
	var items []Item
	for i := 0; i < amount; i++ {
		var item Item
		item.ChrtId = rand.Intn(99999999)
		item.TrackNumber = "WBILMTESTTRACK"
		item.Price = rand.Intn(99999)
		item.Rid = "ab4219087a764ae0btest"
		item.Name = "Mascaras"
		item.Sale = rand.Intn(99)
		item.Size = "0"
		item.TotalPrice = int(item.Price * ((100 - item.Sale) / 100))
		item.NmId = rand.Intn(9999999)
		item.Brand = "Vivienne Sabo"
		item.Status = 202

		items = append(items, item)
	}
	return items
}

func generateRandomDelivery() Delivery {
	rand.Seed(time.Now().UnixNano())

	var delivery Delivery
	// compose random name
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	delivery.Name = firstName + " " + lastName
	// generate random number with 10 digits starting with "+"
	time.Sleep(time.Duration(100) * time.Millisecond)
	delivery.Phone = "+" + strconv.Itoa(rand.Intn(8999999999)+1000000000)
	delivery.Zip = "2639809"
	delivery.City = "Kiryat Mozkin"
	delivery.Address = "Ploshad Mira 15"
	delivery.Region = "Kraiot"
	delivery.Email = strings.ToLower(firstName + lastName + "@gmail.com")

	return delivery
}

func generateRandomPayment(orderingUid string) Payment {
	var payment Payment

	payment.Transaction = orderingUid
	payment.RequestId = uuid.New().String()
	payment.Currency = "USD"
	payment.Provider = "wbpay"
	payment.Amount = rand.Intn(99999)
	payment.PaymentDate = int(time.Now().Unix())
	payment.Bank = "VTB"
	payment.DeliveryCost = 1500
	payment.GoodsTotal = payment.Amount - payment.DeliveryCost
	payment.CustomFee = rand.Intn(23)

	return payment
}

func GenerateRandomOrdering() Ordering {
	var ordering Ordering

	ordering.OrderUid = uuid.New().String()
	ordering.TrackNumber = "WBILMTESTTRACK"
	ordering.Entry = "WBIL"
	ordering.Locale = "en"
	ordering.InternalSignature = "signed"
	ordering.CustomerId = "test"
	ordering.DeliveryService = "meest"
	ordering.ShardKey = "9"
	ordering.SmId = rand.Intn(99)
	ordering.DateOfCreation, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	ordering.OofShard = "1"

	ordering.Delivery = generateRandomDelivery()
	ordering.Payment = generateRandomPayment(ordering.OrderUid)
	// at least 1 and at most 11 items
	ordering.Items = generateRandomItem(rand.Intn(10) + 1)

	return ordering
}

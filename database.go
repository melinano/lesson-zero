package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/melinano/lesson-zero/models"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "melinano"
	password = "1234"
	dbname   = "ordering_db"
)

func startDB(ordering models.Ordering) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	insertOrderingIntoDB(db, ordering)
}

func insertOrderingIntoDB(db *sql.DB, ordering models.Ordering) {

	insertStatementOrdering := `INSERT INTO ordering (order_uid, track_number, 
                      entry, locale, internal_signature, customer_id, 
                      delivery_service, shardkey, sm_id, date_created,
                      oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	err := db.QueryRow(
		insertStatementOrdering,
		ordering.OrderUid,
		ordering.TrackNumber,
		ordering.Entry,
		ordering.Locale,
		ordering.InternalSignature,
		ordering.CustomerId,
		ordering.DeliveryService,
		ordering.ShardKey,
		ordering.SmId,
		ordering.DateOfCreation,
		ordering.OofShard).Err()
	if err != nil {
		log.Fatalf("Error while inserting ordering: %s", err)
		return
	}

	insertDeliveryIntoDB(db, ordering.OrderUid, ordering.Delivery)
	insertPaymentIntoDB(db, ordering.OrderUid, ordering.Payment)
	insertItemsIntoDB(db, ordering.OrderUid, ordering.Items)
}
func insertDeliveryIntoDB(db *sql.DB, orderingUid string, delivery models.Delivery) {

	insertStatementDelivery := `INSERT INTO delivery (name, phone, zip, city, address, region,
                      email, ordering_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	err := db.QueryRow(
		insertStatementDelivery,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
		orderingUid).Err()
	if err != nil {
		log.Fatalf("Error while inserting delivery: %s", err)
	}
}

func insertPaymentIntoDB(db *sql.DB, orderingUid string, payment models.Payment) {

	insertStatementDelivery := `INSERT INTO delivery (request_id, currency, provider, amount, payment_date, bank,
                      delivery_cost, goods_total, custom_fee, transaction) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	err := db.QueryRow(
		insertStatementDelivery,
		payment.RequestId,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDate,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
		orderingUid).Err()
	if err != nil {
		log.Fatalf("Error while inserting payment: %s", err)
	}
}

func insertItemsIntoDB(db *sql.DB, orderingUid string, items []models.Item) {
	for _, item := range items {
		insertStatementItem := `INSERT INTO item (chrt_id, track_number, 
                      price, rid, name, sale, size, total_price, nm_id, 
                      brand, status, ordering_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

		err := db.QueryRow(
			insertStatementItem,
			item.ChrtId,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status,
			orderingUid).Err()
		if err != nil {
			log.Fatalf("Error while inserting item: %s", err)
		}
	}
}

func updateOrdering(db *sql.DB, ordering models.Ordering) {
	updateStatementOrdering := `UPDATE ordering SET 
                    track_number = $2, 
                    entry = $3,
                    locale = $4,
                    internal_signature = $5,
                    customer_id = $6, 
                    delivery_service = $7,
                    shardkey = $8,
                    sm_id = $9,
                    date_created = $10,
                    oof_shard = $11
                    WHERE order_uid = $1`

	result, err := db.Exec(
		updateStatementOrdering,
		ordering.OrderUid,
		ordering.TrackNumber,
		ordering.Entry,
		ordering.Locale,
		ordering.InternalSignature,
		ordering.CustomerId,
		ordering.DeliveryService,
		ordering.ShardKey,
		ordering.SmId,
		ordering.DateOfCreation,
		ordering.OofShard)

	id, err := result.RowsAffected()

	if err != nil {
		log.Fatalf("Error while updating ordering: %s\n %s", id, err)
	}
}

func updateDelivery(db *sql.DB, orderingUid string, delivery models.Delivery) {
	updateStatementDelivery := `UPDATE delivery SET 
                    name = $1, 
                    zip = $3,
                    city = $4,
                    address = $5, 
                    region = $6,
                    email = $7
                    WHERE ordering_id = $8 AND phone = $2`

	result, err := db.Exec(
		updateStatementDelivery,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
		orderingUid)

	id, err := result.RowsAffected()

	if err != nil {
		log.Fatal("Error while updating delivery: %s\n %s", id, err)
	}
}

func updatePayment(db *sql.DB, orderingUid string, payment models.Payment) {
	updateStatementPayment := `UPDATE payment SET 
                    currency = $2,
                    provider = $3,
                    amount = $4,
                    payment_date = $5, 
                    bank = $6,
                    delivery_cost = $7,
                    goods_total = $8,
                    custom_fee = $9
                    WHERE transaction = $10 AND request_id = $1`

	result, err := db.Exec(
		updateStatementPayment,
		payment.RequestId,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDate,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
		orderingUid)

	id, err := result.RowsAffected()

	if err != nil {
		log.Fatal("Error while updating ordering: %s\n %s", id, err)
	}
}

func updateItem(db *sql.DB, orderingUid string, item models.Item) {
	updateStatementItem := `UPDATE item SET 
                    track_number = $2, 
                    price = $3,
                    rid = $4,
                    name = $5,
                    sale = $6, 
                    size = $7,
                    total_price = $8,
                    nm_id = $9,
                    brand = $10,
                    status = $11
                    WHERE ordering_id = $12 AND chrt_id = $1`

	result, err := db.Exec(
		updateStatementItem,
		item.ChrtId,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmId,
		item.Brand,
		item.Status,
		orderingUid)

	id, err := result.RowsAffected()

	if err != nil {
		log.Fatal("Error while updating item: %s\n %s", id, err)
	}
}

func fetchOrderings(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM ordering")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var ordering models.Ordering
		if err := rows.Scan(&ordering.OrderUid, &ordering.TrackNumber, &ordering.Entry,
			&ordering.Locale, &ordering.InternalSignature, &ordering.CustomerId, &ordering.DeliveryService,
			&ordering.ShardKey, &ordering.SmId, &ordering.DateOfCreation, &ordering.OofShard); err != nil {
			log.Fatal(err)
		}
	}
}

func fetchDelivery(db *sql.DB, orderingUid string) models.Delivery {
	rows, err := db.Query(`SELECT * FROM delivery WHERE ordering_id = $1 LIMIT 1`, orderingUid)
	if err != nil {
		log.Fatal(err)
	}
	var delivery models.Delivery
	for rows.Next() {
		if err := rows.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email); err != nil {
			log.Fatal(err)
		}
	}
	return delivery
}

func fetchPayment(db *sql.DB, orderingUid string) models.Payment {
	rows, err := db.Query(`SELECT * FROM payment WHERE transaction = $1 LIMIT 1`, orderingUid)
	if err != nil {
		log.Fatal(err)
	}
	var payment models.Payment
	for rows.Next() {
		if err := rows.Scan(&payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount,
			&payment.PaymentDate, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee, &payment.Transaction); err != nil {
			log.Fatal(err)
		}
	}
	return payment
}

func fetchItems(db *sql.DB, orderingUid string) []models.Item {
	rows, err := db.Query(`SELECT * FROM item WHERE ordering_id = $1`, orderingUid)
	if err != nil {
		log.Fatal(err)
	}
	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand,
			&item.Status); err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}

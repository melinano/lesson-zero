package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
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

func startDB() *pgxpool.Pool {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)

	pool, err := pgxpool.Connect(context.Background(), psqlconn)
	if err != nil {
		log.Print(err)
	}
	return pool
}

func insertOrderingIntoDB(pool *pgxpool.Pool, ordering models.Ordering) error {

	insertStatementOrdering := `INSERT INTO ordering (order_uid, track_number, 
                      entry, locale, internal_signature, customer_id, 
                      delivery_service, shardkey, sm_id, date_created,
                      oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := pool.Exec(
		context.Background(),
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
		ordering.OofShard)
	if err != nil {
		log.Printf("Error while inserting ordering: %s", err)
		return err
	}

	if insertDeliveryIntoDB(pool, ordering.OrderUid, ordering.Delivery) != nil ||
		insertPaymentIntoDB(pool, ordering.OrderUid, ordering.Payment) != nil ||
		insertItemsIntoDB(pool, ordering.OrderUid, ordering.Items) != nil {
		return err
	}
	return nil
}
func insertDeliveryIntoDB(pool *pgxpool.Pool, orderingUid string, delivery models.Delivery) error {

	insertStatementDelivery := `INSERT INTO delivery (name, phone, zip, city, address, region,
                      email, ordering_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := pool.Exec(
		context.Background(),
		insertStatementDelivery,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
		orderingUid)
	if err != nil {
		log.Printf("Error while inserting delivery (phone: %s): %s \n Skipping record ordering(order_uid: %s)", delivery.Phone, err, orderingUid)
		return err
	}
	return nil
}

func insertPaymentIntoDB(pool *pgxpool.Pool, orderingUid string, payment models.Payment) error {

	insertStatementPayment := `INSERT INTO payment (request_id, currency, provider, amount, payment_date, bank,
                      delivery_cost, goods_total, custom_fee, transaction) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := pool.Exec(
		context.Background(),
		insertStatementPayment,
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
	if err != nil {
		log.Printf("Error while inserting payment (request_id: %s): %s\nSkipping record ordering(order_uid: %s)", payment.RequestId, err, orderingUid)
		return err
	}
	return nil
}

func insertItemsIntoDB(pool *pgxpool.Pool, orderingUid string, items []models.Item) error {
	for _, item := range items {
		insertStatementItem := `INSERT INTO item (chrt_id, track_number, 
                      price, rid, name, sale, size, total_price, nm_id, 
                      brand, status, ordering_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

		_, err := pool.Exec(
			context.Background(),
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
			orderingUid)
		if err != nil {
			log.Printf("Error while inserting item(chrt_id: %s): %s\nSkipping record ordering(%s)", item.ChrtId, err, orderingUid)
			return err
		}
	}
	return nil
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
		log.Printf("Error while updating ordering: %s\n %s", id, err)
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
		log.Printf("Error while updating delivery: %s\n %s", id, err)
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
		log.Printf("Error while updating ordering: %s\n %s", id, err)
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
		log.Printf("Error while updating item: %s\n %s", id, err)
	}
}

func fetchOrderings(pool *pgxpool.Pool, orderingsMap *map[string]models.Ordering) error {
	// Sending SQL Query to database
	rows, err := pool.Query(context.Background(), "SELECT * FROM ordering")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	// going through the rows of the Ordering table to get field values
	for rows.Next() {
		var ordering models.Ordering
		if err := rows.Scan(&ordering.OrderUid, &ordering.TrackNumber, &ordering.Entry,
			&ordering.Locale, &ordering.InternalSignature, &ordering.CustomerId, &ordering.DeliveryService,
			&ordering.ShardKey, &ordering.SmId, &ordering.DateOfCreation, &ordering.OofShard); err != nil {
			log.Print(err)
			return err
		} else {
			// fetching the corresponding values for Delivery, Payment and Items
			if err := fetchDelivery(pool, ordering.OrderUid, &ordering.Delivery); err != nil {
				return err
			} else if err = fetchPayment(pool, ordering.OrderUid, &ordering.Payment); err != nil {
				return err
			} else if err = fetchItems(pool, ordering.OrderUid, &ordering.Items); err != nil {
				return err
			}
			// adding to the map
			(*orderingsMap)[ordering.OrderUid] = ordering
		}
	}
	return nil
}

func fetchDelivery(pool *pgxpool.Pool, orderingUid string, delivery *models.Delivery) error {
	// Sending SQL Query to database
	// There is only one Delivery in an Ordering
	rows, err := pool.Query(context.Background(), `SELECT * FROM delivery WHERE ordering_id = $1 LIMIT 1`, orderingUid)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	// going through the rows of the Delivery table to get field values
	for rows.Next() {
		// we don't need orderingId from the DB, but we have to catch the value
		var orderingId *string
		if err := rows.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email, &orderingId); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func fetchPayment(pool *pgxpool.Pool, orderingUid string, payment *models.Payment) error {
	// Sending SQL Query to database
	// there is only one Payment in an Ordering
	rows, err := pool.Query(context.Background(), `SELECT * FROM payment WHERE transaction = $1 LIMIT 1`, orderingUid)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	// going through the rows of the Payment table to get field values
	for rows.Next() {
		if err := rows.Scan(&payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount,
			&payment.PaymentDate, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee, &payment.Transaction); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func fetchItems(pool *pgxpool.Pool, orderingUid string, items *[]models.Item) error {
	// Sending SQL Query to database
	rows, err := pool.Query(context.Background(), `SELECT * FROM item WHERE ordering_id = $1`, orderingUid)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	// going through the rows of the Item table to get field values and appending them
	// into a slice of items
	for rows.Next() {
		var item models.Item
		// we don't need orderingId from the DB, but we have to catch the value
		var ordering_id *string
		if err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand,
			&item.Status, &ordering_id); err != nil {
			log.Print(err)
			return err
		}
		*items = append(*items, item)
	}
	return nil
}

package database

import (
	"database/sql"

	"github.com/luis13005/cleanarchitecture/internal/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (o *OrderRepository) Save(order *entity.Order) error {
	stmt, err := o.DB.Prepare("INSERT INTO orders (id, price, tax, final_price) values ($1,$2,$3,$4)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}

	return nil
}

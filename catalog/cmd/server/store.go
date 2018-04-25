package main

import (
	"database/sql"
	"errors"

	"github.com/darethas/go-microservices/catalog"
)

type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// newStore will prepare our queries on the db instance passed in and return it for use
func newStore(db *sql.DB) (*store, error) {
	// create a map of statements to prepare
	unpreparedStatements := map[string]string{
		"get-products": `
			SELECT id, 
				description, 
				price, 
				quantity 
			FROM catalog.products;
			`,
		"get-product-by-id": `
			SELECT 
				id, 
				description, 
				price,
				quantity 
			FROM catalog.products WHERE id = ?;
		`,
		"decrement-product-quantity-by-id": `
			UPDATE catalog.products 
			SET quantity = quantity - 1 
			WHERE id = ?;
		`,
	}

	prepared := map[string]*sql.Stmt{}
	for k, v := range unpreparedStatements {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, err
		}
		prepared[k] = stmt
	}

	s := &store{
		db:    db,
		stmts: prepared,
	}

	return s, nil
}

// GetProducts will retrieve the entire product catalog
func (s *store) GetProducts() ([]catalog.Product, error) {
	// perform the query
	rows, err := s.stmts["get-products"].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// initialize our slice of products
	products := []catalog.Product{}

	for rows.Next() {
		p := catalog.Product{}
		err := rows.Scan(&p.ID, &p.Description, &p.Price, &p.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID will fetch one product from the database by ID
func (s *store) GetProductByID(id int) (catalog.Product, error) {
	p := catalog.Product{}

	err := s.stmts["get-product-by-id"].QueryRow(id).Scan(&p.ID, &p.Description, &p.Price, &p.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return catalog.Product{}, errors.New("product not found")
		}
		return catalog.Product{}, err
	}

	return p, nil
}

// DecrementProductQuantityByID will decrement the quantity of a product by one, returning an error if the quantity is already zero. You may want to support a batch version of this function (i.e. can decrease quantity by arbitrary amount), but this will suffice for now
func (s *store) DecrementProductQuantityByID(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	p := catalog.Product{}
	// take our prepared statment for getting a product, and re-prepare it under the transaction, and perform the query
	err = tx.Stmt(s.stmts["get-product-by-id"]).QueryRow(id).Scan(&p.ID, &p.Description, &p.Price, &p.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	// if there is no quantity to decrement, return an error
	if p.Quantity == 0 {
		tx.Rollback()
		return errors.New("cannot decrement: no inventory left")
	}

	_, err = tx.Stmt(s.stmts["decrement-product-quantity-by-id"]).Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

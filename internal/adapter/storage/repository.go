package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewReposiory(databaseType, connStr string) *Repository {

	db, err := sql.Open(databaseType, connStr)
	if err != nil {
		return nil
	}

	return &Repository{
		db: db,
	}
}

func (r *Repository) WriteInventoryData(ctx context.Context, data *[]models.InventoryData) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback to ensure transaction is closed in case of an error or panic.
	// If tx.Commit() succeeds later, Rollback() will safely return sql.ErrTxDone.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic during database write: %v", r)
			_ = tx.Rollback()
			panic(r) // Re-throw panic
		}
	}()
	defer func() {
		if err != nil {
			// Only try to rollback if an error occurred before commit
			if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
				log.Printf("Error during transaction rollback: %v", rollbackErr)
				// Prefer returning the original error, but log the rollback failure
			}
		}
	}()

	const insertSQL = `
        INSERT INTO public.inventory (item, quantity, unit, wholesale_price_per_quantity, total_cost_of_product)
        VALUES ($1, $2, $3, $4, $5);
    `

	// 2. Iterate through the data slice and execute the prepared statement within the transaction
	for i, item := range *data {
		_, err = tx.ExecContext(
			ctx,
			insertSQL,
			item.Item,
			item.Quantity,
			item.Unit,
			item.WholesalePricePerQuantity,
			item.TotalCostOfProduct,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item %d (%s): %w", i, item.Item, err)
		}
	}

	// 3. Commit the transaction if all insertions were successful
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) WriteSalesData(ctx context.Context, data *[]models.SalesData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback to ensure transaction is closed in case of an error or panic.
	// If tx.Commit() succeeds later, Rollback() will safely return sql.ErrTxDone.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic during database write: %v", r)
			_ = tx.Rollback()
			panic(r) // Re-throw panic
		}
	}()
	defer func() {
		if err != nil {
			// Only try to rollback if an error occurred before commit
			if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
				log.Printf("Error during transaction rollback: %v", rollbackErr)
				// Prefer returning the original error, but log the rollback failure
			}
		}
	}()

	const insertSQL = `
        INSERT INTO inventory (item, quantity, unit, retail_price_per_quantity, total_selling_price)
        VALUES ($1, $2, $3, $4, $5);
    `

	// 2. Iterate through the data slice and execute the prepared statement within the transaction
	for i, item := range *data {
		_, err = tx.ExecContext(
			ctx,
			insertSQL,
			item.Item,
			item.Quantity,
			item.Unit,
			item.RetailPricePerQuantity,
			item.TotalSellingPrice,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item %d (%s): %w", i, item.Item, err)
		}
	}

	// 3. Commit the transaction if all insertions were successful
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

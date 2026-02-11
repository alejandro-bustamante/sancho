package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// Usando un nombre alternativo a "Store"
type Database interface {
	Querier
	WithTx(tx *sql.Tx) Database
	ExecTx(ctx context.Context, fn func(Database) error) error
}

type sqlProvider struct {
	*Queries
	db *sql.DB
}

func NewDatabase(db *sql.DB) Database {
	return &sqlProvider{
		Queries: New(db),
		db:      db,
	}
}

func (p *sqlProvider) WithTx(tx *sql.Tx) Database {
	return &sqlProvider{
		Queries: p.Queries.WithTx(tx),
		db:      p.db,
	}
}

// ExecTx es una utilidad MUY recomendada.
// Ejecuta una función dentro de una transacción, manejando el commit y rollback automáticamente.
func (p *sqlProvider) ExecTx(ctx context.Context, fn func(Database) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Creamos un Database que usa ESTA transacción
	q := p.WithTx(tx)

	// Ejecutamos la lógica de negocio pasando el Store transaccional
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

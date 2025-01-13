// Package pg provides a PostgreSQL-based implementation for managing order items,
// including operations such as retrieving, listing, creating, updating, and deleting records.
package pg

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mb-feedback/internal/domain/order_detail/model"
	"mb-feedback/internal/errs"
)

// Repo provides methods to interact with the PostgreSQL database for order items operations.
// It holds a connection pool to manage database connections.
type Repo struct {
	Con *pgxpool.Pool
}

// New creates a new instance of Repo with the given PostgreSQL connection pool.
func New(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

// Get retrieves a single order item based on the provided query parameters.
// It returns the item if found, a boolean indicating its existence, and any error encountered.
func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.OrderItem, bool, error) {
	if !pars.IsValid() {
		return nil, false, errs.InvalidInput
	}

	var result model.OrderItem

	queryBuilder := squirrel.Select("*").From("ord_detail")

	if len(pars.ID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if len(pars.OrderID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_id": pars.OrderID})
	}

	if len(pars.ProductCode) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"product_code": pars.ProductCode})
	}

	queryBuilder = queryBuilder.Limit(1)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, err
	}

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.ID, &result.OrderID, &result.ProductCode, &result.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, true, nil
}

// List retrieves multiple order items based on the provided query parameters,
// supporting filters like ID, user OrderID, ProductCode, and timestamps. It returns the list
// of items, the total count, and any error encountered.
func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.OrderItem, int64, error) {
	queryBuilder := squirrel.
		Select("id", "order_id", "product_code", "created_at").
		From("ord_detail")

	if pars.ID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.IDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.IDs})
	}

	if pars.OrderID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_id": pars.OrderID})
	}

	if pars.OrderIDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_id": pars.OrderIDs})
	}

	if pars.ProductCode != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"product_code": pars.ProductCode})
	}

	if pars.ProductCodes != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"product_code": pars.ProductCodes})
	}

	if pars.CreatedBefore != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"created_at": pars.CreatedBefore})
	}

	if pars.CreatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"created_at": pars.CreatedAfter})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.Con.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var result []*model.OrderItem
	for rows.Next() {
		var data model.OrderItem
		err = rows.Scan(&data.ID, &data.OrderID, &data.ProductCode, &data.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, &data)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, int64(len(result)), nil
}

// Create inserts a new order item into the database based on the provided Edit object,
// returning any error encountered.
func (r *Repo) Create(ctx context.Context, obj *model.Edit) error {
	insert := squirrel.Insert("ord_detail").
		Columns("order_id", "product_code").
		Values(obj.OrderID, obj.ProductCode).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := insert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies an existing order item based on the provided query parameters and Edit object,
// returning any error encountered during the operation.
func (r *Repo) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Update("ord_detail")

	if obj.ProductCode != nil {
		queryBuilder = queryBuilder.Set("product_code", obj.ProductCode)
	} else {
		return nil
	}

	if obj.ID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": obj.ID})
	}

	if obj.OrderID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_id": obj.OrderID})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

// Delete removes a order item from the database based on the provided query parameters,
// returning any error encountered during the operation.
func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Delete("ord_detail")

	if pars.ID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.OrderID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_id": pars.OrderID})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

//func (r *Repo) BeginTx(ctx context.Context) (pgx.Tx, error) {
//	tx, err := r.Con.BeginTx(ctx, pgx.TxOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	return tx, nil
//}
//
//func (r *Repo) CommitTx(ctx context.Context, tx pgx.Tx) error {
//	if err := tx.Commit(ctx); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (r *Repo) RollbackTx(ctx context.Context, tx pgx.Tx) error {
//	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
//		return err
//	}
//	return nil
//}
//
//func (r *Repo) HandleTxCompletion(tx pgx.Tx, err *error) {
//	if p := recover(); p != nil {
//		_ = tx.Rollback(context.Background())
//		panic(p)
//	} else if *err != nil {
//		_ = tx.Rollback(context.Background())
//	} else {
//		*err = tx.Commit(context.Background())
//	}
//}

// Package pg provides a PostgreSQL-based implementation for managing order,
// including operations such as retrieving, listing, creating, updating, and deleting records.
package pg

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mb-feedback/internal/domain/order/model"
	"mb-feedback/internal/errs"
)

// Repo provides methods to interact with the PostgreSQL database for order operations.
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

// Get retrieves a single data item based on the provided query parameters.
// It returns the item if found, a boolean indicating its existence, and any error encountered.
func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error) {
	if !pars.IsValid() {
		return nil, false, errs.InvalidInput
	}

	var result model.Order

	queryBuilder := squirrel.Select("*").From("ord")

	if len(pars.ID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if len(pars.ExternalOrderID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"external_order_id": pars.ExternalOrderID})
	}

	if len(pars.UserPhone) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_phone": pars.UserPhone})
	}

	queryBuilder = queryBuilder.Limit(1)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, err
	}

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.ID, &result.ExternalOrderID, &result.UserPhone, &result.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, true, nil
}

// List retrieves multiple order based on the provided query parameters,
// supporting filters like ID, ExternalOrderID, UserPhone, and timestamps. It returns the list
// of items, the total count, and any error encountered.
func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.Order, int64, error) {
	queryBuilder := squirrel.
		Select("id", "external_order_id", "user_phone", "created_at").
		From("ord")

	if pars.ID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.IDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.IDs})
	}

	if pars.ExternalOrderID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"external_order_id": pars.ExternalOrderID})
	}

	if pars.ExternalOrderIDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"external_order_id": pars.ExternalOrderIDs})
	}

	if pars.UserPhone != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_phone": pars.UserPhone})
	}

	if pars.UserPhones != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_phone": pars.UserPhones})
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

	var result []*model.Order
	for rows.Next() {
		var data model.Order
		err = rows.Scan(&data.ID, &data.ExternalOrderID, &data.UserPhone, &data.CreatedAt)
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

// Create inserts a new order into the database based on the provided Edit object,
// returning any error encountered.
func (r *Repo) Create(ctx context.Context, obj *model.Edit) error {
	insert := squirrel.Insert("ord").
		Columns("external_order_id", "user_phone").
		Values(obj.ExternalOrderID, obj.UserPhone).
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

// Update modifies an existing data item based on the provided query parameters and Edit object,
// returning any error encountered during the operation.
func (r *Repo) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Update("ord")

	if obj.UserPhone != nil {
		queryBuilder = queryBuilder.Set("user_phone", obj.UserPhone)
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

// Delete removes a order from the database based on the provided query parameters,
// returning any error encountered during the operation.
func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Delete("ord")

	if pars.ID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.ExternalOrderID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"external_order_id": pars.ExternalOrderID})
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

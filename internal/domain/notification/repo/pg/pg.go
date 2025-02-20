// Package pg provides a PostgreSQL-based implementation for managing notifications,
// including operations such as retrieving, listing, creating, updating, and deleting records.
package pg

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mb-feedback/internal/domain/notification/model"
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
func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.Notification, bool, error) {
	if !pars.IsValid() {
		return nil, false, errs.InvalidInput
	}

	var result model.Notification

	queryBuilder := squirrel.Select("*").From("notification")

	if len(pars.ID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if len(pars.OrderItemID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_item_id": pars.OrderItemID})
	}

	if len(pars.PhoneNumber) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"phone_number": pars.PhoneNumber})
	}

	if len(pars.Status) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Status})
	}

	queryBuilder = queryBuilder.Limit(1)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, err
	}

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.ID, &result.OrderItemID, &result.PhoneNumber, &result.Status, &result.CreatedAt)
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
func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.Notification, int64, error) {
	queryBuilder := squirrel.
		Select("id", "order_item_id", "phone_number", "status", "sent_at", "created_at").
		From("notification")

	if pars.ID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.IDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.IDs})
	}

	if pars.OrderItemID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_item_id": pars.OrderItemID})
	}

	if pars.OrderItemIDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_item_id": pars.OrderItemID})
	}

	if pars.PhoneNumber != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"phone_number": pars.PhoneNumber})
	}

	if pars.PhoneNumbers != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"phone_number": pars.PhoneNumbers})
	}

	if pars.Status != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Status})
	}

	if pars.Statuses != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Status})
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

	var result []*model.Notification
	for rows.Next() {
		var data model.Notification
		err = rows.Scan(&data.ID, &data.OrderItemID, &data.PhoneNumber, &data.Status, &data.CreatedAt)
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
	insert := squirrel.Insert("notification").
		Columns("order_item_id", "phone_number", "status", "sent_at").
		Values(obj.OrderItemID, obj.PhoneNumber, obj.Status, obj.SentAt).
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

	queryBuilder := squirrel.Update("notification")

	if obj.Status != nil {
		queryBuilder = queryBuilder.Set("status", obj.Status)
	}

	if obj.SentAt != nil {
		queryBuilder = queryBuilder.Set("sent_at", obj.SentAt)
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

	queryBuilder := squirrel.Delete("notification")

	if pars.ID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.OrderItemID != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_item_id": pars.OrderItemID})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

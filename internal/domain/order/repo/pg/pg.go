package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mb-feedback/internal/domain/order/model"
	"mb-feedback/internal/errs"
)

type Repo struct {
	Con *pgxpool.Pool
}

func New(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

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

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.ID, &result.ExternalOrderID, &result.UserPhone, &result.UserName, &result.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, true, nil
}

func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.Order, int64, error) {
	queryBuilder := squirrel.
		Select("id", "external_order_id", "user_phone", "user_name", "created_at").
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
		err = rows.Scan(&data.ID, &data.ExternalOrderID, &data.UserPhone, &data.UserName, &data.CreatedAt)
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

func (r *Repo) ListOrdersNotInDetails(ctx context.Context, pars *model.ListPars) ([]*model.Order, error) {
	queryBuilder := squirrel.
		Select("o.id", "o.external_order_id", "o.user_phone", "o.user_name", "o.created_at").
		From("ord o").
		LeftJoin("ord_detail od ON o.id = od.order_id").
		Where("od.order_id IS NULL")
	if pars.CreatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"o.created_at": pars.CreatedAfter})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.Con.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.ExternalOrderID, &order.UserPhone, &order.UserName, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return orders, nil
}

func (r *Repo) Create(ctx context.Context, obj *model.Edit) error {
	insert := squirrel.Insert("ord").
		Columns("external_order_id", "user_phone", "user_name").
		Values(obj.ExternalOrderID, obj.UserPhone, obj.UserName).
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

func (r *Repo) CreateBatch(ctx context.Context, objects []*model.Edit) error {

	query := squirrel.Insert("ord").Columns("external_order_id", "user_phone", "user_name")

	for _, obj := range objects {
		query = query.Values(obj.ExternalOrderID, *obj.UserPhone, obj.UserName)
	}

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := r.Con.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	return nil
}

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

package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"mb-feedback/internal/domain/order_detail/model"
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

func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.OrderDetail, bool, error) {
	if !pars.IsValid() {
		return nil, false, errs.InvalidInput
	}

	var result model.OrderDetail

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

func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetail, int64, error) {
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

	var result []*model.OrderDetail
	for rows.Next() {
		var data model.OrderDetail
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

func (r *Repo) ListDetailNotInNotification(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetailWithUserInfo, error) {
	queryBuilder := squirrel.
		Select(
			"od.id AS order_detail_id",
			"od.product_code",
			"o.external_order_id AS order_id",
			"o.user_phone",
			"o.user_name",
		).
		From("ord_detail od").
		LeftJoin("ord o ON od.order_id = o.id").
		Where("NOT EXISTS (SELECT 1 FROM notification n WHERE n.order_item_id = od.id)")

	if pars.CreatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"od.created_at": pars.CreatedAfter})
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

	var results []*model.OrderDetailWithUserInfo
	for rows.Next() {
		var detail model.OrderDetailWithUserInfo
		if err := rows.Scan(&detail.ID, &detail.ProductCode, &detail.OrderID, &detail.UserPhone, &detail.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, &detail)
	}

	return results, nil
}

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

func (r *Repo) CreateBatch(ctx context.Context, objects []*model.Edit) error {

	query := squirrel.Insert("ord_detail").Columns("order_id", "product_code")

	for _, obj := range objects {
		query = query.Values(obj.OrderID, *obj.ProductCode)
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

func (r *Repo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.Con.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *Repo) HandleTxCompletion(tx pgx.Tx, err *error) {
	if p := recover(); p != nil {
		_ = tx.Rollback(context.Background())
		panic(p)
	} else if *err != nil {
		_ = tx.Rollback(context.Background())
	} else {
		*err = tx.Commit(context.Background())
	}
}

package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"mb-feedback/internal/domain/order_detail/model"
	"mb-feedback/internal/errs"
)

type Service struct {
	repoDB      RepoDBI
	repoFetcher RepoFetcherI
}

func New(repoDB RepoDBI, repoFetcher RepoFetcherI) *Service {
	return &Service{
		repoDB:      repoDB,
		repoFetcher: repoFetcher,
	}
}

type RepoDBI interface {
	Get(ctx context.Context, pars *model.GetPars) (*model.OrderDetail, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetail, int64, error)
	ListDetailNotInNotification(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetailWithUserInfo, error)
	Create(ctx context.Context, obj *model.Edit) error
	CreateBatch(ctx context.Context, objects []*model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
	HandleTxCompletion(tx pgx.Tx, err *error)
}

type RepoFetcherI interface {
	FetchProductCodes(ctx context.Context, externalOrderID string) ([]string, error)
}

func (s *Service) list(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetail, int64, error) {
	return s.repoDB.List(ctx, pars)
}

// ListDetailWithoutNotification retrieves a list of OrderDetails with UserInfo based on the provided filtering parameters,
// delegating the operation to the database repository.
func (s *Service) ListDetailWithoutNotification(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetailWithUserInfo, error) {
	return s.repoDB.ListDetailNotInNotification(ctx, pars)
}

func (s *Service) create(ctx context.Context, obj *model.Edit) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) CreateList(ctx context.Context, objs []*model.Edit) error {
	return s.repoDB.CreateBatch(ctx, objs)
}

func (s *Service) get(ctx context.Context, pars *model.GetPars, errNE bool) (*model.OrderDetail, bool, error) {
	result, found, err := s.repoDB.Get(ctx, pars)
	if err != nil {
		return nil, false, fmt.Errorf("repoDb.Get: %w", err)
	}
	if !found {
		if errNE {
			return nil, false, errs.ObjectNotFound
		}
		return nil, false, nil
	}

	return result, found, nil
}

func (s *Service) update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	return s.repoDB.Update(ctx, pars, obj)
}

func (s *Service) delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Delete(ctx, pars)
}

func (s *Service) FetchProductCodesByOrder(ctx context.Context, externalOrderID string) ([]string, error) {
	return s.repoFetcher.FetchProductCodes(ctx, externalOrderID)
}

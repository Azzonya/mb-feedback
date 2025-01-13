package service

import (
	"context"
	"mb-feedback/internal/client/fetcher"
	"mb-feedback/internal/domain/order_detail/model"
)

// Service provides methods to manage OrderDetail, handling password operations,
// user validation, and CRUD operations through the repository interface.
type Service struct {
	repoDB  RepoDBI
	fetcher fetcher.Fetcher
}

// New creates a new Service instance with the given database repository.
func New(repoDB RepoDBI, fetcher fetcher.Fetcher) *Service {
	return &Service{
		repoDB:  repoDB,
		fetcher: fetcher,
	}
}

// RepoDBI defines the interface for database interactions related to OrderDetails.
// It includes methods for retrieving, listing, creating, updating, deleting, and checking
// the existence of users in the database.
type RepoDBI interface {
	Get(ctx context.Context, pars *model.GetPars) (*model.OrderDetail, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetail, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, pars *model.GetPars) (bool, error)
}

// List retrieves a list of OrderDetails based on the provided filtering parameters,
// delegating the operation to the database repository.
func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.OrderDetail, int64, error) {
	return s.repoDB.List(ctx, pars)
}

// Create stores a new OrderDetail in the database.
func (s *Service) Create(ctx context.Context, obj *model.Edit) error {
	return s.repoDB.Create(ctx, obj)
}

// Get retrieves a OrderDetail from the database based on the provided query parameters.
func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.OrderDetail, bool, error) {
	return s.repoDB.Get(ctx, pars)
}

// Update modifies an existing user account in the database.
func (s *Service) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	return s.repoDB.Update(ctx, pars, obj)
}

// Delete removes a OrderDetail from the database.
func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Delete(ctx, pars)
}

// Exists checks whether a OrderDetail exists in the database based on the provided query parameters.
func (s *Service) Exists(ctx context.Context, pars *model.GetPars) (bool, error) {
	return s.repoDB.Exists(ctx, pars)
}

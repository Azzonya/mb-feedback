package service

import (
	"context"
	"fmt"
	"mb-feedback/internal/client/notifier"
	"mb-feedback/internal/domain/notification/model"
	"mb-feedback/internal/errs"
)

type Service struct {
	repoDB   RepoDBI
	notifier notifier.Notifier
}

func New(repoDB RepoDBI, notifier notifier.Notifier) *Service {
	return &Service{
		repoDB:   repoDB,
		notifier: notifier,
	}
}

type RepoDBI interface {
	Get(ctx context.Context, pars *model.GetPars) (*model.Notification, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.Notification, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
}

func (s *Service) list(ctx context.Context, pars *model.ListPars) ([]*model.Notification, int64, error) {
	return s.repoDB.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *model.Edit) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) get(ctx context.Context, pars *model.GetPars, errNE bool) (*model.Notification, bool, error) {
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

func (s *Service) Notify(ctx context.Context, orderID, userPhone, userName, productCode string) error {
	return s.notifier.SendNotification(ctx, orderID, userPhone, userName, productCode)
}

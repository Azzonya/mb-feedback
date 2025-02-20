package service

import (
	"context"
	"fmt"
	"log/slog"
	"mb-feedback/internal/domain/order/model"
	"mb-feedback/internal/errs"
	"strings"
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
	Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.Order, int64, error)
	ListOrdersNotInDetails(ctx context.Context, pars *model.ListPars) ([]*model.Order, error)
	Create(ctx context.Context, obj *model.Edit) error
	CreateBatch(ctx context.Context, objects []*model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
}

type RepoFetcherI interface {
	FetchOrders(ctx context.Context) ([]*model.Order, error)
}

func (s *Service) list(ctx context.Context, pars *model.ListPars) ([]*model.Order, int64, error) {
	return s.repoDB.List(ctx, pars)
}

// ListOrdersWithoutDetails retrieves a list of Orders based on the provided filtering parameters and without details,
// delegating the operation to the database repository.
func (s *Service) ListOrdersWithoutDetails(ctx context.Context, pars *model.ListPars) ([]*model.Order, error) {
	return s.repoDB.ListOrdersNotInDetails(ctx, pars)
}

func (s *Service) create(ctx context.Context, obj *model.Edit) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) get(ctx context.Context, pars *model.GetPars, errNE bool) (*model.Order, bool, error) {
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

// FetchOrdersFromExternalSource fetch orders from external source, after insert them to DB.
func (s *Service) FetchOrdersFromExternalSource(ctx context.Context) error {
	fetchedOrders, err := s.repoFetcher.FetchOrders(ctx)
	if err != nil {
		return err
	}
	if len(fetchedOrders) == 0 {
		return nil
	}

	externalOrderIDs := make([]string, len(fetchedOrders))
	for i, order := range fetchedOrders {
		externalOrderIDs[i] = order.ExternalOrderID
	}

	existingOrders, _, err := s.repoDB.List(ctx, &model.ListPars{
		ExternalOrderIDs: &externalOrderIDs,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch existing orders from DB: %w", err)
	}

	existingOrderMap := make(map[string]struct{}, len(existingOrders))
	for _, order := range existingOrders {
		existingOrderMap[order.ExternalOrderID] = struct{}{}
	}

	var ordersToInsert []*model.Edit
	for _, order := range fetchedOrders {
		if _, exists := existingOrderMap[order.ExternalOrderID]; exists {
			continue
		}

		var userPhone string
		userPhone, err = s.formatPhoneNumber(order.UserPhone)
		if err != nil {
			slog.Info("Failed to format phone number", "orderID", order.ExternalOrderID, "phoneNumber", order.UserPhone, "error", err.Error())
			continue
		}

		ordersToInsert = append(ordersToInsert, &model.Edit{
			ExternalOrderID: order.ExternalOrderID,
			UserPhone:       &userPhone,
			UserName:        &order.UserName,
		})
	}

	if len(ordersToInsert) > 0 {
		if err = s.repoDB.CreateBatch(ctx, ordersToInsert); err != nil {
			return fmt.Errorf("failed to insert orders to DB: %w", err)
		}
	} else {
		return fmt.Errorf("orders were missed trying to be added")
	}

	return nil
}

// formatPhoneNumber accepts a phone number in various formats
// and returns it in the format +77XXXXXXXXX.
func (s *Service) formatPhoneNumber(phone string) (string, error) {
	normalizedPhone := strings.TrimSpace(phone)

	if len(normalizedPhone) == 10 {
		return "+7" + normalizedPhone, nil
	} else if len(normalizedPhone) == 11 && strings.HasPrefix(normalizedPhone, "7") {
		return "+" + normalizedPhone, nil
	}

	return "", fmt.Errorf("invalid phone number format: %s", phone)
}

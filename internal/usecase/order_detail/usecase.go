package order_detail

import (
	"context"
	"fmt"
	orderModel "mb-feedback/internal/domain/order/model"
	orderDetailModel "mb-feedback/internal/domain/order_detail/model"
)

type OrderServiceI interface {
	ListOrdersWithoutDetails(ctx context.Context, pars *orderModel.ListPars) ([]*orderModel.Order, error)
}

type OrderDetailServiceI interface {
	FetchProductCodesByOrder(ctx context.Context, externalOrderID string) ([]string, error)
	CreateList(ctx context.Context, objs []*orderDetailModel.Edit) error
}

type Usecase struct {
	orderService       OrderServiceI
	orderDetailService OrderDetailServiceI
}

func New(orderService OrderServiceI, orderDetailService OrderDetailServiceI) *Usecase {
	return &Usecase{
		orderService:       orderService,
		orderDetailService: orderDetailService,
	}
}

func (u *Usecase) FetchProductCodes(ctx context.Context) error {
	//createdAfter := time.Now().Add(-time.Hour) // заказы за крайний час

	missingOrders, err := u.orderService.ListOrdersWithoutDetails(ctx, &orderModel.ListPars{
		//CreatedAfter: &createdAfter,
	})
	if err != nil {
		return err
	}

	for _, missingOrder := range missingOrders {
		if err = u.processMissingOrder(ctx, missingOrder); err != nil {
			return fmt.Errorf("failed to process missing order %s: %w", missingOrder.ExternalOrderID, err)
		}
	}

	return nil
}

func (u *Usecase) processMissingOrder(ctx context.Context, missingOrder *orderModel.Order) error {
	productCodes, err := u.orderDetailService.FetchProductCodesByOrder(ctx, missingOrder.ExternalOrderID)
	if err != nil {
		return fmt.Errorf("failed to fetch product codes for order %s: %w", missingOrder.ExternalOrderID, err)
	}

	orderDetailEdit := make([]*orderDetailModel.Edit, 0, len(productCodes))
	for _, productCode := range productCodes {
		orderDetailEdit = append(orderDetailEdit, &orderDetailModel.Edit{
			OrderID:     missingOrder.ID,
			ProductCode: &productCode,
		})
	}

	if err = u.orderDetailService.CreateList(ctx, orderDetailEdit); err != nil {
		return fmt.Errorf("failed to create order details for order %s: %w", missingOrder.ExternalOrderID, err)
	}

	return nil
}

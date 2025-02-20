package fetcher

import (
	"context"
	orderModel "mb-feedback/internal/domain/order/model"
)

type Fetcher interface {
	FetchCompletedOrders(ctx context.Context) ([]*orderModel.Order, error)
	FetchProductCodesByOrder(ctx context.Context, orderID string) ([]string, error)
}

package fetcher

import (
	"context"
	mb_broker "mb-feedback/internal/client/fetcher/mb-broker"
	orderModel "mb-feedback/internal/domain/order/model"
)

type Repo struct {
	client *mb_broker.Client
}

func New(client *mb_broker.Client) *Repo {
	return &Repo{
		client: client,
	}
}

func (r *Repo) FetchOrders(ctx context.Context) ([]*orderModel.Order, error) {
	result, err := r.client.FetchCompletedOrders(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

package fetcher

import (
	"context"
	mb_broker "mb-feedback/internal/client/fetcher/mb-broker"
)

type Repo struct {
	client *mb_broker.Client
}

func New(client *mb_broker.Client) *Repo {
	return &Repo{
		client: client,
	}
}

func (r *Repo) FetchProductCodes(ctx context.Context, externalOrderID string) ([]string, error) {
	result, err := r.client.FetchProductCodes(ctx, externalOrderID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

package order

import "context"

type OrderServiceI interface {
	FetchOrdersFromExternalSource(ctx context.Context) error
}

type Usecase struct {
	orderService OrderServiceI
}

func New(orderService OrderServiceI) *Usecase {
	return &Usecase{
		orderService: orderService,
	}
}

func (u *Usecase) FetchNewOrders(ctx context.Context) error {
	return u.orderService.FetchOrdersFromExternalSource(ctx)
}

package notification

import (
	"context"
	"fmt"
	"mb-feedback/internal/cns"
	notificationModel "mb-feedback/internal/domain/notification/model"
	orderDetail "mb-feedback/internal/domain/order_detail/model"
	"time"
)

type OrderDetailServiceI interface {
	ListDetailWithoutNotification(ctx context.Context, pars *orderDetail.ListPars) ([]*orderDetail.OrderDetailWithUserInfo, error)
}

type NotificationServiceI interface {
	Notify(ctx context.Context, orderID, userPhone, userName, productCode string) error
	Create(ctx context.Context, obj *notificationModel.Edit) error
}

type Usecase struct {
	orderDetailService  OrderDetailServiceI
	notificationService NotificationServiceI
}

func New(orderDetailService OrderDetailServiceI, notificationService NotificationServiceI) *Usecase {
	return &Usecase{
		orderDetailService:  orderDetailService,
		notificationService: notificationService,
	}
}

func (u *Usecase) SendNotification(ctx context.Context) error {
	createdAfter := time.Now().Add(-time.Hour) // заказы за крайний час

	details, err := u.orderDetailService.ListDetailWithoutNotification(ctx, &orderDetail.ListPars{
		CreatedAfter: &createdAfter,
	})
	if err != nil {
		return fmt.Errorf("failed to list details without notification: %w", err)
	}

	for _, detail := range details {
		if err = u.processNotification(ctx, detail); err != nil {
			return fmt.Errorf("failed to process notification for detail ID %d: %w", detail.ID, err)
		}
	}

	return nil
}

func (u *Usecase) processNotification(ctx context.Context, detail *orderDetail.OrderDetailWithUserInfo) error {
	errNotify := u.notificationService.Notify(ctx, detail.OrderID, detail.UserPhone, detail.UserName, detail.ProductCode)

	status := cns.StatusFailed
	if errNotify == nil {
		status = cns.StatusSent
	}

	sentAt := time.Now()

	err := u.notificationService.Create(ctx, &notificationModel.Edit{
		OrderItemID: &detail.ID,
		PhoneNumber: &detail.UserPhone,
		Status:      &status,
		SentAt:      &sentAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create notification log: %w", err)
	}

	return nil
}

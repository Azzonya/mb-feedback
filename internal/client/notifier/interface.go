package notifier

import "context"

type Notifier interface {
	SendNotification(ctx context.Context, orderID, userPhone, userName, productCode string) error
}

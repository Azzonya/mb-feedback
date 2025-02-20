package rest

import (
	"context"
	"log/slog"
	"net/http"
)

// FetchOrdersHandler handles updating the list of orders
func (s *Rest) FetchOrdersHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	go func() {
		slog.Info("Updating orders")
		s.updateOrderMutex.Lock()
		err := s.orderUsc.FetchNewOrders(context.Background())
		if err != nil {
			slog.Error("Error updating orders: ", "error", err)
		} else {
			slog.Info("Updated orders")
		}
		s.updateOrderMutex.Unlock()
	}()
}

// GetProductCodesHandler handles fetching product codes for a given order
func (s *Rest) GetProductCodesHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	go func() {
		slog.Info("Fetching product codes")

		s.getProductCodeMutex.Lock()
		err := s.orderDetailUsc.FetchProductCodes(context.Background())
		if err != nil {
			slog.Error("Error fetching product codes: ", err)
		} else {
			slog.Info("Fetched product codes")
		}
		s.getProductCodeMutex.Unlock()
	}()
}

// SendNotificationHandler handles sending SMS for a given order
func (s *Rest) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	go func() {
		slog.Info("Sending notifications")

		s.sendNotificationMutex.Lock()
		err := s.notificationUsc.SendNotification(context.Background())
		if err != nil {
			slog.Error("Error sending notification: ", err)
		} else {
			slog.Info("Sent notification")
		}
		s.sendNotificationMutex.Unlock()
	}()
}

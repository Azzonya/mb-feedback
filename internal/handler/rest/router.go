package rest

import (
	"context"
	"errors"
	"log/slog"
	notificationUsecase "mb-feedback/internal/usecase/notification"
	orderUsecase "mb-feedback/internal/usecase/order"
	orderDetailUsecase "mb-feedback/internal/usecase/order_detail"
	"net/http"
	"sync"
	"time"
)

type Rest struct {
	httpServer      *http.Server
	orderUsc        *orderUsecase.Usecase
	orderDetailUsc  *orderDetailUsecase.Usecase
	notificationUsc *notificationUsecase.Usecase

	updateOrderMutex      sync.Mutex
	getProductCodeMutex   sync.Mutex
	sendNotificationMutex sync.Mutex

	ErrorChan chan error
}

func New(orderUsc *orderUsecase.Usecase, orderDetailUsc *orderDetailUsecase.Usecase, notificationUsc *notificationUsecase.Usecase) *Rest {
	return &Rest{
		orderUsc:        orderUsc,
		orderDetailUsc:  orderDetailUsc,
		notificationUsc: notificationUsc,

		ErrorChan: make(chan error, 1),
	}
}

func (s *Rest) Start(addr string) {

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("GET /fetch-orders", s.FetchOrdersHandler)
	httpMux.HandleFunc("GET /get-product-codes", s.GetProductCodesHandler)
	httpMux.HandleFunc("GET /send-notification", s.SendNotificationHandler)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: httpMux,
	}

	go func() {
		slog.Info("Server is running on", "address", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server:", "error", err)
		}
	}()

}

func (s *Rest) Stop() error {
	slog.Info("Server is shutting down")
	defer close(s.ErrorChan)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", err)
		return err
	}

	slog.Info("Server stopped gracefully")
	return nil
}

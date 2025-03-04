package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	mb_broker "mb-feedback/internal/client/fetcher/mb-broker"
	"mb-feedback/internal/client/notifier/voximplant"
	"mb-feedback/internal/conf"
	notificationRepoPG "mb-feedback/internal/domain/notification/repo/pg"
	NotificationService "mb-feedback/internal/domain/notification/service"
	orderRepoFetcher "mb-feedback/internal/domain/order/repo/fetcher"
	orderRepoPG "mb-feedback/internal/domain/order/repo/pg"
	OrderService "mb-feedback/internal/domain/order/service"
	orderDetailRepoFetcher "mb-feedback/internal/domain/order_detail/repo/fetcher"
	orderDetailRepoPG "mb-feedback/internal/domain/order_detail/repo/pg"
	OrderDetailService "mb-feedback/internal/domain/order_detail/service"
	"mb-feedback/internal/handler/rest"
	NotificationUsecase "mb-feedback/internal/usecase/notification"
	OrderUsecase "mb-feedback/internal/usecase/order"
	OrderDetailUsecase "mb-feedback/internal/usecase/order_detail"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	pgpool *pgxpool.Pool

	mbBrokerClient   *mb_broker.Client
	voximplantClient *voximplant.Client

	// order
	orderUsc *OrderUsecase.Usecase
	orderSrv *OrderService.Service

	// order-detail
	orderDetailUsc *OrderDetailUsecase.Usecase
	orderDetailSrv *OrderDetailService.Service

	// notification
	notificationUsc *NotificationUsecase.Usecase
	notificationSrv *NotificationService.Service

	httpServer *rest.Rest

	exitCode int
}

func (a *App) Init() {
	var err error

	// logger
	{
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	// pgpool
	{
		a.pgpool, err = pgxpool.New(context.Background(), conf.Conf.PgDsn)
		errCheck(err, "pgxpool.New")
	}

	// mb-broker client
	{
		a.mbBrokerClient = mb_broker.New(conf.Conf.MbBrokerURL, conf.Conf.MbBrokerToken)
	}

	// voximplant
	{
		a.voximplantClient = voximplant.New(
			conf.Conf.VoximplantURL,
			conf.Conf.VoximplantToken,
			conf.Conf.VoximplantDomainName,
			conf.Conf.VoximplantTemplateID,
			conf.Conf.VoximplantChannelID)
	}

	// order
	{
		orderRepoDB := orderRepoPG.New(a.pgpool)
		orderFetcherRepo := orderRepoFetcher.New(a.mbBrokerClient)

		a.orderSrv = OrderService.New(orderRepoDB, orderFetcherRepo)
		a.orderUsc = OrderUsecase.New(a.orderSrv)
	}

	// order-detail
	{
		orderDetailRepoDB := orderDetailRepoPG.New(a.pgpool)
		orderDetailFetcherRepo := orderDetailRepoFetcher.New(a.mbBrokerClient)
		a.orderDetailSrv = OrderDetailService.New(orderDetailRepoDB, orderDetailFetcherRepo)

		a.orderDetailUsc = OrderDetailUsecase.New(a.orderSrv, a.orderDetailSrv)
	}

	// notification
	{
		notificationRepoDB := notificationRepoPG.New(a.pgpool)

		a.notificationSrv = NotificationService.New(notificationRepoDB, a.voximplantClient)
		a.notificationUsc = NotificationUsecase.New(a.orderDetailSrv, a.notificationSrv)
	}

	// http-server
	{
		a.httpServer = rest.New(a.orderUsc, a.orderDetailUsc, a.notificationUsc)
	}
}

func (a *App) Start() {

	// recover
	{
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered in Start", "panic", r)

				a.exitCode = 2
				os.Exit(a.exitCode)
			}
		}()
	}

	slog.Info("Starting")

	// http-server
	{
		a.httpServer.Start(conf.Conf.HTTPListen)
	}
}

func (a *App) Listen() {
	select {
	case <-StopSignal():
	case <-a.httpServer.ErrorChan:
	}
}

func StopSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return ch
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	// http-server
	{
		err := a.httpServer.Stop()
		if err != nil {
			panic(err)
		}
	}
}

func (a *App) Exit() {
	slog.Info("Exit")

	os.Exit(a.exitCode)
}

// errCheck checks if an error occurred and logs it with the specified message.
// If an error is found, the function logs the error and terminates the program.
// If a message is provided, it is included in the logged output.
func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}

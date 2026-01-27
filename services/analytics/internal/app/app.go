package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nats-io/nats.go"

	"github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	"github.com/egor200512/URL_shortener/services/analytics/internal/service"
	"github.com/egor200512/URL_shortener/shared/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	HttpConf *configs.HttpConf
	GrpcConf *configs.GrpcConf
	NatsConf *configs.NatsConf

	AnalyticsHandler desc.AnalyticsServiceServer
	AnalyticsService service.IAnalyticsService
	AnalyticsRepo    repository.IAnalyticsRepo

	httpServer   *http.Server `wire:"-"`
	grpcServer   *grpc.Server `wire:"-"`
	NatsConsumer broker.IConsumer
}

func NewApp(
	httpConf *configs.HttpConf,
	grpcConf *configs.GrpcConf,
	natsConf *configs.NatsConf,
	analyticsHandler desc.AnalyticsServiceServer,
	analyticsService service.IAnalyticsService,
	analyticsRepo repository.IAnalyticsRepo,
	natsConsumer broker.IConsumer,
) *App {
	return &App{
		HttpConf:         httpConf,
		GrpcConf:         grpcConf,
		NatsConf:         natsConf,
		AnalyticsHandler: analyticsHandler,
		AnalyticsService: analyticsService,
		AnalyticsRepo:    analyticsRepo,
		NatsConsumer:     natsConsumer,
	}
}

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	reflection.Register(a.grpcServer)
	desc.RegisterAnalyticsServiceServer(a.grpcServer, a.AnalyticsHandler)
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	httpAddr := a.HttpConf.AnalyticsAddress()
	grpcAddr := a.GrpcConf.AnalyticsAddress()

	router := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := desc.RegisterAnalyticsServiceHandlerFromEndpoint(ctx, router, grpcAddr, opts); err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Handler: router,
		Addr:    httpAddr,
	}
	return nil
}

func (a *App) runGRPCServer() error {
	grpcAddr := a.GrpcConf.AnalyticsAddress()
	log.Printf("GRPC server is running on %s\n", grpcAddr)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}
	if err = a.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (a *App) runHTTPServer() error {
	httpAddr := a.HttpConf.AnalyticsAddress()
	log.Printf("HTTP server is running on %s\n", httpAddr)
	if err := a.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.initGRPCServer(ctx); err != nil {
		return err
	}
	if err := a.initHTTPServer(ctx); err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.runGRPCServer(); err != nil {
			log.Println("grpc error:", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			log.Println("http error:", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			ms, err := a.NatsConsumer.Fetch(
				a.NatsConf.ConsumerBatch(),
				a.NatsConf.ConsumerWait(),
			)
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}
				log.Printf("nats fetch error: %v\n", err)
				continue
			}

			const workerPoolSize = 16

			pool := make(chan struct{}, workerPoolSize)
			for range workerPoolSize {
				pool <- struct{}{}
			}

			for _, msg := range ms {
				<-pool
				go func(m *nats.Msg) {
					defer func() { pool <- struct{}{} }()
					select {
					case <-ctx.Done():
						return
					default:
						if err := a.AnalyticsService.HandleMessage(ctx, m); err != nil {
							log.Printf("failed to handle nats message: %s\n", err.Error())
							return
						}
						if err := m.Ack(); err != nil {
							log.Printf("failed to ack nats message: %s\n", err.Error())
						}
					}
				}(msg)
			}

			for range workerPoolSize {
				<-pool
			}
		}
	}()

	wg.Wait()
	return nil
}

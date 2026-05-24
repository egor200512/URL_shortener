package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

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

	AnalyticsHandler desc.AnalyticsServiceServer
	AnalyticsService service.IAnalyticsService
	AnalyticsRepo    repository.IAnalyticsRepo

	httpServer *http.Server `wire:"-"`
	grpcServer *grpc.Server `wire:"-"`
}

func NewApp(
	httpConf *configs.HttpConf,
	grpcConf *configs.GrpcConf,
	analyticsHandler desc.AnalyticsServiceServer,
	analyticsService service.IAnalyticsService,
	analyticsRepo repository.IAnalyticsRepo,
) *App {
	return &App{
		HttpConf:         httpConf,
		GrpcConf:         grpcConf,
		AnalyticsHandler: analyticsHandler,
		AnalyticsService: analyticsService,
		AnalyticsRepo:    analyticsRepo,
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
	return a.grpcServer.Serve(lis)
}

func (a *App) runHTTPServer() error {
	httpAddr := a.HttpConf.AnalyticsAddress()
	log.Printf("HTTP server is running on %s\n", httpAddr)

	return a.httpServer.ListenAndServe()
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

	wg.Wait()
	return nil
}

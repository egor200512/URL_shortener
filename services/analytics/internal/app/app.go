package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	"github.com/egor200512/URL_shortener/services/analytics/internal/service"
	"github.com/egor200512/URL_shortener/shared/configs"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	"github.com/egor200512/URL_shortener/shared/pkg/events"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type App struct {
	HttpConf *configs.HttpConf
	GrpcConf *configs.GrpcConf

	AnalyticsHandler desc.AnalyticsServiceServer
	AnalyticsService service.IAnalyticsService
	AnalyticsRepo    repository.IAnalyticsRepo
	Consumer         broker.IConsumer

	httpServer *http.Server `wire:"-"`
	grpcServer *grpc.Server `wire:"-"`
}

func NewApp(
	httpConf *configs.HttpConf,
	grpcConf *configs.GrpcConf,
	analyticsHandler desc.AnalyticsServiceServer,
	analyticsService service.IAnalyticsService,
	analyticsRepo repository.IAnalyticsRepo,
	consumer broker.IConsumer,
) *App {
	return &App{
		HttpConf:         httpConf,
		GrpcConf:         grpcConf,
		AnalyticsHandler: analyticsHandler,
		AnalyticsService: analyticsService,
		AnalyticsRepo:    analyticsRepo,
		Consumer:         consumer,
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

func (a *App) runConsumer(ctx context.Context) error {
	if a.Consumer == nil {
		return nil
	}

	log.Println("NATS consumer is running")
	return a.Consumer.Consume(ctx, a.handleLinkEvent)
}

func (a *App) handleLinkEvent(ctx context.Context, payload []byte) error {
	var event events.LinkEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	if event.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if event.UserID == "" {
		return fmt.Errorf("user id is required")
	}
	if event.ShortLink == "" {
		return fmt.Errorf("short link is required")
	}
	if event.OriginalLink == "" {
		return fmt.Errorf("original link is required")
	}

	req := &desc.RecordEventRequest{
		EventType:    event.EventType,
		UserId:       event.UserID,
		ShortLink:    event.ShortLink,
		OriginalLink: event.OriginalLink,
	}
	if !event.ExecutedAt.IsZero() {
		req.ExecutedAt = timestamppb.New(event.ExecutedAt)
	}

	return a.AnalyticsService.RecordEvent(ctx, req)
}

func (a *App) Run(ctx context.Context) error {
	if a.Consumer != nil {
		defer a.Consumer.Close()
	}

	if err := a.initGRPCServer(ctx); err != nil {
		return err
	}
	if err := a.initHTTPServer(ctx); err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := a.runGRPCServer(); err != nil {
			log.Println("grpc error:", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			log.Println("http error:", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runConsumer(ctx); err != nil {
			log.Println("consumer error:", err)
		}
	}()

	wg.Wait()
	return nil
}

package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	descA "github.com/egor200512/URL_shortener/shared/gen/auth"
	descL "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/egor200512/URL_shortener/shared/pkg/broker"
	cache "github.com/egor200512/URL_shortener/shared/pkg/cache"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/egor200512/URL_shortener/shared/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	HttpConf *configs.HttpConf
	GrpcConf *configs.GrpcConf

	LinksHandler descL.LinksServiceServer

	httpServer *http.Server `wire:"-"`
	grpcServer *grpc.Server `wire:"-"`
	AuthClient descA.AuthServiceClient
	Cache      cache.ICache
	Producer   broker.IProducer
}

func NewApp(
	httpConf *configs.HttpConf,
	grpcConf *configs.GrpcConf,
	linksHandler descL.LinksServiceServer,
) *App {
	return &App{
		HttpConf:     httpConf,
		GrpcConf:     grpcConf,
		LinksHandler: linksHandler,
	}
}

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(jwt.AuthIntersepter(a.AuthClient)),
	)
	reflection.Register(a.grpcServer)
	descL.RegisterLinksServiceServer(a.grpcServer, a.LinksHandler)
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	httpAddr := a.HttpConf.LinksAddress()
	grpcAddr := a.GrpcConf.LinksAddress()

	router := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := descL.RegisterLinksServiceHandlerFromEndpoint(ctx, router, grpcAddr, opts); err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Handler: router,
		Addr:    httpAddr,
	}
	return nil
}

func (a *App) runGRPCServer() error {
	grpcAddr := a.GrpcConf.LinksAddress()
	slog.Info("grpc server is running", "addr", grpcAddr)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}
	if err = a.grpcServer.Serve(lis); err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return err
	}
	return nil
}

func (a *App) runHTTPServer() error {
	httpAddr := a.HttpConf.LinksAddress()
	slog.Info("http server is running", "addr", httpAddr)
	if err := a.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (a *App) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
	}

	if a.grpcServer != nil {
		done := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
		case <-shutdownCtx.Done():
			a.grpcServer.Stop()
			return shutdownCtx.Err()
		}
	}

	if a.Producer != nil {
		a.Producer.Close()
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

	errCh := make(chan error, 2)

	go func() {
		if err := a.runGRPCServer(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := a.runHTTPServer(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown links app")
		return a.shutdown()
	case err := <-errCh:
		slog.Error("links app error", "error", err)
		if shutdownErr := a.shutdown(); shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
			slog.Error("links shutdown error", "error", shutdownErr)
		}
		return err
	}
}

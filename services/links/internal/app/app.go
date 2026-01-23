package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

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
	Broker     broker.IProducer
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
	httpAddr := a.HttpConf.LinksAddress()
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
	wg.Add(2)

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

	wg.Wait()
	return nil
}

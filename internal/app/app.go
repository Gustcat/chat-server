package app

import (
	"context"
	"github.com/Gustcat/chat-server/internal/closer"
	"github.com/Gustcat/chat-server/internal/config"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
	}
	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return app.runGRPCServer()
}

func (app *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	log.Println("env-fail are used")
	if err != nil {
		return err
	}
	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(app.grpcServer)

	desc.RegisterChatV1Server(app.grpcServer, app.serviceProvider.Impl(ctx))
	log.Printf("<serviceProvider> - %+v", app.serviceProvider)

	return nil
}

func (app *App) runGRPCServer() error {
	log.Printf("GRPC server listen:  %s", app.serviceProvider.GRPCConfig().Address())
	lis, err := net.Listen("tcp", app.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = app.grpcServer.Serve(lis)

	return nil
}

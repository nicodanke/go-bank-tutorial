package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nicodanke/bankTutorial/api"
	db "github.com/nicodanke/bankTutorial/db/sqlc"
	_ "github.com/nicodanke/bankTutorial/doc/statik"
	"github.com/nicodanke/bankTutorial/gapi"
	"github.com/nicodanke/bankTutorial/pb"
	"github.com/nicodanke/bankTutorial/utils"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("Cannot load configuration")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Error().Err(err).Msg("Cannot connect to db")

	}

	runDBMigrations(config.MigrationUrl, config.DBSource)

	store := db.NewStore(connPool)

	go runGRPCGatewayServer(config, store)
	runGRPCServer(config, store)
}

func runDBMigrations(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("Failed to run migrate up")
	}

	log.Info().Msg("DB migrations runned successfully")
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterBankTutorialServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create listener")
	}

	log.Info().Msgf("gRPC server started at: %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Err(err).Msg("Cannot start gRPC server")
	}
}

func runGRPCGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create server")
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankTutorialHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Error().Err(err).Msg("Cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Error().Err(err).Msg("Cannot create static file system")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("Cannor create listener")
	}

	log.Info().Msgf("HTTP Gateway Server started at: %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Error().Err(err).Msg("Cannor start HTTP Gateway Server")
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("Cannor create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("Cannot start server")
	}
}

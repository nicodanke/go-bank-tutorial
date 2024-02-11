package main

import (
	"context"
	"log"
	"net"
	"net/http"

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration:", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	runDBMigrations(config.MigrationUrl, config.DBSource)

	store := db.NewStore(connPool)

	go runGRPCGatewayServer(config, store)
	runGRPCServer(config, store)
}

func runDBMigrations(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal("Cannor create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange{
		log.Fatal("Failed to run migrate up:", err)
	}

	log.Println("DB migrations runned successfully")
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannor create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankTutorialServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannor create listener:", err)
	}

	log.Printf("gRPC server started at: %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannor start gRPC server:", err)
	}
}

func runGRPCGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannor create server:", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankTutorialHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Cannor register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("Cannot create static file system:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannor create listener:", err)
	}

	log.Printf("HTTP Gateway Server started at: %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Cannor start HTTP Gateway Server:", err)
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannor create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}

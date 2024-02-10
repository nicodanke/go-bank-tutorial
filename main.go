package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nicodanke/bankTutorial/api"
	db "github.com/nicodanke/bankTutorial/db/sqlc"
	"github.com/nicodanke/bankTutorial/gapi"
	"github.com/nicodanke/bankTutorial/pb"
	"github.com/nicodanke/bankTutorial/utils"
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

	store := db.NewStore(connPool)

	go runGRPCGatewayServer(config, store)
	runGRPCServer(config, store)
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

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

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

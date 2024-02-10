package gapi

import (
	"fmt"

	db "github.com/nicodanke/bankTutorial/db/sqlc"
	"github.com/nicodanke/bankTutorial/pb"
	"github.com/nicodanke/bankTutorial/token"
	"github.com/nicodanke/bankTutorial/utils"
)

type Server struct {
	pb.UnimplementedBankTutorialServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}

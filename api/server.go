package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/nicodanke/bankTutorial/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("accounts", server.getAccounts)
	router.GET("accounts/:id", server.getAccount)
	router.POST("accounts", server.createAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

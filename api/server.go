package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/nicodanke/bankTutorial/db/sqlc"
	"github.com/nicodanke/bankTutorial/token"
	"github.com/nicodanke/bankTutorial/utils"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// No Auth Required
	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)

	// Auth Required
	authRoutes.GET("/accounts", server.getAccounts)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.POST("/accounts", server.createAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	authRoutes.GET("/users/:username", server.getUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

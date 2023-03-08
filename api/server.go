package api

import (
	"fmt"

	db "github.com/ghost-codes/simplebank/db/sqlc"
	"github.com/ghost-codes/simplebank/token"
	"github.com/ghost-codes/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	maker, err := token.NewJWTMaker(config.SecretKey)

	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %v", err)
	}
	server := &Server{store: store, tokenMaker: maker, config: config}

	// router
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	fmt.Println("============> Starting server here")
	return server.router.Run()

}

func (server *Server) setupRouter() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	//Add routes to server
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountById)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	authRoutes.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

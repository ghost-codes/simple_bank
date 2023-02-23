package api

import (
	db "github.com/ghost-codes/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}

	// router
	router := gin.Default()

	//Add routes to server
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

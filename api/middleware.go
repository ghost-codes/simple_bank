package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ghost-codes/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authoriztionTypeBearer  = "bearer"
	authorizationPayloadKey = "payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeaderKey)
		if len(authorization) == 0 {
			err := errors.New("Authrization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorization)
		fmt.Println(fields)
		if len(fields) < 2 {
			err := errors.New("Invalid authoriztion header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationtype := strings.ToLower(fields[0])
		if authorizationtype != authoriztionTypeBearer {
			err := fmt.Errorf("unsupported authorizaiton type %v", authorizationtype)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

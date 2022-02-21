package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker, admin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		var adminPayload *token.AdminPayload
		var userPayload *token.UserPayload
		var err error
		if admin {
			adminPayload, err = tokenMaker.VerifyTokenForAdmin(accessToken)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			ctx.Set(authorizationPayloadKey, adminPayload)
			ctx.Next()
		}

		if !admin {
			userPayload, err = tokenMaker.VerifyTokenForUser(accessToken)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			ctx.Set(authorizationPayloadKey, userPayload)
			ctx.Next()

		}

	}
}

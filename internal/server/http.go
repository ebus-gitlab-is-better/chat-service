package server

import (
	"chat-service/internal/biz"
	"chat-service/internal/conf"
	"chat-service/internal/data"
	http1 "net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func ExtractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func AuthMiddleware(api *data.KeycloakAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) < 1 {
			c.JSON(http1.StatusUnauthorized, &gin.H{
				"error": "not token",
			})
			c.Abort()
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.JSON(http1.StatusUnauthorized, &gin.H{
				"error": "not token",
			})
			c.Abort()
			return
		}
		accessToken := authParts[1]

		rptResult, err := api.CheckToken(accessToken)

		if err != nil {
			c.JSON(http1.StatusUnauthorized, &gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		istokenvalid := *rptResult.Active
		if !istokenvalid {
			c.JSON(http1.StatusUnauthorized, &gin.H{
				"error": "token expired",
			})
			c.Abort()
			return
		}
		user, err := api.GetUserInfo(accessToken)

		if err != nil {
			c.JSON(http1.StatusUnauthorized, &gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, uc *biz.ChatUseCase, api *data.KeycloakAPI, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	r := gin.Default()
	r.Use(AuthMiddleware(api))

	srv := http.NewServer(opts...)
	srv.HandlePrefix("/", r)
	return srv
}

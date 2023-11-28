package server

import (
	_ "chat-service/docs"
	"chat-service/internal/biz"
	"chat-service/internal/conf"
	"chat-service/internal/data"
	"chat-service/internal/route"
	slog "log"
	http1 "net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
		slog.Println(rptResult)
		slog.Println(err)
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

//	@title			Chat Service Swagger API
//	@version		1.0
//	@description	This is documentation api for backend

//	@contact.name	Suro
//	@contact.url	https://t.me/suronek
//	@contact.email	suro@hyneo.ru

//	@securityDefinitions.apikey	authorization
//	@in							header
//	@name						authorization

// @host		chats.e-bus.site
// @BasePath	/
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
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route := route.NewChatRoute(uc)
	chatsRoute := r.Group("/chats")
	chatsRoute.Use(AuthMiddleware(api))
	route.Register(chatsRoute)
	srv := http.NewServer(opts...)

	srv.HandlePrefix("/", r)
	return srv
}

package route

import (
	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/a1y/doc-formatter/internal/gateway/clients/auth"
	authhandler "github.com/a1y/doc-formatter/internal/gateway/handler/auth"
	authmanager "github.com/a1y/doc-formatter/internal/gateway/manager/auth"
	"github.com/gin-gonic/gin"

	docs "github.com/a1y/doc-formatter/api/http/gateway/v1"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(config *gateway.Config) (*gin.Engine, error) {
	r := gin.Default()

	docs.SwaggerInfo.Title = "AI Doc Formatter API Gateway"
	docs.SwaggerInfo.Version = "1.0"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Group("/api/v1/")
	setupAPIV1(r, config)

	err := r.Run(config.HTTPAddr)

	return r, err
}

func setupAPIV1(r *gin.Engine, config *gateway.Config) {
	authClient := auth.NewAuthClient(config.AuthGRPCAddr)
	authManager := authmanager.NewAuthManager(authClient)
	authHandler, _ := authhandler.NewAuthHandler(authManager)

	r.Group("/auth")
	{
		r.POST("/signup", authHandler.Signup)
		r.POST("/login", authHandler.Login)
	}
}

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

	v1 := r.Group("/api/v1")
	if err := setupAPIV1(v1, config); err != nil {
		return nil, err
	}

	return r, nil
}

func setupAPIV1(r gin.IRouter, config *gateway.Config) error {
	authClient := auth.NewAuthClient(config.AuthService)
	authManager := authmanager.NewAuthManager(authClient)
	authHandler, err := authhandler.NewAuthHandler(authManager)
	if err != nil {
		return err
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
	}

	return nil
}

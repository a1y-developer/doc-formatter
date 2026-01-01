package route

import (
	"context"

	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/a1y/doc-formatter/internal/gateway/clients/auth"
	storageclient "github.com/a1y/doc-formatter/internal/gateway/clients/storage"
	authhandler "github.com/a1y/doc-formatter/internal/gateway/handler/auth"
	storagehandler "github.com/a1y/doc-formatter/internal/gateway/handler/storage"
	authmanager "github.com/a1y/doc-formatter/internal/gateway/manager/auth"
	storagemanager "github.com/a1y/doc-formatter/internal/gateway/manager/storage"
	"github.com/a1y/doc-formatter/internal/gateway/middleware"
	logutil "github.com/a1y/doc-formatter/internal/gateway/util/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	docs "github.com/a1y/doc-formatter/api/http/gateway/v1"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(config *gateway.Config) (*gin.Engine, error) {
	r := gin.New()
	logger := logutil.GetLogger(context.Background())

	if err := r.SetTrustedProxies(nil); err != nil {
		logger.Error("Failed to set trusted proxies...", zap.Error(err))
		return nil, err
	}

	r.Use(middleware.APILoggerMiddleware(config.Logging, "gateway"))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	docs.SwaggerInfo.Title = "AI Doc Formatter API Gateway"
	docs.SwaggerInfo.Version = "1.0"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	if err := setupAPIV1(v1, config); err != nil {
		logger.Error("Failed to setup API v1...", zap.Error(err))
		return nil, err
	}

	return r, nil
}

func setupAPIV1(r gin.IRouter, config *gateway.Config) error {
	logger := logutil.GetLogger(context.Background())
	logger.Info("Setting up API v1...")

	// Setup clients
	authClient := auth.NewAuthClient(config.AuthService)
	storageClient := storageclient.NewStorageClient(config.StorageService)

	// Setup managers
	authManager := authmanager.NewAuthManager(authClient)
	storageManager := storagemanager.NewStorageManager(storageClient)

	// Setup handlers
	authHandler, err := authhandler.NewAuthHandler(authManager)
	if err != nil {
		logger.Error("Failed to create auth handler...", zap.Error(err))
		return err
	}
	storageHandler, err := storagehandler.NewStorageHandler(storageManager)
	if err != nil {
		logger.Error("Failed to create storage handler...", zap.Error(err))
		return err
	}

	// Setup routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
	}

	storageGroup := r.Group("/storage")
	{
		storageGroup.POST("/upload", storageHandler.UploadFile)
	}

	return nil
}

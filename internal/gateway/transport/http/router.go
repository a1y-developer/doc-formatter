package http

import (
    "github.com/gin-gonic/gin"

    docs "github.com/a1y/ai-doc-formatter/internal/gateway/docs"

    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

func NewRouter(authHandler *AuthHandler) *gin.Engine {
    r := gin.Default()

    docs.SwaggerInfo.Title = "AI Doc Formatter API Gateway"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.BasePath = "/"

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.Group("/api")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/signup", authHandler.Signup)
            auth.POST("/login", authHandler.Login)
        }
    }

    return r
}

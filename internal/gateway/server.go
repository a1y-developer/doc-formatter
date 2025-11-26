package gateway

import (
    "log"

    "github.com/a1y/ai-doc-formatter/internal/gateway/clients"
    httpgw "github.com/a1y/ai-doc-formatter/internal/gateway/transport/http"
)

func StartHTTPServer() {
	authClient := clients.NewAuthClient("localhost:9000")
	authHandler := httpgw.NewAuthHandler(authClient)
	
    r := httpgw.NewRouter(authHandler)
    log.Println("Gateway running at :8080")
    r.Run(":8080")
}

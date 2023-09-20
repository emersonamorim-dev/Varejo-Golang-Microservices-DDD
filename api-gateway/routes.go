package api_gateway

import (
	"github.com/gin-gonic/gin"
)

type Service struct {
	Route string
	URL   string
}

var services = []Service{
	{
		Route: "/customer-service/api/",
		URL:   "http://localhost:8081",
	},

	{
		Route: "/integration-service/api/",
		URL:   "http://localhost:8082",
	},
}

func SetupGatewayRoutes(router *gin.Engine) {
	for _, service := range services {
		proxy := NewServiceProxy(service.URL)
		router.Any(service.Route+"*any", func(c *gin.Context) {
			// Preserve o caminho do URL original
			c.Request.URL.Path = c.Param("any")
			// Define o cabeçalho X-Forwarded-Host
			c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}

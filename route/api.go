package route

import (
	"final/handler/public"
	"final/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAPI() error {
	r := gin.Default()
	r.GET("/", public.Hello)

	apiRoutes := r.Group("/api")

	authorized_routes := apiRoutes.Group("/public", middleware.AuthoriseHeader())
	authorized_routes.GET("/", public.Hello)

	return r.Run()
}

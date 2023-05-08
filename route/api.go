package route

import (
	"final/handler/public"
	"final/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAPI(address string) error {
	r := gin.Default()
	r.GET("/", public.Hello)

	publicRoutes := r.Group("/public")
	{
		publicRoutes.GET("/hello", public.Hello)
		publicRoutes.GET("/sellers", public.ListOfAllSellers)
		publicRoutes.GET("/sellers/:id/products", public.SellerProducts)
		publicRoutes.GET("/sellers/:id/scores", public.SellersScores)
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.GET("/login", middleware.Authorise(), public.Hello)
	}

	buyerRoutes := r.Group("/buyer", middleware.AuthoriseBuyer())
	{
		buyerRoutes.GET("/hello", public.Hello)
	}

	sellerRoutes := r.Group("/seller", middleware.AuthoriseSeller())
	{
		sellerRoutes.GET("/hello", public.Hello)
	}

	adminRoutes := r.Group("/admin", middleware.AuthoriseAdmin())
	{
		adminRoutes.GET("/hello", public.Hello)
	}

	return r.Run(address)
}

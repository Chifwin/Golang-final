package route

import (
	"golang-final/handler/buyer"
	"golang-final/handler/public"
	"golang-final/handler/seller"
	"golang-final/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAPI(address string) error {
	r := gin.Default()
	r.GET("/", public.Hello)

	publicRoutes := r.Group("/public")
	{
		publicRoutes.GET("/sellers", public.ListOfAllSellers)
		publicRoutes.GET("/sellers/:id/products", public.SellerProducts)
		publicRoutes.GET("/sellers/:id/comments", public.SellersComments)

		publicRoutes.GET("/products", public.ListOfAllProducts)
		publicRoutes.GET("/products/search", public.SearchProducts)
		publicRoutes.GET("/products/:id/sellers", public.ProductSellers)
		publicRoutes.GET("/products/:id/comments", public.ProductsComments)
	}

	authRoutes := r.Group("/auth")
	{
		authRoutes.GET("/login", middleware.Authorise(), public.Hello)
	}

	buyerRoutes := r.Group("/buyer", middleware.AuthoriseBuyer())
	{
		buyerRoutes.GET("/purchases", buyer.ListOfAllPurchases)
		buyerRoutes.POST("/purchases/", buyer.AddPurchases)

		buyerRoutes.GET("/comments", buyer.GetComments)
		buyerRoutes.POST("/comments/:id", buyer.AddComment)
		buyerRoutes.PUT("/comments/:id", buyer.UpdateComment)
		buyerRoutes.DELETE("/comments/:id", buyer.DeleteComment)
	}

	sellerRoutes := r.Group("/seller", middleware.AuthoriseSeller())
	{
		sellerRoutes.GET("/products", seller.Products)
		sellerRoutes.GET("/purchases", seller.Purchases)

		sellerRoutes.PUT("/products/:id", seller.UpdateProduct)
		sellerRoutes.DELETE("/products/:id", seller.DeleteProduct)
		sellerRoutes.PUT("products/:id/publish", seller.PublishProduct)
	}

	adminRoutes := r.Group("/admin", middleware.AuthoriseAdmin())
	{
		adminRoutes.GET("/hello", public.Hello)
	}

	return r.Run(address)
}

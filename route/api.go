package route

import (
	"golang-final/handler/admin"
	"golang-final/handler/auth"
	"golang-final/handler/buyer"
	"golang-final/handler/public"
	"golang-final/handler/seller"
	"golang-final/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupAPI(address string) error {
	if os.Getenv("RELEASE") == "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
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
		authRoutes.GET("", middleware.Authorise(), auth.GetUserInfo)
		authRoutes.POST("", auth.RegisterBuyer)
		authRoutes.PUT("", middleware.Authorise(), auth.UpdateUser)
		authRoutes.DELETE("", middleware.Authorise(), auth.DeleteUser)
	}

	buyerRoutes := r.Group("/buyer", middleware.AuthoriseBuyer())
	{
		buyerRoutes.GET("/purchases", buyer.ListOfAllPurchases)
		buyerRoutes.POST("/purchases/", buyer.Buy)

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
		adminRoutes.GET("/users", admin.ListUsers)
		adminRoutes.POST("/users/:role", admin.RegisterUser)
		adminRoutes.DELETE("/users/:id", admin.DeleteUser)

		adminRoutes.POST("/products", admin.AddProduct)
		adminRoutes.PUT("/products/:id", admin.UpdateProduct)

		adminRoutes.GET("/purchases", admin.ListLastPurchases)
	}

	return r.Run(address)
}

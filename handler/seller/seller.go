package seller

import (
	"database/sql"
	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Products(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	sellerID := user_info.ID

	products, err := db.SellerProducts(sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller products not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller products with error: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func Purchases(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	purchases, err := db.SellerPurchases(user_info.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get purchases",
		})
		return
	}

	ctx.JSON(http.StatusOK, purchases)
}

func UpdateProduct(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	sellerID := user_info.ID
	var product db.ProductFromSeller
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}

	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	product, err = db.UpdateSellerProduct(sellerID, productID, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Updating error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func DeleteProduct(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	product, err := db.DeleteSellerProduct(id, user_info.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Deleting error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func PublishProduct(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	product_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	var product db.ProductFromSeller
	query, have := ctx.GetQuery("value")
	value := (query == "true")
	product, err = db.PublishSellerProduct(product_id, user_info.ID, value, !have)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Publishing error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

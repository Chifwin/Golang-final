package seller

import (
	"database/sql"
	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func SellerProducts(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	sellerID := val.ID


	products, err := db.SellerProducts(sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller products not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller products"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}


func UpdateSellerProduct(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	sellerID := val.ID
	var product db.Products
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	value := ctx.GetHeader("value")

	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = db.UpdateSellerProductDB(sellerID, productID, product, value)
	
	ctx.JSON(http.StatusOK, gin.H{
		"error" : err,
	})
}

func DeleteSellerProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	val := ctx.Value("user_info").(db.UserRet)
	sellerID := val.ID

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = db.DeleteSellerProductDB(id, sellerID)

	ctx.JSON(http.StatusOK, gin.H{
		"error" : err,
	})
}

func PurchasesSeller(ctx *gin.Context){
	val := ctx.Value("user_info").(db.UserRet)
	sellerID := val.ID
	purchases, err := db.PurchasesSellerDB(sellerID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get purchases",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"purchases": purchases,
	})
}
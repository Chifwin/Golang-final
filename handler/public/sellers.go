package public

import (
	"database/sql"
	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListOfAllSellers(ctx *gin.Context) {
	sellers, err := db.ListSellers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sellers",
		})
		return
	}

	ctx.JSON(http.StatusOK, sellers)
}

func SellerProducts(ctx *gin.Context) {
	sellerID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

	products, err := db.SellerProducts(sellerID, false)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller products not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller products"})
		}
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func SellersComments(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

	comments, err := db.SellerComments(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller comments not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller comments"})
		}
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

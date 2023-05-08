package public

import (
	"database/sql"
	"final/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ListOfAllSellers(ctx *gin.Context) {
	sellers, err := db.ListSellers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sellers",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sellers": sellers,
	})
}

func SellerProducts(ctx *gin.Context) {
	sellerID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

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

func SellersScores(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

	scores, err := db.SellerScores(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller scores not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller scores"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"scores": scores,
	})
}

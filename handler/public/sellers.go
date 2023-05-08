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

func SellerInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

	seller, err := db.FindSellerById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seller"})
		}
		return
	}

	c.JSON(http.StatusOK, seller)
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

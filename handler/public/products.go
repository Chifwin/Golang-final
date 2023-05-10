package public

import (
	"database/sql"
	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListOfAllProducts(ctx *gin.Context) {
	products, err := db.ListProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get products with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func SearchProducts(ctx *gin.Context) {
	name, _ := ctx.GetQuery("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No name header provided",
		})
		return
	}
	products, err := db.SearchProduct(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get products with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func ProductSellers(ctx *gin.Context) {
	productID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

	sellers, err := db.ProductSellers(productID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product sellers not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product sellers"})
		}
		return
	}

	ctx.JSON(http.StatusOK, sellers)
}

func ProductsComments(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	comments, err := db.ProductComments(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product comments not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product comments"})
		}
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

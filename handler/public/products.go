package public

import (
	"database/sql"
	"golang-final/db"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

func ListOfAllProducts(ctx *gin.Context) {
	name := ctx.GetHeader("name")
	if(name == ""){
		products, err := db.ListProducts()
	
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get products",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	} else {
		products, err := db.ProductSearchDB(name)

		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Products not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	}
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

	ctx.JSON(http.StatusOK, gin.H{
		"sellers": sellers,
	})
}

func ProductsScores(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	scores, err := db.ProductScores(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product scores not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product scores"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"scores": scores,
	})
}

// func ProductSearch(ctx *gin.Context){
// 	name := ctx.GetHeader("name")

// 	products, err := db.ProductSearchDB(name)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": "Products not found"})
// 		} else {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
// 		}
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"products": products,
// 	})
// }

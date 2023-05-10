package admin

import (
	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddProduct(ctx *gin.Context) {
	var product db.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	product, err := db.AddProduct(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add product with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func UpdateProduct(ctx *gin.Context) {
	var product db.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	// Get the ID of the product to update from the URL parameter
	product_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err = db.UpdateProduct(product_id, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func ListUsers(ctx *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func ListLastPurchases(ctx *gin.Context) {
	purchases, err := db.LastPurchases()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get last purchases with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, purchases)
}

func RegisterUser(ctx *gin.Context) {
	var user db.UserCred
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	role, valid := db.ValidRole(ctx.Param("role"))
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role",
		})
		return
	}
	user.Role = role
	if err := db.AddUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add user with error: " + err.Error(),
		})
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Succesfulli created",
	})
}

func DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := db.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user with error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

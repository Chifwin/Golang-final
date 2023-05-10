package buyer

import (
	// "database/sql"

	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListOfAllPurchases(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	purchases, err := db.BuyerPurchases(user_info.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get purchases",
		})
		return
	}

	ctx.JSON(http.StatusOK, purchases)
}

func AddPurchases(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	var purchase db.Purchase
	if err := ctx.ShouldBindJSON(&purchase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	purchase.UserID = user_info.ID
	purchase, err := db.Buy(purchase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get purchases with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, purchase)
}

// Comment
func GetComments(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	comments, err := db.GetBuyerComments(user_info.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get comments with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func AddComment(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	var comment db.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	// Get the ID of the comment to update from the URL parameter
	purchase_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}
	res_comment, err := db.CreateComment(purchase_id, user_info.ID, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get comments with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res_comment)
}

func UpdateComment(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	var comment db.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}

	// Get the ID of the comment to update from the URL parameter
	purchase_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	res_comment, err := db.UpdateComment(purchase_id, user_info.ID, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update comment with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res_comment)
}

func DeleteComment(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := db.DeleteComment(id, user_info.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comment with error: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

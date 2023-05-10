package buyer

import (
	// "database/sql"

	"golang-final/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListOfAllPurchases(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	purchases, err := db.GetBuyerPurchases(val.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get purchases",
		})
		return
	}

	ctx.JSON(http.StatusOK, purchases)
}

func AddPurchases(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	var purchase db.Purchase
	if err := ctx.ShouldBindJSON(&purchase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	purchase.UserID = val.ID
	res_purchase, err := db.CreatePurchase(purchase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get purchases with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"purchase": res_purchase,
	})
}

// Comment
func GetComments(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	comments, err := db.GetBuyerComments(val.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get comments with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func AddComment(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
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
	res_comment, err := db.CreateComment(purchase_id, val.ID, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get comments with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res_comment)
}

func UpdateComment(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
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

	res_comment, err := db.UpdateComment(purchase_id, val.ID, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update comment with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, res_comment)
}

func DeleteComment(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := db.DeleteComment(id, val.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comment with error: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

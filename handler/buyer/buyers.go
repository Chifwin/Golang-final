package buyer

import (
	"database/sql"
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

	ctx.JSON(http.StatusOK, gin.H{
		"purchases": purchases,
	})
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


// Score
func GetComment(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	scores, err := db.GetBuyerComments(val.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get scores",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"scores": scores,
	})
}

func AddComment(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	var score db.Scores
	if err := ctx.ShouldBindJSON(&score); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	score.ProductId = val.ID
	res_score, err := db.CreateComment(score)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get scores with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"score": res_score,
	})
}

func UpdateComment(ctx *gin.Context) {
	val := ctx.Value("user_info").(db.UserRet)
	var score db.Scores
	if err := ctx.ShouldBindJSON(&score); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	score.ProductId = val.ID

	// Get the ID of the comment to update from the URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid comment ID",
        })
        return
    }

	res_score, err := db.UpdateCommentDB(id, score)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update comment with error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"score": res_score,
	})
}



func DeleteComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scores ID"})
		return
	}

	scores, err := db.DeleteCommentDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Seller scores not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get scores scores"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"scores": scores,
	})
}

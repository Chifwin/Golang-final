package auth

import (
	"golang-final/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterBuyer(ctx *gin.Context) {
	var user db.UserCred
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	user.Role = db.BUYER
	new_user, err := db.AddUser(user)
	if err != nil {
		if err == db.ErrUniqueFailed {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "User with this username already exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add user with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, new_user)
}

func GetUserInfo(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	user, err := db.GetUser(user_info.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user info with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

func UpdateUser(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	var user db.UserCred
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided",
		})
		return
	}
	new_user, err := db.UpdateUser(user_info.ID, user)
	if err != nil {
		if err == db.ErrUniqueFailed {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "User with this username already exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, new_user)
}

func DeleteUser(ctx *gin.Context) {
	user_info := ctx.Value("user_info").(db.UserRet)
	user, err := db.DeleteUser(user_info.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user with error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

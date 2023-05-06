package middleware

import (
	"final/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthoriseHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetHeader("username")
		password := ctx.GetHeader("password")
		if username == "" || password == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization headers found",
			})
			return
		}
		user_info, err := db.AuthoriseUser(username, password)
		if err != nil || user_info == nil {
			if err == db.ErrNotFound {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "No user with this credintails",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Internal error: " + err.Error(),
				})
			}
		}
		ctx.Set("user_info", *user_info)
	}
}

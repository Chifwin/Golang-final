package middleware

import (
	"database/sql"
	"golang-final/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

const anyRole db.UserRole = "any"

func authorizeWithRole(role db.UserRole) gin.HandlerFunc {
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
			if err == sql.ErrNoRows {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "No user with this credintails",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Internal error: " + err.Error(),
				})
			}
			return
		}
		if role != anyRole && role != user_info.Role {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User must have role: " + role,
			})
			return
		}
		ctx.Set("user_info", *user_info)
	}
}

func Authorise() gin.HandlerFunc {
	return authorizeWithRole(anyRole)
}

func AuthoriseAdmin() gin.HandlerFunc {
	return authorizeWithRole(db.ADMIN)
}

func AuthoriseSeller() gin.HandlerFunc {
	return authorizeWithRole(db.SELLER)
}

func AuthoriseBuyer() gin.HandlerFunc {
	return authorizeWithRole(db.BUYER)
}
